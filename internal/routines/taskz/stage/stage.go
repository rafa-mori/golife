package stage

import (
	"fmt"
	f "github.com/faelmori/golife/internal/property"
	w "github.com/faelmori/golife/internal/routines/workers"
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
	ID() uuid.UUID
}

type BaseStage[T any] struct {
	StageID   uuid.UUID // Stage identifier
	StageName string    // Stage name
	Data      T         // Stage data
	Type      string    // Stage type
	Desc      string    // Stage description
}

// Stage represents a stage in the lifecycle.
type Stage[T any] struct {
	IStage[T]                                 // Interface for the stage
	BaseStage[T]                              // Base stage properties
	Meta         map[string]f.EventMetadata   // Stage metadata
	Tags         []string                     // Stage tags
	PossibleNext []string                     // Possible next stages
	PossiblePrev []string                     // Possible previous stages
	OnEnterFn    func()                       // Function to execute on entering the stage
	OnExitFn     func()                       // Function to execute on exiting the stage
	EventFns     map[string]func(interface{}) // Event functions
	workerPool   t.IWorkerPool                // Worker pool for the stage
}

// ID returns the stage identifier.
func (s *Stage[T]) ID() uuid.UUID { return s.StageID }

// Name returns the stage name.
func (s *Stage[T]) Name() string { return s.StageName }

// Description returns the stage description.
func (s *Stage[T]) Description() string { return s.Desc }

// OnEnter sets the function to execute on entering the stage.
func (s *Stage[T]) OnEnter(fn func()) IStage[T] {
	s.OnEnterFn = fn
	return s
}

// OnExit sets the function to execute on exiting the stage.
func (s *Stage[T]) OnExit(fn func()) IStage[T] {
	s.OnExitFn = fn
	return s
}

// OnEvent sets the function to execute on a specific event.
func (s *Stage[T]) OnEvent(event string, fn func(interface{})) IStage[T] {
	s.EventFns[event] = fn
	return s
}

// AutoScale sets the worker pool size for the stage.
func (s *Stage[T]) AutoScale(size int) IStage[T] {
	if s.workerPool == nil {
		l.ErrorCtx(fmt.Sprintf("WorkerPool not initialized for stage %s", s.Name()), nil)
		return s
	}

	wkrId := s.workerPool.GetWorkerCount() + 1
	wkr := w.NewWorker(wkrId, s.workerPool.Logger())
	if err := s.workerPool.AddWorker(wkrId, wkr); err != nil {
		l.ErrorCtx(fmt.Sprintf("Error adding worker %d to pool: %v", wkrId, err), nil)
		return s
	}

	return s
}

// Dispatch sends a task to the worker pool.
func (s *Stage[T]) Dispatch(task func()) error {
	if s.workerPool == nil {
		l.ErrorCtx(fmt.Sprintf("WorkerPool not initialized for stage %s", s.Name()), nil)
		return fmt.Errorf("WorkerPool not initialized for stage %s", s.Name())
	}
	//s.workerPool.Tasks <- task
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
func NewStage[T any](name, desc, stageType string, data *T) IStage[T] {
	var defaultData T
	// Initialize the default data for the stage
	if data == nil {
		defaultData = *new(T)
	} else {
		defaultData = *data
	}
	// Initialize the worker pool
	stg := &Stage[T]{
		BaseStage: BaseStage[T]{
			StageID:   uuid.New(),
			StageName: name,
			Desc:      desc,
			Data:      defaultData,
			Type:      stageType,
		},
		Meta:         make(map[string]f.EventMetadata), // Stage metadata
		Tags:         []string{},                       // Stage tags
		PossibleNext: []string{},                       // Possible next stages
		PossiblePrev: []string{},                       // Possible previous stages
		OnEnterFn: func() {

		}, // Function to execute on entering the stage
		OnExitFn: func() {

		}, // Function to execute on exiting the stage
		EventFns:   make(map[string]func(interface{})), // Event functions
		workerPool: nil,
	}
	l.GetLogger("GoLife").InfoCtx(fmt.Sprintf("New stage created: %s", name), nil)
	return stg
}
