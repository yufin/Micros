package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type RcProcessedContent struct {
	Id        int64
	ContentId int64
	Content   string
	UpdatedAt time.Time
	CreatedAt time.Time
}

func (*RcProcessedContent) TableName() string {
	return "rskc_processed_content"
}

type RcProcessedContentRepo interface {
	Get(ctx context.Context, id int64) (*RcProcessedContent, error)
	GetByContentIdUpToDate(ctx context.Context, contentId int64) (*RcProcessedContent, error)
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
func (uc *RcProcessedContentUsecase) GetById(ctx context.Context, id int64) (*RcProcessedContent, error) {
	uc.log.WithContext(ctx).Infof("biz.GetById %d", id)
	return uc.repo.Get(ctx, id)
}

// GetByContentIdUpToDate .
// 使用RcProcessedContentRepo中定义的方法实现具体业务
func (uc *RcProcessedContentUsecase) GetByContentIdUpToDate(ctx context.Context, contentId int64) (*RcProcessedContent, error) {
	uc.log.WithContext(ctx).Infof("biz.GetList %v", contentId)
	return uc.repo.GetByContentIdUpToDate(ctx, contentId)
}
