package data

import (
	"github.com/google/wire"
	"micros-worker/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewGormDb,
	NewMongoDb,
	NewTemporalClient,
	NewRedisClient,
	NewSftpClientPool,

	NewContentRepo,
	NewNoticeRepo,
)

// Data .
type Data struct {
	TemporalClient *TemporalClient
	Db             *Db
	MongoDb        *MongoDb
	Redis          *Rdb
	SftpPool       *SftpClientPool
	DataConf       *conf.Data
}

// NewData .
func NewData(
	db *Db,
	mgo *MongoDb,
	temporalCli *TemporalClient,
	rdb *Rdb,
	sftpPool *SftpClientPool,
	dc *conf.Data,
) *Data {
	return &Data{
		Redis:          rdb,
		Db:             db,
		MongoDb:        mgo,
		TemporalClient: temporalCli,
		SftpPool:       sftpPool,
		DataConf:       dc,
	}
}
