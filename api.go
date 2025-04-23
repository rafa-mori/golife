package golife

//import (
//	_ "github.com/faelmori/golife/internal/routines/taskz"
//)

//// Action is a generic type that represents an action of type T.
//type Action[T any] interface{ t.IAction[T] }
//
//// NewAction creates a new action with the given name, action and type.
//func NewAction[T any](identifier string, actionType string, data *T, ev func(T) error) Action[any] {
//	act := a.NewAction[T](
//		identifier,
//		actionType,
//		data,
//		ev,
//	)
//	if _, ok := act.(a.IAction[T]); ok {
//		return reflect.ValueOf(act).Interface().(Action[any])
//	} else {
//		return nil
//	}
//}
//
//// Job is a generic type that represents a job of type T.
//type Job[T any] interface{ j.IJob[T] }
//
//// NewJob creates a new job with the given name, action and type.
//func NewJob[T any](action Action[any], cancelChanel chan struct{}, doneChanel chan struct{}, Logger l.Logger, data *T) Job[any] {
//	if Logger == nil {
//		Logger = l.GetLogger("Job")
//	}
//	var act a.IAction[T]
//	_, ok := action.(a.IAction[T])
//	if action == nil || !ok {
//		act = &a.Action[T]{
//			ID:         uuid.New().String(),
//			Type:       "Job",
//			Errors:     make([]error, 0),
//			Properties: make(map[string]property.Property[T]),
//		}
//	} else {
//		act = reflect.ValueOf(action).Interface().(a.IAction[T])
//	}
//	if act != nil {
//		jj := j.NewJob[T](act, cancelChanel, doneChanel, Logger, nil)
//		return reflect.ValueOf(jj).Interface().(Job[any])
//	}
//	return nil
//}
//
//// Property is a generic type that represents a property of type T.
//type Property[T any] interface{ property.Property[T] }
//
//// NewProperty creates a new property with the given name and value.
//func NewProperty[T any](name string, value T) Property[any] {
//	// Create a new property with the given name, value and type.
//	// The type is inferred from the value passed to the function.
//	// Type don't will be explicitly exposed but will be used to create the property and validate the value.
//	return property.NewProperty[T](name, &value)
//}
//
//type Routine[T any] = r.IManagedGoroutine[T]
//
//type Channel[T any, N int] = s.IChannel[T, N]
//
//func NewChannel[T any, N int](name string, tp *T, buffers N) Channel[T, N] {
//	return c.NewChannel[T, N](name, tp, buffers)
//}
//
//type WorkerPool interface {
//	w.IWorkerPool
//	DispatchJob(job Job[any]) error
//	Close()
//}
//
//func NewWorkerPool(size int) WorkerPool {
//	return &internalWorkerPool{
//		size:    size,
//		jobs:    make(chan Job[any], size),
//		results: make(chan Result, size),
//	}
//}
//
//type Stage interface {
//	GetName() string
//	GetDescription() string
//	GetStageType() string
//}
//
//func NewStage(name, desc, stageType string) Stage {
//	return &Stage{
//		name:        name,
//		description: desc,
//		stageType:   stageType,
//	}
//}
//
//type LifeCycleManager interface {
//	RegisterEvent(stage string, event string, fn func(interface{})) error
//	RegisterStage(stage Stage) error
//	RegisterProcess(process ManagedProcess) error
//	Trigger(stage string, event string, data interface{})
//}
//
//func NewLifecycleManager(
//	processes map[string]ManagedProcess,
//	stages map[string]Stage,
//	sigChan chan os.Signal,
//	doneChan chan struct{},
//	events []Event,
//	eventsCh chan Event,
//) LifeCycleManager {
//	return &i.LifecycleManager{
//		processes: processes,
//		stages:    stages,
//		sigChan:   sigChan,
//		doneChan:  doneChan,
//		events:    events,
//		eventsCh:  eventsCh,
//	}
//}
//
//type Event interface {
//	Trigger(data interface{})
//	Register(callback func(data interface{}))
//}
//
//func NewEvent(eventFns map[string]func(interface{}), triggerCh chan interface{}) Event {
//	return &internalEvent{
//		eventFns:  eventFns,
//		triggerCh: triggerCh,
//	}
//}
//
////func NewWorkerPool(size int) WorkerPool { return w.NewWorkerPool(size) }
////
////type WorkerPoolGl interface{ i.IWorkerPoolGl }
////
////func NewWorkerPoolGl(size int) WorkerPoolGl { return i.NewWorkerPoolGL(size) }
////
////type LifeCycleManager interface{ i.LifeCycleManager }
////type ManagedProcessEvent interface{ e.IManagedProcessEvents }
////
////func NewLifecycleManager(processes map[string]p.IManagedProcess, stages map[string]t.IStage, sigChan chan os.Signal, doneChan chan struct{}, events []e.IManagedProcessEvents, eventsCh chan e.IManagedProcessEvents) LifeCycleManager {
////	return i.NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh)
////}
////func NewLifecycleMgrSig() (LifeCycleManager, error) {
////	processes := make(map[string]p.IManagedProcess)
////	stages := make(map[string]t.IStage)
////	sigChan := make(chan os.Signal, 2)
////	doneChan := make(chan struct{}, 2)
////	eventsCh := make(chan e.IManagedProcessEvents, 100)
////	events := make([]e.IManagedProcessEvents, 0)
////
////	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
////}
////func NewLifecycleMgrManual(processes map[string]p.IManagedProcess, stages map[string]t.IStage, sigChan chan os.Signal, doneChan chan struct{}, events []e.IManagedProcessEvents, eventsCh chan e.IManagedProcessEvents) (LifeCycleManager, error) {
////	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
////}
////func NewLifecycleMgrDec() (LifeCycleManager, error) {
////	processes := make(map[string]p.IManagedProcess)
////	stages := make(map[string]t.IStage)
////	sigChan := make(chan os.Signal, 2)
////	doneChan := make(chan struct{}, 2)
////	eventsCh := make(chan e.IManagedProcessEvents, 100)
////	events := make([]e.IManagedProcessEvents, 0)
////
////	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
////}
////func NewLifecycleMgrChan(sigChan chan os.Signal, doneChan chan struct{}, eventsCh chan e.IManagedProcessEvents) (LifeCycleManager, error) {
////	processes := make(map[string]p.IManagedProcess)
////	stages := make(map[string]t.IStage)
////	events := make([]e.IManagedProcessEvents, 0)
////	if sigChan == nil {
////		sigChan = make(chan os.Signal, 2)
////	}
////	if doneChan == nil {
////		doneChan = make(chan struct{}, 2)
////	}
////	if eventsCh == nil {
////		eventsCh = make(chan e.IManagedProcessEvents, 100)
////	}
////
////	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
////}
////
////type Stage = t.IStage
////
////func NewStage(name, desc, stageType string) Stage {
////	return t.NewStage(name, desc, stageType)
////}
////
////type ManagedProcess = p.IManagedProcess
////
////func NewManagedProcess(name string, command string, args []string, waitFor bool, customFn func() error) ManagedProcess {
////	return p.NewManagedProcess(name, command, args, waitFor, customFn)
////}
////
////type Event = e.IManagedProcessEvents
////
////func NewEvent(eventFns map[string]func(interface{}), triggerCh chan interface{}) Event {
////	return e.NewManagedProcessEvents(eventFns, triggerCh)
////}
////
////func Trigger(lc LifeCycleManager, stage string, event string, data interface{}) {
////	lc.Trigger(stage, event, data)
////}
////
////func RegisterEvent(lc LifeCycleManager, stage string, event string, fn func(interface{})) error {
////	return lc.RegisterEvent(stage, event, fn)
////}
////
////func RegisterStage(lc LifeCycleManager, stage Stage) error {
////	return lc.RegisterStage(stage)
////}
////
////func RegisterProcess(lc LifeCycleManager, process ManagedProcess) error {
////	return lc.RegisterProcess(process.GetName(), process.GetCommand(), process.GetArgs(), process.WillRestart(), process.GetFunction())
////}
