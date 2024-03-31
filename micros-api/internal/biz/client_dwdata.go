package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "micros-api/api/dwdata/v2"
	pbV3 "micros-api/api/dwdata/v3"
	"micros-api/internal/biz/dto"
)

type ClientDwDataRepo interface {
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

	GetClient(ctx context.Context) pb.DwdataServiceClient
	GetClientV3(ctx context.Context) pbV3.DwdataServiceClient
}

type ClientDwDataUsecase struct {
	repo ClientDwDataRepo
	log  *log.Helper
}

func NewClientDwDataUsecase(repo ClientDwDataRepo, logger log.Logger) *ClientDwDataUsecase {
	return &ClientDwDataUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *ClientDwDataUsecase) GetClientV3(ctx context.Context) pbV3.DwdataServiceClient {
	uc.log.WithContext(ctx).Infof("biz.ClientDwDataUsecase.GetClientV3")
	return uc.repo.GetClientV3(ctx)

}

func (uc *ClientDwDataUsecase) GetClient(ctx context.Context) pb.DwdataServiceClient {
	uc.log.WithContext(ctx).Infof("biz.ClientDwDataUsecase.GetClient")
	return uc.repo.GetClient(ctx)
}

func (uc *ClientDwDataUsecase) GetShareholders(ctx context.Context, uscId string) (*pb.GetShareholdersResp, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDependencyDataUsecase.GetShareholders %s", uscId)
	return uc.repo.GetShareholders(ctx, uscId)
}

func (uc *ClientDwDataUsecase) GetInvestments(ctx context.Context, uscId string) (*pb.GetInvestmentResp, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDependencyDataUsecase.GetInvestments %s", uscId)
	return uc.repo.GetInvestments(ctx, uscId)
}

func (uc *ClientDwDataUsecase) GetBranches(ctx context.Context, uscId string) (*pb.GetBranchesResp, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDependencyDataUsecase.GetBranches %s", uscId)
	return uc.repo.GetBranches(ctx, uscId)
}

func (uc *ClientDwDataUsecase) GetEntIdent(ctx context.Context, name string) (string, error) {
	uc.log.WithContext(ctx).Infof("biz.ClientDwDataUsecase.GetEntIdent %s", name)
	return uc.repo.GetEntIdent(ctx, name)
}

func (uc *ClientDwDataUsecase) GetEntInfo(ctx context.Context, uscId string) (*dto.EnterpriseInfo, error) {
	uc.log.WithContext(ctx).Infof("biz.ClientDwDataUsecase.GetEntInfo %s", uscId)
	return uc.repo.GetEntInfo(ctx, uscId)
}

func (uc *ClientDwDataUsecase) GetEntCredential(ctx context.Context, uscId string) (*[]dto.EnterpriseCertification, error) {
	uc.log.WithContext(ctx).Infof("biz.ClientDwDataUsecase.GetEntCredential %s", uscId)
	return uc.repo.GetEntCredential(ctx, uscId)
}

func (uc *ClientDwDataUsecase) GetEntRankingList(ctx context.Context, uscId string) (*[]dto.EnterpriseRankingList, error) {
	uc.log.WithContext(ctx).Infof("biz.ClientDwDataUsecase.GetEntRankingList %s", uscId)
	return uc.repo.GetEntRankingList(ctx, uscId)
}

func (uc *ClientDwDataUsecase) GetEntIndustry(ctx context.Context, uscId string) (*[]string, error) {
	uc.log.WithContext(ctx).Infof("biz.ClientDwDataUsecase.GetEntIndustry %s", uscId)
	return uc.repo.GetEntIndustry(ctx, uscId)
}

func (uc *ClientDwDataUsecase) GetEntProduct(ctx context.Context, uscId string) (*[]string, error) {
	uc.log.WithContext(ctx).Infof("biz.ClientDwDataUsecase.GetEntProduct %s", uscId)
	return uc.repo.GetEntProduct(ctx, uscId)
}

func (uc *ClientDwDataUsecase) GetEquityTransparency(ctx context.Context, uscId string) (*dto.EnterpriseEquityTransparency, error) {
	uc.log.WithContext(ctx).Infof("biz.ClientDwDataUsecase.GetEquityTransparency %s", uscId)
	return uc.repo.GetEquityTransparency(ctx, uscId)
}
