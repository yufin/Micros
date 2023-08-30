package data

import (
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/biz/dto"
	"brillinkmicros/pkg"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"math"
	"strconv"
)

type MgoRcRepo struct {
	data *Data
	log  *log.Helper
}

func NewMgoRcRepo(data *Data, logger log.Logger) biz.MgoRcRepo {
	return &MgoRcRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *MgoRcRepo) GetProcessedContent(ctx context.Context, contentId int64) (bson.M, error) {
	data := bson.M{}
	coll := repo.data.MgoCli.Client.Database("rc").Collection("processed_content")
	err := coll.FindOne(
		ctx,
		bson.M{"content_id": strconv.FormatInt(contentId, 10)},
		options.FindOne().SetSort(bson.D{{"created_at", -1}}),
	).Decode(&data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	//content := data["content"].(bson.M)
	return data, nil
}

func (repo *MgoRcRepo) GetContentInfos(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentInfosRespV3, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}
	pageNum := int(math.Max(1, float64(page.PageNum)))
	offset := (pageNum - 1) * page.PageSize
	var count int64
	infos := make([]dto.RcOriginContentInfoV3, 0)

	if err := repo.data.Db.Raw(
		`select count(roc.id) as count
			from rc_origin_content roc
					 left join
				 (select *
				  from (select *, row_number() over (partition by content_id order by created_at DESC ) as rn
						from rc_dependency_data
						where content_id is not null) t
				  where t.rn = 1) rdd on roc.id = rdd.content_id
			where rdd.create_by in (?)
			  and rdd.deleted_at is null
			  and roc.deleted_at is null`, dsi.AccessibleIds).
		First(&count).
		Error; err != nil {
		return nil, err
	}

	err = repo.data.Db.Raw(
		`select roc.id as content_id,
			   roc.enterprise_name as enterprise_name,
			   roc.usc_id as usc_id,
			   roc.year_month as data_collect_month,
			   rdd.lh_qylx, 
			   rdd.create_by    as create_by,
			   rdd.id           as dep_id
		from rc_origin_content roc
				 left join
			 (select *
			  from (select *, row_number() over (partition by content_id order by created_at DESC ) as rn
					from rc_dependency_data
					where content_id is not null) t
			  where t.rn = 1) rdd on roc.id = rdd.content_id
		where rdd.create_by in (?)
		  and rdd.deleted_at is null
		  and roc.deleted_at is null
		  limit ? offset ?`, dsi.AccessibleIds, page.PageSize, offset).
		Scan(&infos).
		Error
	if err != nil {
		return nil, err
	}

	for i, info := range infos {
		// get processed content
		processedContent, err := repo.GetProcessedContent(ctx, info.ContentId)
		if err != nil {
			return nil, err
		}
		// if processed content is empty, then skip
		if processedContent == nil {
			docId := processedContent["_id"].(string)
			infos[i].ProcessedId = docId
		}
	}

	return &dto.RcOriginContentInfosRespV3{
		PaginationResp: dto.PaginationResp{
			Total:     count,
			TotalPage: int(math.Ceil(float64(count) / float64(page.PageSize))),
			PageNum:   pageNum,
			PageSize:  page.PageSize,
		},
		Data: &infos,
	}, nil
}
