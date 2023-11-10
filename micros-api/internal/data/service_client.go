package data

import (
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	dwdataV2 "micros-api/api/dwdata/v2"
	pipelineV1 "micros-api/api/pipeline/v1"
	"micros-api/internal/conf"
)

func NewDwdataServiceClient(c *conf.Data) (dwdataV2.DwdataServiceClient, func(), error) {
	conn, err := grpc.Dial(
		c.Service.DwdataUri,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	cleanup := func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("grpc conn closing failed")
		}
	}
	return dwdataV2.NewDwdataServiceClient(conn), cleanup, nil
}

func NewPipelineServiceClient(c *conf.Data) (pipelineV1.PipelineServiceClient, func(), error) {
	conn, err := grpc.Dial(
		c.Service.PipelineUri,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	cleanup := func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("grpc conn closing failed")
		}
	}
	return pipelineV1.NewPipelineServiceClient(conn), cleanup, nil
}
