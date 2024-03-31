package data

import (
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
	dwdataV2 "micros-api/api/dwdata/v2"
	dwdataV3 "micros-api/api/dwdata/v3"
	pipelineV1 "micros-api/api/pipeline/v1"
	"micros-api/internal/conf"
	"micros-api/pkg/miniocli"
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
	NewRedisClient,

	NewDwdataServiceClient,
	NewPipelineServiceClient,

	NewRcProcessedContentRepo,
	NewRcOriginContentRepo,
	NewRcDependencyDataRepo,
	NewGraphRepo,
	NewRcReportOssRepo,
	NewOssMetadataRepo,
	NewRcRdmResultRepo,
	NewRcRdmResDetailRepo,
	NewMgoRcRepo,
	NewClientDwDataRepo,
	NewClientPipelineRepo,
	NewRcDecisionFactorRepo,
	NewRcDecisionFactorV3Repo,
	NewRcContentMetaRepo,
	NewArtifactDataRepo,
	NewUserAuthRepo,
)

type Data struct {
	Db             *gorm.DB
	DbBl           *gorm.DB
	Rdb            *Rdb
	MinioCli       *miniocli.MinioClient
	Nw             *NatsWrap
	Neo            *NeoCli
	MgoCli         *MgoCli
	DwDataClient   dwdataV2.DwdataServiceClient
	DwDataClientV3 dwdataV3.DwdataServiceClient
	PipelineClient pipelineV1.PipelineServiceClient
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

func NewDbs(c *conf.Data) (*Dbs, func(), error) {
	db, err := newGormDB(c.Database.Source)
	if err != nil {
		return nil, nil, err
	}
	dbBl, err := newGormDB(c.BlAuth.Database.Source)
	if err != nil {
		return nil, func() {
			sqlDb, err := db.DB()
			if err != nil {
				log.Errorf("failed to get sqlDb obj while cleanup: %v", err)
			}
			if err := sqlDb.Close(); err != nil {
				log.Errorf("failed to close db: %v", err)
			}
		}, err
	}
	return &Dbs{
			Db:   db,
			DbBl: dbBl,
		}, func() {
			sqlDb1, err := db.DB()
			if err != nil {
				log.Errorf("failed to get sqlDb obj while cleanup: %v", err)
			}
			if err := sqlDb1.Close(); err != nil {
				log.Errorf("failed to close db: %v", err)
			}
			sqlDb2, err := dbBl.DB()
			if err != nil {
				log.Errorf("failed to get sqlDb obj while cleanup: %v", err)
			}
			if err := sqlDb2.Close(); err != nil {
				log.Errorf("failed to close db: %v", err)
			}
		}, nil
}

func NewNatsConn(c *conf.Data) (*NatsWrap, func(), error) {
	uri := c.Nats.Uri
	nc, err := nats.Connect(uri)
	if err != nil {
		return nil, nil, err
	}
	// init js
	cleanUp := func() {
		if err := nc.Drain(); err != nil {
			log.Error(err, "nats close error")
		}
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, cleanUp, err
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "TASK",
		Retention: nats.WorkQueuePolicy,
		Subjects:  []string{"task.rc.>"},
	})
	if err != nil {
		return nil, cleanUp, err
	}
	return &NatsWrap{
		Nc: nc,
		Js: js,
	}, cleanUp, nil
}

func NewNeoCli(c *conf.Data) (*NeoCli, func(), error) {
	d, err := neo4j.NewDriverWithContext(c.Neo4J.Url, neo4j.BasicAuth(c.Neo4J.Username, c.Neo4J.Password, ""))
	if err != nil {
		return nil, nil, err
	}
	cleanUp := func() {
		if err := d.Close(context.Background()); err != nil {
			log.Error(err, "neo4j close error")
		}
	}

	//err = d.VerifyConnectivity(context.Background())
	//if err != nil {
	//	return nil, err
	//}
	return &NeoCli{driver: d}, cleanUp, nil
}

func NewMgoCli(c *conf.Data) (*MgoCli, func(), error) {
	cli, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(c.MongoDb.Uri),
	)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	cleanUp := func() {
		if err := cli.Disconnect(context.Background()); err != nil {
			log.Errorf("failed to close mongo: %v", err)
		}
	}
	err = cli.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, cleanUp, errors.WithStack(err)
	}
	return &MgoCli{Client: cli}, cleanUp, nil
}

// NewData .
func NewData(
	dbs *Dbs,
	nw *NatsWrap,
	neo *NeoCli,
	miCli *miniocli.MinioClient,
	mgoCli *MgoCli,
	dwdataCli *DwDataClients,
	pipeline pipelineV1.PipelineServiceClient,
	rdb *Rdb,
) (*Data, error) {

	return &Data{
		Db:             dbs.Db,
		DbBl:           dbs.DbBl,
		Nw:             nw,
		Neo:            neo,
		MinioCli:       miCli,
		MgoCli:         mgoCli,
		DwDataClient:   dwdataCli.dwDataV2,
		DwDataClientV3: dwdataCli.dwDataV3,
		PipelineClient: pipeline,
		Rdb:            rdb,
	}, nil
}
