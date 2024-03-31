package worker

import (
	"go.temporal.io/sdk/worker"
	"micros-worker/internal/workflow"
)

func NewWorkers(cs *workflow.ContentSyncWorkflow, in *workflow.CommonNoticeWorkflow) (*[]worker.Worker, error) {
	workersQueue := make([]worker.Worker, 0)
	orchestrateWorkers(&workersQueue, cs)
	orchestrateWorkers(&workersQueue, in)
	return &workersQueue, nil
}

type workerRegistrable interface {
	RegistryWorker() worker.Worker
}

func orchestrateWorkers(q *[]worker.Worker, reg workerRegistrable) {
	w := reg.RegistryWorker()
	*q = append(*q, w)
}
