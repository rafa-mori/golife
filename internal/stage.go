package internal

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/faelmori/logz"
)

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
type Stage struct {
	// Stage identifiers
	StageID      string
	StageName    string
	Type         string
	Desc         string
	Tags         []string
	Meta         map[string]interface{}
	Data         interface{}
	PossibleNext []string
	PossiblePrev []string
	OnEnterFn    func()
	OnExitFn     func()
	EventFns     map[string]func(interface{})
	WorkerPool   *WorkerPool
}

func (s *Stage) ID() string          { return s.StageID }
func (s *Stage) Name() string        { return s.StageName }
func (s *Stage) Description() string { return s.Desc }
func (s *Stage) OnEnter(fn func()) IStage {
	s.OnEnterFn = fn
	return s
}
func (s *Stage) OnExit(fn func()) IStage {
	s.OnExitFn = fn
	return s
}
func (s *Stage) OnEvent(event string, fn func(interface{})) IStage {
	s.EventFns[event] = fn
	return s
}
func (s *Stage) AutoScale(size int) IStage {
	s.WorkerPool = NewWorkerPool(size).(*WorkerPool)
	s.WorkerPool.Wg.Add(size)
	return s
}
func (s *Stage) Dispatch(task func()) error {
	if s.WorkerPool == nil {
		logz.Error(fmt.Sprintf("WorkerPool não inicializado para o estágio %s", s.Name()), nil)
		return fmt.Errorf("WorkerPool não inicializado para o estágio %s", s.Name())
	}
	s.WorkerPool.Tasks <- task
	logz.Info(fmt.Sprintf("Tarefa despachada para o estágio %s", s.Name()), nil)
	return nil
}
func (s *Stage) CanTransitionTo(stageID string) bool {
	for _, next := range s.PossibleNext {
		if next == stageID {
			return true
		}
	}
	return false
}
func (s *Stage) GetEvent(event string) func(interface{}) {
	if fn, ok := s.EventFns[event]; ok {
		return fn
	}
	return nil
}
func (s *Stage) GetEventFns() map[string]func(interface{}) { return s.EventFns }
func (s *Stage) GetData() interface{}                      { return s.Data }
func (s *Stage) EventExists(event string) bool {
	if _, ok := s.EventFns[event]; ok {
		return true
	}
	return false
}

func NewStage(name, desc, stageType string) IStage {
	stg := Stage{
		StageID:    uuid.New().String(),
		StageName:  name,
		Type:       stageType,
		Desc:       desc,
		EventFns:   make(map[string]func(interface{})),
		WorkerPool: nil,
	}
	logz.Info(fmt.Sprintf("Novo estágio criado: %s", name), nil)
	return &stg
}
