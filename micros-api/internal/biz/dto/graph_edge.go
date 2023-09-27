package dto

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"google.golang.org/protobuf/types/known/structpb"
	pb "micros-api/api/graph/v1"
)

type Edge struct {
	SourceId string
	TargetId string
	Id       string
	Type     string
	Data     map[string]any
}

func (e *Edge) Gen(r neo4j.Relationship) {
	propsCopy := make(map[string]any)
	for k, v := range r.GetProperties() {
		switch k {
		case "id":
			id, ok := v.(string)
			if ok {
				e.Id = id
			}
		default:
			propsCopy[k] = v
		}
	}
	e.Type = r.Type
	e.Data = propsCopy
}

func (e *Edge) GenPb(pb *pb.Edge) {
	pb.Id = e.Id
	pb.Source = e.SourceId
	pb.Target = e.TargetId
	pb.Label = e.Type
	st, _ := structpb.NewStruct(e.Data)
	pb.Data = st
}
