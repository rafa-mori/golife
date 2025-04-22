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
	SetMeta(key string, value f.EventMetadata) IStage[T]
	GetMeta(key string) (*f.EventMetadata, bool)
	SetTags(tags []string) IStage[T]
	GetTags() []string
	SetWorkerPool(pool t.IWorkerPool) IStage[T]
	GetWorkerPool() t.IWorkerPool
	GetStageID() uuid.UUID
	GetStageName() string
	GetStageType() string
	GetStageDesc() string
	GetStageData() *T
	SetPossibleNext(stages []string) IStage[T]
	GetPossibleNext() []string
	SetPossiblePrev(stages []string) IStage[T]
	GetPossiblePrev() []string
	GetWorkerCount() int
	Dispatch(task func()) error
	Description() string
	Name() string
	ID() uuid.UUID
}

type BaseStage[T any] struct {
	StageID   uuid.UUID // Stage identifier
	StageName string    // Stage name
	Data      *T        // Stage data
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
			Data:      &defaultData,
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
		workerPool: w.NewWorkerPool(1, l.GetLogger("GoLife")),
	}
	return stg
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
		l.Error(fmt.Sprintf("WorkerPool not initialized for stage %s", s.Name()), nil)
		return s
	}

	wkrId := s.workerPool.GetWorkerCount() + 1
	wkr := w.NewWorker(wkrId, s.workerPool.Logger())
	if err := s.workerPool.AddWorker(wkrId, wkr); err != nil {
		l.Error(fmt.Sprintf("Error adding worker %d to pool: %v", wkrId, err), nil)
		return s
	}

	return s
}

// Dispatch sends a task to the worker pool.
func (s *Stage[T]) Dispatch(task func()) error {
	if s.workerPool == nil {
		l.Error(fmt.Sprintf("WorkerPool not initialized for stage %s", s.Name()), nil)
		return fmt.Errorf("WorkerPool not initialized for stage %s", s.Name())
	}
	//s.workerPool.Tasks <- task
	l.Info(fmt.Sprintf("Task dispatched to stage %s", s.Name()), nil)
	return nil
}

// SetMeta sets the metadata for the stage.
func (s *Stage[T]) SetMeta(key string, value f.EventMetadata) IStage[T] {
	s.Meta[key] = value
	return s
}

// GetMeta returns the metadata for the stage.
func (s *Stage[T]) GetMeta(key string) (*f.EventMetadata, bool) {
	if meta, ok := s.Meta[key]; ok {
		return &meta, true
	}
	return nil, false
}

// SetTags sets the tags for the stage.
func (s *Stage[T]) SetTags(tags []string) IStage[T] {
	s.Tags = tags
	return s
}

// GetTags returns the tags for the stage.
func (s *Stage[T]) GetTags() []string { return s.Tags }

// SetWorkerPool sets the worker pool for the stage.
func (s *Stage[T]) SetWorkerPool(pool t.IWorkerPool) IStage[T] {
	s.workerPool = pool
	return s
}

// GetWorkerPool returns the worker pool for the stage.
func (s *Stage[T]) GetWorkerPool() t.IWorkerPool {
	if s.workerPool == nil {
		l.Error(fmt.Sprintf("WorkerPool not initialized for stage %s", s.Name()), nil)
	}
	return s.workerPool
}

// GetStageID returns the stage identifier.
func (s *Stage[T]) GetStageID() uuid.UUID { return s.StageID }

// GetStageName returns the stage name.
func (s *Stage[T]) GetStageName() string { return s.StageName }

// GetStageType returns the stage type.
func (s *Stage[T]) GetStageType() string { return s.Type }

// GetStageDesc returns the stage description.
func (s *Stage[T]) GetStageDesc() string { return s.Desc }

// GetStageData returns the stage data.
func (s *Stage[T]) GetStageData() *T { return s.Data }

// GetPossibleNext returns the possible next stages.
func (s *Stage[T]) GetPossibleNext() []string { return s.PossibleNext }

// SetPossibleNext sets the possible next stages.
func (s *Stage[T]) SetPossibleNext(stages []string) IStage[T] {
	s.PossibleNext = stages
	return s
}

// GetPossiblePrev returns the possible previous stages.
func (s *Stage[T]) GetPossiblePrev() []string { return s.PossiblePrev }

// SetPossiblePrev sets the possible previous stages.
func (s *Stage[T]) SetPossiblePrev(stages []string) IStage[T] {
	s.PossiblePrev = stages
	return s
}

// CanTransitionTo checks if the stage can transition to another stage.
func (s *Stage[T]) CanTransitionTo(stageName string) bool {
	for _, next := range s.PossibleNext {
		if next == stageName {
			return true
		}
	}
	for _, prev := range s.PossiblePrev {
		if prev == stageName {
			return true
		}
	}
	return false
}

// GetWorkerCount returns the number of workers in the pool.
func (s *Stage[T]) GetWorkerCount() int {
	if s.workerPool == nil {
		l.Error(fmt.Sprintf("WorkerPool not initialized for stage %s", s.Name()), nil)
		return 0
	}
	return s.workerPool.GetWorkerCount()
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
