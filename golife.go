package golife

import (
	i "github.com/faelmori/golife/internal"
	"github.com/faelmori/golife/internal/process"
	"github.com/faelmori/golife/internal/routines/taskz"
	"github.com/faelmori/golife/internal/routines/taskz/events"
	"github.com/faelmori/golife/internal/routines/workers"

	"os"
)

type WorkerPool = workers.IWorkerPool

func NewWorkerPool(size int) WorkerPool { return workers.NewWorkerPool(size) }

type WorkerPoolGl interface{ i.IWorkerPoolGl }

func NewWorkerPoolGl(size int) WorkerPoolGl { return i.NewWorkerPoolGL(size) }

type LifeCycleManager interface{ i.LifeCycleManager }
type ManagedProcessEvent interface{ events.IManagedProcessEvents }

func NewLifecycleManager(processes map[string]process.IManagedProcess, stages map[string]taskz.IStage, sigChan chan os.Signal, doneChan chan struct{}, events []events.IManagedProcessEvents, eventsCh chan events.IManagedProcessEvents) LifeCycleManager {
	return i.NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh)
}
func NewLifecycleMgrSig() (LifeCycleManager, error) {
	processes := make(map[string]process.IManagedProcess)
	stages := make(map[string]taskz.IStage)
	sigChan := make(chan os.Signal, 2)
	doneChan := make(chan struct{}, 2)
	eventsCh := make(chan events.IManagedProcessEvents, 100)
	events := make([]events.IManagedProcessEvents, 0)

	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}
func NewLifecycleMgrManual(processes map[string]process.IManagedProcess, stages map[string]taskz.IStage, sigChan chan os.Signal, doneChan chan struct{}, events []events.IManagedProcessEvents, eventsCh chan events.IManagedProcessEvents) (LifeCycleManager, error) {
	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}
func NewLifecycleMgrDec() (LifeCycleManager, error) {
	processes := make(map[string]process.IManagedProcess)
	stages := make(map[string]taskz.IStage)
	sigChan := make(chan os.Signal, 2)
	doneChan := make(chan struct{}, 2)
	eventsCh := make(chan events.IManagedProcessEvents, 100)
	events := make([]events.IManagedProcessEvents, 0)

	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}
func NewLifecycleMgrChan(sigChan chan os.Signal, doneChan chan struct{}, eventsCh chan events.IManagedProcessEvents) (LifeCycleManager, error) {
	processes := make(map[string]process.IManagedProcess)
	stages := make(map[string]taskz.IStage)
	events := make([]events.IManagedProcessEvents, 0)
	if sigChan == nil {
		sigChan = make(chan os.Signal, 2)
	}
	if doneChan == nil {
		doneChan = make(chan struct{}, 2)
	}
	if eventsCh == nil {
		eventsCh = make(chan events.IManagedProcessEvents, 100)
	}

	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}

type Stage = taskz.IStage

func NewStage(name, desc, stageType string) Stage {
	return taskz.NewStage(name, desc, stageType)
}

type ManagedProcess = process.IManagedProcess

func NewManagedProcess(name string, command string, args []string, waitFor bool, customFn func() error) ManagedProcess {
	return process.NewManagedProcess(name, command, args, waitFor, customFn)
}

type Event = events.IManagedProcessEvents

func NewEvent(eventFns map[string]func(interface{}), triggerCh chan interface{}) Event {
	return events.NewManagedProcessEvents(eventFns, triggerCh)
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

//func abc() {
//	var ch = make(chan int)
//
//	// Tentando enviar uma string (vai dar erro de tipo)
//	ch <- "isso não é um número!"
//}
