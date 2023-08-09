package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nats-io/nats.go"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorm.io/gorm"
)

type NatsWrap struct {
	Nc *nats.Conn
	Js nats.JetStreamContext
}

type Dbs struct {
	Db   *gorm.DB
	DbBl *gorm.DB
}

type NeoCli struct {
	driver neo4j.DriverWithContext
}

func (n NeoCli) CypherQuery(ctx context.Context, cypher string, params map[string]any) ([]neo4j.Record, error) {
	//ctxTemp := context.Background()
	session := n.driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})

	defer func(session neo4j.SessionWithContext, ctxTemp context.Context) {
		err := session.Close(ctxTemp)
		if err != nil {
			log.Errorf("Error closing Neo4j session: %v", err)
		}
	}(session, context.TODO())

	result, err := session.Run(context.Background(), cypher, params)
	if err != nil {
		return nil, err
	}
	var output []neo4j.Record
	for result.Next(context.Background()) {
		record := result.Record()
		output = append(output, *record)
	}
	return output, err
}
