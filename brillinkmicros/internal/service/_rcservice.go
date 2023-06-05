package service

import (
	"brillinkmicros/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"

	pb "brillinkmicros/api/rc/v1"
)

type RcServiceService struct {
	pb.UnimplementedRcServiceServer
	log                *log.Helper
	rcProcessedContent *biz.RcProcessedContentUsecase
}

func NewRcServiceService(rpc *biz.RcProcessedContentUsecase, logger log.Logger) *RcServiceService {
	return &RcServiceService{
		rcProcessedContent: rpc,
		log:                log.NewHelper(logger),
	}
}

func (s *RcServiceService) ListReportInfos(ctx context.Context, req *pb.PaginationRequest) (*pb.ReportInfosResponse, error) {
	return &pb.ReportInfosResponse{}, nil
}

func (s *RcServiceService) GetReportContent(ctx context.Context, req *pb.ReportContentRequest) (*pb.ReportContentResponse, error) {
	contentId, err := strconv.ParseInt(req.ReportId, 10, 64)
	if err != nil {
		return nil, err
	}
	dataRpc, err := s.rcProcessedContent.GetByContentIdUpToDate(ctx, contentId)
	if err != nil {
		return nil, err
	}
	if dataRpc == nil {
		return &pb.ReportContentResponse{
			Content:   "",
			Available: false,
		}, nil
	}
	return &pb.ReportContentResponse{
		Content:   dataRpc.Content,
		Available: true,
	}, nil
}
