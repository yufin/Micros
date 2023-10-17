package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-graph/internal/biz/dto"
)

type GraphNetRepo interface {
	GetChildren(ctx context.Context, sourceId int64, relTypeScope []string, labelScope []string, p *dto.PaginationReq) (*dto.Net, int64, error)
}

type GraphNetUsecase struct {
	repo GraphNetRepo
	log  *log.Helper
}

func NewGraphNetUsecase(repo GraphNetRepo, logger log.Logger) *GraphNetUsecase {
	return &GraphNetUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *GraphNetUsecase) GetChildren(ctx context.Context, sourceId int64, relTypeScope []string, labelScope []string, p *dto.PaginationReq) (*dto.Net, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphNetUsecase.GetChildren sourceId=%d, relTypeScope=%v, labelScope=%v", sourceId, relTypeScope, labelScope)
	return uc.repo.GetChildren(ctx, sourceId, relTypeScope, labelScope, p)
}
