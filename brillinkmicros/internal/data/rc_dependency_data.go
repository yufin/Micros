package data

import (
	"brillinkmicros/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type RcDependencyDataRepo struct {
	data *Data
	log  *log.Helper
}

func NewRcDependencyDataRepo(data *Data, logger log.Logger) biz.RcDependencyDataRepo {
	return &RcDependencyDataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *RcDependencyDataRepo) GetByContentId(ctx context.Context, contentId int64) (*biz.RcDependencyData, error) {
	var modelRdd biz.RcDependencyData
	err := repo.data.db.
		Table(modelRdd.TableName()).
		Where("content_id = ?", contentId).
		Order("updated_at desc").
		First(&modelRdd).
		Error
	if err != nil {
		return nil, err
	}
	return &modelRdd, nil
}

func (repo *RcDependencyDataRepo) Get(ctx context.Context, id int64) (*biz.RcDependencyData, error) {
	var modelRdd biz.RcDependencyData
	err := repo.data.db.
		Table(modelRdd.TableName()).
		Where("id = ?", id).
		First(&modelRdd).
		Error
	if err != nil {
		return nil, err
	}
	return &modelRdd, nil
}

func (repo *RcDependencyDataRepo) Insert(ctx context.Context, insertReq *biz.RcDependencyData) (int64, error) {
	var modelRdd biz.RcDependencyData
	insertReq.BaseModel.Gen()
	err := repo.data.db.
		Table(modelRdd.TableName()).
		Create(&insertReq).
		Error
	if err != nil {
		return 0, err
	}
	return insertReq.Id, nil
}

func (repo *RcDependencyDataRepo) Update(ctx context.Context, updateReq *biz.RcDependencyData) (int64, error) {
	var modelRdd biz.RcDependencyData
	if updateReq.Id == 0 {
		return 0, errors.BadRequest("Empty Id", "id is required")
	}
	err := repo.data.db.
		Table(modelRdd.TableName()).
		Updates(&updateReq).Error
	if err != nil {
		return 0, err
	}

	return updateReq.Id, nil
}

func (repo *RcDependencyDataRepo) Delete(ctx context.Context, id int64) error {
	var modelRdd biz.RcDependencyData
	if id == 0 {
		return errors.BadRequest("Empty Id", "id is required")
	}
	err := repo.data.db.
		Table(modelRdd.TableName()).
		Delete(&modelRdd, id).
		Error
	if err != nil {
		return err
	}
	return nil
}
