package internal

type IStage interface {
	OnEnter(fn func()) *Stage
	OnExit(fn func()) *Stage
	OnEvent(event string, fn func(interface{})) *Stage
	AutoScale(size int) *Stage
	Dispatch(task func())
}

type Stage struct {
	// Stage identifiers
	ID   string
	Name string
	Type string
	Desc string
	Tags []string
	Meta map[string]interface{}
	Data interface{}

	// Next and Prev stages
	PossibleNext []string
	PossiblePrev []string

	// Embed methods
	OnEnterFn func()
	OnExitFn  func()
	EventFns  map[string]func(interface{})

	// Worker pool
	WorkerPool *WorkerPool

	// Internals
	// ...
}

func (s *Stage) OnEnter(fn func()) *Stage {
	s.OnEnterFn = fn
	return s
}
func (s *Stage) OnExit(fn func()) *Stage {
	s.OnExitFn = fn
	return s
}
func (s *Stage) OnEvent(event string, fn func(interface{})) *Stage {
	s.EventFns[event] = fn
	return s
}
func (s *Stage) AutoScale(size int) *Stage {
	s.WorkerPool = NewWorkerPool(size).(*WorkerPool)
	s.WorkerPool.Wg.Add(size)
	return s
}
func (s *Stage) Dispatch(task func()) {
	if s.WorkerPool != nil {
		s.WorkerPool.Tasks <- task
	}
}
