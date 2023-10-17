package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-graph/internal/biz/dto"
)

type GraphNodeRepo interface {
	GetNode(ctx context.Context, id int64) (*dto.Node, error)
	GetNodes(ctx context.Context, ids []int64) ([]*dto.Node, error)
	GetNodesByParams(ctx context.Context, labelScope []string, props map[string]interface{}, p dto.PaginationReq) ([]*dto.Node, int64, error)
}

type GraphNodeUsecase struct {
	repo GraphNodeRepo
	log  *log.Helper
}

func NewGraphNodeUsecase(repo GraphNodeRepo, logger log.Logger) *GraphNodeUsecase {
	return &GraphNodeUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *GraphNodeUsecase) GetNode(ctx context.Context, id int64) (*dto.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphNodeUsecase.GetNode id=%d", id)
	return uc.repo.GetNode(ctx, id)
}

func (uc *GraphNodeUsecase) GetNodes(ctx context.Context, ids []int64) ([]*dto.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphNodeUsecase.GetNodes id=%d", ids)
	return uc.repo.GetNodes(ctx, ids)
}

func (uc *GraphNodeUsecase) GetNodesByParams(ctx context.Context, labelScope []string, props map[string]interface{}, p dto.PaginationReq) ([]*dto.Node, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphNodeUsecase.GetNodesByParams labels=%v, props=%v", labelScope, props)
	return uc.repo.GetNodesByParams(ctx, labelScope, props, p)
}
