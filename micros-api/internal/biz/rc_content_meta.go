package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-api/internal/biz/dto"
)

type RcContentMetaRepo interface {
	Get(ctx context.Context, id int64) (*dto.RcContentMeta, error)
	GetContentIdsByUscId(ctx context.Context, uscId string) ([]int64, error)
}

type RcContentMetaUsecase struct {
	repo RcContentMetaRepo
	log  *log.Helper
}

func NewRcContentMetaUsecase(repo RcContentMetaRepo, logger log.Logger) *RcContentMetaUsecase {
	return &RcContentMetaUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *RcContentMetaUsecase) Get(ctx context.Context, id int64) (*dto.RcContentMeta, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.Get %d", id)
	return uc.repo.Get(ctx, id)
}

func (uc *RcContentMetaUsecase) GetContentIdsByUscId(ctx context.Context, uscId string) ([]int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.GetContentIdsByUscId %s", uscId)
	return uc.repo.GetContentIdsByUscId(ctx, uscId)
}
