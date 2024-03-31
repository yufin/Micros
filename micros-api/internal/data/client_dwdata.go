package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	dwdataV2 "micros-api/api/dwdata/v2"
	dwdataV3 "micros-api/api/dwdata/v3"
	"micros-api/internal/biz"
	"micros-api/internal/biz/dto"
)

type ClientDwDataRepo struct {
	data *Data
	log  *log.Helper
}

func NewClientDwDataRepo(data *Data, logger log.Logger) biz.ClientDwDataRepo {
	return &ClientDwDataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *ClientDwDataRepo) GetClient(ctx context.Context) dwdataV2.DwdataServiceClient {
	return repo.data.DwDataClient
}

func (repo *ClientDwDataRepo) GetClientV3(ctx context.Context) dwdataV3.DwdataServiceClient {
	return repo.data.DwDataClientV3
}

func (repo *ClientDwDataRepo) GetEntIdent(ctx context.Context, name string) (string, error) {
	resp, err := repo.data.DwDataClient.GetEnterpriseIdent(context.TODO(), &dwdataV2.GetEntIdentReq{EnterpriseName: name})
	if err != nil {
		return "", err
	}
	return resp.Data.UscId, nil
}

func (repo *ClientDwDataRepo) GetEntInfo(ctx context.Context, uscId string) (*dto.EnterpriseInfo, error) {
	resp, err := repo.data.DwDataClient.GetEnterpriseInfo(context.TODO(), &dwdataV2.GetEntInfoReq{UscId: uscId})
	if err != nil {
		return nil, err
	}
	if resp.Success != true {
		return nil, nil
	}

	return &dto.EnterpriseInfo{
		UscId:                         resp.Data.UscId,
		EnterpriseTitle:               resp.Data.EnterpriseTitle,
		EnterpriseTitleEn:             resp.Data.EnterpriseTitleEn,
		BusinessRegistrationNumber:    resp.Data.BusinessRegistrationNumber,
		EstablishedDate:               resp.Data.EstablishDate,
		Region:                        resp.Data.Region,
		ApprovedDate:                  resp.Data.ApprovedDate,
		RegisteredAddress:             resp.Data.RegisteredAddress,
		RegisteredCapital:             resp.Data.RegisteredCapital,
		PaidInCapital:                 resp.Data.PaidInCapital,
		EnterpriseType:                resp.Data.EnterpriseType,
		StuffSize:                     resp.Data.StuffSize,
		StuffInsuredNumber:            int(resp.Data.StuffInsuredNumber),
		BusinessScope:                 resp.Data.BusinessScope,
		ImportExportQualificationCode: resp.Data.ImportExportQualificationCode,
		LegalRepresentative:           resp.Data.LegalRepresentative,
		RegistrationAuthority:         resp.Data.RegistrationAuthority,
		RegistrationStatus:            resp.Data.RegistrationStatus,
		TaxpayerQualification:         resp.Data.TaxpayerQualification,
		OrganizationCode:              resp.Data.OrganizationCode,
		UrlQcc:                        resp.Data.UrlQcc,
		UrlHomepage:                   resp.Data.UrlHomepage,
		BusinessTermStart:             resp.Data.BusinessTermStart,
		BusinessTermEnd:               resp.Data.BusinessTermEnd,
	}, nil
}

func (repo *ClientDwDataRepo) GetEntCredential(ctx context.Context, uscId string) (*[]dto.EnterpriseCertification, error) {
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

func (repo *ClientDwDataRepo) GetEntRankingList(ctx context.Context, uscId string) (*[]dto.EnterpriseRankingList, error) {
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

func (repo *ClientDwDataRepo) GetEntIndustry(ctx context.Context, uscId string) (*[]string, error) {
	resp, err := repo.data.DwDataClient.GetEnterpriseIndustry(context.TODO(), &dwdataV2.GetEntInfoReq{UscId: uscId})
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (repo *ClientDwDataRepo) GetEntProduct(ctx context.Context, uscId string) (*[]string, error) {
	resp, err := repo.data.DwDataClient.GetEnterpriseProduct(context.TODO(), &dwdataV2.GetEntInfoReq{UscId: uscId})
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (repo *ClientDwDataRepo) GetEquityTransparency(ctx context.Context, uscId string) (*dto.EnterpriseEquityTransparency, error) {
	resp, err := repo.data.DwDataClient.GetEntEquityTransparency(context.TODO(), &dwdataV2.GetEntInfoReq{UscId: uscId})
	if err != nil {
		return nil, err
	}
	if resp.Success != true {
		return nil, nil
	}

	info, err := repo.data.DwDataClient.GetEnterpriseInfo(context.TODO(), &dwdataV2.GetEntInfoReq{UscId: uscId})
	if err != nil {
		return nil, err
	}

	return &dto.EnterpriseEquityTransparency{
		UscId:      resp.UscId,
		Name:       info.Data.EnterpriseTitle,
		Conclusion: resp.Conclusion,
		Data:       resp.Data,
	}, nil
}

func (repo *ClientDwDataRepo) GetShareholders(ctx context.Context, uscId string) (*dwdataV2.GetShareholdersResp, error) {
	resp, err := repo.data.DwDataClient.GetEntShareholders(context.TODO(), &dwdataV2.GetEntInfoReq{
		UscId: uscId,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (repo *ClientDwDataRepo) GetInvestments(ctx context.Context, uscId string) (*dwdataV2.GetInvestmentResp, error) {
	resp, err := repo.data.DwDataClient.GetEntInvestment(context.TODO(), &dwdataV2.GetEntInfoReq{
		UscId: uscId,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (repo *ClientDwDataRepo) GetBranches(ctx context.Context, uscId string) (*dwdataV2.GetBranchesResp, error) {
	resp, err := repo.data.DwDataClient.GetEntBranches(context.TODO(), &dwdataV2.GetEntInfoReq{
		UscId: uscId,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
