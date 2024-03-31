package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"micros-worker/internal/biz"
)

type NoticeRepo struct {
	data *Data
	log  *log.Helper
}

func NewNoticeRepo(data *Data, logger log.Logger) biz.CommonNoticeRepo {
	return &NoticeRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/notice")),
	}
}

func (r *NoticeRepo) SaveConfig(ctx context.Context, b []byte) (*mongo.InsertOneResult, error) {
	var m bson.M
	err := bson.UnmarshalExtJSON(b, false, &m)
	if err != nil {
		return nil, err
	}

	res, err := r.data.MongoDb.Client.
		Database("infra").
		Collection("notice_config").
		InsertOne(ctx, m)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (r *NoticeRepo) GetConfigById(ctx context.Context, idHex string) ([]byte, error) {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	var m bson.M
	err = r.data.MongoDb.Client.
		Database("infra").
		Collection("notice_config").
		FindOne(ctx, bson.M{"_id": id}).
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
