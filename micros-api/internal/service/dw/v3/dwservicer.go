package service

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	dwdataV3 "micros-api/api/dwdata/v3"
	pipelineV1 "micros-api/api/pipeline/v1"
	"micros-api/internal/biz"
	"micros-api/internal/biz/dto"
	"micros-api/internal/data"
	"micros-api/pkg"
	"net"
	"net/http"
	"time"

	pb "micros-api/api/dw/v3"
)

type DwServiceServicer struct {
	pb.UnimplementedDwServiceServer
	log            *log.Helper
	clientDwData   *biz.ClientDwDataUsecase
	mgoRc          *biz.MgoRcUsecase
	data           *data.Data
	pipelineClient *biz.ClientPipelineUsecase
	artifacts      *biz.ArtifactDataUsecase
}

func NewDwServiceServicer(dwe *biz.ClientDwDataUsecase, mgo *biz.MgoRcUsecase, dt *data.Data, plc *biz.ClientPipelineUsecase, art *biz.ArtifactDataUsecase, logger log.Logger) *DwServiceServicer {
	return &DwServiceServicer{
		clientDwData:   dwe,
		mgoRc:          mgo,
		data:           dt,
		pipelineClient: plc,
		artifacts:      art,
		log:            log.NewHelper(logger),
	}
}

func (s *DwServiceServicer) GetEntRelations(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.EnterpriseRelationsResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req0 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}

	info, err := cli.GetEnterpriseInfo(context.TODO(), req0)
	if err != nil {
		return nil, err
	}
	var relData pb.EnterpriseRelationsResp_RelationsData

	if info != nil && info.Success {
		infoM := info.Data.AsMap()
		entName, ok := infoM["enterprise_title"].(string)
		if ok {
			relData.EnterpriseName = entName
		}
	}
	branch, err := cli.GetEntBranches(ctx, req0)
	if err != nil {
		return nil, err
	}
	if branch.Success {
		relData.Branch = branch.Data
		relData.TotalBranch = int32(branch.Total)
	}

	investment, err := cli.GetEntInvestment(ctx, req0)
	if err != nil {
		return nil, err
	}
	if investment.Success {
		relData.Investment = investment.Data
		relData.TotalInvestment = int32(investment.Total)
	}

	shareholder, err := cli.GetEntShareholders(ctx, req0)
	if err != nil {
		return nil, err
	}
	if shareholder.Success {
		relData.Shareholder = shareholder.Data
		relData.TotalShareholder = int32(shareholder.Total)
	}

	httpReq, ok := kratosHttp.RequestFromServerContext(ctx)
	if !ok {
		return nil, errors.New(400, "invalid request context", "")
	}
	lang := httpReq.Header.Get("Accept-Language")
	if lang == "en-US" {
		b, err := protojson.Marshal(&relData)
		if err != nil {
			return nil, err
		}
		var m map[string]interface{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return nil, err
		}
		hash, err := pkg.GenMapCacheKey(m)
		if err != nil {
			return nil, err
		}
		k := "entRelationEn:" + hash
		bEn, err := s.data.Rdb.Client.Get(context.TODO(), k).Bytes()
		if err != nil {
			var netErr net.Error
			if errors.Is(err, redis.Nil) || (errors.As(err, &netErr) && netErr.Timeout()) {
				s.log.Infof("refresh tradeDetailEn key: %v", k)
				transReqSt, err := structpb.NewStruct(m)
				if err != nil {
					return nil, err
				}
				cli := s.pipelineClient.GetClient(context.TODO())
				transResp, err := cli.GetJsonTranslate(context.TODO(), &pipelineV1.GetJsonTranslateReq{Data: transReqSt})
				if err != nil {
					return nil, err
				}
				if !transResp.Success {
					return nil, errors.New(400, "translate failed", "")
				}
				bEn, err = json.Marshal(transResp.Data.AsMap())
				if err != nil {
					return nil, err
				}
				err = s.data.Rdb.Client.Set(context.TODO(), k, bEn, time.Hour*24*3).Err()
				if err != nil {
					return nil, err
				}

			} else {
				return nil, err
			}
		}
		var relDataEn pb.EnterpriseRelationsResp_RelationsData
		err = protojson.Unmarshal(bEn, &relDataEn)
		if err != nil {
			return nil, err
		}
		return &pb.EnterpriseRelationsResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    &relDataEn,
		}, nil

	}

	return &pb.EnterpriseRelationsResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    &relData,
	}, nil
}

