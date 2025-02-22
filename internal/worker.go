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
	wg    sync.WaitGroup
}

func (wp *WorkerPool) worker() {
	for task := range wp.tasks {
		task()
	}
}

func (wp *WorkerPool) Submit(task func()) {
	wp.tasks <- task
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

func NewWorkerPool(size int) IWorkerPool {
	pool := WorkerPool{tasks: make(chan func(), size)}
	for i := 0; i < size; i++ {
		go pool.worker()
	}
	return &pool
}
