package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type RcOriginContent struct {
	Id             int64
	UscId          string
	EnterpriseName string
	YearMonth      string
	Content        string
	StatusCode     int
	UpdatedAt      time.Time
	CreatedAt      time.Time
}

func (*RcOriginContent) TableName() string {
	return "rskc_origin_content"
}

type RcOriginContentGetPageResp struct {
	PaginationResp
	Data *[]RcOriginContent
}

type RcOriginContentInfosResp struct {
	PaginationResp
	Data *[]RcOriginContentInfo
}

type RcOriginContentInfo struct {
	ContentId          int64
	UscId              string
	DataCollectMonth   string
	Content            string
	StatusCode         int
	EnterpriseName     string
	ProcessedId        int64
	ProcessedUpdatedAt time.Time
	UpdatedAt          time.Time
	CreatedAt          time.Time
}

type RcOriginContentRepo interface {
	GetPage(ctx context.Context, page *PaginationReq) (*RcOriginContentGetPageResp, error)
	GetInfos(ctx context.Context, page *PaginationReq) (*RcOriginContentInfosResp, error)
	Get(ctx context.Context, id int64) (*RcOriginContent, error)
}

type RcOriginContentUsecase struct {
	repo RcOriginContentRepo
	log  *log.Helper
}

func NewRcOriginContentUsecase(repo RcOriginContentRepo, logger log.Logger) *RcOriginContentUsecase {
	return &RcOriginContentUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *RcOriginContentUsecase) GetPage(ctx context.Context, page *PaginationReq) (*RcOriginContentGetPageResp, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.GetPage %d", page.PageNum)
	return uc.repo.GetPage(ctx, page)
}

func (uc *RcOriginContentUsecase) GetInfos(ctx context.Context, page *PaginationReq) (*RcOriginContentInfosResp, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.GetInfos %d", page.PageNum)
	return uc.repo.GetInfos(ctx, page)
}

func (uc *RcOriginContentUsecase) Get(ctx context.Context, id int64) (*RcOriginContent, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.Get %d", id)
	return uc.repo.Get(ctx, id)
}
