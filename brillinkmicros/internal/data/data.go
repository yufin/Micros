package data

import (
	"brillinkmicros/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// database wrapper

var NatsConn *nats.Conn

// Data .
// wrapped database client

// ProviderSet is data providers.
// var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)
var ProviderSet = wire.NewSet(
	NewData,
	NewGormDB,
	NewNatsConn,
	NewRcProcessedContentRepo,
	NewRcOriginContentRepo,
	NewRcDependencyDataRepo,
)

type Data struct {
	db *gorm.DB
	js nats.JetStreamContext
}

func NewGormDB(c *conf.Data) (*gorm.DB, error) {
	dsn := c.Database.Source
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// get sql.DB object to set db connection pool options
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	return db, nil
}

func NewNatsConn(c *conf.Data) (nats.JetStreamContext, error) {
	uri := c.Nats.Uri
	nc, err := nats.Connect(uri)
	if err != nil {
		return nil, err
	}
	NatsConn = nc
	// init js
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "TASK",
		Retention: nats.WorkQueuePolicy,
		Subjects:  []string{"task.rskc.>"},
	})
	if err != nil {
		return nil, err
	}
	return js, nil
}

// NewData .
func NewData(logger log.Logger, db *gorm.DB, js nats.JetStreamContext) (*Data, func(), error) {
	ndLog := log.NewHelper(logger)

	cleanup := func() {
		ndLog.Info("Closing the data resources")
		sqlDb, err := db.DB()
		if err != nil {
			ndLog.Errorf("failed to get sqlDb obj while cleanup: %v", err)
		}
		if err := sqlDb.Close(); err != nil {
			ndLog.Errorf("failed to close db: %v", err)
		}
		if err := NatsConn.Drain(); err != nil {
			ndLog.Errorf("failed to drain nats: %v", err)
		}

		ndLog.Info("Data resource Closed")
	}

	return &Data{db: db, js: js}, cleanup, nil
}
