package data

import (
	"brillinkmicros/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// database wrapper

// ProviderSet is data providers.
// var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)
var ProviderSet = wire.NewSet(
	NewData,
	NewGormDB,
	NewRcProcessedContentRepo,
	NewRcOriginContentRepo,
	NewRcDependencyDataRepo,
)

// Data .
// wrapped database client
type Data struct {
	db *gorm.DB
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

// NewData .
func NewData(logger log.Logger, db *gorm.DB) (*Data, func(), error) {
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
		ndLog.Info("Data resource Closed")
	}

	return &Data{db: db}, cleanup, nil
}
