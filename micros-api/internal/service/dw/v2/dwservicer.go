package service

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/structpb"
	dwdataV2 "micros-api/api/dwdata/v2"
	"micros-api/internal/biz"
	"micros-api/internal/data"
	"net/http"

	pb "micros-api/api/dw/v2"
)

type DwServiceServicer struct {
	pb.UnimplementedDwServiceServer
	log          *log.Helper
	clientDwData *biz.ClientDwDataUsecase
	data         *data.Data
}

func NewDwServiceServicer(dwe *biz.ClientDwDataUsecase, logger log.Logger) *DwServiceServicer {
	return &DwServiceServicer{
		clientDwData: dwe,
		log:          log.NewHelper(logger),
	}
}

func (s *DwServiceServicer) GetEntRelations(ctx context.Context, req *pb.GetEntInfoReq) (*pb.EnterpriseRelations, error) {
	info, err := s.clientDwData.GetEntInfo(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return &pb.EnterpriseRelations{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "object enterprise not found",
			Data:    nil,
		}, nil
	}
	var relData pb.EnterpriseRelations_RelationsData
	relData.EnterpriseName = info.EnterpriseTitle

	branch, err := s.clientDwData.GetBranches(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	branchData := make([]*pb.Branches, 0)
	if branch.Success {
		for _, v := range branch.Data {
			branchData = append(branchData, &pb.Branches{
				EnterpriseName: v.EnterpriseName,
				Operator:       v.Operator,
				Area:           v.Area,
				Status:         v.Status,
				StartDate:      v.StartDate,
			})
		}
	}
	relData.Branch = branchData

	investment, err := s.clientDwData.GetInvestments(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	investmentData := make([]*pb.Investment, 0)
	if investment.Success {
		for _, v := range investment.Data {
			investmentData = append(investmentData, &pb.Investment{
				EnterpriseName:    v.EnterpriseName,
				Operator:          v.Operator,
				ShareholdingRatio: v.ShareholdingRatio,
				InvestedAmount:    v.InvestedAmount,
				Status:            v.Status,
				StartDate:         v.StartData,
			})
		}
	}
	relData.Investment = investmentData

	shareholder, err := s.clientDwData.GetShareholders(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	shareholderData := make([]*pb.Shareholders, 0)
	if shareholder.Success {
		for _, v := range shareholder.Data {
			shareholderData = append(shareholderData, &pb.Shareholders{
				ShareholderName: v.ShareholderName,
				ShareholderType: v.ShareholderType,
				CapitalType:     v.CapitalType,
				RealAmount:      v.RealAmount,
				CapitalAmount:   v.CapitalAmount,
				Percent:         v.Percent,
			})
		}
	}
	relData.Shareholder = shareholderData

	return &pb.EnterpriseRelations{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    &relData,
	}, nil
}

func (s *DwServiceServicer) GetEnterpriseIdent(ctx context.Context, req *pb.GetEntIdentReq) (*pb.EntIdentResp, error) {
	res, err := s.clientDwData.GetEntIdent(ctx, req.EnterpriseName)
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

func (s *DwServiceServicer) GetEnterpriseInfo(ctx context.Context, req *pb.GetEntInfoReq) (*pb.EntStructResp, error) {
	res, err := s.clientDwData.GetEntInfo(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &pb.EntStructResp{
			Success: false,
			Code:    0,
			Msg:     "uscId not found",
			Data:    nil,
		}, nil
	}
	b, err := json.Marshal(res)
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
	return &pb.EntStructResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    st,
	}, nil
}

func (s *DwServiceServicer) GetEnterpriseCredential(ctx context.Context, req *pb.GetEntInfoReq) (*pb.EntArrayResp, error) {
	res, err := s.clientDwData.GetEntCredential(ctx, req.UscId)

	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
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

func (s *DwServiceServicer) GetEnterpriseRankingList(ctx context.Context, req *pb.GetEntInfoReq) (*pb.EntArrayResp, error) {
	res, err := s.clientDwData.GetEntRankingList(ctx, req.UscId)
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
	res, err := s.clientDwData.GetEntIndustry(ctx, req.UscId)
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
	res, err := s.clientDwData.GetEntProduct(ctx, req.UscId)
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

func (s *DwServiceServicer) GetEnterpriseEquityTransparency(ctx context.Context, req *pb.GetEntInfoReq) (*pb.EquityTransparencyResp, error) {
	res, err := s.clientDwData.GetEquityTransparency(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &pb.EquityTransparencyResp{
			Success: false,
			Code:    0,
			Msg:     "data not found",
		}, nil
	}
	data := &pb.EquityTransparency{
		Conclusion:  res.Conclusion,
		Shareholder: res.Data,
		KeyNo:       req.UscId,
		Name:        res.Name,
	}
	return &pb.EquityTransparencyResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    data,
	}, nil
}

func (s *DwServiceServicer) GetEntBranches(ctx context.Context, req *pb.GetEntInfoReq) (*pb.BranchesResp, error) {
	res, err := s.clientDwData.GetBranches(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res.Success == false {
		return &pb.BranchesResp{
			Success: true,
			Found:   false,
			Code:    0,
			Msg:     "data not found",
		}, nil
	}
	d := make([]*pb.Branches, 0)
	for _, v := range res.Data {
		d = append(d, &pb.Branches{
			EnterpriseName: v.EnterpriseName,
			Operator:       v.Operator,
			Area:           v.Area,
			Status:         v.Status,
			StartDate:      v.StartDate,
		})
	}

	return &pb.BranchesResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    d,
	}, nil
}

func (s *DwServiceServicer) GetEntInvestment(ctx context.Context, req *pb.GetEntInfoReq) (*pb.InvestmentResp, error) {
	res, err := s.clientDwData.GetInvestments(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.InvestmentResp{
			Success: true,
			Found:   false,
			Code:    0,
			Msg:     "data not found",
		}, nil
	}

	d := make([]*pb.Investment, 0)
	for _, v := range res.Data {
		d = append(d, &pb.Investment{
			EnterpriseName:    v.EnterpriseName,
			Operator:          v.Operator,
			ShareholdingRatio: v.ShareholdingRatio,
			InvestedAmount:    v.InvestedAmount,
			Status:            v.Status,
			StartDate:         v.StartData,
		})
	}
	return &pb.InvestmentResp{
		Success: true,
		Code:    0,
		Msg:     "",
		Found:   true,
		Data:    d,
	}, nil
}

func (s *DwServiceServicer) GetEntShareholders(ctx context.Context, req *pb.GetEntInfoReq) (*pb.ShareholdersResp, error) {
	res, err := s.clientDwData.GetShareholders(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.ShareholdersResp{
			Success: true,
			Found:   false,
			Code:    0,
			Msg:     "data not found",
		}, nil
	}

	d := make([]*pb.Shareholders, 0)
	for _, v := range res.Data {
		d = append(d, &pb.Shareholders{
			ShareholderName: v.ShareholderName,
			ShareholderType: v.ShareholderType,
			CapitalType:     v.CapitalType,
			RealAmount:      v.RealAmount,
			CapitalAmount:   v.CapitalAmount,
			Percent:         v.Percent,
		})
	}
	return &pb.ShareholdersResp{
		Success: true,
		Code:    0,
		Msg:     "",
		Found:   true,
		Data:    d,
	}, nil
}

func (s *DwServiceServicer) GetForeclosureDisposition(ctx context.Context, req *pb.GetEntInfoWithTime) (*pb.EntStructResp, error) {
	cli := s.clientDwData.GetClient(context.TODO())
	req.TimePoint.AsTime()
	resp, err := cli.GetForeclosureDisposition(context.TODO(), &dwdataV2.GetEntInfoWithDurationReq{
		UscId:              req.UscId,
		TimePoint:          req.TimePoint,
		ValidateExtendDate: 60,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &pb.EntStructResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "content not found",
			Data:    nil,
		}, nil
	}
	if !resp.Success {
		return &pb.EntStructResp{
			Success: false,
			Code:    resp.Code,
			Msg:     resp.Msg,
			Data:    nil,
		}, nil
	}
	m := resp.Data.AsMap()
	st, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}
	return &pb.EntStructResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    st,
	}, nil
}

func (s *DwServiceServicer) GetExecutive(ctx context.Context, req *pb.GetEntInfoWithTime) (*pb.EntStructResp, error) {
	cli := s.clientDwData.GetClient(context.TODO())
	req.TimePoint.AsTime()
	resp, err := cli.GetExecutive(context.TODO(), &dwdataV2.GetEntInfoWithDurationReq{
		UscId:              req.UscId,
		TimePoint:          req.TimePoint,
		ValidateExtendDate: 60,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &pb.EntStructResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "content not found",
			Data:    nil,
		}, nil
	}
	if !resp.Success {
		return &pb.EntStructResp{
			Success: false,
			Code:    resp.Code,
			Msg:     resp.Msg,
			Data:    nil,
		}, nil
	}
	m := resp.Data.AsMap()
	st, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}
	return &pb.EntStructResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    st,
	}, nil
}

