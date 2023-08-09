package service

import (
	"context"

	pb "brillinkmicros/api/graph/v1"
)

type NetGraphServiceServicer struct {
	pb.UnimplementedNetGraphServiceServer
}

func NewNetGraphServiceServicer() *NetGraphServiceServicer {
	return &NetGraphServiceServicer{}
}

func (s *NetGraphServiceServicer) GetNetExpand(ctx context.Context, req *pb.NetExpandReq) (*pb.NetResp, error) {

	return &pb.NetResp{}, nil
}
func (s *NetGraphServiceServicer) GetNode(ctx context.Context, req *pb.PgIdReq) (*pb.NodeResp, error) {
	return &pb.NodeResp{}, nil
}
func (s *NetGraphServiceServicer) GetChildren(ctx context.Context, req *pb.PgIdReq) (*pb.NetPaginationResp, error) {
	return &pb.NetPaginationResp{}, nil
}
func (s *NetGraphServiceServicer) GetParents(ctx context.Context, req *pb.PgIdReq) (*pb.NetPaginationResp, error) {
	return &pb.NetPaginationResp{}, nil
}
