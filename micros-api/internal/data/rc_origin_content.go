package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"math"
	"micros-api/internal/biz"
	"micros-api/internal/biz/dto"
	"micros-api/pkg"
)

type RcOriginContentRepo struct {
	data *Data
	log  *log.Helper
}

func NewRcOriginContentRepo(data *Data, logger log.Logger) biz.RcOriginContentRepo {
	return &RcOriginContentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *RcOriginContentRepo) GetPage(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentGetPageResp, error) {
	var modelRoc dto.RcOriginContent
	listRoc := make([]dto.RcOriginContent, 0)
	var count int64
	offset := (page.PageNum - 1) * page.PageSize
	err := repo.data.Db.
		Table(modelRoc.TableName()).
		Find(&listRoc).
		Offset(offset).Limit(page.PageSize).
		Count(&count).
		Error
	if err != nil {
		return nil, err
	}
	return &dto.RcOriginContentGetPageResp{
		PaginationResp: dto.PaginationResp{
			Total:     count,
			TotalPage: int(math.Ceil(float64(count) / float64(page.PageSize))),
			PageNum:   page.PageNum,
			PageSize:  page.PageSize,
		},
		Data: &listRoc,
	}, nil
}

// CheckContentIdAllowed .
func (repo *RcOriginContentRepo) CheckContentIdAllowed(ctx context.Context, contentId int64) (bool, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return false, err
	}
	var allowedDid int64
	err = repo.data.Db.Raw(
		`select rdd.content_id
				from rc_dependency_data rdd
						 INNER Join rc_origin_content roc on rdd.content_id = roc.id
				where rdd.create_by in ?
				  and content_id = ?
				order by rdd.created_at desc`, dsi.AccessibleIds, contentId,
	).First(&allowedDid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (repo *RcOriginContentRepo) GetInfos(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentInfosResp, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}

	pageNum := int(math.Max(1, float64(page.PageNum)))
	offset := (pageNum - 1) * page.PageSize
	var count int64
	if err := repo.data.Db.
		Raw(`WITH rpc_cte AS (SELECT content_id, MAX(created_at) AS max_created_at
							 FROM rc_processed_content
							 WHERE deleted_at IS NULL
							 GROUP BY content_id),
				 rdd_cte AS (SELECT content_id, MAX(created_at) AS max_created_at
							 FROM rc_dependency_data
							 WHERE deleted_at IS NULL
							   AND content_id IS NOT NULL
							 GROUP BY content_id)
				SELECT count(roc.id)
				FROM rc_origin_content roc
						 LEFT JOIN rpc_cte rpc_max ON rpc_max.content_id = roc.id
						 LEFT JOIN rc_processed_content rpc ON rpc.content_id = roc.id
					AND rpc.created_at = rpc_max.max_created_at
					AND rpc.deleted_at IS NULL
						 INNER JOIN rdd_cte rdd_max ON rdd_max.content_id = roc.id
						 INNER JOIN rc_dependency_data rdd ON rdd.content_id = roc.id
					AND rdd.created_at = rdd_max.max_created_at
					AND rdd.deleted_at IS NULL
				WHERE roc.deleted_at IS NULL
				  AND rdd.create_by IN ?`, dsi.AccessibleIds).
		First(&count).
		Error; err != nil {
		return nil, err
	}

	var Infos = make([]dto.RcOriginContentInfo, 0)
	err = repo.data.Db.
		Raw(`WITH rpc_cte AS (SELECT content_id, MAX(created_at) AS max_created_at
							 FROM rc_processed_content
							 WHERE deleted_at IS NULL
							 GROUP BY content_id),
				 rdd_cte AS (SELECT content_id, MAX(created_at) AS max_created_at
							 FROM rc_dependency_data
							 WHERE deleted_at IS NULL
							   AND content_id IS NOT NULL
							 GROUP BY content_id)
				SELECT roc.id         AS content_id,
					   roc.usc_id,
					   roc.enterprise_name,
					   roc.year_month AS data_collect_month,
					   rdd.lh_qylx,
					   rpc.id         AS processed_id,
					   rpc.created_at AS processed_updated_at,
					   rdd.create_by,
					   rdd.id         as dep_id
				FROM rc_origin_content roc
						 LEFT JOIN rpc_cte rpc_max ON rpc_max.content_id = roc.id
						 LEFT JOIN rc_processed_content rpc ON rpc.content_id = roc.id
					AND rpc.created_at = rpc_max.max_created_at
					AND rpc.deleted_at IS NULL
						 INNER JOIN rdd_cte rdd_max ON rdd_max.content_id = roc.id
						 INNER JOIN rc_dependency_data rdd ON rdd.content_id = roc.id
					AND rdd.created_at = rdd_max.max_created_at
					AND rdd.deleted_at IS NULL
				WHERE roc.deleted_at IS NULL
				  AND rdd.create_by IN ?
				LIMIT ? OFFSET ?;`, dsi.AccessibleIds, page.PageSize, offset).
		Scan(&Infos).
		Error
	if err != nil {
		return nil, err
	}
	return &dto.RcOriginContentInfosResp{
		PaginationResp: dto.PaginationResp{
			Total:     count,
			TotalPage: int(math.Ceil(float64(count) / float64(page.PageSize))),
			PageNum:   pageNum,
			PageSize:  page.PageSize,
		},
		Data: &Infos,
	}, nil
}

