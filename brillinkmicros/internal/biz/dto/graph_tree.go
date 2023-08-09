package dto

import (
	pb "brillinkmicros/api/graph/v1"
	"brillinkmicros/internal/biz"
	"brillinkmicros/pkg"
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"sync"
)

type TreeNode struct {
	Node
	RandId        string
	ChildrenCount *int64
	Children      []TreeNode
}

type childrenCountParam struct {
	ParentId       string
	ChildrenCountP *int64
}

func (t *TreeNode) Gen(n neo4j.Node) {
	t.Node.Gen(n)
	t.RandId = pkg.RandUuid()
}

func (t *TreeNode) AutoGen(ctx context.Context, n neo4j.Node, repo biz.GraphRepo, filter PathFilter) error {
	t.Node.Gen(n)
	t.RandId = pkg.RandUuid()
	var count int64
	err := repo.CountChildren(ctx, t.Id, filter, &count)
	if err != nil {
		return err
	}
	t.ChildrenCount = &count
	return nil
}

func (t *TreeNode) setChild(parentId string, neoChild neo4j.Node, chParamP *chan childrenCountParam) bool {
	if t.Id == parentId {
		if t.Children != nil {
			for _, child := range t.Children {
				if child.Id == neoChild.Props["id"].(string) {
					return true
				}
			}
			child := TreeNode{}
			child.Gen(neoChild)
			ccp := new(int64)
			child.ChildrenCount = ccp
			*chParamP <- childrenCountParam{
				ParentId:       child.Id,
				ChildrenCountP: ccp,
			}
			t.Children = append(t.Children, child)
			return true
		} else {
			child := TreeNode{}
			child.Gen(neoChild)
			ccp := new(int64)
			child.ChildrenCount = ccp
			*chParamP <- childrenCountParam{
				ParentId:       child.Id,
				ChildrenCountP: ccp,
			}
			t.Children = []TreeNode{child}
			return true
		}
	} else {
		if t.Children != nil {
			for _, child := range t.Children {
				r := child.setChild(parentId, neoChild, chParamP)
				if r {
					return true
				}
			}
		}
	}
	return false
}

func NewTreeNodeFromPath(ctx context.Context, repo biz.GraphRepo, neoPath *[]neo4j.Path, filter PathFilter) (*TreeNode, error) {
	rootNeo := (*neoPath)[0].Nodes[0]
	root := TreeNode{}
	root.Gen(rootNeo)
	var buffSize int
	for _, path := range *neoPath {
		buffSize += len(path.Relationships)
	}
	chCountingParam := make(chan childrenCountParam, buffSize)

	go func() {
		for _, path := range *neoPath {
			for _, rel := range path.Relationships {
				parent := GetNodeByElementId(&path.Nodes, rel.StartElementId)
				child := GetNodeByElementId(&path.Nodes, rel.EndElementId)
				if parent != nil && child != nil {
					pn := Node{}
					pn.Gen(*parent)
					root.setChild(pn.Id, *child, &chCountingParam)
				}
			}
		}
		close(chCountingParam)
	}()

	var wg sync.WaitGroup
	chErr := make(chan error)

	for ccp := range chCountingParam {
		ccp := ccp
		go func() {
			wg.Add(1)
			defer wg.Done()
			err := repo.CountChildren(ctx, ccp.ParentId, filter, ccp.ChildrenCountP)
			if err != nil {
				chErr <- err
			}
		}()
	}

	select {
	case err := <-chErr:
		return nil, err
	default:
		wg.Wait()
	}
	return &root, nil
}

func (t *TreeNode) GenPb() *pb.TreeNode {
	pbt := pb.TreeNode{}
	pbt.EntityId = t.Id
	pbt.Id = t.RandId
	pbt.Title = t.Title
	pbt.Labels = t.Labels
	pbt.ChildrenCount = int32(*t.ChildrenCount)
	if t.Children != nil {
		if pbt.Children == nil {
			pbt.Children = make([]*pb.TreeNode, 0)
		}
		for _, child := range t.Children {
			pbt.Children = append(pbt.Children, child.GenPb())
		}
	}
	return &pbt
}
