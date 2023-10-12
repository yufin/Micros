package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"micros-api/internal/biz/dto"
	"sync"

	pb "micros-api/api/graph/v1"
	"micros-api/internal/biz"
)

type TreeGraphServiceServicer struct {
	pb.UnimplementedTreeGraphServiceServer
	log   *log.Helper
	graph *biz.GraphUsecase
}

func NewTreeGraphServiceServicer(gn *biz.GraphUsecase, logger log.Logger) *TreeGraphServiceServicer {
	return &TreeGraphServiceServicer{
		graph: gn,
		log:   log.NewHelper(logger),
	}
}

func (s *TreeGraphServiceServicer) GetTreeNode(ctx context.Context, req *pb.IdReq) (*pb.TreeNodeResp, error) {
	n, err := s.graph.GetNode(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	node := dto.TreeNode{}
	if err := node.AutoGen(n, s.graph, nil); err != nil {
		return nil, err
	}
	data := node.GenPb()
	return &pb.TreeNodeResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    data,
	}, nil
}

func (s *TreeGraphServiceServicer) GetChildren(ctx context.Context, req *pb.PgIdReq) (*pb.TreeNodesResp, error) {
	filter := dto.PathFilter{
		NodeLabels: treeGraphLimitNodeLabels(),
		RelLabels:  treeGraphLimitRelLabels(),
	}
	p := dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	children, err := s.graph.GetChildren(ctx, req.Id, filter, p)
	if err != nil {
		return nil, err
	}
	nodes := make([]*dto.TreeNode, 0)
	for _, item := range *children {
		node := dto.TreeNode{}
		node.Gen(item)
		nodes = append(nodes, &node)
	}
	err = dto.CountChildrenParallel(s.graph, &nodes, &filter)
	if err != nil {
		return nil, err
	}
	treeNodes := make([]*pb.TreeNode, 0)
	for _, item := range nodes {
		treeNodes = append(treeNodes, item.GenPb())
	}

	return &pb.TreeNodesResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    treeNodes,
	}, nil
}

func (s *TreeGraphServiceServicer) GetTitleAutoComplete(ctx context.Context, req *pb.TitleAutoCompleteReq) (*pb.TitleAutoCompleteResp, error) {
	var filter dto.PathFilter
	switch req.NodeLabel {
	case pb.LabelType_COMPANY:
		filter = dto.PathFilter{
			RelLabels:  []string{"ATTACH_TO"},
			NodeLabels: []string{"Company"},
		}
	case pb.LabelType_TAG:
		filter = dto.PathFilter{
			RelLabels:  []string{"CLASSIFY_OF"},
			NodeLabels: []string{"Tag"},
		}
	default:
		return nil, errors.New(401, "Invalid nodeType", "")
	}

	p := dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	var resGet *[]dto.TitleAutoCompleteRes
	data := make([]*pb.TitleAutoComplete, 0)
	var (
		count      int64
		errG, errC error
		wg         sync.WaitGroup
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		resGet, errG = s.graph.GetTitleAutoComplete(ctx, filter, p, req.KeyWord)
	}()

	go func() {
		defer wg.Done()
		errC = s.graph.CountTitleAutoComplete(ctx, filter, req.KeyWord, &count)
	}()
	wg.Wait()
	if errG != nil {
		return nil, errG
	}
	if errC != nil {
		return nil, errC
	}

	for _, item := range *resGet {
		data = append(data, &pb.TitleAutoComplete{
			Id:    item.Id,
			Title: item.Title,
		})
	}

	return &pb.TitleAutoCompleteResp{
		Total:    uint32(count),
		Current:  req.PageNum,
		PageSize: req.PageSize,
		Success:  true,
		Code:     200,
		Msg:      "",
		Data:     data,
	}, nil
}

func (s *TreeGraphServiceServicer) GetPathBetween(ctx context.Context, req *pb.GetPathReq) (*pb.TreeNodeResp, error) {
	filter := dto.PathFilter{
		RelLabels:    treeGraphLimitRelLabels(),
		NodeLabels:   treeGraphLimitNodeLabels(),
		MaxPathDepth: 3,
	}
	neoPath, err := s.graph.GetPathTo(ctx, req.Source, req.Target, 3, treeGraphLimitRelLabels())
	if err != nil {
		return nil, err
	}
	if len(*neoPath) == 0 {
		return &pb.TreeNodeResp{
			Success: true,
			Code:    0,
			Msg:     "empty result",
			Data:    nil,
		}, nil
	}

	root, err := dto.NewTreeNodeFromPath(ctx, s.graph, neoPath, filter)
	if err != nil {
		return nil, err
	}
	data := root.GenPb()
	return &pb.TreeNodeResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    data,
	}, nil
}

func (s *TreeGraphServiceServicer) GetConst(ctx context.Context, empty *emptypb.Empty) (*pb.ConstResp, error) {
	data := make(map[string]any)
	data["rootId"] = "3f543cff-5d66-44e9-805f-4d3f8c27ecd2"
	st, err := structpb.NewStruct(data)
	if err != nil {
		return nil, err
	}
	return &pb.ConstResp{
		Success: true,
		Code:    200,
		Msg:     "",
		Data:    st,
	}, nil

}
