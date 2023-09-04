package biz

import (
	"brillinkmicros/internal/biz/dto"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

type MgoRcRepo interface {
	GetProcessedObjIdByContentId(ctx context.Context, contentId int64) (bson.M, error)
	GetProcessedContentByObjId(ctx context.Context, objIdHex string) (bson.M, error)
	GetProcessedContentInfoByObjId(ctx context.Context, objIdHex string) (bson.M, error)
	GetContentInfos(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentInfosRespV3, error)
	GetContentInfosByKwd(ctx context.Context, page *dto.PaginationReq, kwd string) (*dto.RcOriginContentInfosRespV3, error)
}

type MgoRcUsecase struct {
	repo MgoRcRepo
	log  *log.Helper
}

func NewMgoRcUsecase(repo MgoRcRepo, logger log.Logger) *MgoRcUsecase {
	return &MgoRcUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *MgoRcUsecase) GetProcessedObjIdByContentId(ctx context.Context, contentId int64) (bson.M, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetProcessedObjIdByContentId %d", contentId)
	return uc.repo.GetProcessedObjIdByContentId(ctx, contentId)
}

func (uc *MgoRcUsecase) GetContentInfos(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentInfosRespV3, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetContentInfos %d", page.PageNum)
	return uc.repo.GetContentInfos(ctx, page)
}

func (uc *MgoRcUsecase) GetProcessedContentByObjId(ctx context.Context, objIdHex string) (bson.M, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetProcessedContentByObjId %s", objIdHex)
	return uc.repo.GetProcessedContentByObjId(ctx, objIdHex)
}

func (uc *MgoRcUsecase) GetProcessedContentInfoByObjId(ctx context.Context, objIdHex string) (bson.M, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetProcessedContentInfoByObjId %s", objIdHex)
	return uc.repo.GetProcessedContentInfoByObjId(ctx, objIdHex)
}

func (uc *MgoRcUsecase) GetContentInfosByKwd(ctx context.Context, page *dto.PaginationReq, kwd string) (*dto.RcOriginContentInfosRespV3, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.SearchReportInfosByKwd %d", page.PageNum)
	return uc.repo.GetContentInfosByKwd(ctx, page, kwd)
}
