package biz

import (
	"brillinkmicros/internal/biz/dto"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type GraphNodeRepo interface {
	GetNode(ctx context.Context, id string) (*dto.Node, error)
	GetNodes(ctx context.Context, ids []string) ([]*dto.Node, error)
	GetChildren(ctx context.Context, id string, f *dto.PathFilter) ([]*dto.Node, error)
	CountChildren(ctx context.Context, id string, f *dto.PathFilter) (int64, error)
}

type GraphNodeUsecase struct {
	repo GraphNodeRepo
	log  *log.Helper
}

func NewGraphNodeUsecase(repo GraphNodeRepo, logger log.Logger) *GraphNodeUsecase {
	return &GraphNodeUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *GraphNodeUsecase) GetNode(ctx context.Context, id string) (*dto.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphNodeUsecase.GetNode id=%d", id)
	return uc.repo.GetNode(ctx, id)
}

func (uc *GraphNodeUsecase) GetNodes(ctx context.Context, ids []string) ([]*dto.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphNodeUsecase.GetNode")
	return uc.repo.GetNodes(ctx, ids)
}

func (uc *GraphNodeUsecase) GetChildren(ctx context.Context, id string, f *dto.PathFilter) ([]*dto.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphNodeUsecase.GetChildren id=%d", id)
	return uc.repo.GetChildren(ctx, id, f)
}

func (uc *GraphNodeUsecase) CountChildren(ctx context.Context, id string, f *dto.PathFilter) (int64, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphNodeUsecase.CountChildren id=%d", id)
	return uc.repo.CountChildren(ctx, id, f)
}
