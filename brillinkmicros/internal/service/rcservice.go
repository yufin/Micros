package service

import (
	pb "brillinkmicros/api/rc/v1"
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/service/util"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type RcServiceService struct {
	pb.UnimplementedRcServiceServer
	log                *log.Helper
	rcProcessedContent *biz.RcProcessedContentUsecase
	rcOriginContent    *biz.RcOriginContentUsecase
}

func NewRcServiceService(
	rpc *biz.RcProcessedContentUsecase, roc *biz.RcOriginContentUsecase, logger log.Logger) *RcServiceService {
	return &RcServiceService{
		rcOriginContent:    roc,
		rcProcessedContent: rpc,
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
		name, _ := util.ParseEnterpriseName(v.Content)
		available := false
		if v.ProcessedId != 0 {
			available = true
		}

		info := &pb.ReportInfo{
			ContentId:          v.ContentId,
			EnterpriseName:     name,
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
	dataRpc, err := s.rcProcessedContent.GetByContentIdUpToDate(ctx, req.ContentId)
	if err != nil {
		return nil, err
	}
	if dataRpc == nil {
		return &pb.ReportContentResp{
			Content:   "",
			Available: false,
		}, nil
	}
	return &pb.ReportContentResp{
		Content:   dataRpc.Content,
		Available: true,
	}, nil
}

func (s *RcServiceService) RefreshReportContent(ctx context.Context, req *pb.ReportContentReq) (*pb.ReportContentResp, error) {
	return &pb.ReportContentResp{}, nil
}
func (s *RcServiceService) SetReportAdditionData(ctx context.Context, req *pb.SetAdditionDataReq) (*pb.SetAdditionDataResp, error) {
	return &pb.SetAdditionDataResp{}, nil
}
