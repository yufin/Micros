package biz

import (
	"brillinkmicros/internal/biz/dto"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphRepo interface {
	GetNode(ctx context.Context, id string) (*dto.Node, error)
	GetNodes(ctx context.Context, ids []string) ([]*dto.Node, error)
	GetChildren(ctx context.Context, id string, f *dto.PathFilter, p *dto.PaginationReq) ([]*dto.Node, error)
	CountChildren(ctx context.Context, id string, f *dto.PathFilter, amount *int64) error
	GetTitleAutoComplete(ctx context.Context, f *dto.PathFilter, p *dto.PaginationReq, kw string) ([]*dto.TitleAutoCompleteRes, error)
	CountTitleAutoComplete(ctx context.Context, f *dto.PathFilter, kw string, amount *int64) error
	GetPathBetween(ctx context.Context, sourceId string, targetId string, f *dto.PathFilter, p *dto.PaginationReq) ([]neo4j.Path, error)
}

type GraphUsecase struct {
	repo GraphRepo
	log  *log.Helper
}

func NewGraphUsecase(repo GraphRepo, logger log.Logger) *GraphUsecase {
	return &GraphUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *GraphUsecase) GetNode(ctx context.Context, id string) (*dto.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetNode id=%d", id)
	return uc.repo.GetNode(ctx, id)
}

func (uc *GraphUsecase) GetNodes(ctx context.Context, ids []string) ([]*dto.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetNode")
	return uc.repo.GetNodes(ctx, ids)
}

func (uc *GraphUsecase) GetChildren(ctx context.Context, id string, f *dto.PathFilter, p *dto.PaginationReq) ([]*dto.Node, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetChildren id=%d", id)
	return uc.repo.GetChildren(ctx, id, f, p)
}

func (uc *GraphUsecase) CountChildren(ctx context.Context, id string, f *dto.PathFilter, amount *int64) error {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.CountChildren id=%d", id)
	return uc.repo.CountChildren(ctx, id, f, amount)
}

func (uc *GraphUsecase) GetTitleAutoComplete(ctx context.Context, f *dto.PathFilter, p *dto.PaginationReq, kw string) ([]*dto.TitleAutoCompleteRes, error) {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.GetTitleAutoComplete")
	return uc.repo.GetTitleAutoComplete(ctx, f, p, kw)
}

func (uc *GraphUsecase) CountTitleAutoComplete(ctx context.Context, f *dto.PathFilter, kw string, amount *int64) error {
	uc.log.WithContext(ctx).Infof("biz.GraphUsecase.CountTitleAutoComplete")
	return uc.repo.CountTitleAutoComplete(ctx, f, kw, amount)
}
