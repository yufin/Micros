package service

import (
	"brillinkmicros/internal/biz"
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/structpb"

	pb "brillinkmicros/api/dw/v2"
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

func (s *DwServiceServicer) GetEnterpriseInfo(ctx context.Context, req *pb.GetEntInfoReq) (*pb.EntStructResp, error) {
	res, err := s.dwEnterprise.GetEntInfo(ctx, req.UscId)
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
	res, err := s.dwEnterprise.GetEntCredential(ctx, req.UscId)

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

func (s *DwServiceServicer) GetEnterpriseEquityTransparency(ctx context.Context, req *pb.GetEntInfoReq) (*pb.EquityTransparencyResp, error) {
	res, err := s.dwEnterprise.GetEquityTransparency(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &pb.EquityTransparencyResp{
			Success: false,
			Code:    1,
			Msg:     "data not found",
		}, nil
	}
	data := &pb.EquityTransparency{
		Conclusion: res.Conclusion,
		Detail:     res.Data,
		UscId:      req.UscId,
		Name:       res.Name,
	}
	return &pb.EquityTransparencyResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    data,
	}, nil
}
