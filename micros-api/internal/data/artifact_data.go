package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"micros-api/internal/biz"
)

type ArtifactDataRepo struct {
	data *Data
	log  *log.Helper
}

func NewArtifactDataRepo(data *Data, logger log.Logger) biz.ArtifactDataRepo {
	return &ArtifactDataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *ArtifactDataRepo) SaveMany(ctx context.Context, db string, coll string, req []interface{}) error {
	_, err := repo.data.MgoCli.Client.Database(db).Collection(coll).InsertMany(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ArtifactDataRepo) Get(ctx context.Context, db string, coll string, pageSize int64, pageNum int64, orderBy string, filter bson.D, result interface{}) error {
	findOptions := options.Find()
	if orderBy != "" {
		findOptions.SetSort(bson.D{{orderBy, -1}})
	}
	findOptions.SetSkip((pageNum - 1) * pageSize)
	findOptions.SetLimit(pageSize)
	curr, err := repo.data.MgoCli.Client.Database(db).Collection(coll).Find(ctx, filter, findOptions)
	if err != nil {
		return err
	}
	defer curr.Close(ctx)

	if err = curr.All(ctx, result); err != nil {
		return err
	}
	return nil

	//resArr := make([]map[string]any, 0)
	//for _, bm := range results
	//	b, err := bson.MarshalExtJSON(bm, false, false)
	//	if err != nil {
	//		return nil, err
	//	}
	//	var m map[string]any
	//	if err := bson.UnmarshalExtJSON(b, false, &m); err != nil {
	//		return nil, err
	//	}
	//	resArr = append(resArr, m)
	//}
	//return resArr, nil
}

func (repo *ArtifactDataRepo) DeleteOne(ctx context.Context, db string, coll string, id string) error {
	var result bson.M
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	err = repo.data.MgoCli.Client.Database(db).Collection(coll).FindOneAndDelete(ctx, bson.D{{"_id", objId}}).Decode(&result)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ArtifactDataRepo) CountDoc(ctx context.Context, filter bson.D, coll string, db string) (int64, error) {
	count, err := repo.data.MgoCli.Client.Database(db).
		Collection(coll).
		CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}
