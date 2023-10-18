package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-graph/internal/biz"

	pb "micros-graph/api/graph/v1"
)

type GraphServiceService struct {
	pb.UnimplementedGraphServiceServer
	log       *log.Helper
	graphNet  *biz.GraphNetUsecase
	graphNode *biz.GraphNodeUsecase
	graphEdge *biz.GraphEdgeUsecase
}

func NewGraphServiceService(netUc *biz.GraphNetUsecase, nodeUc *biz.GraphNodeUsecase, edgeUc *biz.GraphEdgeUsecase, logger log.Logger) *GraphServiceService {
	return &GraphServiceService{
		graphNet:  netUc,
		graphNode: nodeUc,
		graphEdge: edgeUc,
		log:       log.NewHelper(logger),
	}
}

func (s *GraphServiceService) GetEdge(ctx context.Context, req *pb.GetEdgeReq) (*pb.EdgeResp, error) {
	return &pb.EdgeResp{}, nil
}
func (s *GraphServiceService) GetEdges(ctx context.Context, req *pb.GetEdgesReq) (*pb.EdgesResp, error) {
	return &pb.EdgesResp{}, nil
}
func (s *GraphServiceService) GetEdgesByProps(ctx context.Context, req *pb.GetEdgesByPropsReq) (*pb.EdgesResp, error) {
	return &pb.EdgesResp{}, nil
}
func (s *GraphServiceService) GetNode(ctx context.Context, req *pb.GetNodeReq) (*pb.NodeResp, error) {
	return &pb.NodeResp{}, nil
}
func (s *GraphServiceService) GetNodes(ctx context.Context, req *pb.GetNodesReq) (*pb.NodesResp, error) {
	return &pb.NodesResp{}, nil
}
func (s *GraphServiceService) GetNodesByProps(ctx context.Context, req *pb.GetNodesByPropsReq) (*pb.NodesResp, error) {
	return &pb.NodesResp{}, nil
}
func (s *GraphServiceService) GetChildren(ctx context.Context, req *pb.GetChildrenReq) (*pb.NetResp, error) {
	return &pb.NetResp{}, nil
}
