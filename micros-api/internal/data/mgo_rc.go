package data

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"math"
	"micros-api/internal/biz"
	"micros-api/internal/biz/dto"
	"micros-api/pkg"
	"strconv"
	"time"
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

func (r *MgoRcRepo) SaveReportPrintConfig(ctx context.Context, b []byte) (*mongo.InsertOneResult, error) {
	var m bson.M
	err := bson.UnmarshalExtJSON(b, false, &m)
	if err != nil {
		return nil, err
	}

	res, err := r.data.MgoCli.Client.
		Database("rc").
		Collection("report_print_config").
		InsertOne(ctx, m)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (r *MgoRcRepo) GetReportPrintConfig(ctx context.Context, cond bson.M) ([]byte, error) {
	var m bson.M
	err := r.data.MgoCli.Client.
		Database("rc").
		Collection("report_print_config").
		FindOne(ctx, cond).
		Decode(&m)
	if err != nil {
		return nil, err
	}

	b, err := bson.MarshalExtJSON(m, false, false)
	if err != nil {
		return nil, err
	}
	return b, err
}

func (r *MgoRcRepo) DeleteReportPrintConfig(ctx context.Context, cond bson.M) error {
	_, err := r.data.MgoCli.Client.
		Database("rc").
		Collection("report_print_config").
		DeleteOne(ctx, cond)
	if err != nil {
		return err
	}
	return nil
}

func (repo *MgoRcRepo) GetRdmResultByClaimedId(ctx context.Context, claimId int64) (bson.M, error) {
	var data bson.M
	err := repo.data.MgoCli.Client.Database("rc").Collection("rdm_result").
		FindOne(
			context.TODO(),
			bson.M{"claim_id": claimId},
		).Decode(&data)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return data, nil
}

func (repo *MgoRcRepo) GetNewestDocByContentId(ctx context.Context, contentId int64) (bson.M, error) {
	data := bson.M{}
	err := repo.data.MgoCli.Client.Database("rc").Collection("processed_content").
		FindOne(
			context.TODO(),
			bson.M{"content_id": strconv.FormatInt(contentId, 10)},
			options.FindOne().SetSort(bson.D{{"created_at", -1}}),
		).Decode(&data)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	return data, nil
}

func (repo *MgoRcRepo) GetNewestDocInfoByContentId(ctx context.Context, contentId int64) (string, time.Time, error) {
	data := bson.M{}
	err := repo.data.MgoCli.Client.Database("rc").Collection("processed_content").
		FindOne(
			context.TODO(),
			bson.M{"content_id": strconv.FormatInt(contentId, 10)},
			options.FindOne().SetSort(bson.D{{"created_at", -1}}).SetProjection(bson.D{{"_id", 1}, {"created_at", 1}}),
		).Decode(&data)
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return "", time.Time{}, nil
		}
		return "", time.Time{}, errors.WithStack(err)
	}
	docId := data["_id"].(primitive.ObjectID)
	createdAt := data["created_at"].(primitive.DateTime)
	return docId.Hex(), createdAt.Time(), nil
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

func (repo *MgoRcRepo) GetProcessedContentInfoByObjId(ctx context.Context, objIdHex string) (bson.M, error) {
	var data bson.M
	objId, err := primitive.ObjectIDFromHex(objIdHex)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = repo.data.MgoCli.Client.Database("rc").Collection("processed_content").
		FindOne(
			context.TODO(),
			bson.M{"_id": objId},
			options.FindOne().SetProjection(bson.D{
				{"content", -1},
				{"content_id", 1},
				{"_id", 1},
				{"created_at", 1}}),
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

func (repo *MgoRcRepo) GetContentInfosByKwd(ctx context.Context, page *dto.PaginationReq, kwd string) (*dto.RcOriginContentInfosRespV3, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}
	pageNum := int(math.Max(1, float64(page.PageNum)))
	offset := (pageNum - 1) * page.PageSize
	var count int64
	infos := make([]dto.RcOriginContentInfoV3, 0)
	kwdLike := fmt.Sprintf("%%%s%%", kwd)
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
			And roc.enterprise_name like ?
			  AND rdd.create_by IN ?`, kwdLike, dsi.AccessibleIds).
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
			AND roc.enterprise_name like ?
			  AND rdd.create_by IN ?
			LIMIT ? OFFSET ?;`, kwdLike, dsi.AccessibleIds, page.PageSize, offset).
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

func (repo *MgoRcRepo) GetCountOnDistinctUscIdForColl(ctx context.Context, collName string) (int, error) {
	collection := repo.data.MgoCli.Client.Database("dw2").Collection(collName)

	// Aggregation pipeline
	pipeline := mongo.Pipeline{
		bson.D{{"$group", bson.D{{"_id", "$usc_id"}}}},
		bson.D{{"$group", bson.D{{"_id", nil}, {"count", bson.D{{"$sum", 1}}}}}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}
	if len(results) == 0 {
		return 0, nil
	}
	totalDist := results[0]["count"]
	totalDistInt, ok := totalDist.(int32)
	if !ok {
		return 0, errors.New("type assertion failed")
	}
	return int(totalDistInt), nil
}
