package data

import (
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/biz/dto"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphNodeRepo struct {
	data *Data
	log  *log.Helper
}

func NewGraphNodeRepo(data *Data, logger log.Logger) biz.GraphNodeRepo {
	return &GraphNodeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *GraphNodeRepo) GetNode(ctx context.Context, id string) (*dto.Node, error) {
	cypher := "MATCH (n {id: $id} RETURN n;"
	res, err := repo.data.Neo.CypherQuery(ctx, cypher, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	n, found := res[0].Get("n")
	if !found {
		return nil, errors.New("node with specified id not found")
	}
	var node dto.Node
	node.Gen(n.(*neo4j.Node))
	return &node, nil
}

func (repo *GraphNodeRepo) GetNodes(ctx context.Context, ids []string) ([]*dto.Node, error) {
	return nil, nil
}

func (repo *GraphNodeRepo) GetChildren(ctx context.Context, id string, f *dto.PathFilter) ([]*dto.Node, error) {
	return []*dto.Node{}, nil
}

func (repo *GraphNodeRepo) CountChildren(ctx context.Context, id string, f *dto.PathFilter) (int64, error) {
	return 0, nil
}
