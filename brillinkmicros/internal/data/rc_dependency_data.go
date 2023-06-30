package data

import (
	"brillinkmicros/common"
	"brillinkmicros/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"time"
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
	err := repo.data.Db.
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
	err := repo.data.Db.
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
	dsi, err := common.ParseBlDataScope(ctx)
	if err != nil {
		return 0, err
	}
	if dsi.UserId != 0 {
		insertReq.CreateBy = &dsi.UserId
	}

	var modelRdd biz.RcDependencyData
	insertReq.BaseModel.Gen()
	err = repo.data.Db.
		Table(modelRdd.TableName()).
		Create(&insertReq).
		Error
	if err != nil {
		return 0, err
	}
	return insertReq.Id, nil
}

func (repo *RcDependencyDataRepo) Update(ctx context.Context, updateReq *biz.RcDependencyData) (int64, error) {
	dsi, err := common.ParseBlDataScope(ctx)
	if err != nil {
		return 0, err
	}
	if dsi.UserId != 0 {
		updateReq.UpdateBy = &dsi.UserId
	}
	updatedAt := time.Now()
	updateReq.UpdatedAt = &updatedAt

	var modelRdd biz.RcDependencyData
	if updateReq.Id == 0 {
		return 0, errors.BadRequest("Empty Id", "id is required")
	}
	err = repo.data.Db.
		Table(modelRdd.TableName()).
		Updates(&updateReq).Error
	if err != nil {
		return 0, err
	}
	return updateReq.Id, nil
}

func (repo *RcDependencyDataRepo) Delete(ctx context.Context, id int64) (bool, error) {
	var modelRdd biz.RcDependencyData
	if id == 0 {
		return false, errors.BadRequest("Empty Id", "id is required")
	}
	err := repo.data.Db.
		Table(modelRdd.TableName()).
		Delete(&modelRdd, id).
		Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *RcDependencyDataRepo) CheckIsInsertDepdDataDuplicate(ctx context.Context, uscId string) (bool, error) {
	dsi, err := common.ParseBlDataScope(ctx)
	if err != nil {
		return false, err
	}
	var model biz.RcDependencyData
	var count int64
	err = repo.data.Db.
		Table(model.TableName()).
		Where("usc_id = ?", uscId).
		Where("create_by = ?", dsi.UserId).
		Where("content_id is null").
		Count(&count).
		Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *RcDependencyDataRepo) GetDefaultContentIdForInsertDependencyData(ctx context.Context, uscId string) ([]int64, error) {
	dsi, err := common.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	err = repo.data.Db.
		Raw(`select content_id
				from (select roc.id as content_id, rdd.content_id as cid
					  from rskc_origin_content roc
							   left join (select content_id, id, usc_id from rc_dependency_data where create_by = ?) rdd
										 on roc.id = rdd.content_id and roc.usc_id = rdd.usc_id
					  where roc.usc_id = ?) t
				where t.cid is null`, dsi.UserId, uscId).
		Pluck("content_id", &ids).Error
	if err != nil {
		return []int64{}, err
	}
	return ids, nil
}
