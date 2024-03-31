package v3

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/structpb"
	"math"
	pb "micros-dw/api/dwdata/v3"
	"micros-dw/internal/biz"
	"micros-dw/internal/biz/dto"
	"net/http"
)

type DwdataServiceServicer struct {
	pb.UnimplementedDwdataServiceServer
	log    *log.Helper
	dwData *biz.DwDataUsecase
}

func NewDwdataServiceServicer(dw *biz.DwDataUsecase, logger log.Logger) *DwdataServiceServicer {
	return &DwdataServiceServicer{
		dwData: dw,
		log:    log.NewHelper(logger),
	}
}

// GetUscIdByEnterpriseName 获取企业统一社会信用代码
func (s *DwdataServiceServicer) GetUscIdByEnterpriseName(ctx context.Context, req *pb.GetUscIdByEnterpriseNameReq) (*pb.GetUscIdByEnterpriseNameResp, error) {
	res, err := s.dwData.GetUscIdByEnterpriseName(ctx, req.EnterpriseName)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &pb.GetUscIdByEnterpriseNameResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data: &pb.GetUscIdByEnterpriseNameResp_EntIdent{
				Exists:  false,
				IsLegal: false,
				UscId:   "",
			},
		}, nil
	}
	if res.StatusCode == 9 {
		return &pb.GetUscIdByEnterpriseNameResp{
			Success: true,
			Code:    http.StatusNoContent,
			Data: &pb.GetUscIdByEnterpriseNameResp_EntIdent{
				Exists:  true,
				IsLegal: false,
				UscId:   "",
			},
		}, nil
	}

	return &pb.GetUscIdByEnterpriseNameResp{
		Success: true,
		Data: &pb.GetUscIdByEnterpriseNameResp_EntIdent{
			Exists:  true,
			IsLegal: true,
			UscId:   res.UscId,
		},
	}, nil
}

// GetEnterpriseInfo 工商信息 identical status query
func (s *DwdataServiceServicer) GetEnterpriseInfo(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataDictResp, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"usc_id", req.UscId},
				{"check_date", bson.D{{"$gt", req.TimePoint.AsTime()}}},
				{"create_date", bson.D{{"$lte", req.TimePoint.AsTime()}}},
			}}},
		bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"data", "$gsxx"},
			}},
		},
	}

	res, err := s.dwData.GetDocAggregated(ctx, "basic_status_enterprise_info", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataDictResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	r := (*res)[0]
	data, err := structpb.NewStruct(r["data"].(map[string]any))
	if err != nil {
		return nil, err
	}
	return &pb.GetDataDictResp{
		Success: true,
		Code:    http.StatusOK,
		Data:    data,
		Msg:     "",
	}, nil
}

// GetEntBranches 分支机构 identical status query
func (s *DwdataServiceServicer) GetEntBranches(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"check_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
			{{"usc_id", req.UscId}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "basic_status_branchlist", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
	limit := int64(math.Max(1, float64(req.PageSize)))

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"usc_id", req.UscId},
				{"check_date", bson.D{{"$gt", req.TimePoint.AsTime()}}},
				{"create_date", bson.D{{"$lte", req.TimePoint.AsTime()}}},
			}}},
		bson.D{{"$unwind", "$branch_list"}},
		//bson.D{{"$sort", bson.D{{"invest_info.StartDate", -1}}}},
		bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": "$branch_list"}}}},
		bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
	}

	res, err := s.dwData.GetDocAggregated(ctx, "basic_status_branchlist", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

// GetEntEquityTransparencyConclusion 股东穿透 identical status query
func (s *DwdataServiceServicer) GetEntEquityTransparencyConclusion(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataDictResp, error) {
	cond := bson.D{
		{"usc_id", req.UscId},
		{"check_date", bson.D{{"$gt", req.TimePoint.AsTime()}}},
		{"create_date", bson.D{{"$lte", req.TimePoint.AsTime()}}},
	}

	res, err := s.dwData.GetDoc(ctx, "basic_status_equity_penitration", &cond)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &pb.GetDataDictResp{
				Success: false,
				Code:    http.StatusNoContent,
				Msg:     "not content",
				Data:    nil,
			}, nil
		}
		return nil, err
	}
	m := make(map[string]any)
	var coclu string
	coclu, ok := (*res)["penitration_conclusion"].(string)
	if !ok {
		coclu = ""
	}
	m["conclusion"] = coclu
	st, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}
	return &pb.GetDataDictResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    st,
	}, nil
}

