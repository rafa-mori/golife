package workers

import (
	"context"
)

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

/*type WorkerPool struct {
	IWorkerPoolBase
	IWorkerPool

	ctx      context.Context
	cancel   context.CancelFunc
	Tasks    chan func()
	Wg       sync.WaitGroup
	Shutdown chan struct{}
}*/

/*func (wp *WorkerPool) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.MuDone():
			return
		case task := <-wp.Tasks:
			task()
			wp.Wg.MuDone()
		}
	}
}
func (wp *WorkerPool) Submit(task func()) error {
	select {
	case wp.Tasks <- task:
		wp.Wg.MuAdd(1)
		l.InfoCtx("Task submitted", nil)
		return nil
	default:
		return fmt.Errorf("buffer de tarefas cheio")
	}
}
func (wp *WorkerPool) MuWait() {
	wp.Wg.MuWait()
	l.InfoCtx("All tasks completed", nil)
}
func (wp *WorkerPool) Scale(size int) {
	for i := 0; i < size; i++ {
		// Create a new context for each worker based on the parent context from the pool
		ctxWorker, cancelWorker := context.WithCancel(wp.ctx)
		// MuAdd a task to the wait group
		wp.Wg.MuAdd(1)
		// Define the task for the worker
		wp.Tasks <- func() {
			defer wp.Wg.MuDone()
			defer cancelWorker()
			// Preciso fazer algo logo.. hgahahahahaha
		}
		// MuAdd a task to the wait group
		go wp.worker(ctxWorker)
	}
}
func (wp *WorkerPool) workerWithError(errChan chan error) {
	for task := range wp.Tasks {
		defer wp.Wg.MuDone()
		if err := wp.safeExecute(task); err != nil {
			errChan <- err
		}
	}
}
func (wp *WorkerPool) safeExecute(task func()) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Erro recuperado: %v\n", r)
		}
	}()
	task()
	return nil
}

func NewWorkerPool(size int) IWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := WorkerPool{Tasks: make(chan func(), size), ctx: ctx, cancel: cancel, Shutdown: make(chan struct{})}

	for i := 0; i < size; i++ {
		// Create a new context for each worker based on the parent context from the pool
		ctxWorker, cancelWorker := context.WithCancel(ctx)

		// MuAdd a task to the wait group
		pool.Wg.MuAdd(1)

		// Define the task for the worker
		pool.Tasks <- func() {
			defer pool.Wg.MuDone()
			defer cancelWorker()
			// Preciso fazer algo logo Pode ser aquiiiii...

			// OOOOOU.... logo ali em cima.. kkkkk
		}

		// Start the worker
		go pool.worker(ctxWorker)
	}
	l.InfoCtx("Worker pool created", map[string]interface{}{"size": size})
	return &pool
}
*/
