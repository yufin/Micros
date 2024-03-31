package data

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"micros-dw/internal/biz"
	"micros-dw/internal/biz/dto"
	"time"
)

type DwDataRepo struct {
	data *Data
	log  *log.Helper
}

func NewDwDataRepo(data *Data, logger log.Logger) biz.DwDataRepo {
	return &DwDataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}

}

func (repo *DwDataRepo) GetUscIdByEnterpriseName(ctx context.Context, name string) (*dto.EnterpriseWaitList, error) {
	var data dto.EnterpriseWaitList
	err := repo.data.Dbs.Db.
		Model(&dto.EnterpriseWaitList{}).
		Where("enterprise_name = ?", name).
		Order("created_at desc").
		First(&data).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (repo *DwDataRepo) GetDocInDuration(ctx context.Context, uscId string, tp time.Time, coll string) (bson.M, error) {
	var res bson.M

	filter := bson.D{
		{"$and", []bson.D{
			{{"date", bson.D{{"$lte", tp}}}},
			{{"check_date", bson.D{{"$gte", tp}}}},
			{{"usc_id", uscId}},
		}},
	}

	err := repo.data.Mongo.Client.Database("dw").
		Collection(coll).
		FindOne(context.TODO(), filter).
		Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func (repo *DwDataRepo) GetDocInExtendDuration(ctx context.Context, uscId string, tp time.Time, collName string, extendDate int) (bson.M, error) {
	var res bson.M

	pipeline := mongo.Pipeline{
		{
			{"$addFields", bson.D{
				{"valid_time_end", bson.D{{"$add", bson.A{"$check_date", extendDate * 24 * 60 * 60000}}}},
			}},
		},
		{
			{"$match", bson.D{
				{"$and", bson.A{
					bson.D{{"date", bson.D{{"$lte", tp}}}},
					bson.D{{"valid_time_end", bson.D{{"$gte", tp}}}},
					bson.D{{"usc_id", uscId}},
				}},
			}},
		},
		{
			{"$project", bson.D{
				{"valid_time_end", 0},
			}},
		},
	}
	cur, err := repo.data.Mongo.Client.Database("dw").Collection(collName).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	if cur.TryNext(context.TODO()) {
		if err := cur.Decode(&res); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, nil
			}
			return nil, err
		}
	}
	return res, nil
}

func (repo *DwDataRepo) GetDocWithDuration(ctx context.Context, uscId string, tp time.Time, coll string, extendDate int) (map[string]any, error) {
	res, err := repo.GetDocInDuration(ctx, uscId, tp, coll)
	if err != nil {
		return nil, err
	}
	if res == nil {
		res, err = repo.GetDocInExtendDuration(ctx, uscId, tp, coll, extendDate)
		if err != nil {
			return nil, err
		}
	}
	var m map[string]any
	if res != nil {
		b, err := bson.MarshalExtJSON(res, false, false)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (repo *DwDataRepo) GetDocAggregated(ctx context.Context, coll string, pipeline *mongo.Pipeline) (*[]map[string]any, error) {
	var res []bson.M
	cur, err := repo.data.Mongo.Client.Database("dw2").
		Collection(coll).
		Aggregate(ctx, *pipeline)

	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())
	if err = cur.All(context.TODO(), &res); err != nil {
		return nil, err
	}
	resArr := make([]map[string]any, 0)
	for _, bm := range res {
		b, err := bson.MarshalExtJSON(bm, false, false)
		if err != nil {
			return nil, err
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, err
		}
		resArr = append(resArr, m)
	}
	return &resArr, nil
}

func (repo *DwDataRepo) GetDoc(ctx context.Context, coll string, cond *bson.D) (*map[string]any, error) {
	var res bson.M
	err := repo.data.Mongo.Client.Database("dw2").
		Collection(coll).
		FindOne(context.TODO(), *cond).Decode(&res)

	if err != nil {
		return nil, err
	}

	b, err := bson.MarshalExtJSON(res, false, false)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (repo *DwDataRepo) CheckStatusInDuration(ctx context.Context, uscId string, tp time.Time, coll string) (bool, error) {
	var res bson.M
	filter := bson.D{
		{"$and", []bson.D{
			{{"create_date", bson.D{{"$lte", tp}}}},
			{{"check_date", bson.D{{"$gte", tp}}}},
			{{"usc_id", uscId}},
		}},
	}
	err := repo.data.Mongo.Client.Database("dw").
		Collection(coll).
		FindOne(context.TODO(), filter).
		Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (repo *DwDataRepo) CountDoc(ctx context.Context, filter bson.D, coll string, db string) (int64, error) {
	//filter := bson.D{
	//	{"$and", []bson.D{
	//		{{"create_date", bson.D{{"$lte", tp}}}},
	//		{{"check_date", bson.D{{"$gte", tp}}}},
	//		{{"usc_id", uscId}},
	//	}},
	//}
	count, err := repo.data.Mongo.Client.Database(db).
		Collection(coll).
		CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *DwDataRepo) GetDocs(ctx context.Context, db string, coll string, orderBy string, pageSize int64, pageNum int64) ([]map[string]any, error) {
	findOptions := options.Find()
	if orderBy != "" {
		findOptions.SetSort(bson.D{{orderBy, -1}})
	}
	findOptions.SetSkip((pageNum - 1) * pageSize)
	findOptions.SetLimit(pageSize)

	var results []bson.M

	curr, err := repo.data.Mongo.Client.Database(db).Collection(coll).Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer curr.Close(ctx)
	if err = curr.All(ctx, &results); err != nil {
		return nil, err
	}
	resArr := make([]map[string]any, 0)
	for _, bm := range results {
		b, err := bson.MarshalExtJSON(bm, false, false)
		if err != nil {
			return nil, err
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, err
		}
		resArr = append(resArr, m)
	}
	return resArr, nil
}
