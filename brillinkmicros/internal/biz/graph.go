package biz

import (
	dto2 "brillinkmicros/internal/biz/dto"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphRepo interface {
	GetNode(ctx context.Context, id string) (*dto2.Node, error)
	GetNodes(ctx context.Context, ids []string) (*[]dto2.Node, error)
	GetChildren(ctx context.Context, id string, f dto2.PathFilter, p dto2.PaginationReq) (*[]dto2.Node, error)

	CountChildren(ctx context.Context, id string, f dto2.PathFilter, amount *int64) error

	GetTitleAutoComplete(ctx context.Context, f dto2.PathFilter, p dto2.PaginationReq, kw string) (*[]dto2.TitleAutoCompleteRes, error)
	CountTitleAutoComplete(ctx context.Context, f dto2.PathFilter, kw string, amount *int64) error

	GetPathBetween(ctx context.Context, sourceId string, targetId string, f dto2.PathFilter) (*[]neo4j.Path, error)
	GetPathBetweenByIds(ctx context.Context, sourceId string, targetIds []string, f *dto2.PathFilter) (*[]neo4j.Path, error)

	GetPathExpand(ctx context.Context, sourceId string, depth uint32, limit uint32, f *dto2.PathFilter) (*[]neo4j.Path, error)
	GetPathToChildren(ctx context.Context, sourceId string, p dto2.PaginationReq, f *dto2.PathFilter) (*[]neo4j.Path, int64, error)
	GetPathToParent(ctx context.Context, targetId string, p dto2.PaginationReq, f *dto2.PathFilter) (*[]neo4j.Path, int64, error)
}

type GraphUsecase struct {
	repo GraphRepo
	log  *log.Helper
}

func NewGraphUsecase(repo GraphRepo, logger log.Logger) *GraphUsecase {
	return &GraphUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *GraphUsecase) GetNode(ctx context.Context, id string) (*dto2.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetNode id=%d", id)
	return uc.repo.GetNode(ctx, id)
}

func (uc *GraphUsecase) GetNodes(ctx context.Context, ids []string) (*[]dto2.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetNode")
	return uc.repo.GetNodes(ctx, ids)
}

func (uc *GraphUsecase) GetChildren(ctx context.Context, id string, f dto2.PathFilter, p dto2.PaginationReq) (*[]dto2.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetChildren id=%d", id)
	return uc.repo.GetChildren(ctx, id, f, p)
}

func (uc *GraphUsecase) CountChildren(ctx context.Context, id string, f dto2.PathFilter, amount *int64) error {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.CountChildren id=%d", id)
	return uc.repo.CountChildren(ctx, id, f, amount)
}

func (uc *GraphUsecase) GetTitleAutoComplete(ctx context.Context, f dto2.PathFilter, p dto2.PaginationReq, kw string) (*[]dto2.TitleAutoCompleteRes, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetTitleAutoComplete")
	return uc.repo.GetTitleAutoComplete(ctx, f, p, kw)
}

func (uc *GraphUsecase) CountTitleAutoComplete(ctx context.Context, f dto2.PathFilter, kw string, amount *int64) error {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.CountTitleAutoComplete")
	return uc.repo.CountTitleAutoComplete(ctx, f, kw, amount)
}

func (uc *GraphUsecase) GetPathBetween(ctx context.Context, sourceId string, targetId string, f dto2.PathFilter) (*[]neo4j.Path, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetPathBetween")
	return uc.repo.GetPathBetween(ctx, sourceId, targetId, f)
}

func (uc *GraphUsecase) GetPathBetweenByIds(ctx context.Context, sourceId string, targetIds []string, f *dto2.PathFilter) (*[]neo4j.Path, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetPathBetweenByIds")
	return uc.repo.GetPathBetweenByIds(ctx, sourceId, targetIds, f)
}

func (uc *GraphUsecase) GetPathExpand(ctx context.Context, sourceId string, depth uint32, limit uint32, f *dto2.PathFilter) (*[]neo4j.Path, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetPathExpand")
	return uc.repo.GetPathExpand(ctx, sourceId, depth, limit, f)
}

func (uc *GraphUsecase) GetPathToChildren(ctx context.Context, sourceId string, p dto2.PaginationReq, f *dto2.PathFilter) (*[]neo4j.Path, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetPathToChildren")
	return uc.repo.GetPathToChildren(ctx, sourceId, p, f)
}

func (uc *GraphUsecase) GetPathToParent(ctx context.Context, targetId string, p dto2.PaginationReq, f *dto2.PathFilter) (*[]neo4j.Path, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetPathToParent")
	return uc.repo.GetPathToParent(ctx, targetId, p, f)
}
