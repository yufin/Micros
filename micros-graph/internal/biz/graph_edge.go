package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-graph/internal/biz/dto"
)

type GraphEdgeRepo interface {
	GetEdge(ctx context.Context, sourceId int64, targetId int64, relType string, rank int64) (*dto.Edge, error)
	GetEdges(ctx context.Context, sourceId int64, targetId int64, relType string, p dto.PaginationReq) ([]*dto.Edge, int64, error)
	GetEdgesByProps(ctx context.Context, sourceId int64, targetId int64, relType string, props map[string]interface{}, p dto.PaginationReq) ([]*dto.Node, int64, error)
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

func (uc *GraphEdgeUsecase) GetEdgesByProps(ctx context.Context, sourceId int64, targetId int64, relType string, props map[string]interface{}, p dto.PaginationReq) ([]*dto.Node, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphEdgeUsecase.GetEdgesByParams sourceId=%d, targetId=%d, label=%s", sourceId, targetId, relType)
	return uc.repo.GetEdgesByProps(ctx, sourceId, targetId, relType, props, p)
}
