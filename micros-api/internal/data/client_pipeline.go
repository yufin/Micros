package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pipelineV1 "micros-api/api/pipeline/v1"
	"micros-api/internal/biz"
)

type ClientPipelineRepo struct {
	data *Data
	log  *log.Helper
}

func NewClientPipelineRepo(data *Data, logger log.Logger) biz.ClientPipelineRepo {
	return &ClientPipelineRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *ClientPipelineRepo) GetClient(ctx context.Context) pipelineV1.PipelineServiceClient {
	return repo.data.PipelineClient
}
