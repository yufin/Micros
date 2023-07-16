package biz

import (
	"brillinkmicros/internal/biz/dto"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type RcProcessedContentRepo interface {
	Get(ctx context.Context, id int64) (*dto.RcProcessedContent, error)
	GetByContentIdUpToDate(ctx context.Context, contentId int64) (*dto.RcProcessedContent, error)
	RefreshReportContent(ctx context.Context, contentId int64) (bool, error)
	GetByContentIdUpToDateByUser(ctx context.Context, contentId int64, userId int64, allowedUserId int64) (*dto.RcProcessedContent, error)
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

// GetByContentIdUpToDateByUser .
// 使用RcProcessedContentRepo中定义的方法实现具体业务
func (uc *RcProcessedContentUsecase) GetByContentIdUpToDateByUser(ctx context.Context, contentId int64, userId int64, allowedUserId int64) (*dto.RcProcessedContent, error) {
	uc.log.WithContext(ctx).Infof("biz.GetByContentIdUpToDateByUser %v", contentId)
	return uc.repo.GetByContentIdUpToDateByUser(ctx, contentId, userId, allowedUserId)
}

// RefreshReportContent .
func (uc *RcProcessedContentUsecase) RefreshReportContent(ctx context.Context, contentId int64) (bool, error) {
	uc.log.WithContext(ctx).Infof("biz.RefreshReportContent %v", contentId)
	return uc.repo.RefreshReportContent(ctx, contentId)
}
