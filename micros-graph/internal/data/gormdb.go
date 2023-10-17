package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"micros-graph/internal/conf"
	"time"
)

type Db struct {
	Db *gorm.DB
}

func NewGormDb(c *conf.Data) (*Db, func(), error) {
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		sqlDB, _ := db.DB()
		if err := sqlDB.Close(); err != nil {
			log.Errorf("Failed to close database: %v", err)
		}
	}
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	return &Db{Db: db}, cleanup, nil
}
