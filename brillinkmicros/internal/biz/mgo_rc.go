package biz

import (
	"brillinkmicros/internal/biz/dto"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

type MgoRcRepo interface {
	GetProcessedContent(ctx context.Context, contentId int64) (bson.M, error)
	GetContentInfos(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentInfosRespV3, error)
}

type MgoRcUsecase struct {
	repo MgoRcRepo
	log  *log.Helper
}

func NewMgoRcUsecase(repo MgoRcRepo, logger log.Logger) *MgoRcUsecase {
	return &MgoRcUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *MgoRcUsecase) GetProcessedContent(ctx context.Context, contentId int64) (bson.M, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetProcessedContent %d", contentId)
	return uc.repo.GetProcessedContent(ctx, contentId)
}

func (uc *MgoRcUsecase) GetContentInfos(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentInfosRespV3, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetContentInfos %d", page.PageNum)
	return uc.repo.GetContentInfos(ctx, page)
}
