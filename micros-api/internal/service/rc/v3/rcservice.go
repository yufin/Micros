package v3

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/types/known/structpb"
	"micros-api/internal/biz"
	"micros-api/internal/biz/dto"
	"net/http"
	"sync"

	//"github.com/gogo/protobuf/proto/protojson"
	//"github.com/gogo/protobuf/proto/structpb"
	pb "micros-api/api/rc/v3"
)

type RcServiceServicer struct {
	pb.UnimplementedRcServiceServer
	log                *log.Helper
	rcProcessedContent *biz.RcProcessedContentUsecase
	rcOriginContent    *biz.RcOriginContentUsecase
	rcDependencyData   *biz.RcDependencyDataUsecase
	rcReportOss        *biz.RcReportOssUsecase
	rcDecisionFactor   *biz.RcDecisionFactorUsecase
	ossMetadata        *biz.OssMetadataUsecase
	mgoRc              *biz.MgoRcUsecase
}

func NewRcServiceServicer(
	rpc *biz.RcProcessedContentUsecase,
	roc *biz.RcOriginContentUsecase,
	rdd *biz.RcDependencyDataUsecase,
	omd *biz.OssMetadataUsecase,
	rro *biz.RcReportOssUsecase,
	rdf *biz.RcDecisionFactorUsecase,
	mgo *biz.MgoRcUsecase,
	logger log.Logger) *RcServiceServicer {
	return &RcServiceServicer{
		rcOriginContent:    roc,
		rcProcessedContent: rpc,
		rcDependencyData:   rdd,
		rcReportOss:        rro,
		ossMetadata:        omd,
		rcDecisionFactor:   rdf,
		mgoRc:              mgo,
		log:                log.NewHelper(logger),
	}
}

// GetReportDecisionFactor 查询企业风控参数
func (s *RcServiceServicer) GetReportDecisionFactor(ctx context.Context, req *pb.GetDecisionFactorReq) (*pb.GetDecisionFactorResp, error) {
	factor, err := s.rcDecisionFactor.GetByContentIdWithDataScope(ctx, req.ContentId)
	if err != nil {
		return nil, err
	}
	if factor == nil {
		return &pb.GetDecisionFactorResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "record not found/data not accessible",
			Data:    nil,
		}, nil
	}

	return &pb.GetDecisionFactorResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data: &pb.GetDecisionFactorResp_DecisionFactorData{
			UscId:     factor.UscId,
			LhQylx:    int32(factor.LhQylx),
			LhCylwz:   int32(factor.LhCylwz),
			LhGdct:    int32(factor.LhGdct),
			LhYhsx:    int32(factor.LhYhsx),
			LhSfsx:    int32(factor.LhSfsx),
			CreatedAt: factor.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy: factor.UserId,
			ClaimId:   factor.ClaimId,
		},
	}, nil
}

