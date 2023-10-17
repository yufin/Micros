package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-graph/internal/biz/dto"
)

type GraphEdgeRepo interface {
	GetEdge(ctx context.Context, sourceId int64, targetId int64, relType string, rank int64) (*dto.Edge, error)
	GetEdges(ctx context.Context, sourceId int64, targetId int64, relType string, p dto.PaginationReq) ([]*dto.Edge, int64, error)
	GetEdgesByParams(ctx context.Context, sourceId int64, targetId int64, relType string, rank int64, props map[string]interface{}, p dto.PaginationReq) ([]*dto.Node, int64, error)
}

type GraphEdgeUsecase struct {
	repo GraphEdgeRepo
	log  *log.Helper
}

func NewGraphEdgeUsecase(repo GraphEdgeRepo, logger log.Logger) *GraphEdgeUsecase {
	return &GraphEdgeUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *GraphEdgeUsecase) GetEdge(ctx context.Context, sourceId int64, targetId int64, relType string, rank int64) (*dto.Edge, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphEdgeUsecase.GetEdge sourceId=%d, targetId=%d, label=%s, rank=%d", sourceId, targetId, relType, rank)
	return uc.repo.GetEdge(ctx, sourceId, targetId, relType, rank)
}

func (uc *GraphEdgeUsecase) GetEdges(ctx context.Context, sourceId int64, targetId int64, relType string, p dto.PaginationReq) ([]*dto.Edge, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphEdgeUsecase.GetEdges sourceId=%d, targetId=%d, label=%s", sourceId, targetId, relType)
	return uc.repo.GetEdges(ctx, sourceId, targetId, relType, p)
}

func (uc *GraphEdgeUsecase) GetEdgesByParams(ctx context.Context, sourceId int64, targetId int64, relType string, rank int64, props map[string]interface{}, p dto.PaginationReq) ([]*dto.Node, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphEdgeUsecase.GetEdgesByParams sourceId=%d, targetId=%d, label=%s, rank=%d, props=%v", sourceId, targetId, relType, rank, props)
	return uc.repo.GetEdgesByParams(ctx, sourceId, targetId, relType, rank, props, p)
}
