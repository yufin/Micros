package workflow

import (
	"fmt"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"micros-worker/internal/biz"
	"micros-worker/internal/biz/dto"
	"micros-worker/internal/data"
	"time"
)

type ContentSyncWorkflow struct {
	uc       *biz.ContentUsecase
	dataRepo *data.Data
}

// NewContentSyncWorkflow new a ContentSyncWorkflow.
func NewContentSyncWorkflow(uc *biz.ContentUsecase, data *data.Data) *ContentSyncWorkflow {
	return &ContentSyncWorkflow{uc: uc, dataRepo: data}
}

func (s *ContentSyncWorkflow) SyncNewContent(ctx workflow.Context) ([]dto.RcContentMeta, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("SyncNewContent Workflow Invoked")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 1,
	})
	newFileMetas := make([]dto.SftpContentFileMeta, 0)
	logger.Info("GetNewContentFileMetaActivity start")
	err := workflow.ExecuteActivity(ctx, s.uc.GetNewContentFileMetaActivity).Get(ctx, &newFileMetas)
	if err != nil {
		return nil, err
	}

	if len(newFileMetas) > 0 {
		err = workflow.ExecuteActivity(ctx, s.uc.NoticeOnFoundNewContent, newFileMetas).Get(ctx, nil)
		if err != nil {
			logger.Error("NoticeOnFoundNewContent error", "error", err)
		}
	}

	var cwFutures []workflow.Future
	for _, nfm := range newFileMetas {
		cwo := workflow.ChildWorkflowOptions{
			ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			WorkflowID:        nfm.SourcePath,
		}
		cwCtx := workflow.WithChildOptions(ctx, cwo)
		logger.Info("ChildWorkflow ParseContent start", "sourcePath", nfm.SourcePath)
		future := workflow.ExecuteChildWorkflow(cwCtx, s.ParseContent, nfm)
		cwFutures = append(cwFutures, future)
	}

	resContentMeta := make([]dto.RcContentMeta, 0)
	cwNoticeFutures := make([]workflow.Future, 0)
	for _, f := range cwFutures {
		var resMeta dto.RcContentMeta
		err := f.Get(ctx, &resMeta)
		if err != nil {
			logger.Error("ChildWorkflow ParseContent error", "sourcePath", resMeta.SourcePath, "error", err)
		} else {
			resContentMeta = append(resContentMeta, resMeta)
			logger.Info("ChildWorkflow ParseContent complete", "sourcePath", resMeta.SourcePath, "contentId", resMeta.Id, "enterpriseName", resMeta.EnterpriseName)
			cwoNotice := workflow.ChildWorkflowOptions{
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
				WorkflowID:        fmt.Sprintf("Notice:%s|%s|%s|%v", resMeta.EnterpriseName, resMeta.AttributeMonth, resMeta.Version, resMeta.Id),
			}
			cwNoticeCtx := workflow.WithChildOptions(ctx, cwoNotice)
			futureNotice := workflow.ExecuteChildWorkflow(cwNoticeCtx, s.NoticeOnNewContentSynced, resMeta)
			cwNoticeFutures = append(cwNoticeFutures, futureNotice)
		}
	}

	metasDone := make([]dto.RcContentMeta, 0)
	for _, f := range cwNoticeFutures {
		var noticeRes dto.RcContentMeta
		err := f.Get(ctx, &noticeRes)
		if err != nil {
			logger.Error("ChildWorkflow NoticeOnNewContentSynced error", "error", err)
		} else {
			logger.Info("ChildWorkflow NoticeOnNewContentSynced complete", "contentId", noticeRes.Id)
			metasDone = append(metasDone, noticeRes)
		}
	}

	return metasDone, nil
}

func (s *ContentSyncWorkflow) ParseContent(ctx workflow.Context, fileMeta dto.SftpContentFileMeta) (dto.RcContentMeta, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ParseContent Workflow Invoked")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 1,
	})

	var meta dto.RcContentMeta
	logger.Info("InsertContentActivity start")
	err := workflow.ExecuteActivity(ctx, s.uc.InsertContentActivity, fileMeta).Get(ctx, &meta)
	if err != nil {
		return meta, err
	}

	childWorkflowId := func(rcm dto.RcContentMeta) string {
		return fmt.Sprintf("SyncDepd:%s|%s|%s|%v", rcm.EnterpriseName, rcm.AttributeMonth, rcm.Version, rcm.Id)
	}(meta)

	cwo := workflow.ChildWorkflowOptions{
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
		WorkflowID:        childWorkflowId,
	}

	cwCtx := workflow.WithChildOptions(ctx, cwo)
	logger.Info("ChildWorkflow SyncDependencyData start", "contentId", meta.Id)
	err = workflow.ExecuteChildWorkflow(cwCtx, s.SyncDependencyData, meta).Get(cwCtx, nil)
	if err != nil {
		return meta, errors.WithStack(err)
	}

	return meta, nil
}

