package data

import (
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/biz/dto"
	"brillinkmicros/pkg"
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

func (repo *RcDependencyDataRepo) GetByContentId(ctx context.Context, contentId int64) (*dto.RcDependencyData, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}
	var modelRdd dto.RcDependencyData
	err = repo.data.Db.
		Raw(
			`WITH user_records AS (SELECT *, 1 AS sort_priority
								  FROM rc_dependency_data
								  WHERE content_id = ?
									AND create_by = ?
									and deleted_at is null),
				 valid_records AS (SELECT *, 2 AS sort_priority
								   FROM rc_dependency_data
								   WHERE content_id = ?
									 AND create_by IN (?)
									 and deleted_at is null),
				 combined_records AS (SELECT *
									  FROM user_records
									  UNION ALL
									  SELECT *
									  FROM valid_records
									  WHERE NOT EXISTS (SELECT 1 FROM user_records))
			SELECT *
			FROM combined_records
			ORDER BY sort_priority,
					CASE WHEN sort_priority = 1 THEN created_at END DESC, 
					CASE WHEN sort_priority = 2 THEN created_at END ASC
			limit 1;`, contentId, dsi.UserId, contentId, dsi.AccessibleIds).
		Scan(&modelRdd).Error
	if err != nil {
		return nil, err
	}
	return &modelRdd, nil
}

func (repo *RcDependencyDataRepo) Get(ctx context.Context, id int64) (*dto.RcDependencyData, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}
	var modelRdd dto.RcDependencyData
	err = repo.data.Db.
		Table(modelRdd.TableName()).
		Where("id = ?", id).
		Where("create_by IN (?)", dsi.AccessibleIds).
		First(&modelRdd).
		Error
	if err != nil {
		return nil, err
	}
	return &modelRdd, nil
}

func (repo *RcDependencyDataRepo) Insert(ctx context.Context, insertReq *dto.RcDependencyData) (int64, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return 0, err
	}
	if dsi.UserId != 0 {
		insertReq.CreateBy = &dsi.UserId
	}

	var modelRdd dto.RcDependencyData
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

func (repo *RcDependencyDataRepo) Update(ctx context.Context, updateReq *dto.RcDependencyData) (int64, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return 0, err
	}
	if dsi.UserId != 0 {
		updateReq.UpdateBy = &dsi.UserId
	}
	//updatedAt := time.Now()
	//updateReq.UpdatedAt = &updatedAt
	var modelRdd dto.RcDependencyData
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
	var modelRdd dto.RcDependencyData
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
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return false, err
	}
	var model dto.RcDependencyData
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
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	err = repo.data.Db.
		Raw(`select content_id
				from (select roc.id as content_id, rdd.content_id as cid
					  from rc_origin_content roc
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
