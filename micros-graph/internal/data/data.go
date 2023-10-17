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
	NewNebulaDb,

	NewGraphEdgeRepo,
	NewGraphNetRepo,
	NewGraphNodeRepo,
)

type Data struct {
	Dbs      *Db
	Mongo    *MongoDb
	nebulaDb *NebulaDb
}

// NewData .
func NewData(db *Db, mgo *MongoDb, nebulaDb *NebulaDb) *Data {
	return &Data{
		Dbs:      db,
		Mongo:    mgo,
		nebulaDb: nebulaDb,
	}
}