func (s *ContentSyncWorkflow) SyncDependencyData(ctx workflow.Context, meta dto.RcContentMeta) (dto.RcContentMeta, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("SyncDependencyActivity start")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 1,
	})

	err := workflow.ExecuteActivity(ctx, s.uc.SyncDependencyActivity, meta.Id).Get(ctx, nil)
	if err != nil {
		return meta, errors.WithStack(err)
	}
	return meta, nil
}

func (s *ContentSyncWorkflow) NoticeOnNewContentSynced(ctx workflow.Context, meta dto.RcContentMeta) (dto.RcContentMeta, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("NoticeOnParseContentDone start")
	noticeRetryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Minute, // Set the initial retry interval to 1 minute
		BackoffCoefficient: 1.0,         // Keep the retry interval constant at 1 minute
		MaximumInterval:    time.Minute, // Maximum interval also set to 1 minute to ensure it doesn't increase
		MaximumAttempts:    300,         // Maximum number of retry attempts
	}
	ao := workflow.ActivityOptions{
		RetryPolicy:         noticeRetryPolicy,
		StartToCloseTimeout: time.Hour,
	}

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 1,
	})
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, ao),
		s.uc.NoticeOnParseContentDone,
		meta,
	).Get(ctx, nil)
	if err != nil {
		return meta, errors.WithStack(err)
	}
	return meta, nil
}

func (s *ContentSyncWorkflow) ReSyncDependencyData(ctx workflow.Context, req dto.PaginationReq) ([]dto.RcContentMeta, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ReSyncDependencyData Workflow Invoked")
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 1,
	})

	var metas []dto.RcContentMeta

	err := workflow.ExecuteActivity(ctx, s.uc.GetAllContentMetas, req).Get(ctx, &metas)
	if err != nil {
		return nil, err
	}

	var cwFutures []workflow.Future
	for _, meta := range metas {
		cwo := workflow.ChildWorkflowOptions{
			ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			WorkflowID:        meta.SourcePath,
		}
		cwCtx := workflow.WithChildOptions(ctx, cwo)
		future := workflow.ExecuteChildWorkflow(cwCtx, s.SyncDependencyData, meta)
		cwFutures = append(cwFutures, future)
	}

	resContentMeta := make([]dto.RcContentMeta, 0)
	for _, f := range cwFutures {
		var resMeta dto.RcContentMeta
		err := f.Get(ctx, &resMeta)
		if err != nil {
			logger.Error("ChildWorkflow ParseContent error", "sourcePath", resMeta.SourcePath, "error", err)
		} else {
			resContentMeta = append(resContentMeta, resMeta)
			logger.Info("ChildWorkflow ParseContent complete", "sourcePath", resMeta.SourcePath, "contentId", resMeta.Id, "enterpriseName", resMeta.EnterpriseName)
		}
	}

	return resContentMeta, nil
}

func (s *ContentSyncWorkflow) RegistryWorker() worker.Worker {
	w := worker.New(*s.dataRepo.TemporalClient.Client, "content-sync", worker.Options{})
	w.RegisterWorkflow(s.SyncNewContent)
	w.RegisterActivity(s.uc.GetNewContentFileMetaActivity)
	w.RegisterActivity(s.uc.NoticeOnFoundNewContent)

	w.RegisterWorkflow(s.ParseContent)
	w.RegisterWorkflow(s.NoticeOnNewContentSynced)
	w.RegisterActivity(s.uc.InsertContentActivity)

	w.RegisterWorkflow(s.SyncDependencyData)
	w.RegisterActivity(s.uc.SyncDependencyActivity)
	w.RegisterActivity(s.uc.NoticeOnParseContentDone)

	w.RegisterWorkflow(s.ReSyncDependencyData)
	w.RegisterActivity(s.uc.GetAllContentMetas)
	return w
}
