package dto

import (
	pb "brillinkmicros/api/graph/v1"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"google.golang.org/protobuf/types/known/structpb"
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
	s.Data = propsCopy
}

func (s *Node) GenPb(pb *pb.Node) {
	pb.Id = s.Id
	pb.Labels = s.Labels
	pb.Title = s.Title
	st, _ := structpb.NewStruct(s.Data)
	pb.Data = st
}
