package data

import (
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/biz/dto"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

func NewRcReportOssRepo(data *Data, logger log.Logger) biz.RcReportOssRepo {
	return &RcReportOssRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

type RcReportOssRepo struct {
	data *Data
	log  *log.Helper
}

func (repo *RcReportOssRepo) GetOssIdUptoDateByDepId(ctx context.Context, depId int64) (int64, error) {
	var reportOss *dto.RcReportOss
	err := repo.data.Db.Model(&reportOss).
		Where("dep_id = ?", depId).
		Order("created_at DESC").
		First(&reportOss).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return reportOss.OssId, nil
}
