package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"micros-api/internal/biz/dto"
)

type GraphRepo interface {
	GetNode(ctx context.Context, id string) (neo4j.Node, error)
	GetNodes(ctx context.Context, ids []string) (*[]neo4j.Node, error)
	GetChildren(ctx context.Context, id string, f dto.PathFilter, p dto.PaginationReq) (*[]neo4j.Node, error)

	GetRelTypeAvailable(ctx context.Context, id string, targetType int) ([]string, error)

	CountChildren(ctx context.Context, id string, f dto.PathFilter, amount *int64) error

	GetTitleAutoComplete(ctx context.Context, f dto.PathFilter, p dto.PaginationReq, kw string) (*[]dto.TitleAutoCompleteRes, error)
	CountTitleAutoComplete(ctx context.Context, f dto.PathFilter, kw string, amount *int64) error

	GetPathTo(ctx context.Context, sourceId string, targetId string, maxDepth int, relScope []string) (*[]neo4j.Path, error)
	GetPathBetween(ctx context.Context, sourceId string, targetId string, maxDepth int, relScope []string) (*[]neo4j.Path, error)
	GetPathBetweenByIds(ctx context.Context, sourceId string, targetIds []string, f *dto.PathFilter) (*[]neo4j.Path, error)

	GetPathExpand(ctx context.Context, sourceId string, depth uint32, limit uint32, relScope []string) (*[]neo4j.Path, error)
	GetPathToChildren(ctx context.Context, sourceId string, p dto.PaginationReq, scopeRelType []string) (*[]neo4j.Path, int64, error)
	GetPathToParent(ctx context.Context, targetId string, p dto.PaginationReq, scopeRelType []string) (*[]neo4j.Path, int64, error)
}

type GraphUsecase struct {
	repo GraphRepo
	log  *log.Helper
}

func NewGraphUsecase(repo GraphRepo, logger log.Logger) *GraphUsecase {
	return &GraphUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *GraphUsecase) GetNode(ctx context.Context, id string) (neo4j.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetNode id=%d", id)
	return uc.repo.GetNode(ctx, id)
}

func (uc *GraphUsecase) GetNodes(ctx context.Context, ids []string) (*[]neo4j.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetNode")
	return uc.repo.GetNodes(ctx, ids)
}

func (uc *GraphUsecase) GetChildren(ctx context.Context, id string, f dto.PathFilter, p dto.PaginationReq) (*[]neo4j.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetChildren id=%d", id)
	return uc.repo.GetChildren(ctx, id, f, p)
}

func (uc *GraphUsecase) GetRelTypeAvailable(ctx context.Context, id string, targetType int) ([]string, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetRelTypeAvailableToChildren id=%d", id)
	return uc.repo.GetRelTypeAvailable(ctx, id, targetType)
}

func (uc *GraphUsecase) CountChildren(ctx context.Context, id string, f dto.PathFilter, amount *int64) error {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.CountChildren id=%d", id)
	return uc.repo.CountChildren(ctx, id, f, amount)
}

func (uc *GraphUsecase) GetTitleAutoComplete(ctx context.Context, f dto.PathFilter, p dto.PaginationReq, kw string) (*[]dto.TitleAutoCompleteRes, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetTitleAutoComplete")
	return uc.repo.GetTitleAutoComplete(ctx, f, p, kw)
}

func (uc *GraphUsecase) CountTitleAutoComplete(ctx context.Context, f dto.PathFilter, kw string, amount *int64) error {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.CountTitleAutoComplete")
	return uc.repo.CountTitleAutoComplete(ctx, f, kw, amount)
}

func (uc *GraphUsecase) GetPathTo(ctx context.Context, sourceId string, targetId string, maxDepth int, relScope []string) (*[]neo4j.Path, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetPathTo")
	return uc.repo.GetPathTo(ctx, sourceId, targetId, maxDepth, relScope)
}

func (uc *GraphUsecase) GetPathBetween(ctx context.Context, sourceId string, targetId string, maxDepth int, relScope []string) (*[]neo4j.Path, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetPathTo")
	return uc.repo.GetPathTo(ctx, sourceId, targetId, maxDepth, relScope)
}

func (uc *GraphUsecase) GetPathBetweenByIds(ctx context.Context, sourceId string, targetIds []string, f *dto.PathFilter) (*[]neo4j.Path, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetPathBetweenByIds")
	return uc.repo.GetPathBetweenByIds(ctx, sourceId, targetIds, f)
}

func (uc *GraphUsecase) GetPathExpand(ctx context.Context, sourceId string, depth uint32, limit uint32, relScope []string) (*[]neo4j.Path, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetPathExpand")
	return uc.repo.GetPathExpand(ctx, sourceId, depth, limit, relScope)
}

func (uc *GraphUsecase) GetPathToChildren(ctx context.Context, sourceId string, p dto.PaginationReq, scopeRelType []string) (*[]neo4j.Path, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetPathToChildren")
	return uc.repo.GetPathToChildren(ctx, sourceId, p, scopeRelType)
}

func (uc *GraphUsecase) GetPathToParent(ctx context.Context, targetId string, p dto.PaginationReq, scopeRelType []string) (*[]neo4j.Path, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetPathToParent")
	return uc.repo.GetPathToParent(ctx, targetId, p, scopeRelType)
}
