package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"micros-graph/internal/biz"
	"micros-graph/internal/biz/dto"
)

type GraphNetRepo struct {
	data *Data
	log  *log.Helper
}

func NewGraphNetRepo(data *Data, logger log.Logger) biz.GraphNetRepo {
	return &GraphNetRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *GraphNetRepo) GetChildren(ctx context.Context, sourceId int64, relTypeScope []string, labelScope []string, p *dto.PaginationReq) (*dto.Net, int64, error) {
	session, err := repo.data.nebulaDb.getSession()
	if err != nil {
		return nil, 0, err
	}
	defer session.Release()

	//resultSet, err := session.Execute("USE demo_sns; match (n) return n limit 2;")
	//resultSet, err := session.Execute("USE demo_sns; match p=(n)-[r]->(c) return p limit 3")
	res, err := session.ExecuteJson("USE demo_sns; match p=(n)-[r]->(c) return p limit 3")
	if err != nil {
		return nil, 0, err
	}
	fmt.Println(res)
	//if !resultSet.IsSucceed() {
	//	return nil, 0, errors.New("query failed")
	//}
	//
	//colNames := resultSet.GetColNames()
	//fmt.Printf("column names: %s\n", strings.Join(colNames, ", "))
	//
	//// Get a row from resultSet
	//record, err := resultSet.GetRowValuesByIndex(0)
	//if err != nil {
	//	log.Error(err.Error())
	//}
	//// Print whole row
	//fmt.Printf("row elements: %s\n", record.String())
	//// Get a value in the row by column index
	//valueWrapper, err := record.GetValueByIndex(0)
	//if err != nil {
	//	log.Error(err.Error())
	//}
	//// Get type of the value
	//fmt.Printf("valueWrapper type: %s \n", valueWrapper.GetType())
	//// Check if valueWrapper is a string type
	//if valueWrapper.IsString() {
	//	// Convert valueWrapper to a string value
	//	v1Str, err := valueWrapper.AsString()
	//	if err != nil {
	//		log.Error(err.Error())
	//	}
	//	fmt.Printf("Result of ValueWrapper.AsString(): %s\n", v1Str)
	//}
	//// Print ValueWrapper using String()
	//fmt.Printf("Print using ValueWrapper.String(): %s", valueWrapper.String())

	return nil, 0, nil
}
