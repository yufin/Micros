//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"micros-worker/cmd"
	"micros-worker/internal/biz"
	"micros-worker/internal/conf"
	"micros-worker/internal/data"
	"micros-worker/internal/server"
	"micros-worker/internal/service"
	"micros-worker/internal/worker"
	"micros-worker/internal/workflow"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*cmd.WorkerApp, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		service.ProviderSet,
		workflow.ProviderSet,
		worker.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		newWorkerApp,
		newApp,
	))
}
