package biz

import (
	"brillinkmicros/internal/biz/dto"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type RcRdmResDetailRepo interface {
	GetByResId(ctx context.Context, resId int64) (*[]dto.RcRdmResDetail, error)
	GetByResIdAndLevel(ctx context.Context, resId int64, level int) (*[]dto.RcRdmResDetail, error)
}

type RcRdmResDetailUsecase struct {
	repo RcRdmResDetailRepo
	log  *log.Helper
}

func NewRcRdmResDetailUsecase(repo RcRdmResDetailRepo, logger log.Logger) *RcRdmResDetailUsecase {
	return &RcRdmResDetailUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *RcRdmResDetailUsecase) GetByResId(ctx context.Context, resId int64) (*[]dto.RcRdmResDetail, error) {
	uc.log.WithContext(ctx).Infof("biz.RcRdmResDetailUsecase.Get %d", resId)
	return uc.repo.GetByResId(ctx, resId)
}

func (uc *RcRdmResDetailUsecase) GetByResIdAndLevel(ctx context.Context, resId int64, level int) (*[]dto.RcRdmResDetail, error) {
	uc.log.WithContext(ctx).Infof("biz.RcRdmResDetailUsecase.Get %d", resId)
	return uc.repo.GetByResIdAndLevel(ctx, resId, level)
}