func (s *DwServiceServicer) GetUscIdByEnterpriseName(ctx context.Context, req *pb.GetUscIdByEnterpriseNameReq) (*pb.GetUscIdByEnterpriseNameResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetUscIdByEnterpriseNameReq{
		EnterpriseName: req.EnterpriseName,
	}
	res, err := cli.GetUscIdByEnterpriseName(ctx, req2)
	if err != nil {
		return nil, err
	}
	return &pb.GetUscIdByEnterpriseNameResp{
		Success: res.Success,
		Code:    res.Code,
		Msg:     res.Msg,
		Data: &pb.GetUscIdByEnterpriseNameResp_EntIdent{
			Exists:  res.Data.Exists,
			IsLegal: res.Data.IsLegal,
			UscId:   res.Data.UscId,
		},
	}, nil

}

func (s *DwServiceServicer) GetEnterpriseInfo(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataDictResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}

	res, err := cli.GetEnterpriseInfo(ctx, req2)
	if err != nil {
		return nil, err
	}
	return &pb.GetDataDictResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
	}, nil
}

func (s *DwServiceServicer) GetEnterpriseEquityTransparency(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetEntEquityTransparency(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetEnterpriseEquityTransparencyConclusion(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataDictResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetEntEquityTransparencyConclusion(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataDictResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
		}, nil
	}
	return &pb.GetDataDictResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
	}, nil
}

