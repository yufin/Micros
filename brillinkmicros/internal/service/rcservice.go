package service

import (
	"brillinkmicros/internal/biz"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/gorm"
	//"github.com/gogo/protobuf/proto/protojson"
	//"github.com/gogo/protobuf/proto/structpb"
	pb "brillinkmicros/api/rc/v1"
)

type RcServiceService struct {
	pb.UnimplementedRcServiceServer
	log                *log.Helper
	rcProcessedContent *biz.RcProcessedContentUsecase
	rcOriginContent    *biz.RcOriginContentUsecase
	rcDependencyData   *biz.RcDependencyDataUsecase
}

func NewRcServiceService(
	rpc *biz.RcProcessedContentUsecase,
	roc *biz.RcOriginContentUsecase,
	rdd *biz.RcDependencyDataUsecase,
	logger log.Logger,
) *RcServiceService {
	return &RcServiceService{
		rcOriginContent:    roc,
		rcProcessedContent: rpc,
		rcDependencyData:   rdd,
		log:                log.NewHelper(logger),
	}
}

func (s *RcServiceService) ListReportInfos(ctx context.Context, req *pb.PaginationReq) (*pb.ReportInfosResp, error) {
	pageReq := &biz.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	infosResp, err := s.rcOriginContent.GetInfos(ctx, pageReq)
	if err != nil {
		return nil, err
	}

	pbInfos := make([]*pb.ReportInfo, 0)
	for _, v := range *infosResp.Data {
		v := v
		available := false
		if v.ProcessedId != 0 {
			available = true
		}
		info := &pb.ReportInfo{
			ContentId:          v.ContentId,
			EnterpriseName:     v.EnterpriseName,
			UnifiedCreditId:    v.UscId,
			DataCollectMonth:   v.DataCollectMonth,
			Available:          available,
			ContentUpdatedTime: v.ProcessedUpdatedAt.Format("2006-01-02 15:04:05"),
			ImpStatus:          int32(v.StatusCode),
			// TODO: add i18n info
		}
		pbInfos = append(pbInfos, info)

	}

	return &pb.ReportInfosResp{
		PageNum:     uint32(infosResp.PageNum),
		PageSize:    uint32(infosResp.PageSize),
		Total:       uint32(infosResp.Total),
		TotalPage:   uint32(infosResp.TotalPage),
		ReportInfos: pbInfos,
	}, nil
}

func (s *RcServiceService) GetReportContent(ctx context.Context, req *pb.ReportContentReq) (*pb.ReportContentResp, error) {
	rpcData, err := s.rcProcessedContent.GetByContentIdUpToDate(ctx, req.ContentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.ReportContentResp{
				Content:   nil,
				Available: false,
			}, nil
		}
		return nil, err
	}
	m := make(map[string]interface{})
	if err = json.Unmarshal([]byte(rpcData.Content), &m); err != nil {
		return nil, err
	}
	var st *structpb.Struct
	st, err = structpb.NewStruct(m)
	return &pb.ReportContentResp{
		Content:   st,
		Available: true,
	}, nil
}

func (s *RcServiceService) RefreshReportContent(ctx context.Context, req *pb.ReportContentReq) (*pb.RefreshReportContentResp, error) {
	
	return &pb.RefreshReportContentResp{}, nil
}

func (s *RcServiceService) SetReportDependencyData(ctx context.Context, req *pb.SetDependencyDataReq) (*pb.SetDependencyDataResp, error) {
	dataRoc, err := s.rcOriginContent.Get(ctx, req.ContentId)
	if err != nil {
		return nil, err
	}
	insertReq := biz.RcDependencyData{
		ContentId:       req.ContentId,
		AttributedMonth: dataRoc.YearMonth,
		UscId:           dataRoc.UscId,
		LhQylx:          int(req.LhQylx),
		LhCylwz:         int(req.LhCylwz),
		LhGdct:          int(req.LhGdct),
		LhQybq:          int(req.LhQybq),
		LhYhsx:          int(req.LhYhsx),
		LhSfsx:          int(req.LhSfsx),
		AdditionData:    req.AdditionData,
	}
	_, err = s.rcDependencyData.Insert(ctx, &insertReq)
	if err != nil {
		return nil, err
	}
	return &pb.SetDependencyDataResp{
		ContentId:    insertReq.ContentId,
		LhQylx:       req.LhQylx,
		LhCylwz:      req.LhCylwz,
		LhGdct:       req.LhGdct,
		LhQybq:       req.LhQybq,
		LhYhsx:       req.LhYhsx,
		LhSfsx:       req.LhSfsx,
		AdditionData: req.AdditionData,
		Success:      true,
	}, nil
}
func (s *RcServiceService) UpdateReportDependencyData(ctx context.Context, req *pb.SetDependencyDataReq) (*pb.SetDependencyDataResp, error) {
	if req.ContentId == 0 {
		return nil, errors.BadRequest("Empty ContentId", "contentId is required")
	}
	dataRdd, err := s.rcDependencyData.GetByContentId(ctx, req.ContentId)
	if err != nil {
		return nil, err
	}

	updateReq := biz.RcDependencyData{
		BaseModel: biz.BaseModel{
			Id: dataRdd.Id,
		},
		LhQylx:       int(req.LhQylx),
		LhCylwz:      int(req.LhCylwz),
		LhGdct:       int(req.LhGdct),
		LhQybq:       int(req.LhQybq),
		LhYhsx:       int(req.LhYhsx),
		LhSfsx:       int(req.LhSfsx),
		AdditionData: req.AdditionData,
	}
	newRddId, err := s.rcDependencyData.Update(ctx, &updateReq)
	if err != nil {
		return nil, err
	}
	fmt.Println(newRddId)
	newRdd, err := s.rcDependencyData.Get(ctx, newRddId)
	if err != nil {
		return nil, err
	}
	return &pb.SetDependencyDataResp{
		ContentId:    newRdd.ContentId,
		LhQylx:       int32(newRdd.LhQylx),
		LhCylwz:      int32(newRdd.LhCylwz),
		LhGdct:       int32(newRdd.LhGdct),
		LhQybq:       int32(newRdd.LhQybq),
		LhYhsx:       int32(newRdd.LhYhsx),
		LhSfsx:       int32(newRdd.LhSfsx),
		AdditionData: newRdd.AdditionData,
		Success:      true,
	}, nil
}

func (s *RcServiceService) GetReportDependencyData(ctx context.Context, req *pb.GetDependencyDataReq) (*pb.GetDependencyDataResp, error) {
	dataRoc, err := s.rcDependencyData.GetByContentId(ctx, req.ContentId)
	if err != nil {
		return nil, err
	}
	return &pb.GetDependencyDataResp{
		ContentId:    dataRoc.ContentId,
		LhQylx:       int32(dataRoc.LhQylx),
		LhCylwz:      int32(dataRoc.LhCylwz),
		LhGdct:       int32(dataRoc.LhGdct),
		LhQybq:       int32(dataRoc.LhQybq),
		LhYhsx:       int32(dataRoc.LhYhsx),
		LhSfsx:       int32(dataRoc.LhSfsx),
		AdditionData: dataRoc.AdditionData,
	}, nil
}
