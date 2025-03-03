package internal

type IStage interface {
	OnEnter(fn func()) IStage
	OnExit(fn func()) IStage
	OnEvent(event string, fn func(interface{})) IStage
	AutoScale(size int) IStage
	Dispatch(task func())
	Name() string
}

type Stage struct {
	// Stage identifiers
	ID           string
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

func (s *Stage) Name() string {
	return s.StageName
}
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
func (s *Stage) Dispatch(task func()) {
	if s.WorkerPool != nil {
		s.WorkerPool.Tasks <- task
	}
}

func NewStage(id, name, desc, stageType string) IStage {
	stg := Stage{
		ID:         id,
		StageName:  name,
		Type:       stageType,
		Desc:       desc,
		EventFns:   make(map[string]func(interface{})),
		WorkerPool: nil,
	}
	return &stg
}
