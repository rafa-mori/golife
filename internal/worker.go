package internal

import (
	"github.com/rafa-mori/logz"
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
		logz.Info("Task completed", nil)
	}
}
func (wp *WorkerPool) Submit(task func()) {
	wp.Wg.Add(1)
	wp.Tasks <- task
	logz.Info("Task submitted", nil)
}
func (wp *WorkerPool) Wait() {
	wp.Wg.Wait()
	logz.Info("All tasks completed", nil)
}

func NewWorkerPool(size int) IWorkerPool {
	pool := WorkerPool{Tasks: make(chan func(), size)}
	for i := 0; i < size; i++ {
		go pool.worker()
	}
	logz.Info("Worker pool created", map[string]interface{}{"size": size})
	return &pool
}
