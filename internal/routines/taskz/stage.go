package taskz

import (
	"fmt"
	a "github.com/faelmori/golife/internal/routines/taskz/actions"
	j "github.com/faelmori/golife/internal/routines/taskz/jobs"
	"github.com/faelmori/golife/internal/routines/workers"
	t "github.com/faelmori/golife/internal/types"
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
)

// IStage represents the interface for a stage in the lifecycle.
type IStage[T any] interface {
	EventExists(event string) bool
	GetEvent(event string) func(interface{})
	GetEventFns() map[string]func(interface{})
	GetData() interface{}
	CanTransitionTo(stageID string) bool
	OnEnter(fn func()) IStage[T]
	OnExit(fn func()) IStage[T]
	OnEvent(event string, fn func(interface{})) IStage[T]
	AutoScale(size int) IStage[T]
	Dispatch(task func()) error
	Description() string
	Name() string
	ID() string
}

// Stage represents a stage in the lifecycle.
type Stage[T any] struct {
	StageID      string                      // Stage identifier
	StageName    string                      // Stage name
	StageObj     T                           // Stage object
	Type         string                      // Stage type
	Desc         string                      // Stage description
	Tags         []string                    // Stage tags
	Meta         t.Metadata                  // Stage metadata
	Data         t.Property[T]               // Stage data
	PossibleNext []string                    // Possible next stages
	PossiblePrev []string                    // Possible previous stages
	OnEnterFn    t.GenericChannelCallback[T] // Function to execute on entering the stage
	OnExitFn     t.ChangeListener[T]         // Function to execute on exiting the stage
	EventFns     t.ChangeListener[T]         // Event functions
	WorkerPool   t.IWorkerPool               // Worker pool for the stage
}

// ID returns the stage identifier.
func (s *Stage[T]) ID() string { return s.StageID }

// Name returns the stage name.
func (s *Stage[T]) Name() string { return s.StageName }

// Description returns the stage description.
func (s *Stage[T]) Description() string { return s.Desc }

// OnEnter sets the function to execute on entering the stage.
func (s *Stage[T]) OnEnter(fn func()) IStage[T] {
	//s.OnEnterFn = fn
	return s
}

// OnExit sets the function to execute on exiting the stage.
func (s *Stage[T]) OnExit(fn func()) IStage[T] {
	//s.OnExitFn = fn
	return s
}

// OnEvent sets the function to execute on a specific event.
func (s *Stage[T]) OnEvent(event string, fn func(any)) IStage[T] {
	//s.EventFns[event] = fn
	return s
}

// AutoScale sets the worker pool size for the stage.
func (s *Stage[T]) AutoScale(size int) IStage[T] {
	s.WorkerPool = workers.NewWorkerPool(size, l.GetLogger("GoLife"))
	if s.WorkerPool == nil {
		l.ErrorCtx(fmt.Sprintf("WorkerPool not initialized for stage %s", s.Name()), nil)
		return nil
	} else {
		l.InfoCtx(fmt.Sprintf("WorkerPool initialized for stage %s", s.Name()), nil)
		if scaleErr := s.WorkerPool.AddWorker(size, workers.NewWorker(s.WorkerPool.GetWorkerCount()+1, l.GetLogger("GoLife"))); scaleErr != nil {
			l.ErrorCtx(fmt.Sprintf("Error adding worker to stage %s: %v", s.Name(), scaleErr), nil)
			return nil
		}
		l.InfoCtx(fmt.Sprintf("Worker added to stage %s", s.Name()), nil)
	}
	return s
}

// Dispatch sends a task to the worker pool.
func (s *Stage[T]) Dispatch(task func()) error {
	if s.WorkerPool == nil {
		l.ErrorCtx(fmt.Sprintf("WorkerPool not initialized for stage %s", s.Name()), nil)
		return fmt.Errorf("WorkerPool not initialized for stage %s", s.Name())
	}
	act := a.NewAction("Teste do pezinho")
	cancelChan := make(chan struct{})
	doneChan := make(chan struct{})
	job := j.NewJob(act, cancelChan, doneChan, l.GetLogger("GoLife"))
	if sendJobErr := s.WorkerPool.SendToWorker(0, job); sendJobErr != nil {
		return sendJobErr
	}
	l.InfoCtx(fmt.Sprintf("Task dispatched to stage %s", s.Name()), nil)
	return nil
}

// CanTransitionTo checks if the stage can transition to another stage.
func (s *Stage[T]) CanTransitionTo(stageID string) bool {
	for _, next := range s.PossibleNext {
		if next == stageID {
			return true
		}
	}
	return false
}

// GetEvent returns the function for a specific event.
func (s *Stage[T]) GetEvent(event string) func(interface{}) {
	if fn, ok := s.EventFns[event]; ok {
		return fn
	}
	return nil
}

// GetEventFns returns all event functions.
func (s *Stage[T]) GetEventFns() map[string]func(interface{}) { return s.EventFns }

// GetData returns the stage data.
func (s *Stage[T]) GetData() interface{} { return s.Data }

// EventExists checks if an event exists in the stage.
func (s *Stage[T]) EventExists(event string) bool {
	if _, ok := s.EventFns[event]; ok {
		return true
	}
	return false
}

// NewStage creates a new stage with the given name, description, and type.
func NewStage[T any](name string, desc string, stageType string) IStage[T] {
	stg := Stage[T]{
		StageID:    uuid.New().String(),
		StageName:  name,
		Type:       stageType,
		Desc:       desc,
		EventFns:   make(map[string]func(interface{})),
		WorkerPool: nil,
	}
	l.GetLogger("GoLife").InfoCtx(fmt.Sprintf("New stage created: %s", name), nil)
	return &stg
}
