package biz

import (
	dto2 "brillinkmicros/internal/biz/dto"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type RcOriginContentRepo interface {
	GetPage(ctx context.Context, page *dto2.PaginationReq) (*dto2.RcOriginContentGetPageResp, error)
	GetInfos(ctx context.Context, page *dto2.PaginationReq) (*dto2.RcOriginContentInfosResp, error)
	Get(ctx context.Context, id int64) (*dto2.RcOriginContent, error)
}

type RcOriginContentUsecase struct {
	repo RcOriginContentRepo
	log  *log.Helper
}

func NewRcOriginContentUsecase(repo RcOriginContentRepo, logger log.Logger) *RcOriginContentUsecase {
	return &RcOriginContentUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *RcOriginContentUsecase) GetPage(ctx context.Context, page *dto2.PaginationReq) (*dto2.RcOriginContentGetPageResp, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.GetPage %d", page.PageNum)
	return uc.repo.GetPage(ctx, page)
}

func (uc *RcOriginContentUsecase) GetInfos(ctx context.Context, page *dto2.PaginationReq) (*dto2.RcOriginContentInfosResp, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.GetInfos %d", page.PageNum)
	return uc.repo.GetInfos(ctx, page)
}

func (uc *RcOriginContentUsecase) Get(ctx context.Context, id int64) (*dto2.RcOriginContent, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.Get %d", id)
	return uc.repo.Get(ctx, id)
}
