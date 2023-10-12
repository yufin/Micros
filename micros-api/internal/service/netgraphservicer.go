package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"micros-api/internal/biz"
	"micros-api/internal/biz/dto"

	pb "micros-api/api/graph/v1"
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

func (s *NetGraphServiceServicer) GetConst(ctx context.Context, empty *emptypb.Empty) (*pb.NetConstResp, error) {
	edge := dto.Edge{}
	relDict := edge.RelTypeAliasDict()
	availableRels := make(map[string]any)
	for k, v := range relDict {
		availableRels[v] = k
	}
	st, err := structpb.NewStruct(availableRels)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &pb.NetConstResp{
		AvailableRelScope: st,
		NetDemoDefaultId:  "",
	}, nil
}

func (s *NetGraphServiceServicer) GetPathBetween(ctx context.Context, req *pb.GetPathBetweenReq) (*pb.NetResp, error) {
	//res, err := s.graph.GetPathTo(ctx, req.SourceId, req.TargetId, int(req.MaxDepth), req.RelScope)
	res, err := s.graph.GetPathBetween(ctx, req.SourceId, req.TargetId, int(req.MaxDepth), req.RelScope)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(*res) == 0 {
		return &pb.NetResp{
			Success: true,
			Code:    0,
			Msg:     "没有更多啦",
			Data:    nil,
		}, nil
	}
	net := dto.Net{}
	net.Gen(res)
	data := pb.Net{}
	net.GenPb(&data)

	return &pb.NetResp{
		Success: true,
		Code:    200,
		Msg:     "success",
		Data:    &data,
	}, nil

}

func (s *NetGraphServiceServicer) GetNetExpands(ctx context.Context, req *pb.NetExpandsReq) (*pb.NetResp, error) {
	finalNet := pb.Net{}
	for _, id := range req.Ids {
		res, err := s.graph.GetPathExpand(ctx, id, req.Depth, req.Limit, req.RelScope)
		if err != nil {
			return nil, err
		}
		if len(*res) > 0 {
			net := dto.Net{}
			net.Gen(res)
			data := pb.Net{}
			net.GenPb(&data)

			finalNet.Nodes = append(finalNet.Nodes, data.Nodes...)
			finalNet.Edges = append(finalNet.Edges, data.Edges...)
		}
	}
	return &pb.NetResp{
		Success: true,
		Code:    200,
		Msg:     "success",
		Data:    &finalNet,
	}, nil
}

func (s *NetGraphServiceServicer) GetNetExpand(ctx context.Context, req *pb.NetExpandReq) (*pb.NetResp, error) {
	res, err := s.graph.GetPathExpand(ctx, req.Id, req.Depth, req.Limit, req.RelScope)
	if err != nil {
		return nil, err
	}
	if len(*res) == 0 {
		return &pb.NetResp{
			Success: true,
			Code:    0,
			Msg:     "没有更多啦",
			Data:    nil,
		}, nil
	}
	net := dto.Net{}
	net.Gen(res)
	data := pb.Net{}
	net.GenPb(&data)

	return &pb.NetResp{
		Success: true,
		Code:    200,
		Msg:     "success",
		Data:    &data,
	}, nil
}

func (s *NetGraphServiceServicer) GetNode(ctx context.Context, req *pb.GetNodeReq) (*pb.NodeResp, error) {
	n, err := s.graph.GetNode(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	//data := pb.Node{}
	node := dto.Node{}
	node.Gen(n)
	data := node.GenPb()
	return &pb.NodeResp{
		Success: true,
		Code:    200,
		Msg:     "success",
		Data:    data,
	}, nil
}

func (s *NetGraphServiceServicer) GetChildrenNet(ctx context.Context, req *pb.GetPaginationNodeReq) (*pb.NetPaginationResp, error) {
	pReq := dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	res, count, err := s.graph.GetPathToChildren(ctx, req.Id, pReq, req.ScopeRelType)
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
		Code:     200,
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
	res, count, err := s.graph.GetPathToParent(ctx, req.Id, pReq, req.ScopeRelType)
	if err != nil {
		return nil, err
	}
	if len(*res) == 0 {
		return &pb.NetPaginationResp{
			Success:  true,
			Code:     0,
			Msg:      "没有更多啦",
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
		Code:     200,
		Msg:      "",
		Total:    int32(count),
		Current:  int32(req.PageNum),
		PageSize: int32(req.PageSize),
		Data:     &data,
	}, nil
}

func (s *NetGraphServiceServicer) GetAvailableRelTypeToParents(ctx context.Context, req *pb.GetNodeReq) (*pb.AvailableRelTypeResp, error) {
	res, err := s.graph.GetRelTypeAvailable(ctx, req.Id, 0)
	if err != nil {
		return nil, err
	}
	var edge dto.Edge
	data := make(map[string]any)
	for _, rel := range res {
		data[edge.GetRelTypeAlias(rel)] = rel
	}

	st, err := structpb.NewStruct(data)
	if err != nil {
		return nil, err
	}
	return &pb.AvailableRelTypeResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    st,
	}, nil
}

func (s *NetGraphServiceServicer) GetAvailableRelTypeToChildren(ctx context.Context, req *pb.GetNodeReq) (*pb.AvailableRelTypeResp, error) {
	res, err := s.graph.GetRelTypeAvailable(ctx, req.Id, 1)
	if err != nil {
		return nil, err
	}
	var edge dto.Edge
	data := make(map[string]any)
	for _, rel := range res {
		data[edge.GetRelTypeAlias(rel)] = rel
	}

	st, err := structpb.NewStruct(data)
	if err != nil {
		return nil, err
	}
	return &pb.AvailableRelTypeResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    st,
	}, nil
}
