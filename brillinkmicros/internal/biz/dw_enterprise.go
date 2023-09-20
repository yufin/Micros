package biz

import (
	"brillinkmicros/internal/biz/dto"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type DwEnterpriseRepo interface {
	GetEntIdent(ctx context.Context, name string) (string, error)
	GetEntInfo(ctx context.Context, uscId string) (*dto.EnterpriseInfo, error)
	GetEntCredential(ctx context.Context, uscId string) (*[]dto.EnterpriseCertification, error)
	GetEntRankingList(ctx context.Context, uscId string) (*[]dto.EnterpriseRankingList, error)
	GetEntIndustry(ctx context.Context, uscId string) (*[]string, error)
	GetEntProduct(ctx context.Context, uscId string) (*[]string, error)
	GetEquityTransparency(ctx context.Context, uscId string) (*dto.EnterpriseEquityTransparency, error)
}

type DwEnterpriseUsecase struct {
	repo DwEnterpriseRepo
	log  *log.Helper
}

func NewDwEnterpriseUsecase(repo DwEnterpriseRepo, logger log.Logger) *DwEnterpriseUsecase {
	return &DwEnterpriseUsecase{repo: repo, log: log.NewHelper(logger)}
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
