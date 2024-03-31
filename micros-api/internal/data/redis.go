package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"micros-api/internal/conf"
)

type Rdb struct {
	Client *redis.Client
}

func NewRedisClient(c *conf.Data) (*Rdb, func(), error) {
	opt := redis.Options{
		Addr:         c.Redis.Addr,
		DB:           int(c.Redis.Db),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		Password:     c.Redis.Password,
	}
	rdb := redis.NewClient(&opt)
	// check if redis connected success
	if err := rdb.Ping(context.TODO()).Err(); err != nil {
		return nil, nil, err
	}

	return &Rdb{Client: rdb}, func() {
			if err := rdb.Close(); err != nil {
				log.Error(err, "redis close error")
			}
		},
		nil
}
