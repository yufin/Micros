package main

import (
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2"
	"go.temporal.io/sdk/worker"
	"os"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"
	"micros-worker/cmd"
	"micros-worker/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf    string
	hostName, _ = os.Hostname()
	id          string
)

func init() {
	Name = "worker.micros"
	flag.StringVar(&flagconf, "conf", "../../configs/config.yaml", "config path, eg: -conf config.yaml")
	flag.StringVar(&Version, "version", "v1.0.0", "version of this service")
	id = fmt.Sprintf("%s-%s-%s", Name, Version, hostName)
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func newWorkerApp(workers *[]worker.Worker, srvApp *kratos.App, log log.Logger) *cmd.WorkerApp {
	return cmd.NewWorkerApp(
		workers,
		srvApp,
		log,
	)
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
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
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
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}

	// start and wait for stop signal
	//if err := workerApp.Run(); err != nil {
	//	panic(err)
	//}
}

//tctl --ad 192.168.44.169:7233 workflow start \
//--workflow_type SyncNewContent \
//--taskqueue content-sync \
//--workflow_id sync_new_content_01 \

//tctl --ad 192.168.44.169:7233 workflow start \
//--workflow_type ReSyncDependencyData \
//--taskqueue content-sync \
//--workflow_id re_sync_dependency_data_01

//tctl --ad 10.0.232.121:7233 workflow start \
//--workflow_type ReSyncDependencyData \
//--taskqueue content-sync \
//--workflow_id re_sync_dependency_data_01
//-- input_json '{"f1": 1}'

//tctl --ad 192.168.44.169:7233 --namespace infra namespace register
//tctl --ad 10.0.232.121:7233 --namespace infra namespace register
