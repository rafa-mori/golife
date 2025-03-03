package internal

import (
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
	}
}
func (wp *WorkerPool) Submit(task func()) {
	wp.Wg.Add(1)
	wp.Tasks <- task
}
func (wp *WorkerPool) Wait() {
	wp.Wg.Wait()
}

func NewWorkerPool(size int) IWorkerPool {
	pool := WorkerPool{Tasks: make(chan func(), size)}
	for i := 0; i < size; i++ {
		go pool.worker()
	}
	return &pool
}
