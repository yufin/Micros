package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type RcReportOssRepo interface {
	GetOssIdUptoDateByDepId(ctx context.Context, depId int64) (int64, error)
}

type RcReportOssUsecase struct {
	repo RcReportOssRepo
	log  *log.Helper
}

func NewRcReportOssUsecase(repo RcReportOssRepo, logger log.Logger) *RcReportOssUsecase {
	return &RcReportOssUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *RcReportOssUsecase) GetOssIdUptoDateByDepId(ctx context.Context, depId int64) (int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcReportOssUsecase.GetOssIdByDepId %d", depId)
	return uc.repo.GetOssIdUptoDateByDepId(ctx, depId)
}
