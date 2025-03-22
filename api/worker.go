package api

import i "github.com/faelmori/golife/internal"

type WorkerPool = i.IWorkerPool

func NewWorkerPool(size int) WorkerPool {
	return i.NewWorkerPool(size)
}
