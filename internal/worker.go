package internal

import (
	l "github.com/faelmori/logz"

	"sync"
)

type IWorkerPool interface {
	Submit(task func())
	Wait()
	worker()
}

type WorkerPool struct {
	Tasks chan func()
	Wg    sync.WaitGroup
}

func (wp *WorkerPool) worker() {
	for task := range wp.Tasks {
		task()
		wp.Wg.Done()
		l.Info("Task completed", nil)
	}
}
func (wp *WorkerPool) Submit(task func()) {
	wp.Wg.Add(1)
	wp.Tasks <- task
	l.Info("Task submitted", nil)
}
func (wp *WorkerPool) Wait() {
	wp.Wg.Wait()
	l.Info("All tasks completed", nil)
}

func NewWorkerPool(size int) IWorkerPool {
	pool := WorkerPool{Tasks: make(chan func(), size)}
	for i := 0; i < size; i++ {
		go pool.worker()
	}
	l.Info("Worker pool created", map[string]interface{}{"size": size})
	return &pool
}
