package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/structpb"
	pb "micros-dw/api/dwdata/v2"
	"micros-dw/internal/biz"
	"net/http"
)

type DwdataServiceServicer struct {
	pb.UnimplementedDwdataServiceServer
	log          *log.Helper
	dwEnterprise *biz.DwEnterpriseUsecase
}

func NewDwdataServiceServicer(dwe *biz.DwEnterpriseUsecase, logger log.Logger) *DwdataServiceServicer {
	return &DwdataServiceServicer{
		dwEnterprise: dwe,
		log:          log.NewHelper(logger),
	}
}

func (s *DwdataServiceServicer) GetEnterpriseIdent(ctx context.Context, req *pb.GetEntIdentReq) (*pb.GetEntIdentResp, error) {
	res, err := s.dwEnterprise.GetEntIdent(ctx, req.EnterpriseName)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &pb.GetEntIdentResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data: &pb.EntIdent{
				Exists:  false,
				IsLegal: false,
				UscId:   "",
			},
		}, nil
	}
	if res.StatusCode == 9 {
		return &pb.GetEntIdentResp{
			Success: false,
			Code:    http.StatusNoContent,
			Data: &pb.EntIdent{
				Exists:  true,
				IsLegal: false,
				UscId:   "",
			},
		}, nil
	}

	return &pb.GetEntIdentResp{
		Success: true,
		Data: &pb.EntIdent{
			Exists:  true,
			IsLegal: true,
			UscId:   res.UscId,
		},
	}, nil
}

func (s *DwdataServiceServicer) GetEnterpriseInfo(ctx context.Context, req *pb.GetEntInfoReq) (*pb.GetEntInfoResp, error) {
	res, err := s.dwEnterprise.GetEntInfo(ctx, req.UscId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if res == nil {
		return &pb.GetEntInfoResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
		}, nil
	}
	data := &pb.EntInfo{
		UscId:                         res.UscId,
		EnterpriseTitle:               res.EnterpriseTitle,
		EnterpriseTitleEn:             res.EnterpriseTitleEn,
		BusinessRegistrationNumber:    res.BusinessRegistrationNumber,
		Region:                        res.Region,
		ApprovedDate:                  res.ApprovedDate.Format("2006-01-02"),
		RegisteredAddress:             res.RegisteredAddress,
		RegisteredCapital:             res.RegisteredCapital,
		PaidInCapital:                 res.PaidInCapital,
		EnterpriseType:                res.EnterpriseType,
		StuffSize:                     res.StuffSize,
		EstablishDate:                 res.EstablishedDate.Format("2006-01-02"),
		StuffInsuredNumber:            int32(res.StuffInsuredNumber),
		BusinessScope:                 res.BusinessScope,
		ImportExportQualificationCode: res.ImportExportQualificationCode,
		LegalRepresentative:           res.LegalRepresentative,
		RegistrationAuthority:         res.RegistrationAuthority,
		RegistrationStatus:            res.RegistrationStatus,
		TaxpayerQualification:         res.TaxpayerQualification,
		OrganizationCode:              res.OrganizationCode,
		UrlQcc:                        res.UrlQcc,
		UrlHomepage:                   res.UrlHomepage,
		BusinessTermStart:             res.BusinessTermStart.Format("2006-01-02"),
		BusinessTermEnd:               res.BusinessTermEnd.Format("2006-01-02"),
		Id:                            res.InfoId,
		CreatedAt:                     res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:                     res.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return &pb.GetEntInfoResp{
		Success: true,
		Code:    0,
		Msg:     "",
		Data:    data,
	}, nil
}

func (s *DwdataServiceServicer) GetEnterpriseCredential(ctx context.Context, req *pb.GetEntInfoReq) (*pb.GetEntCredentialResp, error) {
	res, err := s.dwEnterprise.GetEntCredential(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &pb.GetEntCredentialResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
		}, nil
	}
	data := make([]*pb.EntCredential, 0)
	for _, v := range *res {
		data = append(data, &pb.EntCredential{
			Id:                     v.CertId,
			UscId:                  v.UscId,
			CertificationTitle:     v.CertificationTitle,
			CertificationCode:      v.CertificationCode,
			CertificationLevel:     v.CertificationLevel,
			CertificationType:      v.CertificationType,
			CertificationSource:    v.CertificationSource,
			CertificationDate:      v.CertificationDate.Format("2006-01-02"),
			CertificationTermStart: v.CertificationTermStart.Format("2006-01-02"),
			CertificationTermEnd:   v.CertificationTermEnd.Format("2006-01-02"),
			CertificationAuthority: v.CertificationAuthority,
			CreatedAt:              v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:              v.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &pb.GetEntCredentialResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
	}, nil
}

func (s *DwdataServiceServicer) GetEnterpriseRankingList(ctx context.Context, req *pb.GetEntInfoReq) (*pb.GetEntRankingListResp, error) {
	res, err := s.dwEnterprise.GetEntRankingList(ctx, req.UscId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if res == nil {
		return &pb.GetEntRankingListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
		}, nil
	}
	stArray := make([]*pb.EnterpriseRankingList, 0)
	for _, v := range *res {
		stArray = append(
			stArray,
			&pb.EnterpriseRankingList{
				UscId:                 v.UscId,
				RankingPosition:       int32(v.RankingPosition),
				ListTitle:             v.ListTitle,
				ListType:              v.ListType,
				ListSource:            v.ListSource,
				ListParticipantsTotal: int32(v.ListParticipantsTotal),
				ListPublishedDate:     v.ListPublishedDate.Format("2006-01-02"),
				ListUrlQcc:            v.ListUrlQcc,
				ListUrlOrigin:         v.ListUrlOrigin,
			},
		)
	}
	return &pb.GetEntRankingListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    stArray,
	}, nil
}

func (s *DwdataServiceServicer) GetEnterpriseIndustry(ctx context.Context, req *pb.GetEntInfoReq) (*pb.GetEntStrArrayResp, error) {
	res, err := s.dwEnterprise.GetEntIndustry(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &pb.GetEntStrArrayResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
		}, nil
	}
	return &pb.GetEntStrArrayResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    *res,
	}, nil
}

