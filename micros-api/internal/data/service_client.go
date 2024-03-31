package data

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/pkg/errors"
	dwdataV2 "micros-api/api/dwdata/v2"
	dwdataV3 "micros-api/api/dwdata/v3"
	pipelineV1 "micros-api/api/pipeline/v1"
	"micros-api/internal/conf"
	"os"
	"time"
)

type DwDataClients struct {
	dwDataV2 dwdataV2.DwdataServiceClient
	dwDataV3 dwdataV3.DwdataServiceClient
}

func NewDwdataServiceClient(c *conf.Data, reg *consul.Registry) (*DwDataClients, func(), error) {
	cert, err := tls.LoadX509KeyPair(
		c.Tls.FolderPath+"/client_cert.pem",
		c.Tls.FolderPath+"/client_key.pem",
	)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	certBytes, err := os.ReadFile(c.Tls.FolderPath + "/ca_cert.pem")
	if err != nil {
		panic("Unable to read cert.pem")
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse root certificate")
	}
	tlsConfig := &tls.Config{
		RootCAs:      clientCertPool,
		Certificates: []tls.Certificate{cert},
		//InsecureSkipVerify:    false,
		//VerifyPeerCertificate: nil,
	}

	conn, err := grpc.Dial(
		context.Background(),
		grpc.WithEndpoint("discovery:///dw.micros"),
		grpc.WithMiddleware(
			metadata.Client(),
		),
		grpc.WithDiscovery(reg),
		grpc.WithTLSConfig(tlsConfig),
	)

	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	cleanup := func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("grpc conn closing failed")
		}
	}
	return &DwDataClients{
		dwDataV2: dwdataV2.NewDwdataServiceClient(conn),
		dwDataV3: dwdataV3.NewDwdataServiceClient(conn),
	}, cleanup, nil
}

func NewPipelineServiceClient(c *conf.Data, reg *consul.Registry) (pipelineV1.PipelineServiceClient, func(), error) {
	cert, err := tls.LoadX509KeyPair(c.Tls.FolderPath+"/client_cert.pem", c.Tls.FolderPath+"/client_key.pem")
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	certBytes, err := os.ReadFile(c.Tls.FolderPath + "/ca_cert.pem")
	if err != nil {
		panic("Unable to read cert.pem")
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse root certificate")
	}
	tlsConfig := &tls.Config{
		RootCAs:      clientCertPool,
		Certificates: []tls.Certificate{cert},
		//InsecureSkipVerify:    false,
		//VerifyPeerCertificate: nil,
	}

	conn, err := grpc.Dial(
		context.Background(),
		grpc.WithEndpoint(c.Service.PipelineUri),
		grpc.WithTimeout(time.Duration(10)*time.Second),
		grpc.WithMiddleware(
			metadata.Client(),
		),
		grpc.WithTLSConfig(tlsConfig),
	)
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
