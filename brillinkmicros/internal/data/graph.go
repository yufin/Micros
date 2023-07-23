package data

import (
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/biz/dto"
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphRepo struct {
	data *Data
	log  *log.Helper
}

func NewGraphRepo(data *Data, logger log.Logger) biz.GraphRepo {
	return &GraphRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *GraphRepo) GetPathBetween(ctx context.Context, sourceId string, targetId string, f *dto.PathFilter, p *dto.PaginationReq) ([]neo4j.Path, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *GraphRepo) GetNode(ctx context.Context, id string) (*dto.Node, error) {
	cypher := "MATCH (n {id: $id}) RETURN n;"
	res, err := CypherQuery(repo.data.Neo, ctx, cypher, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	n, found := res[0].Get("n")
	if !found {
		return nil, errors.New("node with specified id not found")
	}
	var node dto.Node
	node.Gen(n.(neo4j.Node))
	return &node, nil
}

func (repo *GraphRepo) GetNodes(ctx context.Context, ids []string) ([]*dto.Node, error) {
	cypher := "MATCH (n) where n.id in $ids RETURN n;"
	res, err := CypherQuery(repo.data.Neo, ctx, cypher, map[string]interface{}{"ids": ids})
	if err != nil {
		return nil, err
	}
	nodes := make([]*dto.Node, 0)
	for _, item := range res {
		item := item
		n, _ := item.Get("n")
		var node dto.Node
		node.Gen(n.(neo4j.Node))
		nodes = append(nodes, &node)
	}
	return nodes, nil
}

func (repo *GraphRepo) GetChildren(ctx context.Context, id string, f *dto.PathFilter, p *dto.PaginationReq) ([]*dto.Node, error) {
	offset := (p.PageNum - 1) * p.PageSize
	cypher := `MATCH (p)-[r]->(c) 
			WHERE p.id = $nodeId 
			AND any(label IN labels(c) WHERE label IN $childLabels) 
			AND type(r) IN $relTypes 
			RETURN c 
			SKIP $offset LIMIT $pageSize;`

	res, err := CypherQuery(repo.data.Neo, ctx, cypher, map[string]any{
		"nodeId":      id,
		"childLabels": f.NodeLabels,
		"relTypes":    f.RelLabels,
		"offset":      offset,
		"pageSize":    p.PageSize,
	})
	if err != nil {
		return nil, err
	}
	nodes := make([]*dto.Node, 0)
	for _, item := range res {
		item := item
		n, _ := item.Get("c")
		var node dto.Node
		node.Gen(n.(neo4j.Node))
		nodes = append(nodes, &node)
	}
	return nodes, nil
}

func (repo *GraphRepo) CountChildren(ctx context.Context, id string, f *dto.PathFilter, amount *int64) error {
	cypher := `MATCH (p)-[r]->(c) 
			WHERE p.id = $nodeId 
			AND any(label IN labels(c) WHERE label IN $childLabels) 
			AND type(r) IN $relTypes 
			RETURN count(c) as childrenCount;`
	res, err := CypherQuery(repo.data.Neo, ctx, cypher, map[string]any{
		"nodeId":      id,
		"childLabels": f.NodeLabels,
		"relTypes":    f.RelLabels,
	})
	if err != nil {
		return err
	}
	if len(res) == 0 {
		return errors.New("empty result")
	}
	cc, found := res[0].Get("childrenCount")
	if !found {
		return errors.New("key not found")
	}
	cci, ok := cc.(int64)
	if !ok {
		return errors.New("result type assert error")
	}
	*amount = cci
	return nil
}

func (repo *GraphRepo) GetTitleAutoComplete(ctx context.Context, f *dto.PathFilter, p *dto.PaginationReq, kw string) ([]*dto.TitleAutoCompleteRes, error) {
	//var relLabel string
	//if limitLabel == "Company" {
	//	relLabel = "ATTACH_TO"
	//} else if limitLabel == "Tag" {
	//	relLabel = "CLASSIFY_OF"
	//} else {
	//	return nil, errors.New("invalid limit label")
	//}
	kwPattern := fmt.Sprintf("(?i).*%s.*", kw)
	offset := (p.PageNum - 1) * p.PageSize
	cypher := `match ()-[r]->(n) 
			where any(label IN labels(n) WHERE label IN $limitLabels) 
			AND type(r) in $relTypes 
			AND n.title =~ $kwPattern 
			WITH distinct n 
			WITH {title: n.title, id: n.id} as res skip $offset limit $pageSize 
			WITH collect(res) as propList 
			RETURN propList;`

	tac := make([]*dto.TitleAutoCompleteRes, 0)
	res, err := CypherQuery(repo.data.Neo, ctx, cypher, map[string]any{
		"limitLabels": f.NodeLabels,
		"relTypes":    f.RelLabels,
		"kwPattern":   kwPattern,
		"offset":      offset,
		"pageSize":    p.PageSize,
	})
	if err != nil {
		return nil, err
	}
	props, _ := res[0].Get("propList")
	for _, item := range props.([]any) {
		item := item.(map[string]any)
		tac = append(tac, &dto.TitleAutoCompleteRes{
			Title: item["title"].(string),
			Id:    item["id"].(string),
		})
	}
	return tac, nil
}

func (repo *GraphRepo) CountTitleAutoComplete(ctx context.Context, f *dto.PathFilter, kw string, amount *int64) error {
	//var relLabel string
	//if limitLabel == "Company" {
	//	relLabel = "ATTACH_TO"
	//} else if limitLabel == "Tag" {
	//	relLabel = "CLASSIFY_OF"
	//} else {
	//	return errors.New("invalid limit label")
	//}
	kwPattern := fmt.Sprintf("(?i).*%s.*", kw)
	cypher := `match ()-[r]->(n) 
			where any(label IN labels(n) WHERE label IN $limitLabels) 
			AND type(r) in $relTypes 
			AND n.title =~ $kwPattern 
			WITH distinct n 
			return count(n) as counts;`
	res, err := CypherQuery(repo.data.Neo, ctx, cypher, map[string]any{
		"limitLabels": f.NodeLabels,
		"relTypes":    f.RelLabels,
		"kwPattern":   kwPattern,
	})
	if err != nil {
		return err
	}
	if len(res) == 0 {
		return errors.New("empty result")
	}
	cc, found := res[0].Get("counts")
	if !found {
		return errors.New("key not found")
	}
	cci, ok := cc.(int64)
	if !ok {
		return errors.New("result type assert error")
	}
	*amount = cci
	return nil
}
