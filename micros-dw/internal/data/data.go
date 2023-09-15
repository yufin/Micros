package data

import (
	"github.com/google/wire"
)

// ProviderSet is data providers.
// var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)
var ProviderSet = wire.NewSet(
	NewMongoDb,
	NewGormDb,
	NewData,
	NewDwEnterpriseRepo,
)

type Data struct {
	Dbs   *Db
	Mongo *MongoDb
}

// NewData .
func NewData(db *Db, mgo *MongoDb) *Data {
	return &Data{
		Dbs:   db,
		Mongo: mgo,
	}
}
