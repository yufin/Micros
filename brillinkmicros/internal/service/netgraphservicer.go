package service

import (
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/biz/dto"
	"context"
	"github.com/go-kratos/kratos/v2/log"

	pb "brillinkmicros/api/graph/v1"
)

type NetGraphServiceServicer struct {
	pb.UnimplementedNetGraphServiceServer
	log   *log.Helper
	graph *biz.GraphUsecase
}

func NewNetGraphServiceServicer(gn *biz.GraphUsecase, logger log.Logger) *NetGraphServiceServicer {
	return &NetGraphServiceServicer{
		graph: gn,
		log:   log.NewHelper(logger),
	}
}

func (s *NetGraphServiceServicer) GetNetExpand(ctx context.Context, req *pb.NetExpandReq) (*pb.NetResp, error) {
	f := dto.PathFilter{
		RelLabels:    nil,
		NodeLabels:   nil,
		MaxPathDepth: 0,
	}
	res, err := s.graph.GetPathExpand(ctx, req.Id, req.Depth, req.Limit, &f)
	if err != nil {
		return nil, err
	}
	if len(*res) == 0 {
		return &pb.NetResp{
			Success: true,
			Code:    0,
			Msg:     "empty result",
			Data:    nil,
		}, nil
	}
	net := dto.Net{}
	net.Gen(res)
	data := pb.Net{}
	net.GenPb(&data)

	return &pb.NetResp{
		Success: true,
		Code:    0,
		Msg:     "success",
		Data:    &data,
	}, nil
}

func (s *NetGraphServiceServicer) GetNode(ctx context.Context, req *pb.GetNodeReq) (*pb.NodeResp, error) {
	n, err := s.graph.GetNode(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	data := pb.Node{}
	n.GenPb(&data)

	return &pb.NodeResp{
		Success: true,
		Code:    0,
		Msg:     "success",
		Data:    &data,
	}, nil
}

func (s *NetGraphServiceServicer) GetChildrenNet(ctx context.Context, req *pb.GetPaginationNodeReq) (*pb.NetPaginationResp, error) {
	pReq := dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	res, count, err := s.graph.GetPathToChildren(ctx, req.Id, pReq, nil)
	if err != nil {
		return nil, err
	}
	if len(*res) == 0 {
		return &pb.NetPaginationResp{
			Success:  true,
			Code:     0,
			Msg:      "",
			Total:    int32(count),
			Current:  int32(req.PageNum),
			PageSize: int32(req.PageSize),
			Data:     nil,
		}, nil
	}
	net := dto.Net{}
	net.Gen(res)
	data := pb.Net{}
	net.GenPb(&data)

	return &pb.NetPaginationResp{
		Success:  true,
		Code:     0,
		Msg:      "",
		Total:    int32(count),
		Current:  int32(req.PageNum),
		PageSize: int32(req.PageSize),
		Data:     &data,
	}, nil
}
func (s *NetGraphServiceServicer) GetParentsNet(ctx context.Context, req *pb.GetPaginationNodeReq) (*pb.NetPaginationResp, error) {
	pReq := dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	res, count, err := s.graph.GetPathToParent(ctx, req.Id, pReq, nil)
	if err != nil {
		return nil, err
	}
	if len(*res) == 0 {
		return &pb.NetPaginationResp{
			Success:  true,
			Code:     0,
			Msg:      "",
			Total:    int32(count),
			Current:  int32(req.PageNum),
			PageSize: int32(req.PageSize),
			Data:     nil,
		}, nil
	}
	net := dto.Net{}
	net.Gen(res)
	data := pb.Net{}
	net.GenPb(&data)

	return &pb.NetPaginationResp{
		Success:  true,
		Code:     0,
		Msg:      "",
		Total:    int32(count),
		Current:  int32(req.PageNum),
		PageSize: int32(req.PageSize),
		Data:     &data,
	}, nil
}
