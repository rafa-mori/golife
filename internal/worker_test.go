package internal

import (
	"sync"
	"testing"
)

func TestSubmit(t *testing.T) {
	pool := NewWorkerPool(5).(*WorkerPool)
	var wg sync.WaitGroup
	wg.Add(1)

	pool.Submit(func() {
		defer wg.Done()
		// Simulate task
	})

	wg.Wait()
	if len(pool.Tasks) != 0 {
		t.Errorf("Expected task to be submitted and completed")
	}
}

func TestWait(t *testing.T) {
	pool := NewWorkerPool(5).(*WorkerPool)
	var wg sync.WaitGroup
	wg.Add(1)

	pool.Submit(func() {
		defer wg.Done()
		// Simulate task
	})

	pool.Wait()
	if len(pool.Tasks) != 0 {
		t.Errorf("Expected all tasks to be completed")
	}
}
