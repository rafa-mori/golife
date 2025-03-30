package golife

import (
	i "github.com/faelmori/golife/internal"
	"os"
)

type LifeCycleManager = i.LifeCycleManager
type ManagedProcessEvent = i.IManagedProcessEvents

func NewLifecycleManager(processes map[string]i.IManagedProcess, stages map[string]i.IStage, sigChan chan os.Signal, doneChan chan struct{}, events []i.IManagedProcessEvents, eventsCh chan i.IManagedProcessEvents) LifeCycleManager {
	return i.NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh)
}
func NewLifecycleMgrSig() (LifeCycleManager, error) {
	processes := make(map[string]i.IManagedProcess)
	stages := make(map[string]i.IStage)
	sigChan := make(chan os.Signal, 2)
	doneChan := make(chan struct{}, 2)
	eventsCh := make(chan i.IManagedProcessEvents, 100)
	events := make([]i.IManagedProcessEvents, 0)

	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}
func NewLifecycleMgrManual(processes map[string]i.IManagedProcess, stages map[string]i.IStage, sigChan chan os.Signal, doneChan chan struct{}, events []i.IManagedProcessEvents, eventsCh chan i.IManagedProcessEvents) (LifeCycleManager, error) {
	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}
func NewLifecycleMgrDec() (LifeCycleManager, error) {
	processes := make(map[string]i.IManagedProcess)
	stages := make(map[string]i.IStage)
	sigChan := make(chan os.Signal, 2)
	doneChan := make(chan struct{}, 2)
	eventsCh := make(chan i.IManagedProcessEvents, 100)
	events := make([]i.IManagedProcessEvents, 0)

	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}
func NewLifecycleMgrChan(sigChan chan os.Signal, doneChan chan struct{}, eventsCh chan i.IManagedProcessEvents) (LifeCycleManager, error) {
	processes := make(map[string]i.IManagedProcess)
	stages := make(map[string]i.IStage)
	events := make([]i.IManagedProcessEvents, 0)
	if sigChan == nil {
		sigChan = make(chan os.Signal, 2)
	}
	if doneChan == nil {
		doneChan = make(chan struct{}, 2)
	}
	if eventsCh == nil {
		eventsCh = make(chan i.IManagedProcessEvents, 100)
	}

	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}

type Stage = i.IStage

func NewStage(name, desc, stageType string) Stage {
	return i.NewStage(name, desc, stageType)
}

type WorkerPool = i.IWorkerPool

func NewWorkerPool(size int) WorkerPool {
	return i.NewWorkerPool(size)
}

type ManagedProcess = i.IManagedProcess

func NewManagedProcess(name string, command string, args []string, waitFor bool, customFn func() error) ManagedProcess {
	return i.NewManagedProcess(name, command, args, waitFor, customFn)
}

type Event = i.IManagedProcessEvents

func NewEvent(eventFns map[string]func(interface{}), triggerCh chan interface{}) Event {
	return i.NewManagedProcessEvents(eventFns, triggerCh)
}

func Trigger(lc LifeCycleManager, stage string, event string, data interface{}) {
	lc.Trigger(stage, event, data)
}

func RegisterEvent(lc LifeCycleManager, stage string, event string, fn func(interface{})) error {
	return lc.RegisterEvent(stage, event, fn)
}

func RegisterStage(lc LifeCycleManager, stage Stage) error {
	return lc.RegisterStage(stage)
}

func RegisterProcess(lc LifeCycleManager, process ManagedProcess) error {
	return lc.RegisterProcess(process.GetName(), process.GetCommand(), process.GetArgs(), process.WillRestart(), process.GetCustomFunc())
}
