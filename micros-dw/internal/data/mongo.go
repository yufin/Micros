package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"micros-dw/internal/conf"
)

type MongoDb struct {
	Client *mongo.Client
}

func NewMongoDb(c *conf.Data) (*MongoDb, func(), error) {
	cli, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(c.MongoDb.Uri),
	)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	err = cli.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	cleanup := func() {
		err := cli.Disconnect(context.Background())
		if err != nil {
			log.Errorf("Error closing MongoDb client: %v", err)
		}
	}
	return &MongoDb{Client: cli}, cleanup, nil
}
