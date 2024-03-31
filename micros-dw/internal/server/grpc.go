package server

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	dwV2 "micros-dw/api/dwdata/v2"
	dwV3 "micros-dw/api/dwdata/v3"
	"micros-dw/internal/conf"
	v2 "micros-dw/internal/service/v2"
	v3 "micros-dw/internal/service/v3"
	"os"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, cfg *conf.Data, dwsV2 *v2.DwdataServiceServicer, dwsV3 *v3.DwdataServiceServicer, logger log.Logger) *grpc.Server {
	//https://www.jianshu.com/p/4cf92c5a386d
	cert, err := tls.LoadX509KeyPair(
		cfg.Tls.FolderPath+"/server_cert.pem",
		cfg.Tls.FolderPath+"/server_key.pem",
	)
	if err != nil {
		panic(err)
	}
	certBytes, err := os.ReadFile(cfg.Tls.FolderPath + "/client_ca_cert.pem")
	if err != nil {
		panic("Unable to read cert.pem")
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse root certificate")
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
		//VerifyPeerCertificate: nil,
	}
	//_ = &tls.Config{
	//	Certificates:          []tls.Certificate{cert},
	//	ClientAuth:            tls.RequireAndVerifyClientCert,
	//	ClientCAs:             clientCertPool,
	//	VerifyPeerCertificate: nil,
	//}

	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			metadata.Server(),
			logging.Server(logger),
		),
		grpc.TLSConfig(tlsConfig),
	}

	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	dwV2.RegisterDwdataServiceServer(srv, dwsV2)
	dwV3.RegisterDwdataServiceServer(srv, dwsV3)
	return srv
}
