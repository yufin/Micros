package workflow

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"micros-worker/internal/biz"
	"micros-worker/internal/biz/dto"
	"micros-worker/internal/data"
	"time"
)

type CommonNoticeWorkflow struct {
	dataRepo *data.Data
	uc       *biz.CommonNoticeUsecase
}

func NewCommonNoticeWorkflow(uc *biz.CommonNoticeUsecase, data *data.Data) *CommonNoticeWorkflow {
	return &CommonNoticeWorkflow{
		dataRepo: data,
		uc:       uc,
	}
}

func (n *CommonNoticeWorkflow) PubNoticeByWechatBotMarkdown(ctx workflow.Context, req dto.PubNoticeByWechatBotMarkdownReq) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("PubNotice workflow started")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Hour * 1,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Minute,
			BackoffCoefficient: 1.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    300,
		},
	})

	logger.Info("SendWechatNoticeInMarkdown Activity start")
	err := workflow.ExecuteActivity(ctx, n.uc.SendWechatNoticeInMarkdown, req).Get(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (n *CommonNoticeWorkflow) RegistryWorker() worker.Worker {

	w := worker.New(*n.dataRepo.TemporalClient.InfraClient, "common-notice", worker.Options{})

	w.RegisterWorkflow(n.PubNoticeByWechatBotMarkdown)
	w.RegisterActivity(n.uc.SendWechatNoticeInMarkdown)

	return w
}
