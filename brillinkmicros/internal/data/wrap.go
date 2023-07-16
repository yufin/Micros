package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nats-io/nats.go"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorm.io/gorm"
)

type NatsWrap struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

type Dbs struct {
	db   *gorm.DB
	dbBl *gorm.DB
}

type NeoCli struct {
	Neo neo4j.DriverWithContext
}

func (t *NeoCli) CypherQuery(ctx context.Context, cypher string, params map[string]any) ([]neo4j.Record, error) {
	session := t.Neo.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Errorf("Error closing Neo4j session: %v", err)
		}
	}(session, ctx)
	result, err := session.Run(ctx, cypher, params)
	if err != nil {
		return nil, err
	}
	var output []neo4j.Record
	for result.Next(ctx) {
		record := result.Record()
		output = append(output, *record)
	}
	return output, err
}