func (s *DwServiceServicer) GetEquityFrozen(ctx context.Context, req *pb.GetEntInfoWithTime) (*pb.EntStructResp, error) {
	cli := s.clientDwData.GetClient(context.TODO())
	req.TimePoint.AsTime()
	resp, err := cli.GetEquityFrozen(context.TODO(), &dwdataV2.GetEntInfoWithDurationReq{
		UscId:              req.UscId,
		TimePoint:          req.TimePoint,
		ValidateExtendDate: 60,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &pb.EntStructResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "content not found",
			Data:    nil,
		}, nil
	}
	if !resp.Success {
		return &pb.EntStructResp{
			Success: false,
			Code:    resp.Code,
			Msg:     resp.Msg,
			Data:    nil,
		}, nil
	}
	m := resp.Data.AsMap()
	st, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}
	return &pb.EntStructResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    st,
	}, nil
}

func (s *DwServiceServicer) GetHighConsumptionRestriction(ctx context.Context, req *pb.GetEntInfoWithTime) (*pb.EntStructResp, error) {
	cli := s.clientDwData.GetClient(context.TODO())
	req.TimePoint.AsTime()
	resp, err := cli.GetHighConsumptionRestriction(context.TODO(), &dwdataV2.GetEntInfoWithDurationReq{
		UscId:              req.UscId,
		TimePoint:          req.TimePoint,
		ValidateExtendDate: 60,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &pb.EntStructResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "content not found",
			Data:    nil,
		}, nil
	}
	if !resp.Success {
		return &pb.EntStructResp{
			Success: false,
			Code:    resp.Code,
			Msg:     resp.Msg,
			Data:    nil,
		}, nil
	}
	m := resp.Data.AsMap()
	st, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}
	return &pb.EntStructResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    st,
	}, nil
}

