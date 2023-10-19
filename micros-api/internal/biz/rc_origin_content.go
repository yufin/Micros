package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-api/internal/biz/dto"
)

type RcOriginContentRepo interface {
	GetPage(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentGetPageResp, error)
	GetInfos(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentInfosResp, error)
	GetInfosByKwd(ctx context.Context, page *dto.PaginationReq, kwd string) (*dto.RcOriginContentInfosResp, error)
	Get(ctx context.Context, id int64) (*dto.RcOriginContent, error)
	CheckContentIdAllowed(ctx context.Context, contentId int64) (bool, error)

	GetContentIdsByUscId(ctx context.Context, uscId string) ([]int64, error)
}

type RcOriginContentUsecase struct {
	repo RcOriginContentRepo
	log  *log.Helper
}

func NewRcOriginContentUsecase(repo RcOriginContentRepo, logger log.Logger) *RcOriginContentUsecase {
	return &RcOriginContentUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *RcOriginContentUsecase) GetPage(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentGetPageResp, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.GetPage %d", page.PageNum)
	return uc.repo.GetPage(ctx, page)
}

func (uc *RcOriginContentUsecase) GetInfos(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentInfosResp, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.GetInfos %d", page.PageNum)
	return uc.repo.GetInfos(ctx, page)
}

func (uc *RcOriginContentUsecase) GetInfosByKwd(ctx context.Context, page *dto.PaginationReq, kwd string) (*dto.RcOriginContentInfosResp, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.GetInfosByKwd %d", page.PageNum)
	return uc.repo.GetInfosByKwd(ctx, page, kwd)
}

func (uc *RcOriginContentUsecase) Get(ctx context.Context, id int64) (*dto.RcOriginContent, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.Get %d", id)
	return uc.repo.Get(ctx, id)
}

func (uc *RcOriginContentUsecase) CheckContentIdAllowed(ctx context.Context, contentId int64) (bool, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.CheckContentIdAllowed %d", contentId)
	return uc.repo.CheckContentIdAllowed(ctx, contentId)
}

func (uc *RcOriginContentUsecase) GetContentIdsByUscId(ctx context.Context, uscId string) ([]int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.GetContentIdsByUscId %s", uscId)
	return uc.repo.GetContentIdsByUscId(ctx, uscId)
}
