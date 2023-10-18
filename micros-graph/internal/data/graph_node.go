package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-graph/internal/biz"
	"micros-graph/internal/biz/dto"
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

func (repo *GraphNodeRepo) GetNode(ctx context.Context, id int64) (*dto.Node, error) {
	return nil, nil
}
func (repo *GraphNodeRepo) GetNodes(ctx context.Context, ids []int64) ([]*dto.Node, error) {
	return nil, nil
}
func (repo *GraphNodeRepo) GetNodesByProps(ctx context.Context, labelScope []string, props map[string]interface{}, p dto.PaginationReq) ([]*dto.Node, int64, error) {
	return nil, 0, nil
}
