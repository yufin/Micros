package dto

import "github.com/neo4j/neo4j-go-driver/v5/neo4j"

type PathFilter struct {
	RelLabels    []string
	NodeLabels   []string
	MaxPathDepth int
}

type TitleAutoCompleteRes struct {
	Title string
	Id    string
}

func GetNodeByElementId(neoNodes *[]neo4j.Node, elementId string) *neo4j.Node {
	for _, neoNode := range *neoNodes {
		if neoNode.GetElementId() == elementId {
			return &neoNode
		}
	}
	return nil
}
