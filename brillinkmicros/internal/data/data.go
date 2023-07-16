package data

import (
	"brillinkmicros/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
	NewNeoCli,
	NewRcProcessedContentRepo,
	NewRcOriginContentRepo,
	NewRcDependencyDataRepo,
	NewGraphNodeRepo,
)

type Data struct {
	Db   *gorm.DB
	DbBl *gorm.DB
	Nw   *NatsWrap
	Neo  *NeoCli
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

func NewNeoCli(c *conf.Data) (*NeoCli, error) {
	d, err := neo4j.NewDriverWithContext(c.Neo4J.Url, neo4j.BasicAuth(c.Neo4J.Username, c.Neo4J.Password, ""))
	if err != nil {
		return nil, err
	}
	return &NeoCli{Neo: d}, nil
}

// NewData .
func NewData(logger log.Logger, dbs *Dbs, nw *NatsWrap, neo *NeoCli) (*Data, func(), error) {
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

	return &Data{
		Db:   dbs.db,
		DbBl: dbs.dbBl,
		Nw:   nw,
		Neo:  neo,
	}, cleanup, nil
}
