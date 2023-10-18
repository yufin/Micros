package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-graph/internal/biz"
	"micros-graph/internal/biz/dto"
)

type GraphEdgeRepo struct {
	data *Data
	log  *log.Helper
}

func NewGraphEdgeRepo(data *Data, logger log.Logger) biz.GraphEdgeRepo {
	return &GraphEdgeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *GraphEdgeRepo) GetEdge(ctx context.Context, sourceId int64, targetId int64, relType string, rank int64) (*dto.Edge, error) {
	return nil, nil
}
func (repo *GraphEdgeRepo) GetEdges(ctx context.Context, sourceId int64, targetId int64, relType string, p dto.PaginationReq) ([]*dto.Edge, int64, error) {
	return nil, 0, nil
}
func (repo *GraphEdgeRepo) GetEdgesByProps(ctx context.Context, sourceId int64, targetId int64, relType string, props map[string]interface{}, p dto.PaginationReq) ([]*dto.Node, int64, error) {
	return nil, 0, nil
}
