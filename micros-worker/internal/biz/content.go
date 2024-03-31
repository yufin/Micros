package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.temporal.io/sdk/activity"
	"micros-worker/internal/biz/dto"
	"micros-worker/internal/conf"
	"micros-worker/pkg"
	"strings"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

// ContentRepo is a Greater repo.
type ContentRepo interface {
	InsertRawContent(ctx context.Context, meta *dto.RcContentMeta, content *map[string]any) (*mongo.InsertOneResult, error)
	InsertContentMeta(ctx context.Context, meta *dto.RcContentMeta) error
	CountContentMetaByPath(ctx context.Context, sftpPath string) (int64, error)
	CountRawContentByPath(ctx context.Context, sftpPath string) (int64, error)
	GetContentMetasByPath(ctx context.Context, sftpPath string) ([]dto.RcContentMeta, error)
	ReadContentByPath(ctx context.Context, path string) (*[]byte, error)
	PubPostLoanFileInfo(ctx context.Context, metaCh *chan dto.SftpContentFileMeta) error
	PubPreLoanFileMeta(ctx context.Context, metaCh *chan dto.SftpContentFileMeta) error
	SaveContentTransaction(ctx context.Context, meta *dto.RcContentMeta, content *map[string]any) error
	DeleteRawContentById(ctx context.Context, id primitive.ObjectID) error
	SaveRiskIndexesIdempotent(ctx context.Context, req []dto.RcRiskIndex) error
	GetRawContentByContentId(ctx context.Context, contentId int64) (*map[string]any, error)
	GetConf() *conf.Data
	DeleteRiskIndexByContentId(ctx context.Context, contentId int64) error
	GetRcContentMetas(ctx context.Context, pageReq dto.PaginationReq) ([]dto.RcContentMeta, error)
}

// ContentUsecase is a Greeter usecase.
type ContentUsecase struct {
	repo ContentRepo
	log  *log.Helper
}

// NewContentUsecase new a Greeter usecase.
func NewContentUsecase(repo ContentRepo, logger log.Logger) *ContentUsecase {
	return &ContentUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *ContentUsecase) GetNewContentFileMetaActivity(ctx context.Context) ([]dto.SftpContentFileMeta, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("start SaveNewContentActivity")
	var (
		metaCh   = make(chan dto.SftpContentFileMeta, 200)
		errCh    = make(chan error)
		done     = make(chan struct{})
		newMetas = make([]dto.SftpContentFileMeta, 0)
		wg       sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(metaCh)
		err := uc.repo.PubPreLoanFileMeta(ctx, &metaCh)
		if err != nil {
			errCh <- err
		}
		err = uc.repo.PubPostLoanFileInfo(ctx, &metaCh)
		if err != nil {
			errCh <- err
		}
		return
	}()

	wg.Add(1)
	go func(newM *[]dto.SftpContentFileMeta) {
		defer wg.Done()
		for m := range metaCh {
			c, err := uc.repo.CountContentMetaByPath(ctx, m.SourcePath)
			if err != nil {
				errCh <- err
				return
			}
			if c == 0 {
				*newM = append(*newM, m)
			}
		}
	}(&newMetas)

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return newMetas, nil
	case err := <-errCh:
		return newMetas, err
	}
}

func (uc *ContentUsecase) GetAllContentMetas(ctx context.Context, pageReq dto.PaginationReq) ([]dto.RcContentMeta, error) {
	return uc.repo.GetRcContentMetas(ctx, pageReq)
}

