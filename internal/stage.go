package internal

import (
	"fmt"
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
)

// IStage represents the interface for a stage in the lifecycle.
type IStage interface {
	EventExists(event string) bool
	GetEvent(event string) func(interface{})
	GetEventFns() map[string]func(interface{})
	GetData() interface{}
	CanTransitionTo(stageID string) bool
	OnEnter(fn func()) IStage
	OnExit(fn func()) IStage
	OnEvent(event string, fn func(interface{})) IStage
	AutoScale(size int) IStage
	Dispatch(task func()) error
	Description() string
	Name() string
	ID() string
}

// Stage represents a stage in the lifecycle.
type Stage struct {
	StageID      string                       // Stage identifier
	StageName    string                       // Stage name
	Type         string                       // Stage type
	Desc         string                       // Stage description
	Tags         []string                     // Stage tags
	Meta         map[string]interface{}       // Stage metadata
	Data         interface{}                  // Stage data
	PossibleNext []string                     // Possible next stages
	PossiblePrev []string                     // Possible previous stages
	OnEnterFn    func()                       // Function to execute on entering the stage
	OnExitFn     func()                       // Function to execute on exiting the stage
	EventFns     map[string]func(interface{}) // Event functions
	WorkerPool   *WorkerPool                  // Worker pool for the stage
}

// ID returns the stage identifier.
func (s *Stage) ID() string { return s.StageID }

// Name returns the stage name.
func (s *Stage) Name() string { return s.StageName }

// Description returns the stage description.
func (s *Stage) Description() string { return s.Desc }

// OnEnter sets the function to execute on entering the stage.
func (s *Stage) OnEnter(fn func()) IStage {
	s.OnEnterFn = fn
	return s
}

// OnExit sets the function to execute on exiting the stage.
func (s *Stage) OnExit(fn func()) IStage {
	s.OnExitFn = fn
	return s
}

// OnEvent sets the function to execute on a specific event.
func (s *Stage) OnEvent(event string, fn func(interface{})) IStage {
	s.EventFns[event] = fn
	return s
}

// AutoScale sets the worker pool size for the stage.
func (s *Stage) AutoScale(size int) IStage {
	s.WorkerPool = NewWorkerPool(size).(*WorkerPool)
	s.WorkerPool.Wg.Add(size)
	return s
}

// Dispatch sends a task to the worker pool.
func (s *Stage) Dispatch(task func()) error {
	if s.WorkerPool == nil {
		l.ErrorCtx(fmt.Sprintf("WorkerPool not initialized for stage %s", s.Name()), nil)
		return fmt.Errorf("WorkerPool not initialized for stage %s", s.Name())
	}
	s.WorkerPool.Tasks <- task
	l.InfoCtx(fmt.Sprintf("Task dispatched to stage %s", s.Name()), nil)
	return nil
}

// CanTransitionTo checks if the stage can transition to another stage.
func (s *Stage) CanTransitionTo(stageID string) bool {
	for _, next := range s.PossibleNext {
		if next == stageID {
			return true
		}
	}
	return false
}

// GetEvent returns the function for a specific event.
func (s *Stage) GetEvent(event string) func(interface{}) {
	if fn, ok := s.EventFns[event]; ok {
		return fn
	}
	return nil
}

// GetEventFns returns all event functions.
func (s *Stage) GetEventFns() map[string]func(interface{}) { return s.EventFns }

// GetData returns the stage data.
func (s *Stage) GetData() interface{} { return s.Data }

// EventExists checks if an event exists in the stage.
func (s *Stage) EventExists(event string) bool {
	if _, ok := s.EventFns[event]; ok {
		return true
	}
	return false
}

// NewStage creates a new stage with the given name, description, and type.
func NewStage(name, desc, stageType string) IStage {
	stg := Stage{
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
