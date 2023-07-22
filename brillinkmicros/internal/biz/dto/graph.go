package dto

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Node struct {
	Id     string
	Labels []string
	Title  string
	Data   map[string]any
}

func (s *Node) Gen(n neo4j.Node) {
	propsCopy := make(map[string]any)
	for k, v := range n.Props {
		switch k {
		case "id":
			id, ok := v.(string)
			if ok {
				s.Id = id
			}
		case "title":
			title, ok := v.(string)
			if ok {
				s.Title = title
			}
		default:
			propsCopy[k] = v
		}
	}
	s.Labels = n.Labels
	s.Data = n.Props
}

type Edge struct {
	SourceId string
	TargetId string
	Id       string
	Label    string
	Data     map[string]any
}

func (s *Edge) Gen (r neo4j.Relationship) {
	return
}

type Net struct {
	Nodes []Node
	Edges []Edge
}

type TreeNode struct {
	Node
	RandId        string
	ChildrenCount int64
	Children      []TreeNode
}

type PathFilter struct {
	RelLabels  []string
	NodeLabels []string
}

type TitleAutoCompleteRes struct {
	Title string
	Id    string
}
