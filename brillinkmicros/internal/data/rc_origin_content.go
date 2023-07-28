package data

import (
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/biz/dto"
	"brillinkmicros/pkg"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"math"
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

func (repo *RcOriginContentRepo) GetInfos(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentInfosResp, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}

	var modelRoc dto.RcOriginContent
	var Infos = make([]dto.RcOriginContentInfo, 0)
	pageNum := int(math.Max(1, float64(page.PageNum)))
	offset := (pageNum - 1) * page.PageSize
	var count int64
	if err := repo.data.Db.
		Table(modelRoc.TableName()).
		Scopes(pkg.ApplyBlDataScope(dsi)).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	err = repo.data.Db.
		Raw(`SELECT roc.id         AS content_id,
					   roc.usc_id,
					   roc.enterprise_name,
					   roc.year_month AS data_collect_month,
					   rdd.lh_qylx,
					   rpc.id         AS processed_id,
					   rpc.created_at AS processed_updated_at,
					   rdd.content_id,
					   rdd.create_by,
					   rdd.id         as dep_id
				FROM rc_origin_content roc
						 LEFT JOIN
					 (SELECT *,
							 ROW_NUMBER() OVER (PARTITION BY content_id ORDER BY created_at DESC) AS rn
					  FROM rc_processed_content
					  WHERE deleted_at IS NULL) rpc ON rpc.content_id = roc.id AND rpc.rn = 1
						 INNER JOIN
					 (select *
					  from (select *, row_number() over (partition by content_id order by created_at DESC ) as rn
							from rc_dependency_data
							where content_id is not null) t
					  where t.rn = 1) rdd ON rdd.content_id = roc.id
				WHERE roc.deleted_at IS NULL
				  AND rdd.deleted_at IS NULL
					AND rdd.create_by IN (?)
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