func (s *DwServiceServicer) GetJudicialStatics(ctx context.Context, req *pb.GetEntInfoWithTime) (*pb.EntStructResp, error) {
	cli := s.clientDwData.GetClient(context.TODO())
	req.TimePoint.AsTime()
	resp, err := cli.GetJudicialStatics(context.TODO(), &dwdataV2.GetEntInfoWithDurationReq{
		UscId:              req.UscId,
		TimePoint:          req.TimePoint,
		ValidateExtendDate: 730,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &pb.EntStructResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "content not found",
			Data:    nil,
		}, nil
	}
	if !resp.Success {
		return &pb.EntStructResp{
			Success: false,
			Code:    resp.Code,
			Msg:     resp.Msg,
			Data:    nil,
		}, nil
	}
	m := resp.Data.AsMap()
	st, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}
	return &pb.EntStructResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    st,
	}, nil
}

func (s *DwServiceServicer) GetCourtAnnouncement(ctx context.Context, req *pb.GetEntInfoWithTime) (*pb.EntStructResp, error) {
	cli := s.clientDwData.GetClient(context.TODO())
	req.TimePoint.AsTime()
	resp, err := cli.GetCourtAnnouncement(context.TODO(), &dwdataV2.GetEntInfoWithDurationReq{
		UscId:              req.UscId,
		TimePoint:          req.TimePoint,
		ValidateExtendDate: 730,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &pb.EntStructResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "content not found",
			Data:    nil,
		}, nil
	}
	if !resp.Success {
		return &pb.EntStructResp{
			Success: false,
			Code:    resp.Code,
			Msg:     resp.Msg,
			Data:    nil,
		}, nil
	}
	m := resp.Data.AsMap()
	st, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}
	return &pb.EntStructResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    st,
	}, nil
}

func (s *DwServiceServicer) GetDiscreditedDebtor(ctx context.Context, req *pb.GetEntInfoWithTime) (*pb.EntStructResp, error) {
	cli := s.clientDwData.GetClient(context.TODO())
	req.TimePoint.AsTime()
	resp, err := cli.GetDiscreditedDebtor(context.TODO(), &dwdataV2.GetEntInfoWithDurationReq{
		UscId:              req.UscId,
		TimePoint:          req.TimePoint,
		ValidateExtendDate: 60,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &pb.EntStructResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "content not found",
			Data:    nil,
		}, nil
	}
	if !resp.Success {
		return &pb.EntStructResp{
			Success: false,
			Code:    resp.Code,
			Msg:     resp.Msg,
			Data:    nil,
		}, nil
	}
	m := resp.Data.AsMap()
	st, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}
	return &pb.EntStructResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    st,
	}, nil
}