func (s *DwServiceServicer) GetEntShareholders(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetEntShareholders(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetEntInvestment(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetEntInvestment(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetEntBranches(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetEntBranches(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetCaseRegistrationInfo(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetCaseRegistrationInfo(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
			Total:   int32(res.Total),
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetForeclosureDisposition(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetForeclosureDisposition(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
			Total:   int32(res.Total),
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetExecutive(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetExecutive(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
			Total:   int32(res.Total),
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetEquityFrozen(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetEquityFrozen(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetHighConsumptionRestriction(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetHighConsumptionRestriction(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetCourtAnnouncement(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetCourtAnnouncement(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetEnterpriseManagerInfo(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetEnterpriseManagerInfo(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetDiscreditedDebtor(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetDiscreditedDebtor(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetEnterpriseCredential(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetEnterpriseCredential(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetCollStat(ctx context.Context, req *pb.GetCollStatReq) (*pb.GetCollStatResp, error) {
	totalDist, err := s.mgoRc.GetCountOnDistinctUscIdForColl(ctx, req.CollName.String())
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	m["totalDist"] = totalDist
	st, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}
	return &pb.GetCollStatResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    st,
	}, nil
}

func (s *DwServiceServicer) GetEnterpriseRankingList(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetEntRankingListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetDataBeforeTimePointReq{
		UscId:     req.UscId,
		TimePoint: req.TimePoint,
		PageSize:  req.PageSize,
		PageNum:   req.PageNum,
	}
	res, err := cli.GetEnterpriseRankingList(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetEntRankingListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
		}, nil
	}

	rankings := make([]*pb.GetEntRankingListResp_EnterpriseRankingList, 0)
	for _, v := range res.Data {
		rankings = append(rankings, &pb.GetEntRankingListResp_EnterpriseRankingList{
			UscId:                 v.UscId,
			RankingPosition:       v.RankingPosition,
			ListTitle:             v.ListTitle,
			ListType:              v.ListType,
			ListSource:            v.ListSource,
			ListParticipantsTotal: v.ListParticipantsTotal,
			ListPublishedDate:     v.ListPublishedDate,
			ListUrlQcc:            v.ListUrlQcc,
			ListUrlOrigin:         v.ListUrlOrigin,
		})

	}

	return &pb.GetEntRankingListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    rankings,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) GetMacroEconomyData(ctx context.Context, req *pb.GetMacroEconomyDataReq) (*pb.GetDataListResp, error) {
	cli := s.clientDwData.GetClientV3(ctx)
	req2 := &dwdataV3.GetMacroEconomyDataReq{
		Item:     req.Item,
		SortBy:   req.SortBy,
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
	}
	res, err := cli.GetMacroEconomyData(ctx, req2)

	if err != nil {
		return nil, err
	}
	if !res.Success {
		return &pb.GetDataListResp{
			Success: res.Success,
			Code:    res.Code,
			Msg:     res.Msg,
			Data:    res.Data,
		}, nil
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    res.Data,
		Total:   int32(res.Total),
	}, nil
}

func (s *DwServiceServicer) InsertArtifactData(ctx context.Context, req *pb.InsertArtifactDataReq) (*pb.InsertArtifactDataResp, error) {
	switch req.Item {
	case pb.ArtifactDataItem_ENTERPRISE_COMMENT:
		insertReq := make([]dto.EnterpriseCommentInsertReq, 0)
		for _, v := range req.Data {
			m := v.AsMap()
			b, err := json.Marshal(m)
			if err != nil {
				return nil, err
			}
			var ec dto.EnterpriseCommentInsertReq
			if err = json.Unmarshal(b, &ec); err != nil {
				return nil, err
			}
			insertReq = append(insertReq, ec)
		}
		if err := s.artifacts.InsertEnterpriseComment(ctx, insertReq); err != nil {
			return nil, err
		}
	case pb.ArtifactDataItem_PRODUCT_EVAL_RULE:
		insertReq := make([]dto.ProductEvalRuleInsertResp, 0)
		for _, v := range req.Data {
			m := v.AsMap()
			b, err := json.Marshal(m)
			if err != nil {
				return nil, err
			}
			var per dto.ProductEvalRuleInsertResp
			if err = json.Unmarshal(b, &per); err != nil {
				return nil, err
			}
			insertReq = append(insertReq, per)
		}
		if err := s.artifacts.InsertProductEvalRule(ctx, insertReq); err != nil {
			return nil, err
		}
	default:
		return &pb.InsertArtifactDataResp{
			Success: false,
			Code:    http.StatusConflict,
			Msg:     "Item not exist",
		}, nil
	}

	return &pb.InsertArtifactDataResp{
		Success: true,
		Code:    http.StatusCreated,
		Msg:     "",
	}, nil
}

func (s *DwServiceServicer) GetArtifactData(ctx context.Context, req *pb.GetArtifactDataReq) (*pb.GetDataListResp, error) {

	stArr := make([]*structpb.Struct, 0)
	f := bson.D{}
	if req.MatchTarget != "" {
		f = bson.D{{req.FieldMatch, req.MatchTarget}}
	}
	var total int64

	switch req.Item {
	case pb.ArtifactDataItem_ENTERPRISE_COMMENT:
		res, t, err := s.artifacts.GetEnterpriseComment(ctx, req.PageSize, req.PageNum, f)
		if err != nil {
			return nil, err
		}
		for _, v := range res {
			st, err := structpb.NewStruct(v)
			if err != nil {
				return nil, err
			}
			stArr = append(stArr, st)
		}
		total = t
	case pb.ArtifactDataItem_PRODUCT_EVAL_RULE:
		res, t, err := s.artifacts.GetProductEvalRule(ctx, req.PageSize, req.PageNum, f)
		if err != nil {
			return nil, err
		}
		for _, v := range res {
			st, err := structpb.NewStruct(v)
			if err != nil {
				return nil, err
			}
			stArr = append(stArr, st)
		}
		total = t
	default:
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusConflict,
			Msg:     "Item not exist",
		}, nil
	}
	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusCreated,
		Msg:     "",
		Data:    stArr,
		Total:   int32(total),
	}, nil
}

func (s *DwServiceServicer) DeleteArtifactData(ctx context.Context, req *pb.DeleteArtifactDataReq) (*pb.DeleteArtifactDataResp, error) {
	switch req.Item {
	case pb.ArtifactDataItem_ENTERPRISE_COMMENT:
		err := s.artifacts.DeleteEnterpriseComment(context.TODO(), req.CollId)
		if err != nil {
			return nil, err
		}
	case pb.ArtifactDataItem_PRODUCT_EVAL_RULE:
		err := s.artifacts.DeleteProductEvalRuleById(context.TODO(), req.CollId)
		if err != nil {
			return nil, err
		}
	default:
		return &pb.DeleteArtifactDataResp{
			Success: false,
			Code:    http.StatusConflict,
			Msg:     "Item not exist",
		}, nil
	}
	return &pb.DeleteArtifactDataResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
	}, nil
}
