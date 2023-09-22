package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-api/internal/biz/dto"
	"net/url"
)

type OssMetadataRepo interface {
	GetById(ctx context.Context, id int64) (*dto.OssMetaData, error)
	GetDownloadUrlByObjName(ctx context.Context, fileName string, metadata *dto.OssMetaData) (*url.URL, error)
}

type OssMetadataUsecase struct {
	repo OssMetadataRepo
	log  *log.Helper
}

func NewOssMetadataUsecase(repo OssMetadataRepo, logger log.Logger) *OssMetadataUsecase {
	return &OssMetadataUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *OssMetadataUsecase) GetById(ctx context.Context, id int64) (*dto.OssMetaData, error) {
	uc.log.WithContext(ctx).Infof("biz.OssMetadataUsecase.GetById %d", id)
	return uc.repo.GetById(ctx, id)
}

func (uc *OssMetadataUsecase) GetDownloadUrlByObjName(ctx context.Context, fileName string, metadata *dto.OssMetaData) (*url.URL, error) {
	uc.log.WithContext(ctx).Infof("biz.OssMetadataUsecase.GetDownloadUrlByObjName %v", metadata.Id)
	return uc.repo.GetDownloadUrlByObjName(ctx, fileName, metadata)
}
