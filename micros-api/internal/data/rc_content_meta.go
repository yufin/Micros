package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-api/internal/biz"
	"micros-api/internal/biz/dto"
)

type RcContentMetaRepo struct {
	data *Data
	log  *log.Helper
}

func NewRcContentMetaRepo(data *Data, logger log.Logger) biz.RcContentMetaRepo {
	return &RcContentMetaRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *RcContentMetaRepo) Get(ctx context.Context, id int64) (*dto.RcContentMeta, error) {
	var modelRcm dto.RcContentMeta
	err := repo.data.Db.Table(modelRcm.TableName()).Where("id = ?", id).First(&modelRcm).Error
	if err != nil {
		return nil, err
	}
	return &modelRcm, nil
}

func (repo *RcContentMetaRepo) GetContentIdsByUscId(ctx context.Context, uscId string) ([]int64, error) {
	var ids []int64
	err := repo.data.Db.
		Model(&dto.RcContentMeta{}).
		Select("id").
		Where("usc_id = ? ", uscId).
		Pluck("id", &ids).
		Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}
