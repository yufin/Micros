package dto

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	pb "micros-api/api/graph/v1"
	"micros-api/pkg"
	"sync"
)

type TreeNodeGenRepo interface {
	CountChildren(ctx context.Context, id string, f PathFilter, amount *int64) error
}

func (t *TreeNode) TreeDefaultPathFilter() *PathFilter {
	return &PathFilter{
		NodeLabels: []string{"Tag", "Company", "Classification", "Application"},
		RelLabels:  []string{"ATTACH_TO", "CLASSIFY_OF", "APPLICATION_OF"},
	}
}

type TreeNode struct {
	Node
	RandId        string
	ChildrenCount *int64
	Children      []*TreeNode
}

func (t *TreeNode) Gen(n neo4j.Node) {
	t.Node.Gen(n)
	t.RandId = pkg.RandUuid()
}

func (t *TreeNode) AutoGen(n neo4j.Node, repo TreeNodeGenRepo, filter *PathFilter) error {
	t.Node.Gen(n)
	err := t.CountChildren(repo, filter)
	if err != nil {
		return err
	}
	return nil
}

func (t *TreeNode) CountChildren(repo TreeNodeGenRepo, filter *PathFilter) error {
	var count int64
	if filter == nil {
		filter = t.TreeDefaultPathFilter()
	}
	err := repo.CountChildren(context.TODO(), t.Id, *filter, &count)
	if err != nil {
		return err
	}
	t.ChildrenCount = &count
	return nil
}

func (t *TreeNode) setChild(parentId string, neoChild neo4j.Node) bool {
	if t.Id == parentId {
		if t.Children != nil {
			for _, child := range t.Children {
				if child.Id == neoChild.Props["id"].(string) {
					return true
				}
			}
			child := TreeNode{}
			child.Gen(neoChild)
			t.Children = append(t.Children, &child)
			return true
		} else {
			child := TreeNode{}
			child.Gen(neoChild)
			ccp := new(int64)
			child.ChildrenCount = ccp

			t.Children = []*TreeNode{&child}
			return true
		}
	} else {
		if t.Children != nil {
			for _, child := range t.Children {
				r := child.setChild(parentId, neoChild)
				if r {
					return true
				}
			}
		}
	}
	return false
}

func (t *TreeNode) CountChildrenRecursively(repo TreeNodeGenRepo, filter *PathFilter) error {
	waitingCount := make([]*TreeNode, 0)
	t.collectWaitingCount(&waitingCount)
	return CountChildrenParallel(repo, &waitingCount, filter)
}

func CountChildrenParallel(repo TreeNodeGenRepo, waitingCount *[]*TreeNode, filter *PathFilter) error {
	var wg sync.WaitGroup
	chErr := make(chan error)
	done := make(chan struct{})
	for _, node := range *waitingCount {
		node := node
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := node.CountChildren(repo, filter)
			if err != nil {
				chErr <- err
			}
		}()
	}
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case err := <-chErr:
		return err
	}
}

func (t *TreeNode) collectWaitingCount(waiting *[]*TreeNode) {
	*waiting = append(*waiting, t)
	if t.Children == nil {
		return
	}
	for _, child := range t.Children {
		if child != nil {
			child.collectWaitingCount(waiting)
		}
	}
}

func NewTreeNodeFromPath(ctx context.Context, repo TreeNodeGenRepo, neoPath *[]neo4j.Path, filter PathFilter) (*TreeNode, error) {
	rootNeo := (*neoPath)[0].Nodes[0]
	root := TreeNode{}
	root.Gen(rootNeo)

	for _, path := range *neoPath {
		for _, rel := range path.Relationships {
			parent := GetNodeByElementId(&path.Nodes, rel.StartElementId)
			child := GetNodeByElementId(&path.Nodes, rel.EndElementId)
			if parent != nil && child != nil {
				pn := Node{}
				pn.Gen(*parent)
				root.setChild(pn.Id, *child)
			}
		}
	}

	err := root.CountChildrenRecursively(repo, &filter)
	if err != nil {
		return nil, err
	}
	return &root, nil
}

func (t *TreeNode) GenPb() *pb.TreeNode {
	pbt := pb.TreeNode{}
	pbt.EntityId = t.Id
	pbt.Id = t.RandId
	pbt.Title = t.Title
	pbt.Labels = t.Labels
	if t.ChildrenCount != nil {
		pbt.ChildrenCount = int32(*t.ChildrenCount)
	} else {
		pbt.ChildrenCount = 0
	}
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