// GetEntEquityTransparency 股东穿透 identical status query
func (s *DwdataServiceServicer) GetEntEquityTransparency(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"check_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
			{{"usc_id", req.UscId}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "basic_status_equity_penitration", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
	limit := int64(math.Max(1, float64(req.PageSize)))

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"usc_id", req.UscId},
				{"check_date", bson.D{{"$gt", req.TimePoint.AsTime()}}},
				{"create_date", bson.D{{"$lte", req.TimePoint.AsTime()}}},
			}}},
		bson.D{{"$unwind", "$penitration_tree"}},
		//bson.D{{"$sort", bson.D{{"invest_info.StartDate", -1}}}},
		bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": "$penitration_tree"}}}},
		bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
	}
	res, err := s.dwData.GetDocAggregated(ctx, "basic_status_equity_penitration", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}
	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}
	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

// GetEntInvestment 对外投资 identical status query
func (s *DwdataServiceServicer) GetEntInvestment(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"check_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
			{{"usc_id", req.UscId}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "basic_status_investments", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
	limit := int64(math.Max(1, float64(req.PageSize)))

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"usc_id", req.UscId},
				{"check_date", bson.D{{"$gt", req.TimePoint.AsTime()}}},
				{"create_date", bson.D{{"$lte", req.TimePoint.AsTime()}}},
			}}},
		bson.D{{"$unwind", "$invest_info"}},
		bson.D{{"$sort", bson.D{{"invest_info.StartDate", -1}}}},

		bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": "$invest_info"}}}},
		bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
	}

	res, err := s.dwData.GetDocAggregated(ctx, "basic_status_investments", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

// GetEnterpriseManagerInfo 高管信息 identical status query
func (s *DwdataServiceServicer) GetEnterpriseManagerInfo(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"check_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
			{{"usc_id", req.UscId}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "basic_status_members", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
	limit := int64(math.Max(1, float64(req.PageSize)))
	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"usc_id", req.UscId},
				{"check_date", bson.D{{"$gt", req.TimePoint.AsTime()}}},
				{"create_date", bson.D{{"$lte", req.TimePoint.AsTime()}}},
			}}},
		bson.D{{"$unwind", "$members"}},
		bson.D{{"$sort", bson.D{{"members.No", 1}}}},
		bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": "$members"}}}},
		bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
	}
	res, err := s.dwData.GetDocAggregated(ctx, "basic_status_members", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}
	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}
	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

