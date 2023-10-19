package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-api/internal/biz/dto"
)

type RcDecisionFactorRepo interface {
	CountByUscIdAndUserId(ctx context.Context, uscId string) (int64, error)
	Insert(ctx context.Context, data *dto.RcDecisionFactor) (int64, error)
	InsertClaimNoDupe(ctx context.Context, data *dto.RcContentFactorClaim) (int64, error)
	ListReportClaimed(ctx context.Context, page *dto.PaginationReq, kwd string) (*[]dto.ListReportInfo, dto.PaginationInfo, error)
}

type RcDecisionFactorUsecase struct {
	repo RcDecisionFactorRepo
	log  *log.Helper
}

func NewRcDecisionFactorUsecase(repo RcDecisionFactorRepo, logger log.Logger) *RcDecisionFactorUsecase {
	return &RcDecisionFactorUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *RcDecisionFactorUsecase) CountByUscIdAndUserId(ctx context.Context, uscId string) (int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.CountByUscIdAndUserId %s", uscId)
	return uc.repo.CountByUscIdAndUserId(ctx, uscId)
}

func (uc *RcDecisionFactorUsecase) Insert(ctx context.Context, data *dto.RcDecisionFactor) (int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.Insert %v", data)
	return uc.repo.Insert(ctx, data)
}

func (uc *RcDecisionFactorUsecase) InsertClaimNoDupe(ctx context.Context, data *dto.RcContentFactorClaim) (int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.InsertClaimNoDupe %v", data)
	return uc.repo.InsertClaimNoDupe(ctx, data)
}

func (uc *RcDecisionFactorUsecase) ListReportClaimed(ctx context.Context, page *dto.PaginationReq, kwd string) (*[]dto.ListReportInfo, dto.PaginationInfo, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.ListReportClaimed %v", page)
	return uc.repo.ListReportClaimed(ctx, page, kwd)
}