func (uc *ContentUsecase) SyncDependencyActivity(ctx context.Context, contentId int64) error {
	content, err := uc.repo.GetRawContentByContentId(ctx, contentId)
	if err != nil {
		return errors.WithStack(err)
	}

	// parse riskIndex and save
	if err = uc.repo.DeleteRiskIndexByContentId(ctx, contentId); err != nil {
		return errors.WithStack(err)
	}
	if err = uc.parseRiskIndexFromContentAndSave(ctx, content, contentId); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (uc *ContentUsecase) deleteRiskIndexByContentId(ctx context.Context, contentId int64) error {
	err := uc.repo.DeleteRiskIndexByContentId(ctx, contentId)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (uc *ContentUsecase) parseRiskIndexFromContentAndSave(ctx context.Context, content *map[string]any, contentId int64) error {
	riskIndexes, ok := (*content)["impExpEntReport"].(map[string]any)["riskIndexes"].([]any)
	if !ok {
		return errors.WithStack(errors.New("riskIndexes not found"))
	}

	idxs := make([]dto.RcRiskIndex, 0)

	for _, v := range riskIndexes {
		v := v.(map[string]any)
		insertReq := dto.RcRiskIndex{
			ContentId:  contentId,
			RiskDec:    safeGetString(v, "RISK_DEC"),
			IndexDec:   safeGetString(v, "INDEX_DEC"),
			IndexValue: safeGetString(v, "INDEX_VALUE"),
			IndexFlag:  safeGetString(v, "INDEX_FLAG"),
		}
		insertReq.Gen()
		idxs = append(idxs, insertReq)
	}
	err := uc.repo.SaveRiskIndexesIdempotent(ctx, idxs)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (uc *ContentUsecase) InsertContentActivity(ctx context.Context, fileMeta dto.SftpContentFileMeta) (dto.RcContentMeta, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("start InsertContentActivity", "sourcePath", fileMeta.SourcePath)

	b, err := uc.repo.ReadContentByPath(ctx, fileMeta.SourcePath)
	if err != nil {
		return dto.RcContentMeta{}, err
	}
	contentMap := make(map[string]any)
	if err := json.Unmarshal(*b, &contentMap); err != nil {
		return dto.RcContentMeta{}, err
	}

	entName, err := uc.parseEnterpriseName(ctx, &contentMap)
	if err != nil {
		return dto.RcContentMeta{}, err
	}
	version := uc.parseVersion(fileMeta.SourcePath)

	meta := dto.RcContentMeta{
		UscId:          fileMeta.UscId,
		EnterpriseName: entName,
		AttributeMonth: fileMeta.AttributeMonth,
		SourcePath:     fileMeta.SourcePath,
		Version:        version,
	}
	meta.Gen()

	err = uc.repo.SaveContentTransaction(ctx, &meta, &contentMap)
	if err != nil {
		return meta, err
	}
	logger.Info("InsertContentActivity success", "contentId", meta.BaseModel.Id)
	return meta, nil
}

func (uc *ContentUsecase) parseEnterpriseName(ctx context.Context, content *map[string]any) (string, error) {
	enterpriseName := (*content)["impExpEntReport"].(map[string]any)["businessInfo"].(map[string]any)["QYMC"]
	// determine enterprise name is nil or string
	if _, ok := enterpriseName.(string); ok {
		return enterpriseName.(string), nil
	}
	return "", errors.New("enterpriseName parse failed")
}

func (uc *ContentUsecase) parseVersion(path string) string {
	if strings.HasSuffix(path, "_V2.0.json") {
		return "V2.0"
	}
	return "V1.0"
}

func (uc *ContentUsecase) NoticeOnParseContentDone(ctx context.Context, meta dto.RcContentMeta) error {
	md := fmt.Sprintf(
		`<font color="warning">新的风控报告数据完成入库</font>:
		>env: <font color="comment">%s</font>
		>企业名称: <font color="comment">%s</font>
		>企业统一社会信用代码: <font color="comment">%s</font>
		>数据更新月份: <font color="comment">%s</font>
		>录入时间: <font color="comment">%s</font>
		>contentId: <font color="comment">%d</font>
		>sftpPath: <font color="comment">%s</font>`,
		uc.repo.GetConf().WechatBot.EnvName,
		meta.EnterpriseName,
		meta.UscId,
		meta.AttributeMonth,
		meta.CreatedAt.Format("2006-01-02 15:04:05"),
		meta.Id,
		meta.SourcePath,
	)

	wb := pkg.WechatBot{Hook: uc.repo.GetConf().WechatBot.HookContentSync}

	err := wb.NoticeInMarkdown(md)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (uc *ContentUsecase) NoticeOnFoundNewContent(ctx context.Context, fileMetas []dto.SftpContentFileMeta) error {
	md := fmt.Sprintf("检索到新的报文<font color=\"warning\">%v例</font>，请关注。(env:%s)", len(fileMetas), uc.repo.GetConf().WechatBot.EnvName)
	wb := pkg.WechatBot{Hook: uc.repo.GetConf().WechatBot.HookContentSync}
	err := wb.NoticeInMarkdown(md)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func safeGetString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		return v.(string)
	}
	return ""
}