// GetEntFinancing 融资信息
func (s *DwdataServiceServicer) GetEntFinancing(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"usc_id", req.UscId}},
			{{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "develop_records_financing", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	pipeline := func(unwind string, sortBy string, req *pb.GetDataBeforeTimePointReq) mongo.Pipeline {
		skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
		limit := int64(math.Max(1, float64(req.PageSize)))
		return mongo.Pipeline{
			bson.D{
				{"$match", bson.D{
					{"usc_id", req.UscId},
					{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}},
					{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}},
				}}},
			bson.D{{"$unwind", fmt.Sprintf("$%s", unwind)}},
			bson.D{{"$match", bson.D{{sortBy, bson.D{{"$lte", req.TimePoint.AsTime().Format("2006-01-02")}}}}}},
			bson.D{{"$sort", bson.D{{sortBy, -1}}}},

			bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": fmt.Sprintf("$%s", unwind)}}}},
			bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
		}
	}("financing", "financing.Date", req)

	res, err := s.dwData.GetDocAggregated(ctx, "develop_records_financing", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

// GetCaseRegistrationInfo 立案信息 identical records query
func (s *DwdataServiceServicer) GetCaseRegistrationInfo(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"usc_id", req.UscId}},
			{{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "law_records_new_cases_", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	pipeline := func(unwind string, sortBy string, req *pb.GetDataBeforeTimePointReq) mongo.Pipeline {
		skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
		limit := int64(math.Max(1, float64(req.PageSize)))
		return mongo.Pipeline{
			bson.D{
				{"$match", bson.D{
					{"usc_id", req.UscId},
					{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}},
					{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}},
				}}},
			bson.D{{"$unwind", fmt.Sprintf("$%s", unwind)}},
			bson.D{{"$match", bson.D{{sortBy, bson.D{{"$lte", req.TimePoint.AsTime().Unix()}}}}}},
			bson.D{{"$sort", bson.D{{sortBy, -1}}}},

			bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": fmt.Sprintf("$%s", unwind)}}}},
			bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
		}
	}("data", "data.PunishDate", req)

	res, err := s.dwData.GetDocAggregated(ctx, "law_records_new_cases_", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

// GetEnterpriseRankingList  identical records query
func (s *DwdataServiceServicer) GetEnterpriseRankingList(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetEntRankingListResp, error) {
	//$new_case_info.cases, new_case_info.cases.CaseDate, law_records_new_cases
	filter := bson.D{
		{"$and", []bson.D{
			{{"usc_id", req.UscId}},
			{{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "develop_records_ranking_tags", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetEntRankingListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
	limit := int64(math.Max(1, float64(req.PageSize)))
	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"usc_id", req.UscId},
				{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}},
				{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}},
			}}},
		bson.D{{"$unwind", "$ranking_tags"}},
		bson.D{{"$match", bson.D{{"ranking_tags.PublishDate", bson.D{{"$lte", req.TimePoint.AsTime().Format("2006-01-02")}}}}}},
		bson.D{{"$sort", bson.D{{"ranking_tags.PublishDate", -1}}}},

		bson.D{{"$lookup", bson.D{
			{"from", "develop_records_bd_detail"},
			{"localField", "ranking_tags.Id"},
			{"foreignField", "bd_id"},
			{"as", "detail"},
		}}},

		bson.D{{"$project", bson.M{
			"list_detail":  bson.M{"$arrayElemAt": bson.A{"$detail.bd_detail.BdDetail", 0}},
			"ranking_tags": 1,
		}}},

		bson.D{{"$addFields", bson.M{
			"data": bson.M{
				"ranking_tags": "$ranking_tags",
				"list_detail":  "$list_detail",
			},
		}}},

		bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": "$data"}}}},
		bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
	}

	res, err := s.dwData.GetDocAggregated(ctx, "develop_records_ranking_tags", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetEntRankingListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*pb.GetEntRankingListResp_EnterpriseRankingList, 0),
		}, nil
	}

	resSt := dto.RankingTagsRes{}
	b, err := json.Marshal((*res)[0])
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &resSt); err != nil {
		return nil, err
	}

	data := make([]*pb.GetEntRankingListResp_EnterpriseRankingList, 0)
	for _, v := range resSt.Data {
		detail := pb.GetEntRankingListResp_EnterpriseRankingList{
			UscId:           req.UscId,
			RankingPosition: int32(v.RankingTags.Ranking),
			ListTitle:       v.RankingTags.Title,
			ListType:        v.ListDetail.BdType,
			ListSource:      v.ListDetail.InstitutionName,
			//ListParticipantsTotal: v.,
			ListPublishedDate: v.ListDetail.PublishDate.Date.Format("2006-01-02"),
			ListUrlQcc:        v.ListDetail.Url,
			//ListUrlOrigin:         "",
		}
		data = append(data, &detail)
	}

	return &pb.GetEntRankingListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(resSt.Total),
	}, nil
}

