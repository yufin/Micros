package data

import (
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/biz/dto"
	"brillinkmicros/pkg"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (repo *MgoRcRepo) GetProcessedObjIdByContentId(ctx context.Context, contentId int64) (bson.M, error) {
	data := bson.M{}
	err := repo.data.MgoCli.Client.Database("rc").Collection("processed_content").
		FindOne(
			context.TODO(),
			bson.M{"content_id": strconv.FormatInt(contentId, 10)},
			options.FindOne().SetSort(bson.D{{"created_at", -1}}).SetProjection(bson.D{{"_id", 1}}),
		).Decode(&data)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return data, nil
}

func (repo *MgoRcRepo) GetProcessedContentByObjId(ctx context.Context, objIdHex string) (bson.M, error) {
	var data bson.M
	objId, err := primitive.ObjectIDFromHex(objIdHex)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = repo.data.MgoCli.Client.Database("rc").Collection("processed_content").
		FindOne(
			context.TODO(),
			bson.M{"_id": objId},
		).Decode(&data)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
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
		`WITH rdd_cte AS (SELECT content_id, MAX(created_at) AS max_created_at
                 FROM rc_dependency_data
                 WHERE deleted_at IS NULL
                   AND content_id IS NOT NULL
                 GROUP BY content_id)
			SELECT count(roc.id)
			FROM rc_origin_content roc
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

	err = repo.data.Db.Raw(
		`WITH rdd_cte AS (SELECT content_id, MAX(created_at) AS max_created_at
                 FROM rc_dependency_data
                 WHERE deleted_at IS NULL
                   AND content_id IS NOT NULL
                 GROUP BY content_id)
			SELECT roc.id         AS content_id,
				   roc.usc_id,
				   roc.enterprise_name,
				   roc.year_month AS data_collect_month,
				   rdd.lh_qylx,
				   rdd.create_by,
				   rdd.id         as dep_id
			FROM rc_origin_content roc
					 INNER JOIN rdd_cte rdd_max ON rdd_max.content_id = roc.id
					 INNER JOIN rc_dependency_data rdd ON rdd.content_id = roc.id
				AND rdd.created_at = rdd_max.max_created_at
				AND rdd.deleted_at IS NULL
			WHERE roc.deleted_at IS NULL
			  AND rdd.create_by IN ?
			LIMIT ? OFFSET ?;`, dsi.AccessibleIds, page.PageSize, offset).
		Scan(&infos).
		Error
	if err != nil {
		return nil, err
	}

	for i, info := range infos {
		// get processed content
		processedContent, err := repo.GetProcessedObjIdByContentId(ctx, info.ContentId)
		if err != nil {
			return nil, err
		}
		// if processed content is empty, then skip
		if processedContent != nil {
			docId := processedContent["_id"].(primitive.ObjectID)
			infos[i].ProcessedId = docId.Hex()
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
