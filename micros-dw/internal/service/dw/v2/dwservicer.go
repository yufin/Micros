package v2

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/structpb"
	pb "micros-dw/api/dw/v2"
	"micros-dw/internal/biz"
)

type DwServiceServicer struct {
	pb.UnimplementedDwServiceServer
	log          *log.Helper
	dwEnterprise *biz.DwEnterpriseUsecase
}

func NewDwServiceServicer(dwe *biz.DwEnterpriseUsecase, logger log.Logger) *DwServiceServicer {
	return &DwServiceServicer{
		dwEnterprise: dwe,
		log:          log.NewHelper(logger),
	}
}

func (s *DwServiceServicer) GetEnterpriseIdent(ctx context.Context, req *pb.GetEntIdentReq) (*pb.EntIdentResp, error) {
	res, err := s.dwEnterprise.GetEntIdent(ctx, req.EnterpriseName)
	if err != nil {
		return nil, err
	}
	if res == "" {
		return &pb.EntIdentResp{
			Success: false,
			Code:    0,
			Msg:     "enterpriseName not found",
			UscId:   "",
		}, nil
	}
	return &pb.EntIdentResp{
		Success: true,
		Code:    200,
		Msg:     "",
		UscId:   res,
	}, nil
}

func (s *DwServiceServicer) GetEnterpriseInfo(ctx context.Context, req *pb.GetEntInfoReq) (*pb.EntInfoResp, error) {
	res, err := s.dwEnterprise.GetEntInfo(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.WithStack(errors.New("record not found"))
	}
	return &pb.EntInfoResp{
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
	}, nil
}
func (s *DwServiceServicer) GetEnterpriseCredential(ctx context.Context, req *pb.GetEntInfoReq) (*pb.EntCredentialResp, error) {
	res, err := s.dwEnterprise.GetEntCredential(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.WithStack(errors.New("record not found"))
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
	return &pb.EntCredentialResp{Data: data}, nil
}

func (s *DwServiceServicer) GetEnterpriseRankingList(ctx context.Context, req *pb.GetEntInfoReq) (*pb.EntArrayResp, error) {
	res, err := s.dwEnterprise.GetEntRankingList(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &pb.EntArrayResp{
			Success: false,
			Code:    0,
			Msg:     "uscId not found",
			Data:    nil,
		}, nil
	}

	stArray := make([]*structpb.Struct, 0)
	for _, v := range *res {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		m := make(map[string]interface{})
		err = json.Unmarshal(b, &m)
		if err != nil {
			return nil, err
		}
		st, err := structpb.NewStruct(m)
		if err != nil {
			return nil, err
		}
		stArray = append(stArray, st)
	}
	return &pb.EntArrayResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    stArray,
	}, nil
}
func (s *DwServiceServicer) GetEnterpriseIndustry(ctx context.Context, req *pb.GetEntInfoReq) (*pb.EntStrArrayResp, error) {
	res, err := s.dwEnterprise.GetEntIndustry(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &pb.EntStrArrayResp{
			Success: false,
			Code:    0,
			Msg:     "uscId not found",
			Data:    nil,
		}, nil
	}
	return &pb.EntStrArrayResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    *res,
	}, nil
}
func (s *DwServiceServicer) GetEnterpriseProduct(ctx context.Context, req *pb.GetEntInfoReq) (*pb.EntStrArrayResp, error) {
	res, err := s.dwEnterprise.GetEntProduct(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &pb.EntStrArrayResp{
			Success: false,
			Code:    0,
			Msg:     "uscId not found",
			Data:    nil,
		}, nil
	}
	return &pb.EntStrArrayResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    *res,
	}, nil
}
