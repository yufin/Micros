package service

import (
	"context"

	pb "brillinkmicros/api/graph/v1"
)

type TreeGraphService struct {
	pb.UnimplementedTreeGraphServer
}

func NewTreeGraphService() *TreeGraphService {
	return &TreeGraphService{}
}

func (s *TreeGraphService) GetNodeById(ctx context.Context, req *pb.IdReq) (*pb.TreeNodeResp, error) {
	return &pb.TreeNodeResp{}, nil
}
func (s *TreeGraphService) GetChildren(ctx context.Context, req *pb.PgIdReq) (*pb.TreeNodesResp, error) {
	return &pb.TreeNodesResp{}, nil
}
func (s *TreeGraphService) GetTitleAutoComplete(ctx context.Context, req *pb.TitleAutoCompleteReq) (*pb.TitleAutoCompleteResp, error) {
	return &pb.TitleAutoCompleteResp{}, nil
}
func (s *TreeGraphService) GetPathBetween(ctx context.Context, req *pb.GetPathReq) (*pb.TreeNodeResp, error) {
	return &pb.TreeNodeResp{}, nil
}
