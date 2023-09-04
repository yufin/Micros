package data

import (
	"brillinkmicros/internal/conf"
	"brillinkmicros/pkg/miniocli"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nats-io/nats.go"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	NewMinioClient,
	NewMgoCli,
	NewRcProcessedContentRepo,
	NewRcOriginContentRepo,
	NewRcDependencyDataRepo,
	NewGraphRepo,
	NewRcReportOssRepo,
	NewOssMetadataRepo,
	NewRcRdmResultRepo,
	NewRcRdmResDetailRepo,
	NewMgoRcRepo,
	NewDwEnterpriseRepo,
)

type Data struct {
	Db       *gorm.DB
	DbBl     *gorm.DB
	Nw       *NatsWrap
	Neo      *NeoCli
	MinioCli *miniocli.MinioClient
	MgoCli   *MgoCli
}

func newGormDB(dsn string) (*gorm.DB, error) {
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

func NewMinioClient(c *conf.Data) (*miniocli.MinioClient, error) {
	moCli, err := minio.New(c.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKey, c.Minio.SecretKey, ""),
		Secure: c.Minio.UseSsl,
	})
	if err != nil {
		return nil, err
	}
	return &miniocli.MinioClient{
		Cli: moCli,
	}, nil
}

func NewDbs(c *conf.Data) (*Dbs, error) {
	db, err := newGormDB(c.Database.Source)
	if err != nil {
		return nil, err
	}
	dbBl, err := newGormDB(c.BlAuth.Database.Source)
	if err != nil {
		return nil, err
	}
	return &Dbs{
		Db:   db,
		DbBl: dbBl,
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
		Subjects:  []string{"task.rc.>"},
	})
	if err != nil {
		return nil, err
	}
	return &NatsWrap{
		Nc: nc,
		Js: js,
	}, nil
}

func NewNeoCli(c *conf.Data) (*NeoCli, error) {
	d, err := neo4j.NewDriverWithContext(c.Neo4J.Url, neo4j.BasicAuth(c.Neo4J.Username, c.Neo4J.Password, ""))
	if err != nil {
		return nil, err
	}
	err = d.VerifyConnectivity(context.Background())
	if err != nil {
		return nil, err
	}
	return &NeoCli{driver: d}, nil
}

func NewMgoCli(c *conf.Data) (*MgoCli, error) {
	cli, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(c.MongoDb.Uri),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = cli.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &MgoCli{Client: cli}, nil
}

// NewData .
func NewData(logger log.Logger, dbs *Dbs, nw *NatsWrap, neo *NeoCli, miCli *miniocli.MinioClient, mgoCli *MgoCli) (*Data, func(), error) {
	ndLog := log.NewHelper(logger)

	cleanup := func() {
		ndLog.Info("Closing the data resources")
		for _, db := range []*gorm.DB{dbs.Db, dbs.DbBl} {
			db := db
			sqlDb, err := db.DB()
			if err != nil {
				ndLog.Errorf("failed to get sqlDb obj while cleanup: %v", err)
			}
			if err := sqlDb.Close(); err != nil {
				ndLog.Errorf("failed to close db: %v", err)
			}
		}

		if err := nw.Nc.Drain(); err != nil {
			ndLog.Errorf("failed to drain nats: %v", err)
		}

		if err := mgoCli.close(); err != nil {
			ndLog.Errorf("failed to disconnect mongodb: %v", err)
		}

		ndLog.Info("Data resource Closed")
	}

	return &Data{
		Db:       dbs.Db,
		DbBl:     dbs.DbBl,
		Nw:       nw,
		Neo:      neo,
		MinioCli: miCli,
		MgoCli:   mgoCli,
	}, cleanup, nil
}
