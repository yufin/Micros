package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pipelineV1 "micros-api/api/pipeline/v1"
)

type ClientPipelineRepo interface {
	GetClient(ctx context.Context) pipelineV1.PipelineServiceClient
}

type ClientPipelineUsecase struct {
	repo ClientPipelineRepo
	log  *log.Helper
}

func NewClientPipelineUsecase(repo ClientPipelineRepo, logger log.Logger) *ClientPipelineUsecase {
	return &ClientPipelineUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *ClientPipelineUsecase) GetClient(ctx context.Context) pipelineV1.PipelineServiceClient {
	//uc.log.WithContext(ctx).Infof("biz.ClientPipelineUsecase.GetClient")
	return uc.repo.GetClient(ctx)
}
