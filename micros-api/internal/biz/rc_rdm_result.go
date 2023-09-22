package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-api/internal/biz/dto"
)

type RcRdmResultRepo interface {
	GetUpToDate(ctx context.Context, depId int64) (*dto.RcRdmResult, error)
}

type RcRdmResultUsecase struct {
	repo RcRdmResultRepo
	log  *log.Helper
}

func NewRcRdmResultUsecase(repo RcRdmResultRepo, logger log.Logger) *RcRdmResultUsecase {
	return &RcRdmResultUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *RcRdmResultUsecase) GetUpToDate(ctx context.Context, depId int64) (*dto.RcRdmResult, error) {
	uc.log.WithContext(ctx).Infof("biz.RcRdmResultUsecase.Get %d", depId)
	return uc.repo.GetUpToDate(ctx, depId)
}
