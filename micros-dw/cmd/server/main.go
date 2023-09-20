package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2/config/file"
	clientv3 "go.etcd.io/etcd/client/v3"
	grpc2 "google.golang.org/grpc"
	"os"
	"time"

	"micros-dw/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagConf is the config flag.
	//flagConf string
	flagConfCenter         string
	flagConfCenterRootPath string
	flagConf               string
	id, _                  = os.Hostname()
)

func init() {
	flag.StringVar(&flagConf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&flagConfCenter, "conf-center", "http://192.168.44.169:2389", "conf center address, eg: -conf-center http://localhost:2389")
	flag.StringVar(&flagConfCenterRootPath, "conf-center-root", "/micros/dw", "conf center root path, eg: -conf-center-root /dw-config")
}

func newApp(logger log.Logger, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			//hs,
		),
	)
}

func newEtcd(endpoint string) (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 3 * time.Second,
		DialOptions: []grpc2.DialOption{grpc2.WithBlock()},
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	//etcdCli, err := newEtcd(flagConfCenter)
	//defer etcdCli.Close()
	//if err != nil {
	//	log.Error(err)
	//	panic(err)
	//}
	//source, err := cfg.New(etcdCli, cfg.WithPath(flagConfCenterRootPath), cfg.WithPrefix(true))
	//if err != nil {
	//	panic(err)
	//}

	c := config.New(
		config.WithSource(
			file.NewSource(flagConf),
			//source,
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer app.Stop()
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}

}
