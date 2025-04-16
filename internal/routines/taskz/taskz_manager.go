package taskz

import (
	"context"
	"errors"
	"sync"
)

// ExecTaskzHandler defines the signature of a function that executes a serial task.
// Receives a context and a pointer to the task request.
// Returns an error, if any.
type ExecTaskzHandler[T any] func(ctx context.Context, taskReq *T) error

// ExecParallelHandler defines the signature of a function that executes a parallel task.
// Receives a context, a pointer to the task request, and a read/write mutex.
// Returns an error, if any.
type ExecParallelHandler[T any] func(ctx context.Context, taskReq *T, mu *sync.RWMutex) error

// Taskz is a struct that manages the execution of tasks.
type Taskz[T any] struct {
	ctx context.Context // Execution context
	mu  sync.RWMutex    // Mutex for concurrent access

	taskReq T // Task request

	tasks         []ExecTaskzHandler[T]    // Serial tasks list
	parallelTasks []ExecParallelHandler[T] // Parallel tasks list
}

// NewTaskz creates a new instance of Taskz.
// Receives a context and a task request.
// Returns a pointer to Taskz.
func NewTaskz[T any](ctx context.Context, taskReq T) *Taskz[T] {
	return &Taskz[T]{
		ctx:     ctx,
		taskReq: taskReq,
		mu:      sync.RWMutex{},
	}
}

// Then adds a serial task to the task list.
// Receives a function that executes the serial task.
// Returns a pointer to Taskz.
func (s *Taskz[T]) Then(taskExec ExecTaskzHandler[T]) *Taskz[T] {
	s.tasks = append(s.tasks, taskExec)
	return s
}

// Parallel adds a parallel task to the parallel task list.
// Receives a function that executes the parallel task.
// Returns a pointer to Taskz.
func (s *Taskz[T]) Parallel(parallelExec ExecParallelHandler[T]) *Taskz[T] {
	s.parallelTasks = append(s.parallelTasks, parallelExec)
	return s
}

// Result executes all tasks and returns the final result.
// Returns the task request and an error, if any.
func (s *Taskz[T]) Result() (T, error) {
	err := s.ExecSerial()
	err = errors.Join(err, s.ExecParallel())
	return s.taskReq, err
}

// ExecSerial executes the serial tasks in order.
// Returns an error, if any.
func (s *Taskz[T]) ExecSerial() error {
	for _, task := range s.tasks {
		err := task(s.ctx, &s.taskReq)
		if err != nil {
			return err
		}
	}
	return nil
}

// ExecParallel executes the parallel tasks simultaneously.
// Returns an error, if any.
func (s *Taskz[T]) ExecParallel() error {
	errChan := make(chan error)
	wg := sync.WaitGroup{}
	var err error
	for _, task := range s.parallelTasks {
		wg.Add(1)
		go func(ctx context.Context, mu *sync.RWMutex, taskReq *T) {
			defer wg.Done()
			err := task(ctx, taskReq, mu)
			errChan <- err
		}(s.ctx, &s.mu, &s.taskReq)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect errors from errChan
	for goErr := range errChan {
		if goErr != nil {
			err = errors.Join(err, goErr)
		}
	}
	return err
}
