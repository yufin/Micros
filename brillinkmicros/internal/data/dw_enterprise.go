package data

import (
	dwdataV2 "brillinkmicros/api/dwdata/v2"
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/biz/dto"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type DwEnterpriseDataRepo struct {
	data *Data
	log  *log.Helper
}

func NewDwEnterpriseRepo(data *Data, logger log.Logger) biz.DwEnterpriseRepo {
	return &DwEnterpriseDataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *DwEnterpriseDataRepo) GetEntIdent(ctx context.Context, name string) (string, error) {
	resp, err := repo.data.DwDataClient.GetEnterpriseIdent(context.TODO(), &dwdataV2.GetEntIdentReq{EnterpriseName: name})
	if err != nil {
		return "", err
	}
	return resp.UscId, nil
}

func (repo *DwEnterpriseDataRepo) GetEntInfo(ctx context.Context, uscId string) (*dto.EnterpriseInfo, error) {
	resp, err := repo.data.DwDataClient.GetEnterpriseInfo(context.TODO(), &dwdataV2.GetEntInfoReq{UscId: uscId})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	return &dto.EnterpriseInfo{
		UscId:                         resp.UscId,
		EnterpriseTitle:               resp.EnterpriseTitle,
		EnterpriseTitleEn:             resp.EnterpriseTitleEn,
		BusinessRegistrationNumber:    resp.BusinessRegistrationNumber,
		EstablishedDate:               resp.EstablishDate,
		Region:                        resp.Region,
		ApprovedDate:                  resp.ApprovedDate,
		RegisteredAddress:             resp.RegisteredAddress,
		RegisteredCapital:             resp.RegisteredCapital,
		PaidInCapital:                 resp.PaidInCapital,
		EnterpriseType:                resp.EnterpriseType,
		StuffSize:                     resp.StuffSize,
		StuffInsuredNumber:            int(resp.StuffInsuredNumber),
		BusinessScope:                 resp.BusinessScope,
		ImportExportQualificationCode: resp.ImportExportQualificationCode,
		LegalRepresentative:           resp.LegalRepresentative,
		RegistrationAuthority:         resp.RegistrationAuthority,
		RegistrationStatus:            resp.RegistrationStatus,
		TaxpayerQualification:         resp.TaxpayerQualification,
		OrganizationCode:              resp.OrganizationCode,
		UrlQcc:                        resp.UrlQcc,
		UrlHomepage:                   resp.UrlHomepage,
		BusinessTermStart:             resp.BusinessTermStart,
		BusinessTermEnd:               resp.BusinessTermEnd,
	}, nil
}

func (repo *DwEnterpriseDataRepo) GetEntCredential(ctx context.Context, uscId string) (*[]dto.EnterpriseCertification, error) {
	resp, err := repo.data.DwDataClient.GetEnterpriseCredential(context.TODO(), &dwdataV2.GetEntInfoReq{UscId: uscId})
	if err != nil {
		return nil, err
	}
	data := make([]dto.EnterpriseCertification, 0)
	for _, v := range resp.Data {
		data = append(data, dto.EnterpriseCertification{
			UscId:                  v.UscId,
			CertificationTitle:     v.CertificationTitle,
			CertificationCode:      v.CertificationCode,
			CertificationLevel:     v.CertificationLevel,
			CertificationType:      v.CertificationType,
			CertificationSource:    v.CertificationSource,
			CertificationDate:      v.CertificationDate,
			CertificationTermStart: v.CertificationTermStart,
			CertificationTermEnd:   v.CertificationTermEnd,
			CertificationAuthority: v.CertificationAuthority,
		})
	}
	return &data, nil
}

func (repo *DwEnterpriseDataRepo) GetEntRankingList(ctx context.Context, uscId string) (*[]dto.EnterpriseRankingList, error) {
	resp, err := repo.data.DwDataClient.GetEnterpriseRankingList(context.TODO(), &dwdataV2.GetEntInfoReq{UscId: uscId})
	if err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, nil
	}
	data := make([]dto.EnterpriseRankingList, 0)
	for _, v := range resp.Data {
		data = append(data, dto.EnterpriseRankingList{
			UscId:                 v.UscId,
			RankingPosition:       int(v.RankingPosition),
			ListTitle:             v.ListTitle,
			ListType:              v.ListType,
			ListSource:            v.ListSource,
			ListParticipantsTotal: int(v.ListParticipantsTotal),
			ListPublishedDate:     v.ListPublishedDate,
			ListUrlQcc:            v.ListUrlQcc,
			ListUrlOrigin:         v.ListUrlOrigin,
		})
	}

	return &data, nil
}

func (repo *DwEnterpriseDataRepo) GetEntIndustry(ctx context.Context, uscId string) (*[]string, error) {
	resp, err := repo.data.DwDataClient.GetEnterpriseIndustry(context.TODO(), &dwdataV2.GetEntInfoReq{UscId: uscId})
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (repo *DwEnterpriseDataRepo) GetEntProduct(ctx context.Context, uscId string) (*[]string, error) {
	resp, err := repo.data.DwDataClient.GetEnterpriseProduct(context.TODO(), &dwdataV2.GetEntInfoReq{UscId: uscId})
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (repo *DwEnterpriseDataRepo) GetEquityTransparency(ctx context.Context, uscId string) (*dto.EnterpriseEquityTransparency, error) {
	resp, err := repo.data.DwDataClient.GetEntEquityTransparency(context.TODO(), &dwdataV2.GetEntInfoReq{UscId: uscId})
	if err != nil {
		return nil, err
	}
	if resp.UscId == "" {
		return nil, nil
	}
	info, err := repo.data.DwDataClient.GetEnterpriseInfo(context.TODO(), &dwdataV2.GetEntInfoReq{UscId: uscId})
	if err != nil {
		return nil, err
	}

	return &dto.EnterpriseEquityTransparency{
		UscId:      resp.UscId,
		Name:       info.EnterpriseTitle,
		Conclusion: resp.Conclusion,
		Data:       resp.Data,
	}, nil
}
