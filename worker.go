package golife

import (
	"sync"
)

type WorkerPool struct {
	tasks chan func()
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

func NewWorkerPool(size int) *WorkerPool {
	pool := &WorkerPool{tasks: make(chan func(), size)}
	for i := 0; i < size; i++ {
		go pool.worker()
	}
	return pool
}
