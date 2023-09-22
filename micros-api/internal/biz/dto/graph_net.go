package dto

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	pb "micros-api/api/graph/v1"
	"strings"
)

type Net struct {
	Nodes *[]Node
	Edges *[]Edge
}

func (n *Net) Gen(p *[]neo4j.Path) {
	nodes := make([]Node, 0)
	edges := make([]Edge, 0)
	var nodeIds []string

	for _, path := range *p {
		for _, node := range path.Nodes {
			nId := node.Props["id"].(string)
			if !strings.Contains(strings.Join(nodeIds, ","), nId) {
				tempN := Node{}
				tempN.Gen(node)
				nodes = append(nodes, tempN)
				nodeIds = append(nodeIds, nId)
			}
		}

		for _, rel := range path.Relationships {
			pStartNode := GetNodeByElementId(&path.Nodes, rel.StartElementId)
			pEndNode := GetNodeByElementId(&path.Nodes, rel.EndElementId)
			e := Edge{}
			e.Gen(rel)
			if pStartNode != nil {
				e.SourceId = (*pStartNode).Props["id"].(string)
			}
			if pEndNode != nil {
				e.TargetId = (*pEndNode).Props["id"].(string)
			}
			edges = append(edges, e)
		}
	}
	n.Nodes = &nodes
	n.Edges = &edges
}

func (n *Net) GenPb(pbNet *pb.Net) {
	if pbNet.Nodes == nil {
		pbNet.Nodes = make([]*pb.Node, 0)
	}
	for _, node := range *n.Nodes {
		pbNode := node.GenPb()
		pbNet.Nodes = append(pbNet.Nodes, pbNode)
	}

	if pbNet.Edges == nil {
		pbNet.Edges = make([]*pb.Edge, 0)
	}
	for _, edge := range *n.Edges {
		pbEdge := &pb.Edge{}
		edge.GenPb(pbEdge)
		pbNet.Edges = append(pbNet.Edges, pbEdge)
	}
}
