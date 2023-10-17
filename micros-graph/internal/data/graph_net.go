package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-graph/internal/biz"
	"micros-graph/internal/biz/dto"
)

type GraphNetRepo struct {
	data *Data
	log  *log.Helper
}

func NewGraphNetRepo(data *Data, logger log.Logger) biz.GraphNetRepo {
	return &GraphNetRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *GraphNetRepo) GetChildren(ctx context.Context, sourceId int64, relTypeScope []string, labelScope []string, p *dto.PaginationReq) (*dto.Net, int64, error) {
	return nil, 0, nil
}
