package service

import (
	"context"

	pb "brillinkmicros/api/helloworld/v1"
)

type RskcService struct {
	pb.UnimplementedRskcServer
}

func NewRskcService() *RskcService {
	return &RskcService{}
}

func (s *RskcService) ListReportInfos(ctx context.Context, req *pb.PaginationRequest) (*pb.ReportInfosResponse, error) {
	return &pb.ReportInfosResponse{}, nil
}
func (s *RskcService) GetReportContent(ctx context.Context, req *pb.ReportContentRequest) (*pb.ReportContentResponse, error) {
	return &pb.ReportContentResponse{}, nil
}
