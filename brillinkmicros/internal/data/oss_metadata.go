package data

import (
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/biz/dto"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"net/url"
	"time"
)

type OssMetadataRepo struct {
	data *Data
	log  *log.Helper
}

func NewOssMetadataRepo(data *Data, logger log.Logger) biz.OssMetadataRepo {
	return &OssMetadataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *OssMetadataRepo) GetById(ctx context.Context, id int64) (*dto.OssMetaData, error) {
	var oss *dto.OssMetaData
	err := repo.data.Db.Model(&oss).
		First(&oss, id).
		Error
	if err != nil {
		return nil, err
	}
	return oss, nil
}

func (repo *OssMetadataRepo) GetDownloadUrlByObjName(ctx context.Context, fileName string, metadata *dto.OssMetaData) (*url.URL, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	preUrl, err := repo.data.MinioCli.Cli.PresignedGetObject(ctx, metadata.BucketName, metadata.ObjName, time.Second*120, reqParams)
	if err != nil {
		return nil, err
	}
	return preUrl, nil
}