// InsertReportDecisionFactor 录入企业风控参数（认领企业报告）
func (s *RcServiceServicer) InsertReportDecisionFactor(ctx context.Context, req *pb.InsertReportDecisionFactorReq) (*pb.InsertReportDecisionFactorResp, error) {

	countRdf, err := s.rcDecisionFactor.CountByUscIdAndUserId(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if countRdf > 0 {
		return &pb.InsertReportDecisionFactorResp{
			Success: false,
			Code:    0,
			Msg:     "该企业风控参数已录入.",
		}, nil
	}

	insertReq := dto.RcDecisionFactor{
		UscId:   req.UscId,
		LhQylx:  int(req.LhQylx),
		LhCylwz: int(req.LhCylwz),
		LhGdct:  int(req.LhGdct),
		LhYhsx:  int(req.LhYhsx),
		LhSfsx:  int(req.LhSfsx),
	}
	rdfId, err := s.rcDecisionFactor.Insert(ctx, &insertReq)
	if err != nil {
		return nil, err
	}
	contentIds, err := s.rcOriginContent.GetContentIdsByUscId(ctx, req.UscId)
	if err != nil {
		return nil, err
	}

	for _, contentId := range contentIds {
		contentId := contentId
		claimReq := dto.RcContentFactorClaim{
			ContentId: contentId,
			FactorId:  rdfId,
		}
		_, err := s.rcDecisionFactor.InsertClaimNoDupe(ctx, &claimReq)
		if err != nil {
			return nil, err
		}
	}

	return &pb.InsertReportDecisionFactorResp{
		Success: true,
		Code:    http.StatusAccepted,
		Msg:     "",
	}, nil
}

// UpdateReportDecisionFactor 更新企业风控参数
// Todo: implement me
func (s *RcServiceServicer) UpdateReportDecisionFactor(ctx context.Context, req *pb.UpdateReportDecisionFactorReq) (*pb.InsertReportDecisionFactorResp, error) {
	// logic: insert to rdf, got newFactorId
	// get contentId by claimId, insert new row to claim table with contentId and newFactorId
	return nil, nil
}

// ListReport 获取报告列表
func (s *RcServiceServicer) ListReport(ctx context.Context, req *pb.ListReportKwdSearchReq) (*pb.ListReportResp, error) {
	page := &dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	list, pageInfo, err := s.rcDecisionFactor.ListReportClaimed(ctx, page, req.NameKwd)
	if err != nil {
		return nil, err
	}
	infos := make([]*pb.ListReportResp_ReportInfo, 0)
	for _, item := range *list {
		infos = append(infos, &pb.ListReportResp_ReportInfo{
			UscId:            item.UscId,
			ContentId:        item.ContentId,
			EnterpriseName:   item.EnterpriseName,
			DataCollectMonth: item.DataCollectMonth,
		})
	}

	sem := make(chan struct{}, 10)
	errCh := make(chan error, 1)
	done := make(chan bool, 1)
	var wg sync.WaitGroup
	var checkProcessedFunc func(item *pb.ListReportResp_ReportInfo)

	switch req.Version {
	case pb.ReportVersion_V3:
		checkProcessedFunc = func(item *pb.ListReportResp_ReportInfo) {
			defer wg.Done()
			defer func() { <-sem }()
			processedId, createdAt, err := s.mgoRc.GetNewestDocInfoByContentId(ctx, item.ContentId)
			if err != nil {
				errCh <- err
				return
			}
			if processedId != "" {
				item.Available = true
				item.ContentUpdatedTime = createdAt.Format("2006-01-02 15:04:05")
			} else {
				item.Available = false
			}
		}
	case pb.ReportVersion_V2:
		checkProcessedFunc = func(item *pb.ListReportResp_ReportInfo) {
			defer wg.Done()
			defer func() { <-sem }()
			processedId, createdAt, err := s.rcProcessedContent.GetNewestRowInfoByContentId(ctx, item.ContentId)
			if err != nil {
				errCh <- err
				return
			}
			if processedId != 0 {
				item.Available = true
				item.ContentUpdatedTime = createdAt.Format("2006-01-02 15:04:05")
			} else {
				item.Available = false
			}
		}
	default:
		return nil, errors.New(400, "invalid version", string(req.Version))
	}

	for _, item := range infos {
		wg.Add(1)
		sem <- struct{}{}
		go checkProcessedFunc(item)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case err := <-errCh:
		return nil, err
	case <-done:
		break
	}

	return &pb.ListReportResp{
		Success: true,
		Msg:     "",
		Code:    0,
		Total:   uint32(pageInfo.Total),
		Offset:  uint32(pageInfo.Offset),
		Data:    infos,
	}, nil
}

// GetReportContent 获取企业报告内容
func (s *RcServiceServicer) GetReportContent(ctx context.Context, req *pb.GetReportContentReq) (*pb.GetReportContentResp, error) {
	accessible, err := s.rcDecisionFactor.CheckContentIdAccessible(ctx, req.ContentId)
	if err != nil {
		return nil, err
	}
	if !accessible {
		return &pb.GetReportContentResp{
			Success: false,
			Code:    http.StatusForbidden,
			Msg:     "data not accessible",
			Data:    nil,
		}, nil
	}

	m := make(map[string]interface{})

	// get processed content
	switch req.Version {
	case pb.ReportVersion_V2:
		data, err := s.rcProcessedContent.GetNewestByContentId(ctx, req.ContentId)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return &pb.GetReportContentResp{
				Success: false,
				Code:    http.StatusNoContent,
				Msg:     "record not found",
				Data:    nil,
			}, nil
		}
		err = json.Unmarshal([]byte(data.Content), &m)
		if err != nil {
			return nil, err
		}
	case pb.ReportVersion_V3:
		data, err := s.mgoRc.GetNewestDocByContentId(ctx, req.ContentId)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return &pb.GetReportContentResp{
				Success: false,
				Code:    http.StatusNoContent,
				Msg:     "record not found",
				Data:    nil,
			}, nil
		}
		b, err := bson.MarshalExtJSON(data, true, false)
		if err != nil {
			return nil, err
		}
		var temp map[string]interface{}
		if err := json.Unmarshal(b, &temp); err != nil {
			return nil, err
		}
		m = temp["content"].(map[string]interface{})

	default:
		return nil, errors.New(400, "invalid version", string(req.Version))
	}
	st, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}

	return &pb.GetReportContentResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    st,
	}, nil
}
