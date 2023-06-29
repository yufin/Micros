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

// Data .
// wrapped database client

// ProviderSet is data providers.
// var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)
var ProviderSet = wire.NewSet(
	NewData,
	NewDbs,
	NewNatsConn,
	NewRcProcessedContentRepo,
	NewRcOriginContentRepo,
	NewRcDependencyDataRepo,
)

type Data struct {
	Db   *gorm.DB
	DbBl *gorm.DB
	Nw   *NatsWrap
}
type NatsWrap struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

type Dbs struct {
	db   *gorm.DB
	dbBl *gorm.DB
}

func NewGormDB(dsn string) (*gorm.DB, error) {
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

func NewDbs(c *conf.Data) (*Dbs, error) {
	db, err := NewGormDB(c.Database.Source)
	if err != nil {
		return nil, err
	}
	dbBl, err := NewGormDB(c.BlAuth.Database.Source)
	if err != nil {
		return nil, err
	}
	return &Dbs{
		db:   db,
		dbBl: dbBl,
	}, nil
}

func NewNatsConn(c *conf.Data) (*NatsWrap, error) {
	uri := c.Nats.Uri
	nc, err := nats.Connect(uri)
	if err != nil {
		return nil, err
	}
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
	return &NatsWrap{
		nc: nc,
		js: js,
	}, nil
}

// NewData .
func NewData(logger log.Logger, dbs *Dbs, nw *NatsWrap) (*Data, func(), error) {
	ndLog := log.NewHelper(logger)

	cleanup := func() {
		ndLog.Info("Closing the data resources")
		for _, db := range []*gorm.DB{dbs.db, dbs.dbBl} {
			db := db
			sqlDb, err := db.DB()
			if err != nil {
				ndLog.Errorf("failed to get sqlDb obj while cleanup: %v", err)
			}
			if err := sqlDb.Close(); err != nil {
				ndLog.Errorf("failed to close db: %v", err)
			}
		}

		if err := nw.nc.Drain(); err != nil {
			ndLog.Errorf("failed to drain nats: %v", err)
		}

		ndLog.Info("Data resource Closed")
	}

	return &Data{Db: dbs.db, DbBl: dbs.dbBl, Nw: nw}, cleanup, nil
}
