package data

import (
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	dwdataV2 "micros-api/api/dwdata/v2"
	"micros-api/internal/conf"
)

func NewGrpcConn(c *conf.Data) (*grpc.ClientConn, func(), error) {
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
	return conn, cleanup, nil
}

func NewDwdataServiceClient(conn *grpc.ClientConn) dwdataV2.DwdataServiceClient {
	return dwdataV2.NewDwdataServiceClient(conn)
}