// GetEnterpriseCredential  identical records query
func (s *DwdataServiceServicer) GetEnterpriseCredential(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	//$new_case_info.cases, new_case_info.cases.CaseDate, law_records_new_cases
	filter := bson.D{
		{"$and", []bson.D{
			{{"usc_id", req.UscId}},
			{{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "develop_records_honour_tags", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	pipeline := func(unwind string, sortBy string, req *pb.GetDataBeforeTimePointReq) mongo.Pipeline {
		skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
		limit := int64(math.Max(1, float64(req.PageSize)))
		return mongo.Pipeline{
			bson.D{
				{"$match", bson.D{
					{"usc_id", req.UscId},
					{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}},
					{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}},
				}}},
			bson.D{{"$unwind", fmt.Sprintf("$%s", unwind)}},
			bson.D{{"$match", bson.D{{sortBy, bson.D{{"$lte", req.TimePoint.AsTime().Unix()}}}}}},
			bson.D{{"$sort", bson.D{{sortBy, -1}}}},

			bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": fmt.Sprintf("$%s", unwind)}}}},
			bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
		}
	}("honour_tags", "honour_tags.PublishDate", req)

	res, err := s.dwData.GetDocAggregated(ctx, "develop_records_honour_tags", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

// GetHistoricalExecutive 历史执法信息 identical records query
func (s *DwdataServiceServicer) GetHistoricalExecutive(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"usc_id", req.UscId}},
			{{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "law_records_his_execution", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	pipeline := func(unwind string, sortBy string, req *pb.GetDataBeforeTimePointReq) mongo.Pipeline {
		skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
		limit := int64(math.Max(1, float64(req.PageSize)))
		return mongo.Pipeline{
			bson.D{
				{"$match", bson.D{
					{"usc_id", req.UscId},
					{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}},
					{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}},
				}}},
			bson.D{{"$unwind", fmt.Sprintf("$%s", unwind)}},
			bson.D{{"$match", bson.D{{sortBy, bson.D{{"$lte", req.TimePoint.AsTime().Unix()}}}}}},
			bson.D{{"$sort", bson.D{{sortBy, -1}}}},

			bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": fmt.Sprintf("$%s", unwind)}}}},
			bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
		}
	}("zhixing_info.cases", "zhixing_info.cases.LiAnDate", req)

	res, err := s.dwData.GetDocAggregated(ctx, "law_records_his_execution", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

// GetExecutive 执法信息 identical records query
func (s *DwdataServiceServicer) GetExecutive(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"usc_id", req.UscId}},
			{{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "law_records_execution", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	pipeline := func(unwind string, sortBy string, req *pb.GetDataBeforeTimePointReq) mongo.Pipeline {
		skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
		limit := int64(math.Max(1, float64(req.PageSize)))
		return mongo.Pipeline{
			bson.D{
				{"$match", bson.D{
					{"usc_id", req.UscId},
					{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}},
					{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}},
				}}},
			bson.D{{"$unwind", fmt.Sprintf("$%s", unwind)}},
			bson.D{{"$match", bson.D{{sortBy, bson.D{{"$lte", req.TimePoint.AsTime().Unix()}}}}}},
			bson.D{{"$sort", bson.D{{sortBy, -1}}}},

			bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": fmt.Sprintf("$%s", unwind)}}}},
			bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
		}
	}("zhixing_info.cases", "zhixing_info.cases.LiAnDate", req)

	res, err := s.dwData.GetDocAggregated(ctx, "law_records_execution", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

// GetCourtAnnouncement 案件信息 identical records query
func (s *DwdataServiceServicer) GetCourtAnnouncement(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"usc_id", req.UscId}},
			{{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "law_records_cases", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	pipeline := func(unwind string, sortBy string, req *pb.GetDataBeforeTimePointReq) mongo.Pipeline {
		skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
		limit := int64(math.Max(1, float64(req.PageSize)))
		return mongo.Pipeline{
			bson.D{
				{"$match", bson.D{
					{"usc_id", req.UscId},
					{"update_date", bson.D{{"$gte", req.TimePoint.AsTime()}}},
					{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}},
				}}},
			bson.D{{"$unwind", fmt.Sprintf("$%s", unwind)}},
			bson.D{{"$match", bson.D{{sortBy, bson.D{{"$lte", req.TimePoint.AsTime().Unix()}}}}}},
			bson.D{{"$sort", bson.D{{sortBy, -1}}}},

			bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": fmt.Sprintf("$%s", unwind)}}}},
			bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
		}
	}("recent_cases", "recent_cases.SubmitDate", req)

	res, err := s.dwData.GetDocAggregated(ctx, "law_records_cases", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

// GetEntShareholders 股东信息 identical status query
func (s *DwdataServiceServicer) GetEntShareholders(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"check_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
			{{"usc_id", req.UscId}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "basic_status_shareholders", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
	limit := int64(math.Max(1, float64(req.PageSize)))

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"usc_id", req.UscId},
				{"check_date", bson.D{{"$gt", req.TimePoint.AsTime()}}},
				{"create_date", bson.D{{"$lte", req.TimePoint.AsTime()}}},
			}}},
		bson.D{{"$unwind", "$shareholders"}},
		bson.D{{"$sort", bson.D{{"shareholders.ShouldDate", -1}}}},

		bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": "$shareholders"}}}},
		bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
	}

	res, err := s.dwData.GetDocAggregated(ctx, "basic_status_shareholders", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

func (s *DwdataServiceServicer) GetForeclosureDisposition(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"check_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
			{{"usc_id", req.UscId}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "law_status_sales", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
	limit := int64(math.Max(1, float64(req.PageSize)))

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"usc_id", req.UscId},
				{"check_date", bson.D{{"$gt", req.TimePoint.AsTime()}}},
				{"create_date", bson.D{{"$lte", req.TimePoint.AsTime()}}},
			}}},
		bson.D{{"$unwind", "$sales_info.cases"}},
		//bson.D{{"$sort", bson.D{{"invest_info.StartDate", -1}}}},

		bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": "$sales_info.cases"}}}},
		bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
	}

	res, err := s.dwData.GetDocAggregated(ctx, "law_status_sales", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

func (s *DwdataServiceServicer) GetHighConsumptionRestriction(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"check_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
			{{"usc_id", req.UscId}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "law_status_consumption_restriction", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
	limit := int64(math.Max(1, float64(req.PageSize)))

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"usc_id", req.UscId},
				{"check_date", bson.D{{"$gt", req.TimePoint.AsTime()}}},
				{"create_date", bson.D{{"$lte", req.TimePoint.AsTime()}}},
			}}},
		bson.D{{"$unwind", "$restriction_info.cases"}},
		//bson.D{{"$sort", bson.D{{"invest_info.StartDate", -1}}}},

		bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": "$restriction_info.cases"}}}},
		bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
	}

	res, err := s.dwData.GetDocAggregated(ctx, "law_status_consumption_restriction", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

// GetEquityFrozen 获取股权冻结信息
func (s *DwdataServiceServicer) GetEquityFrozen(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"check_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
			{{"usc_id", req.UscId}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "basic_status_equity_freezes", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
	limit := int64(math.Max(1, float64(req.PageSize)))

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"usc_id", req.UscId},
				{"check_date", bson.D{{"$gt", req.TimePoint.AsTime()}}},
				{"create_date", bson.D{{"$lte", req.TimePoint.AsTime()}}},
			}}},
		bson.D{{"$unwind", "$equity_freezes.cases"}},
		//bson.D{{"$sort", bson.D{{"invest_info.StartDate", -1}}}},

		bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": "equity_freezes.cases"}}}},
		bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
	}

	res, err := s.dwData.GetDocAggregated(ctx, "basic_status_equity_freezes", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

// GetDiscreditedDebtor law_status_lose_credit
func (s *DwdataServiceServicer) GetDiscreditedDebtor(ctx context.Context, req *pb.GetDataBeforeTimePointReq) (*pb.GetDataListResp, error) {
	filter := bson.D{
		{"$and", []bson.D{
			{{"check_date", bson.D{{"$gte", req.TimePoint.AsTime()}}}},
			{{"create_date", bson.D{{"$lt", req.TimePoint.AsTime()}}}},
			{{"usc_id", req.UscId}},
		}},
	}
	c, err := s.dwData.CountDoc(ctx, filter, "law_status_lose_credit", "dw2")
	if err != nil {
		return nil, err
	}
	if c == 0 {
		return &pb.GetDataListResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "not content",
			Data:    nil,
		}, nil
	}

	skip := (int64(math.Max(1, float64(req.PageNum))) - 1) * req.PageSize
	limit := int64(math.Max(1, float64(req.PageSize)))
	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"usc_id", req.UscId},
				{"check_date", bson.D{{"$gt", req.TimePoint.AsTime()}}},
				{"create_date", bson.D{{"$lte", req.TimePoint.AsTime()}}},
			}}},
		bson.D{{"$unwind", "$shixin_info.cases"}},
		bson.D{{"$sort", bson.D{{"shixin_info.cases.LiAnDate", -1}}}},

		bson.D{{"$group", bson.M{"_id": nil, "total": bson.M{"$sum": 1}, "data": bson.M{"$push": "$shixin_info.cases"}}}},
		bson.D{{"$project", bson.M{"_id": 0, "total": 1, "data": bson.M{"$slice": bson.A{"$data", skip, limit}}}}},
	}

	res, err := s.dwData.GetDocAggregated(ctx, "law_status_lose_credit", &pipeline)
	if err != nil {
		return nil, err
	}
	if res == nil || len(*res) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
		}, nil
	}
	total := (*res)[0]["total"].(float64)
	arr, ok := (*res)[0]["data"].([]any)
	if !ok {
		return nil, errors.New(500, "data is not array", "")
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range arr {
		st, err := structpb.NewStruct(v.(map[string]any))
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}

	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   int64(total),
	}, nil
}

func (s *DwdataServiceServicer) GetMacroEconomyData(ctx context.Context, req *pb.GetMacroEconomyDataReq) (*pb.GetDataListResp, error) {
	resArr, count, err := s.dwData.GetDocsByPagination(ctx, "macro_data_spider", req.Item, req.PageSize, req.PageNum, req.SortBy)
	if err != nil {
		return nil, err
	}
	if resArr == nil || len(resArr) == 0 {
		return &pb.GetDataListResp{
			Success: true,
			Code:    http.StatusOK,
			Msg:     "",
			Data:    make([]*structpb.Struct, 0),
			Total:   count,
		}, nil
	}

	data := make([]*structpb.Struct, 0)
	for _, v := range resArr {
		st, err := structpb.NewStruct(v)
		if err != nil {
			return nil, err
		}
		data = append(data, st)
	}
	return &pb.GetDataListResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    data,
		Total:   count,
	}, nil

}
