package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "micros-api/api/dwdata/v2"
	"micros-api/internal/biz/dto"
)

type DwEnterpriseRepo interface {
	GetEntIdent(ctx context.Context, name string) (string, error)
	GetEntInfo(ctx context.Context, uscId string) (*dto.EnterpriseInfo, error)
	GetEntCredential(ctx context.Context, uscId string) (*[]dto.EnterpriseCertification, error)
	GetEntRankingList(ctx context.Context, uscId string) (*[]dto.EnterpriseRankingList, error)
	GetEntIndustry(ctx context.Context, uscId string) (*[]string, error)
	GetEntProduct(ctx context.Context, uscId string) (*[]string, error)
	GetEquityTransparency(ctx context.Context, uscId string) (*dto.EnterpriseEquityTransparency, error)

	GetShareholders(ctx context.Context, uscId string) (*pb.GetShareholdersResp, error)
	GetInvestments(ctx context.Context, uscId string) (*pb.GetInvestmentResp, error)
	GetBranches(ctx context.Context, uscId string) (*pb.GetBranchesResp, error)
}

type DwEnterpriseUsecase struct {
	repo DwEnterpriseRepo
	log  *log.Helper
}

func NewDwEnterpriseUsecase(repo DwEnterpriseRepo, logger log.Logger) *DwEnterpriseUsecase {
	return &DwEnterpriseUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *DwEnterpriseUsecase) GetShareholders(ctx context.Context, uscId string) (*pb.GetShareholdersResp, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDependencyDataUsecase.GetShareholders %s", uscId)
	return uc.repo.GetShareholders(ctx, uscId)
}

func (uc *DwEnterpriseUsecase) GetInvestments(ctx context.Context, uscId string) (*pb.GetInvestmentResp, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDependencyDataUsecase.GetInvestments %s", uscId)
	return uc.repo.GetInvestments(ctx, uscId)
}

func (uc *DwEnterpriseUsecase) GetBranches(ctx context.Context, uscId string) (*pb.GetBranchesResp, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDependencyDataUsecase.GetBranches %s", uscId)
	return uc.repo.GetBranches(ctx, uscId)
}

func (uc *DwEnterpriseUsecase) GetEntIdent(ctx context.Context, name string) (string, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetEntIdent %s", name)
	return uc.repo.GetEntIdent(ctx, name)
}

func (uc *DwEnterpriseUsecase) GetEntInfo(ctx context.Context, uscId string) (*dto.EnterpriseInfo, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetEntInfo %s", uscId)
	return uc.repo.GetEntInfo(ctx, uscId)
}

func (uc *DwEnterpriseUsecase) GetEntCredential(ctx context.Context, uscId string) (*[]dto.EnterpriseCertification, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetEntCredential %s", uscId)
	return uc.repo.GetEntCredential(ctx, uscId)
}

func (uc *DwEnterpriseUsecase) GetEntRankingList(ctx context.Context, uscId string) (*[]dto.EnterpriseRankingList, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetEntRankingList %s", uscId)
	return uc.repo.GetEntRankingList(ctx, uscId)
}

func (uc *DwEnterpriseUsecase) GetEntIndustry(ctx context.Context, uscId string) (*[]string, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetEntIndustry %s", uscId)
	return uc.repo.GetEntIndustry(ctx, uscId)
}

func (uc *DwEnterpriseUsecase) GetEntProduct(ctx context.Context, uscId string) (*[]string, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetEntProduct %s", uscId)
	return uc.repo.GetEntProduct(ctx, uscId)
}

func (uc *DwEnterpriseUsecase) GetEquityTransparency(ctx context.Context, uscId string) (*dto.EnterpriseEquityTransparency, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetEquityTransparency %s", uscId)
	return uc.repo.GetEquityTransparency(ctx, uscId)
}