func (repo *RcOriginContentRepo) Get(ctx context.Context, id int64) (*dto.RcOriginContent, error) {
	var modelRoc dto.RcOriginContent
	err := repo.data.Db.Table(modelRoc.TableName()).Where("id = ?", id).First(&modelRoc).Error
	if err != nil {
		return nil, err
	}
	return &modelRoc, nil
}

func (repo *RcOriginContentRepo) GetInfosByKwd(ctx context.Context, page *dto.PaginationReq, kwd string) (*dto.RcOriginContentInfosResp, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}
	kwdLike := fmt.Sprintf("%%%s%%", kwd)
	pageNum := int(math.Max(1, float64(page.PageNum)))
	offset := (pageNum - 1) * page.PageSize
	var count int64
	if err := repo.data.Db.
		Raw(`WITH rpc_cte AS (SELECT content_id, MAX(created_at) AS max_created_at
							 FROM rc_processed_content
							 WHERE deleted_at IS NULL
							 GROUP BY content_id),
				 rdd_cte AS (SELECT content_id, MAX(created_at) AS max_created_at
							 FROM rc_dependency_data
							 WHERE deleted_at IS NULL
							   AND content_id IS NOT NULL
							 GROUP BY content_id)
				SELECT count(roc.id)
				FROM rc_origin_content roc
						 LEFT JOIN rpc_cte rpc_max ON rpc_max.content_id = roc.id
						 LEFT JOIN rc_processed_content rpc ON rpc.content_id = roc.id
					AND rpc.created_at = rpc_max.max_created_at
					AND rpc.deleted_at IS NULL
						 INNER JOIN rdd_cte rdd_max ON rdd_max.content_id = roc.id
						 INNER JOIN rc_dependency_data rdd ON rdd.content_id = roc.id
					AND rdd.created_at = rdd_max.max_created_at
					AND rdd.deleted_at IS NULL
				WHERE roc.deleted_at IS NULL
				And roc.enterprise_name like ?
				AND rdd.create_by IN ?`, kwdLike, dsi.AccessibleIds).
		First(&count).
		Error; err != nil {
		return nil, err
	}

	var Infos = make([]dto.RcOriginContentInfo, 0)
	err = repo.data.Db.
		Raw(`WITH rpc_cte AS (SELECT content_id, MAX(created_at) AS max_created_at
							 FROM rc_processed_content
							 WHERE deleted_at IS NULL
							 GROUP BY content_id),
				 rdd_cte AS (SELECT content_id, MAX(created_at) AS max_created_at
							 FROM rc_dependency_data
							 WHERE deleted_at IS NULL
							   AND content_id IS NOT NULL
							 GROUP BY content_id)
				SELECT roc.id         AS content_id,
					   roc.usc_id,
					   roc.enterprise_name,
					   roc.year_month AS data_collect_month,
					   rdd.lh_qylx,
					   rpc.id         AS processed_id,
					   rpc.created_at AS processed_updated_at,
					   rdd.create_by,
					   rdd.id         as dep_id
				FROM rc_origin_content roc
						 LEFT JOIN rpc_cte rpc_max ON rpc_max.content_id = roc.id
						 LEFT JOIN rc_processed_content rpc ON rpc.content_id = roc.id
					AND rpc.created_at = rpc_max.max_created_at
					AND rpc.deleted_at IS NULL
						 INNER JOIN rdd_cte rdd_max ON rdd_max.content_id = roc.id
						 INNER JOIN rc_dependency_data rdd ON rdd.content_id = roc.id
					AND rdd.created_at = rdd_max.max_created_at
					AND rdd.deleted_at IS NULL
				WHERE roc.deleted_at IS NULL
				AND rdd.create_by IN ?
				And roc.enterprise_name like ?
				LIMIT ? OFFSET ?;`, dsi.AccessibleIds, kwdLike, page.PageSize, offset).
		Scan(&Infos).
		Error
	if err != nil {
		return nil, err
	}
	return &dto.RcOriginContentInfosResp{
		PaginationResp: dto.PaginationResp{
			Total:     count,
			TotalPage: int(math.Ceil(float64(count) / float64(page.PageSize))),
			PageNum:   pageNum,
			PageSize:  page.PageSize,
		},
		Data: &Infos,
	}, nil
}

func (repo *RcOriginContentRepo) GetContentIdsByUscId(ctx context.Context, uscId string) ([]int64, error) {
	var ids []int64
	err := repo.data.Db.
		Model(&dto.RcOriginContent{}).
		Select("id").
		Where("usc_id = ? ", uscId).
		Pluck("id", &ids).
		Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}