func (s *DwdataServiceServicer) GetEnterpriseProduct(ctx context.Context, req *pb.GetEntInfoReq) (*pb.GetEntStrArrayResp, error) {
	res, err := s.dwEnterprise.GetEntProduct(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &pb.GetEntStrArrayResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}
	return &pb.GetEntStrArrayResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    *res,
	}, nil
}

func (s *DwdataServiceServicer) GetEntEquityTransparency(ctx context.Context, req *pb.GetEntInfoReq) (*pb.GetEquityTransparencyResp, error) {
	res, err := s.dwEnterprise.GetEntEquityTransparency(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &pb.GetEquityTransparencyResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
		}, nil
	}
	stArr := make([]*structpb.Struct, 0)
	for _, v := range res.Data {
		st, err := structpb.NewStruct(v)
		if err != nil {
			return nil, err
		}
		stArr = append(stArr, st)
	}

	return &pb.GetEquityTransparencyResp{
		Success:    true,
		Code:       http.StatusOK,
		Conclusion: res.Conclusion,
		Data:       stArr,
		UscId:      res.UscId,
	}, nil
}

func (s *DwdataServiceServicer) GetEntShareholders(ctx context.Context, req *pb.GetEntInfoReq) (*pb.GetShareholdersResp, error) {
	res, err := s.dwEnterprise.GetShareholders(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetShareholdersResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "no content",
			Data:    nil,
		}, nil
	}
	data := make([]*pb.Shareholders, 0)
	for _, v := range *res {
		data = append(data, &pb.Shareholders{
			ShareholderName: v.ShareholderName,
			ShareholderType: v.ShareholderType,
			CapitalAmount:   v.CapitalAmount,
			CapitalType:     v.CapitalType,
			Percent:         v.Percent,
		})
	}
	return &pb.GetShareholdersResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
	}, nil
}

func (s *DwdataServiceServicer) GetEntInvestment(ctx context.Context, req *pb.GetEntInfoReq) (*pb.GetInvestmentResp, error) {
	res, err := s.dwEnterprise.GetInvestments(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetInvestmentResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
		}, nil
	}
	data := make([]*pb.Investment, 0)
	for _, v := range *res {
		data = append(data, &pb.Investment{
			EnterpriseName:    v.EnterpriseName,
			Operator:          v.Operator,
			ShareholdingRatio: v.ShareholdingRatio,
			InvestedAmount:    v.InvestedAmount,
			StartData:         v.StartDate,
			Status:            v.Status,
		})
	}
	return &pb.GetInvestmentResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
	}, nil
}

func (s *DwdataServiceServicer) GetEntBranches(ctx context.Context, req *pb.GetEntInfoReq) (*pb.GetBranchesResp, error) {
	res, err := s.dwEnterprise.GetBranches(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetBranchesResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
		}, nil
	}
	data := make([]*pb.Branches, 0)
	for _, v := range *res {
		data = append(data, &pb.Branches{
			EnterpriseName: v.EnterpriseName,
			Operator:       v.Operator,
			Area:           v.Area,
			StartDate:      v.StartDate.Format("2006-01-02"),
			Status:         v.Status,
		})
	}
	return &pb.GetBranchesResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
	}, nil
}
