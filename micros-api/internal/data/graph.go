package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"math"
	"micros-api/internal/biz"
	"micros-api/internal/biz/dto"
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

func (repo *GraphRepo) GetNode(ctx context.Context, id string) (neo4j.Node, error) {
	cypher := "MATCH (n {id: $id}) RETURN n;"
	res, err := repo.data.Neo.CypherQuery(ctx, cypher, map[string]interface{}{"id": id})
	if err != nil {
		return neo4j.Node{}, err
	}
	if len(res) == 0 {
		return neo4j.Node{}, nil
	}
	n, found := res[0].Get("n")
	if !found {
		return neo4j.Node{}, errors.New("node with specified id not found")
	}
	var node dto.Node
	node.Gen(n.(neo4j.Node))
	return n.(neo4j.Node), nil
}

func (repo *GraphRepo) GetNodes(ctx context.Context, ids []string) (*[]neo4j.Node, error) {
	cypher := "MATCH (n) where n.id in $ids RETURN n;"
	res, err := repo.data.Neo.CypherQuery(ctx, cypher, map[string]interface{}{"ids": ids})
	if err != nil {
		return nil, err
	}
	nodes := make([]neo4j.Node, 0)
	for _, item := range res {
		n, _ := item.Get("n")
		nodes = append(nodes, n.(neo4j.Node))
	}
	return &nodes, nil
}

func (repo *GraphRepo) GetChildren(ctx context.Context, id string, f dto.PathFilter, p dto.PaginationReq) (*[]neo4j.Node, error) {
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
	nodes := make([]neo4j.Node, 0)
	for _, item := range res {
		item := item
		n, _ := item.Get("c")
		nodes = append(nodes, n.(neo4j.Node))
	}
	return &nodes, nil
}

func (repo *GraphRepo) CountChildren(ctx context.Context, id string, f dto.PathFilter, amount *int64) error {
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

func (repo *GraphRepo) GetTitleAutoComplete(ctx context.Context, f dto.PathFilter, p dto.PaginationReq, kw string) (*[]dto.TitleAutoCompleteRes, error) {
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

	tac := make([]dto.TitleAutoCompleteRes, 0)
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
		tac = append(tac, dto.TitleAutoCompleteRes{
			Title: item["title"].(string),
			Id:    item["id"].(string),
		})
	}
	return &tac, nil
}

func (repo *GraphRepo) CountTitleAutoComplete(ctx context.Context, f dto.PathFilter, kw string, amount *int64) error {
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

func (repo *GraphRepo) GetPathBetween(ctx context.Context, sourceId string, targetId string, f dto.PathFilter) (*[]neo4j.Path, error) {
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

func (repo *GraphRepo) GetPathExpand(ctx context.Context, sourceId string, depth uint32, limit uint32, f *dto.PathFilter) (*[]neo4j.Path, error) {
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

func (repo *GraphRepo) GetPathBetweenByIds(ctx context.Context, sourceId string, targetIds []string, f *dto.PathFilter) (*[]neo4j.Path, error) {
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

func (repo *GraphRepo) GetPathToChildren(ctx context.Context, sourceId string, p dto.PaginationReq, ScopeRelType []string) (*[]neo4j.Path, int64, error) {
	var cypher, cypherCount string
	paramCount := map[string]any{
		"sourceId": sourceId,
	}
	offset := math.Max(float64((p.PageNum-1)*p.PageSize), float64(0))
	param := map[string]any{
		"sourceId": sourceId,
		"offset":   int(offset),
		"pageSize": p.PageSize,
	}
	if len(ScopeRelType) > 0 {
		//cypher =
		//	`MATCH p = (s {id:$sourceId})-[r]->(c)
		//	where any(label IN labels(c) WHERE label IN $nodeLabels)
		//	AND type(r) in $relTypes
		//	return p skip $offset limit $pageSize;`
		//param["nodeLabels"] = f.NodeLabels
		//param["relTypes"] = f.NodeLabels
		cypher =
			`MATCH p = (s {id:$sourceId})-[r]->(c) 
			where type(r) in $relTypes 
			return p skip $offset limit $pageSize;`
		//param["nodeLabels"] = f.NodeLabels
		cypherCount =
			`MATCH (p {id:$sourceId})-[r]->(c)
			where type(r) in $relTypes 
			return count(c) as total;`

		param["relTypes"] = ScopeRelType
		paramCount["relTypes"] = ScopeRelType
	} else {
		cypher = `MATCH p = (s {id:$sourceId})-[r]->(c) return p skip $offset limit $pageSize;`
		cypherCount = `MATCH (p {id:$sourceId})-[r]->(c) return count(c) as total;`
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
	resTotal, err := repo.data.Neo.CypherQuery(ctx, cypherCount, paramCount)
	if err != nil {
		return nil, 0, err
	}
	respTotal, _ := resTotal[0].Get("total")
	return &resp, respTotal.(int64), nil
}

func (repo *GraphRepo) GetPathToParent(ctx context.Context, sourceId string, p dto.PaginationReq, f *dto.PathFilter) (*[]neo4j.Path, int64, error) {
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

// GetRelTypeAvailable returns the relationship types available to the specified node targetType(1: toChildren, 0:toParent)
func (repo *GraphRepo) GetRelTypeAvailable(ctx context.Context, id string, targetType int) ([]string, error) {

	var cypher string
	switch targetType {
	case 0:
		cypher = `MATCH (p)-[r]->(c {id: $id}) return collect(distinct type(r)) as relType;`
	case 1:
		cypher = `MATCH (p {id: $id})-[r]->(c) return collect(distinct type(r)) as relType;`
	default:
		return nil, errors.New("invalid target type")
	}
	res, err := repo.data.Neo.CypherQuery(ctx, cypher, map[string]any{
		"id": id,
	})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("empty result")
	}
	relType, found := res[0].Get("relType")
	if !found {
		return nil, errors.New("key not found")
	}

	relTypeStr := make([]string, 0)
	for _, s := range relType.([]interface{}) {
		relTypeStr = append(relTypeStr, s.(string))
	}

	return relTypeStr, nil
}
