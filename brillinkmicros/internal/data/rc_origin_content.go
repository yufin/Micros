package data

import (
	"brillinkmicros/internal/biz"
	"context"
	"fmt"
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
	err := repo.data.db.
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
	var modelRoc biz.RcOriginContent
	var modelRpc biz.RcProcessedContent
	var Infos = make([]biz.RcOriginContentInfo, 0)
	pageNum := int(math.Max(1, float64(page.PageNum)))
	offset := (pageNum - 1) * page.PageSize
	var count int64
	err := repo.data.db.Table(modelRoc.TableName()).Count(&count).Error
	if err != nil {
		return nil, err
	}

	err = repo.data.db.
		Raw(
			fmt.Sprintf("select id as content_id, usc_id, enterprise_name,`year_month` as data_collect_month, content as content, status_code, processed_id, processed_updated_at "+
				"from %s roc left join (select id as processed_id, content_id as processed_content_id, updated_at as processed_updated_at "+
				"from (select *, row_number() over (partition by content_id order by updated_at desc ) as rn from %s) t where rn = 1) rpc "+
				"on rpc.processed_content_id = roc.id limit ? offset ?", modelRoc.TableName(), modelRpc.TableName(),
			), page.PageSize, offset).
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
	err := repo.data.db.Table(modelRoc.TableName()).Where("id = ?", id).First(&modelRoc).Error
	if err != nil {
		return nil, err
	}
	return &modelRoc, nil
}
