package interfaces

import "context"

type IWorkerPoolBase interface {
	worker(ctx context.Context)
	safeExecute(task func()) error
	workerWithError(errChan chan error)
}

type IWorkerPool interface {
	worker(ctx context.Context)
	Submit(task func()) error
	Wait()
	Scale(size int)
	workerWithError(errChan chan error)
	safeExecute(task func()) error
}
