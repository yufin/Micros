package data

import (
	"brillinkmicros/internal/biz"
	dto2 "brillinkmicros/internal/biz/dto"
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

func (repo *GraphRepo) GetNode(ctx context.Context, id string) (*dto2.Node, error) {
	cypher := "MATCH (n {id: $id}) RETURN n;"
	res, err := repo.data.Neo.CypherQuery(ctx, cypher, map[string]interface{}{"id": id})
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
	var node dto2.Node
	node.Gen(n.(neo4j.Node))
	return &node, nil
}

func (repo *GraphRepo) GetNodes(ctx context.Context, ids []string) (*[]dto2.Node, error) {
	cypher := "MATCH (n) where n.id in $ids RETURN n;"
	res, err := repo.data.Neo.CypherQuery(ctx, cypher, map[string]interface{}{"ids": ids})
	if err != nil {
		return nil, err
	}
	nodes := make([]dto2.Node, 0)
	for _, item := range res {
		item := item
		n, _ := item.Get("n")
		var node dto2.Node
		node.Gen(n.(neo4j.Node))
		nodes = append(nodes, node)
	}
	return &nodes, nil
}

func (repo *GraphRepo) GetChildren(ctx context.Context, id string, f dto2.PathFilter, p dto2.PaginationReq) (*[]dto2.Node, error) {
	offset := (p.PageNum - 1) * p.PageSize
	cypher := `MATCH (p)-[r]->(c) 
			WHERE p.id = $nodeId 
			AND any(label IN labels(c) WHERE label IN $childLabels) 
			AND type(r) IN $relTypes 
			RETURN c 
			SKIP $offset LIMIT $pageSize;`

	res, err := repo.data.Neo.CypherQuery(ctx, cypher, map[string]any{
		"nodeId":      id,
		"childLabels": f.NodeLabels,
		"relTypes":    f.RelLabels,
		"offset":      offset,
		"pageSize":    p.PageSize,
	})
	if err != nil {
		return nil, err
	}
	nodes := make([]dto2.Node, 0)
	for _, item := range res {
		item := item
		n, _ := item.Get("c")
		var node dto2.Node
		node.Gen(n.(neo4j.Node))
		nodes = append(nodes, node)
	}
	return &nodes, nil
}

func (repo *GraphRepo) CountChildren(ctx context.Context, id string, f dto2.PathFilter, amount *int64) error {
	cypher := `MATCH (p)-[r]->(c) 
			WHERE p.id = $nodeId 
			AND any(label IN labels(c) WHERE label IN $childLabels) 
			AND type(r) IN $relTypes 
			RETURN count(c) as childrenCount;`
	res, err := repo.data.Neo.CypherQuery(ctx, cypher, map[string]any{
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

func (repo *GraphRepo) GetTitleAutoComplete(ctx context.Context, f dto2.PathFilter, p dto2.PaginationReq, kw string) (*[]dto2.TitleAutoCompleteRes, error) {
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

	tac := make([]dto2.TitleAutoCompleteRes, 0)
	res, err := repo.data.Neo.CypherQuery(ctx, cypher, map[string]any{
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
		tac = append(tac, dto2.TitleAutoCompleteRes{
			Title: item["title"].(string),
			Id:    item["id"].(string),
		})
	}
	return &tac, nil
}

func (repo *GraphRepo) CountTitleAutoComplete(ctx context.Context, f dto2.PathFilter, kw string, amount *int64) error {
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
	res, err := repo.data.Neo.CypherQuery(ctx, cypher, map[string]any{
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

func (repo *GraphRepo) GetPathBetween(ctx context.Context, sourceId string, targetId string, f dto2.PathFilter) (*[]neo4j.Path, error) {
	if f.MaxPathDepth == 0 {
		f.MaxPathDepth = 6
	}
	cypher :=
		fmt.Sprintf(
			`MATCH (s {id: $sourceId}) 
		MATCH (t {id: $targetId}) 
		MATCH p = (s)-[r*..%d]->(t) 
		WHERE any(label IN labels(t) WHERE label IN $targetLabels) 
		AND all(rel in relationships(p) where type(rel) in $relTypes) 
		RETURN p;`, f.MaxPathDepth)
	param := map[string]any{
		"sourceId":     sourceId,
		"targetId":     targetId,
		"targetLabels": f.NodeLabels,
		"relTypes":     f.RelLabels,
	}
	res, err := repo.data.Neo.CypherQuery(ctx, cypher, param)
	if err != nil {
		return nil, err
	}
	var resp []neo4j.Path
	for _, path := range res {
		p, found := path.Get("p")
		if found {
			resp = append(resp, p.(neo4j.Path))
		}
	}
	return &resp, nil
}

func (repo *GraphRepo) GetPathExpand(ctx context.Context, sourceId string, depth uint32, limit uint32, f *dto2.PathFilter) (*[]neo4j.Path, error) {
	cypher :=
		fmt.Sprintf(
			`MATCH p = (s {id: $sourceId})-[*%d]->(t) 
		RETURN p LIMIT $limit;`, int(depth))
	param := map[string]any{
		"sourceId": sourceId,
		"limit":    int(limit),
	}
	res, err := repo.data.Neo.CypherQuery(ctx, cypher, param)
	if err != nil {
		return nil, err
	}
	var resp []neo4j.Path
	for _, path := range res {
		p, found := path.Get("p")
		if found {
			resp = append(resp, p.(neo4j.Path))
		}
	}
	return &resp, nil
}

func (repo *GraphRepo) GetPathBetweenByIds(ctx context.Context, sourceId string, targetIds []string, f *dto2.PathFilter) (*[]neo4j.Path, error) {
	if f.MaxPathDepth == 0 {
		f.MaxPathDepth = 6
	}
	cypher :=
		fmt.Sprintf(
			`MATCH (s {id: $sourceId}) 
		MATCH (t) WHERE t.id IN $targetIds 
		AND any(label IN labels(t) WHERE label IN $NodeLabels) 
		MATCH p = (s)-[r*..%d]->(t) 
		WHERE type(r) IN $relLabels 
		RETURN p;`, f.MaxPathDepth)

	res, err := repo.data.Neo.CypherQuery(ctx, cypher, map[string]any{
		"sourceId":   sourceId,
		"targetIds":  targetIds,
		"NodeLabels": f.NodeLabels,
		"$relLabels": f.RelLabels,
	})
	if err != nil {
		return nil, err
	}
	var resp []neo4j.Path
	for _, path := range res {
		p, found := path.Get("p")
		if found {
			resp = append(resp, p.(neo4j.Path))
		}
	}
	return &resp, nil
}

func (repo *GraphRepo) GetPathToChildren(ctx context.Context, sourceId string, p dto2.PaginationReq, f *dto2.PathFilter) (*[]neo4j.Path, int64, error) {
	var cypher string
	param := map[string]any{
		"sourceId": sourceId,
		"offset":   (p.PageNum - 1) * p.PageSize,
		"pageSize": p.PageSize,
	}
	if f != nil {
		cypher =
			`MATCH p = (s {id:$sourceId})-[r]->(c) 
			where any(label IN labels(c) WHERE label IN $nodeLabels) 
			AND type(r) in $relTypes 
			return p skip $offset limit $pageSize;`
		param["nodeLabels"] = f.NodeLabels
		param["relTypes"] = f.NodeLabels
	} else {
		cypher = `MATCH p = (s {id:$sourceId})-[r]->(c) return p skip $offset limit $pageSize;`
	}
	res, err := repo.data.Neo.CypherQuery(ctx, cypher, param)
	if err != nil {
		return nil, 0, err
	}
	var resp []neo4j.Path
	for _, path := range res {
		path, found := path.Get("p")
		if found {
			resp = append(resp, path.(neo4j.Path))
		}
	}
	cypherCount := `MATCH (p {id:$sourceId})-[r]->(c) return count(c) as total;`
	resTotal, err := repo.data.Neo.CypherQuery(ctx, cypherCount, map[string]any{
		"sourceId": sourceId,
	})
	if err != nil {
		return nil, 0, err
	}
	respTotal, _ := resTotal[0].Get("total")
	return &resp, respTotal.(int64), nil
}

func (repo *GraphRepo) GetPathToParent(ctx context.Context, sourceId string, p dto2.PaginationReq, f *dto2.PathFilter) (*[]neo4j.Path, int64, error) {
	cypher := `MATCH p = (s {id:$sourceId})<-[r]-(c) return p skip $offset limit $pageSize;`
	res, err := repo.data.Neo.CypherQuery(ctx, cypher, map[string]any{
		"sourceId": sourceId,
		"offset":   (p.PageNum - 1) * p.PageSize,
		"pageSize": p.PageSize,
	})
	if err != nil {
		return nil, 0, err
	}

	cypherCount := `MATCH (s {id:$sourceId})<-[r]-(c) return count(c) as total;`
	resCount, err := repo.data.Neo.CypherQuery(ctx, cypherCount, map[string]any{
		"sourceId": sourceId,
	})
	if err != nil {
		return nil, 0, err
	}
	respTotal, _ := resCount[0].Get("total")
	var resp []neo4j.Path
	for _, path := range res {
		path, found := path.Get("p")
		if found {
			resp = append(resp, path.(neo4j.Path))
		}
	}
	return &resp, respTotal.(int64), nil
}
