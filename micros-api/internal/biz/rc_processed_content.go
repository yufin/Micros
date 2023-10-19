package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-api/internal/biz/dto"
	"time"
)

type RcProcessedContentRepo interface {
	Get(ctx context.Context, id int64) (*dto.RcProcessedContent, error)
	GetByContentIdUpToDate(ctx context.Context, contentId int64) (*dto.RcProcessedContent, error)
	RefreshReportContent(ctx context.Context, contentId int64) (bool, error)
	GetContentUpToDateByDepId(ctx context.Context, depId int64, allowedUserId int64) (*dto.RcProcessedContent, error)
	GetNewestRowInfoByContentId(ctx context.Context, contentId int64) (int64, time.Time, error)
}

type RcProcessedContentUsecase struct {
	repo RcProcessedContentRepo
	log  *log.Helper
}

func NewRcProcessedContentUsecase(repo RcProcessedContentRepo, logger log.Logger) *RcProcessedContentUsecase {
	return &RcProcessedContentUsecase{repo: repo, log: log.NewHelper(logger)}
}

// GetById .
// 使用RcProcessedContentRepo中定义的方法实现具体业务
func (uc *RcProcessedContentUsecase) GetById(ctx context.Context, id int64) (*dto.RcProcessedContent, error) {
	uc.log.WithContext(ctx).Infof("biz.GetById %d", id)
	return uc.repo.Get(ctx, id)
}

// GetByContentIdUpToDate .
// 使用RcProcessedContentRepo中定义的方法实现具体业务
func (uc *RcProcessedContentUsecase) GetByContentIdUpToDate(ctx context.Context, contentId int64) (*dto.RcProcessedContent, error) {
	uc.log.WithContext(ctx).Infof("biz.GetByContentIdUpToDate %v", contentId)
	return uc.repo.GetByContentIdUpToDate(ctx, contentId)
}

// GetContentUpToDateByDepId .
// 使用RcProcessedContentRepo中定义的方法实现具体业务
func (uc *RcProcessedContentUsecase) GetContentUpToDateByDepId(ctx context.Context, depId int64, allowedUserId int64) (*dto.RcProcessedContent, error) {
	uc.log.WithContext(ctx).Infof("biz.GetContentUpToDateByDepId %v", depId)
	return uc.repo.GetContentUpToDateByDepId(ctx, depId, allowedUserId)
}

// RefreshReportContent .
func (uc *RcProcessedContentUsecase) RefreshReportContent(ctx context.Context, contentId int64) (bool, error) {
	uc.log.WithContext(ctx).Infof("biz.RefreshReportContent %v", contentId)
	return uc.repo.RefreshReportContent(ctx, contentId)
}

// GetNewestRowInfoByContentId .
func (uc *RcProcessedContentUsecase) GetNewestRowInfoByContentId(ctx context.Context, contentId int64) (int64, time.Time, error) {
	uc.log.WithContext(ctx).Infof("biz.GetNewestProcessedIdByContentId %v", contentId)
	return uc.repo.GetNewestRowInfoByContentId(ctx, contentId)
}
