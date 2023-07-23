package service

import (
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/biz/dto"
	"brillinkmicros/pkg"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"sync"

	pb "brillinkmicros/api/graph/v1"
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

func (s *TreeGraphServiceServicer) GetNodeById(ctx context.Context, req *pb.IdReq) (*pb.TreeNodeResp, error) {
	var (
		n                *dto.Node
		count            int64
		errGet, errCount error
	)
	filter := &dto.PathFilter{
		NodeLabels: treeGraphLimitNodeLabels(),
		RelLabels:  treeGraphLimitRelLabels(),
	}
	var wg sync.WaitGroup
	wg.Add(2)
	func() {
		defer wg.Done()
		n, errGet = s.graph.GetNode(ctx, req.Id)
	}()
	go func() {
		defer wg.Done()
		errCount = s.graph.CountChildren(ctx, req.Id, filter, &count)
	}()
	wg.Wait()
	if errGet != nil {
		return nil, errGet
	}
	if errCount != nil {
		return nil, errCount
	}
	treeNode := &pb.TreeNode{
		EntityId:      n.Id,
		Id:            pkg.RandUuid(),
		Title:         n.Title,
		Labels:        n.Labels,
		ChildrenCount: int32(count),
		Children:      nil,
	}
	return &pb.TreeNodeResp{
		Success: true,
		Code:    0,
		Msg:     "",
		Data:    treeNode,
	}, nil
}

func (s *TreeGraphServiceServicer) GetChildren(ctx context.Context, req *pb.PgIdReq) (*pb.TreeNodesResp, error) {
	var (
		children []*dto.Node
		errGet   error
	)
	filter := &dto.PathFilter{
		NodeLabels: treeGraphLimitNodeLabels(),
		RelLabels:  treeGraphLimitRelLabels(),
	}
	p := &dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	children, errGet = s.graph.GetChildren(ctx, req.Id, filter, p)
	if errGet != nil {
		return nil, errGet
	}

	treeNodes := make([]*pb.TreeNode, 0)
	var mutex sync.Mutex
	errCh := make(chan error, len(children))
	for _, node := range children {
		node := node
		go func() {
			var count int64
			errCh <- s.graph.CountChildren(ctx, node.Id, filter, &count)
			mutex.Lock()
			treeNodes = append(treeNodes, &pb.TreeNode{
				EntityId:      node.Id,
				Id:            pkg.RandUuid(),
				Title:         node.Title,
				Labels:        node.Labels,
				ChildrenCount: int32(count),
				Children:      nil,
			})
			mutex.Unlock()
		}()
	}

	for range children {
		err := <-errCh
		if err != nil {
			return nil, err
		}
	}

	return &pb.TreeNodesResp{
		Success: true,
		Code:    0,
		Msg:     "",
		Data:    treeNodes,
	}, nil
}

func (s *TreeGraphServiceServicer) GetTitleAutoComplete(ctx context.Context, req *pb.TitleAutoCompleteReq) (*pb.TitleAutoCompleteResp, error) {
	var relLabel string
	if req.LimitLabel == "Company" {
		relLabel = "ATTACH_TO"
	} else if req.LimitLabel == "Tag" {
		relLabel = "CLASSIFY_OF"
	} else {
		return nil, errors.New(401, "Invalid limit label", "")
	}
	filter := &dto.PathFilter{
		RelLabels:  []string{relLabel},
		NodeLabels: []string{req.LimitLabel},
	}
	p := &dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	resGet := make([]*dto.TitleAutoCompleteRes, 0)
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

	for _, item := range resGet {
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
		Code:     0,
		Msg:      "",
		Data:     data,
	}, nil
}

func (s *TreeGraphServiceServicer) GetPathBetween(ctx context.Context, req *pb.GetPathReq) (*pb.TreeNodeResp, error) {
	return &pb.TreeNodeResp{}, nil
}
