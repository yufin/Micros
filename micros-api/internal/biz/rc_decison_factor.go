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
	GetWithDataScope(ctx context.Context, id int64) (*dto.RcDecisionFactor, error)
	CheckContentIdAccessible(ctx context.Context, contentId int64) (bool, error)
	GetByContentIdWithDataScope(ctx context.Context, contentId int64) (*dto.RcDecisionFactorClaimed, error)
	GetClaimRecord(ctx context.Context, claimId int64) (*dto.RcContentFactorClaim, error)
	InsertClaim(ctx context.Context, data *dto.RcContentFactorClaim) (int64, error)
}

type RcDecisionFactorUsecase struct {
	repo RcDecisionFactorRepo
	log  *log.Helper
}

func NewRcDecisionFactorUsecase(repo RcDecisionFactorRepo, logger log.Logger) *RcDecisionFactorUsecase {
	return &RcDecisionFactorUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *RcDecisionFactorUsecase) GetClaimRecord(ctx context.Context, claimId int64) (*dto.RcContentFactorClaim, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.GetClaimRecord %d", claimId)
	return uc.repo.GetClaimRecord(ctx, claimId)
}

func (uc *RcDecisionFactorUsecase) GetByContentIdWithDataScope(ctx context.Context, contentId int64) (*dto.RcDecisionFactorClaimed, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.GetByContentIdWithDataScope %d", contentId)
	return uc.repo.GetByContentIdWithDataScope(ctx, contentId)
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

func (uc *RcDecisionFactorUsecase) GetWithDataScope(ctx context.Context, id int64) (*dto.RcDecisionFactor, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.GetWithDataScope %v", id)
	return uc.repo.GetWithDataScope(ctx, id)
}

func (uc *RcDecisionFactorUsecase) CheckContentIdAccessible(ctx context.Context, contentId int64) (bool, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.CheckContentIdAccessible %v", contentId)
	return uc.repo.CheckContentIdAccessible(ctx, contentId)
}
