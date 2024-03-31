package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-api/internal/biz/dto"
)

type RcDecisionFactorV3Repo interface {
	CountByUscIdAndUserId(ctx context.Context, uscId string) (int64, error)
	Insert(ctx context.Context, data *dto.RcDecisionFactor) (int64, error)
	InsertClaimNoDupe(ctx context.Context, data *dto.RcContentFactorClaimV3) (int64, error)
	ListReportClaimed(ctx context.Context, page *dto.PaginationReq, kwd string) (*[]dto.ListReportInfo, dto.PaginationInfo, error)
	GetWithDataScope(ctx context.Context, id int64) (*dto.RcDecisionFactor, error)
	CheckContentIdAccessible(ctx context.Context, contentId int64) (bool, error)
	GetByContentIdWithDataScope(ctx context.Context, contentId int64) (*dto.RcDecisionFactorClaimed, error)
	GetClaimRecord(ctx context.Context, claimId int64) (*dto.RcContentFactorClaimV3, error)
	InsertClaim(ctx context.Context, data *dto.RcContentFactorClaimV3) (int64, error)
	GetLatestIdByUscIdAndUserId(ctx context.Context, uscId string) (int64, error)
	ListCompanies(ctx context.Context, page *dto.PaginationReq, kwd string, version string) (*[]dto.ListCompaniesLatest, int64, error)
	ListReportClaimedByUscId(ctx context.Context, page *dto.PaginationReq, uscId string, version string) (*[]dto.ListReportInfo, int64, error)
	SyncClaimed(ctx context.Context, uscId string, version string) error
	ListCompaniesWaiting(ctx context.Context, page *dto.PaginationReq, kwd string, version string) (*[]dto.ListCompaniesWaitingResp, int64, error)
}

type RcDecisionFactorV3Usecase struct {
	repo RcDecisionFactorV3Repo
	log  *log.Helper
}

func NewRcDecisionFactorV3Usecase(repo RcDecisionFactorV3Repo, logger log.Logger) *RcDecisionFactorV3Usecase {
	return &RcDecisionFactorV3Usecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *RcDecisionFactorV3Usecase) GetLatestIdByUscIdAndUserId(ctx context.Context, uscId string) (int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.GetLatestIdByUscIdAndUserId %s", uscId)
	return uc.repo.GetLatestIdByUscIdAndUserId(ctx, uscId)
}

func (uc *RcDecisionFactorV3Usecase) GetClaimRecord(ctx context.Context, claimId int64) (*dto.RcContentFactorClaimV3, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.GetClaimRecord %d", claimId)
	return uc.repo.GetClaimRecord(ctx, claimId)
}

func (uc *RcDecisionFactorV3Usecase) GetByContentIdWithDataScope(ctx context.Context, contentId int64) (*dto.RcDecisionFactorClaimed, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.GetByContentIdWithDataScope %d", contentId)
	return uc.repo.GetByContentIdWithDataScope(ctx, contentId)
}

func (uc *RcDecisionFactorV3Usecase) CountByUscIdAndUserId(ctx context.Context, uscId string) (int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.CountByUscIdAndUserId %s", uscId)
	return uc.repo.CountByUscIdAndUserId(ctx, uscId)
}

func (uc *RcDecisionFactorV3Usecase) Insert(ctx context.Context, data *dto.RcDecisionFactor) (int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.Insert %v", data)
	return uc.repo.Insert(ctx, data)
}

func (uc *RcDecisionFactorV3Usecase) InsertClaimNoDupe(ctx context.Context, data *dto.RcContentFactorClaimV3) (int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.InsertClaimNoDupe %v", data)
	return uc.repo.InsertClaimNoDupe(ctx, data)
}

func (uc *RcDecisionFactorV3Usecase) ListReportClaimed(ctx context.Context, page *dto.PaginationReq, kwd string) (*[]dto.ListReportInfo, dto.PaginationInfo, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.ListReportClaimed %v", page)
	return uc.repo.ListReportClaimed(ctx, page, kwd)
}

func (uc *RcDecisionFactorV3Usecase) GetWithDataScope(ctx context.Context, id int64) (*dto.RcDecisionFactor, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.GetWithDataScope %v", id)
	return uc.repo.GetWithDataScope(ctx, id)
}

func (uc *RcDecisionFactorV3Usecase) CheckContentIdAccessible(ctx context.Context, contentId int64) (bool, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.CheckContentIdAccessible %v", contentId)
	return uc.repo.CheckContentIdAccessible(ctx, contentId)
}

func (uc *RcDecisionFactorV3Usecase) InsertClaim(ctx context.Context, data *dto.RcContentFactorClaimV3) (int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.InsertClaim %v", data)
	return uc.repo.InsertClaim(ctx, data)
}

func (uc *RcDecisionFactorV3Usecase) ListCompanies(ctx context.Context, page *dto.PaginationReq, kwd string, version string) (*[]dto.ListCompaniesLatest, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.ListCompanies %v", page)
	return uc.repo.ListCompanies(ctx, page, kwd, version)
}

func (uc *RcDecisionFactorV3Usecase) ListReportClaimedByUscId(ctx context.Context, page *dto.PaginationReq, uscId string, version string) (*[]dto.ListReportInfo, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.ListReportClaimedByUscId %v", uscId)
	return uc.repo.ListReportClaimedByUscId(ctx, page, uscId, version)
}

func (uc *RcDecisionFactorV3Usecase) SyncClaimed(ctx context.Context, uscId string, version string) error {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.SyncClaimed %v, %v", uscId, version)
	return uc.repo.SyncClaimed(ctx, uscId, version)
}

func (uc *RcDecisionFactorV3Usecase) ListCompaniesWaiting(ctx context.Context, page *dto.PaginationReq, kwd string, version string) (*[]dto.ListCompaniesWaitingResp, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDecisionFactorUsecase.ListCompaniesWaiting %v", version)
	return uc.repo.ListCompaniesWaiting(ctx, page, kwd, version)
}
