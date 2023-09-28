package dto

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"google.golang.org/protobuf/types/known/structpb"
	pb "micros-api/api/graph/v1"
	"time"
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
	pb.Label = e.modifyRelTypeToAlias()
	e.modifyTimeToStringFromData()
	st, err := structpb.NewStruct(e.Data)
	if err != nil {
		panic(err)
	}
	pb.Data = st
}

func (e *Edge) modifyRelTypeToAlias() string {
	trans, ok := e.RelTypeAliasDict()[e.Type]
	if ok {
		return trans
	}
	return "关联"
}

func (e *Edge) modifyTimeToStringFromData() {
	for k, v := range e.Data {
		if t, ok := v.(time.Time); ok {
			e.Data[k] = t.Format("2006-01-02")
			//e.Data[k] = t.Format(time.RFC3339)
		}
	}
}

func (*Edge) RelTypeAliasDict() map[string]string {
	return map[string]string{
		"ATTACH_TO":      "标签",
		"CLASSIFY_OF":    "归属",
		"APPLICATION_OF": "应用",
		"INVOICED_TO":    "供货",
	}
}

func (e *Edge) GetRelTypeAlias(relType string) string {
	alias, ok := e.RelTypeAliasDict()[relType]
	if ok {
		return alias
	}
	return relType
}
