package data

import (
	"brillinkmicros/common"
	"brillinkmicros/internal/biz"
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

func (repo *RcOriginContentRepo) GetPage(ctx context.Context, page *biz.PaginationReq) (*biz.RcOriginContentGetPageResp, error) {
	var modelRoc biz.RcOriginContent
	listRoc := make([]biz.RcOriginContent, 0)
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
	return &biz.RcOriginContentGetPageResp{
		PaginationResp: biz.PaginationResp{
			Total:     count,
			TotalPage: int(math.Ceil(float64(count) / float64(page.PageSize))),
			PageNum:   page.PageNum,
			PageSize:  page.PageSize,
		},
		Data: &listRoc,
	}, nil
}

func (repo *RcOriginContentRepo) GetInfos(ctx context.Context, page *biz.PaginationReq) (*biz.RcOriginContentInfosResp, error) {
	dsi, err := common.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}

	var modelRoc biz.RcOriginContent
	var Infos = make([]biz.RcOriginContentInfo, 0)
	pageNum := int(math.Max(1, float64(page.PageNum)))
	offset := (pageNum - 1) * page.PageSize
	var count int64
	if err := repo.data.Db.
		Table(modelRoc.TableName()).
		Scopes(ApplyBlDataScope(dsi)).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	err = repo.data.Db.
		Raw(`SELECT
					roc.id AS content_id,
					roc.usc_id,
					roc.enterprise_name,
					roc.year_month AS data_collect_month,
					roc.status_code,
					rpc.id AS processed_id,
					rpc.updated_at AS processed_updated_at,
					rdd.content_id,
					rdd.create_by
				FROM
					rskc_origin_content roc
				LEFT JOIN
					(
						SELECT
							*,
							ROW_NUMBER() OVER (PARTITION BY content_id ORDER BY updated_at DESC) AS rn
						FROM
							rskc_processed_content
						WHERE
							deleted_at IS NULL
					) rpc ON rpc.content_id = roc.id AND rpc.rn = 1
				INNER JOIN
					rc_dependency_data rdd ON rdd.content_id = roc.id
				WHERE
					roc.deleted_at IS NULL
					AND rdd.deleted_at IS NULL
					AND rdd.create_by IN (?)
				LIMIT ? OFFSET ?;`, dsi.AccessibleIds, page.PageSize, offset).
		Scan(&Infos).
		Error
	if err != nil {
		return nil, err
	}
	return &biz.RcOriginContentInfosResp{
		PaginationResp: biz.PaginationResp{
			Total:     count,
			TotalPage: int(math.Ceil(float64(count) / float64(page.PageSize))),
			PageNum:   pageNum,
			PageSize:  page.PageSize,
		},
		Data: &Infos,
	}, nil
}

func (repo *RcOriginContentRepo) Get(ctx context.Context, id int64) (*biz.RcOriginContent, error) {
	var modelRoc biz.RcOriginContent
	err := repo.data.Db.Table(modelRoc.TableName()).Where("id = ?", id).First(&modelRoc).Error
	if err != nil {
		return nil, err
	}
	return &modelRoc, nil
}
