package cmd

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go.temporal.io/sdk/worker"
	"sync"
)

type WorkerApp struct {
	workers *[]worker.Worker
	srvApp  *kratos.App
	ctx     context.Context
	mu      sync.Mutex
	cancel  func()
	logger  log.Logger
}

func NewWorkerApp(workers *[]worker.Worker, srv *kratos.App, log log.Logger) *WorkerApp {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerApp{
		srvApp:  srv,
		ctx:     ctx,
		cancel:  cancel,
		workers: workers,
		logger:  log,
	}
}

func (a *WorkerApp) Run() error {
	var errCh chan error

	defer a.srvApp.Stop()
	go func(ch *chan error) {
		if err := a.srvApp.Run(); err != nil {
			fmt.Printf("srv app run error: %+v", err)
			*ch <- err
			return
		}
	}(&errCh)

	for _, wi := range *a.workers {
		go func(w worker.Worker) {
			err := w.Run(worker.InterruptCh())
			if err != nil {
				errCh <- err
				return
			}
		}(wi)
	}
	select {
	case err := <-errCh:
		return err
	case <-worker.InterruptCh():
		return nil
	}
}
