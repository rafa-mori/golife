package internal

import (
	l "github.com/faelmori/logz"

	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type LifeCycleManager interface {
	Start() error
	Stop() error
	Restart() error
	Status() string

	RegisterProcess(name string, command string, args []string, restart bool, customFn func() error) error
	RegisterStage(stage IStage) error
	RegisterEvent(event, stage string, callback func(interface{})) error

	RemoveEvent(event, stage string) error
	StopEvents() error

	StartProcess(proc IManagedProcess) error
	StartAll() error
	StopAll() error

	Trigger(stage, event string, data interface{})
	DefineStage(name string) error
	GetCurrentStage() IStage
	GetStage(name string) IStage
	GetStages() []IStage
	UpdateStage(stage IStage) error

	IsStageAllowed(stage string) bool

	Send(stage string, msg string)
	Receive(stage string) interface{}
	ListenForSignals() error

	getStageIDByName(name string) string
}
type LifeCycle struct {
	processes map[string]IManagedProcess
	stages    map[string]IStage
	events    []IManagedProcessEvents

	currentStage string
	currentEvent string

	lastEvent string
	lastStage string

	sigChan  chan os.Signal
	doneChan chan struct{}

	eventsMu  sync.Mutex
	eventsCh  chan IManagedProcessEvents
	triggerCh chan interface{}

	mu sync.Mutex
}

func (lm *LifeCycle) Trigger(stageName, eventName string, data interface{}) {
	stage := lm.GetStage(stageName)
	if stage == nil {
		l.Error(fmt.Sprintf("Stage %s not found when triggering event %s", stageName, eventName), nil)
		return
	}

	// Execute the event callback
	callback := stage.GetEvent(eventName)
	if callback == nil {
		l.Error(fmt.Sprintf("Event %s not found in stage %s", eventName, stageName), nil)
		return
	}

	// Log the trigger
	l.Info(fmt.Sprintf("Triggering event %s in %s...", eventName, stageName), nil)

	// Execute the callback
	callback(data)

	// Send to the channel, if necessary
	if lm.eventsCh != nil {
		lm.eventsCh <- NewManagedProcessEvents(map[string]func(interface{}){eventName: callback}, lm.triggerCh)
	}
}
func (lm *LifeCycle) DefineStage(name string) error {
	if id := lm.getStageIDByName(name); id == "" {
		return fmt.Errorf("stage %s not found", name)
	} else {
		lm.currentStage = id
	}
	return nil
}
func (lm *LifeCycle) IsStageAllowed(stage string) bool {
	return lm.currentStage == stage
}
func (lm *LifeCycle) GetCurrentStage() IStage {
	if s, ok := lm.stages[lm.currentStage]; ok {
		return s
	}
	return nil
}
func (lm *LifeCycle) GetStage(name string) IStage {
	if id := lm.getStageIDByName(name); id != "" {
		if s, ok := lm.stages[id]; ok {
			return s
		}
	}
	return nil
}
func (lm *LifeCycle) GetStages() []IStage {
	stages := make([]IStage, 0, len(lm.stages))
	for _, stage := range lm.stages {
		stages = append(stages, stage)
	}
	return stages
}
func (lm *LifeCycle) UpdateStage(stage IStage) error {
	if id := lm.getStageIDByName(stage.Name()); id != "" {
		if s, ok := lm.stages[id]; ok {
			lm.stages[s.ID()] = stage
			return nil
		}
	}
	return fmt.Errorf("stage %s not found", stage.Name())
}

func (lm *LifeCycle) Send(stageName string, msg string) {
	if id := lm.getStageIDByName(stageName); id != "" {
		if s, ok := lm.stages[id]; ok {
			s.Dispatch(func() {
				sB := s.(*Stage)
				if sB.OnEnterFn != nil {
					s.OnEnter(sB.OnEnterFn)
				}
				if sB.OnExitFn != nil {
					s.OnExit(sB.OnExitFn)
				}
				if sB.EventFns != nil {
					for event, fn := range sB.EventFns {
						s.OnEvent(event, fn)
					}
				}
			})

			l.Info(msg, map[string]interface{}{
				"context": "GoLife",
				"stage":   stageName,
				"message": msg,
			})
		}
	} else {
		l.Error(fmt.Sprintf("Stage %s not found", stageName), map[string]interface{}{
			"context": "GoLife",
			"stage":   stageName,
		})
	}
}
func (lm *LifeCycle) Receive(stageName string) interface{} {
	fmt.Printf("Receiving message from %s...\n", stageName)
	if id := lm.getStageIDByName(stageName); id != "" {
		if s, ok := lm.stages[id]; ok {
			fmt.Printf("Message received from %s!\n", stageName)
			sB := s.(*Stage)
			return sB.Data
		} else {
			fmt.Printf("No message received from %s!\n", stageName)
			return nil
		}
	} else {
		fmt.Printf("Stage %s not found!\n", stageName)
		return nil
	}
}
func (lm *LifeCycle) Start() error {
	l.Info("Starting processes...", map[string]interface{}{
		"context":  "GoLife",
		"showData": false,
	})
	for _, proc := range lm.processes {
		if err := proc.Start(); err != nil {
			l.Error(fmt.Sprintf("Error starting process %s: %v", proc.String(), err), map[string]interface{}{
				"context":  "GoLife",
				"process":  proc.String(),
				"showData": true,
			})
			return err
		}
	}
	l.Info(fmt.Sprintf("Processes started successfully!"), map[string]interface{}{
		"context":   "GoLife",
		"processes": len(lm.processes),
		"showData":  false,
	})
	return nil
}
func (lm *LifeCycle) Stop() error {
	l.Info("Stopping processes...", map[string]interface{}{
		"context":  "GoLife",
		"showData": false,
	})
	for _, proc := range lm.processes {
		l.Info(fmt.Sprintf("Stopping process %s...", proc.String()), map[string]interface{}{
			"context":  "GoLife",
			"process":  proc.String(),
			"showData": false,
		})
		if err := proc.Stop(); err != nil {
			l.Error(fmt.Sprintf("Error stopping process %s: %v", proc.String(), err), map[string]interface{}{
				"context":  "GoLife",
				"process":  proc.String(),
				"showData": true,
			})
			return err
		}
	}
	l.Info(fmt.Sprintf("Processes stopped successfully!"), map[string]interface{}{
		"context":   "GoLife",
		"processes": len(lm.processes),
		"showData":  false,
	})
	return nil
}
func (lm *LifeCycle) Restart() error {
	for _, proc := range lm.processes {
		if err := proc.Restart(); err != nil {
			return err
		}
	}
	return nil
}
func (lm *LifeCycle) Status() string {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	l.Info("Checking process status...", map[string]interface{}{"context": "GoLife", "showData": false})
	var status string
	for name, proc := range lm.processes {
		l.Info(fmt.Sprintf("Process %s (PID %d) is running: %t", name, proc.Pid(), proc.IsRunning()), map[string]interface{}{
			"context":  "GoLife",
			"process":  name,
			"pid":      proc.Pid(),
			"running":  proc.IsRunning(),
			"showData": false,
		})
		status += fmt.Sprintf("Process %s (PID %d) is running: %t\n", name, proc.Pid(), proc.IsRunning())
	}
	l.Info("Process status checked successfully!", map[string]interface{}{"context": "GoLife", "showData": false})
	return status
}

func (lm *LifeCycle) RegisterProcess(name string, command string, args []string, restart bool, customFn func() error) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	l.Info(fmt.Sprintf("Registering process %s...", name), map[string]interface{}{"context": "GoLife", "process": name, "showData": false})
	lm.processes[name] = NewManagedProcess(name, command, args, restart, customFn)

	l.Info(fmt.Sprintf("Process %s registered successfully!", name), map[string]interface{}{"context": "GoLife", "process": name, "showData": false})
	return nil
}
func (lm *LifeCycle) RegisterStage(stage IStage) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if id := lm.getStageIDByName(stage.Name()); id != "" {
		l.Error(fmt.Sprintf("Stage %s already registered", stage.Name()), map[string]interface{}{"context": "GoLife", "stage": stage.Name(), "showData": false})
		return nil
	}

	l.Info(fmt.Sprintf("Registering stage %s...", stage.Name()), map[string]interface{}{"context": "GoLife", "stage": stage.Name(), "showData": false})
	lm.stages[stage.ID()] = stage
	l.Info(fmt.Sprintf("Stage %s registered successfully!", stage.Name()), map[string]interface{}{"context": "GoLife", "stage": stage.Name(), "showData": false})

	return nil
}
func (lm *LifeCycle) RegisterEvent(event, stageName string, callback func(interface{})) error {
	lm.eventsMu.Lock()
	defer lm.eventsMu.Unlock()

	// Validate if the stage exists
	stage := lm.GetStage(stageName)
	if stage == nil {
		return fmt.Errorf("stage %s not found when registering event %s", stageName, event)
	}

	// Add the event to the stage
	stage.OnEvent(event, callback)

	// Log success
	l.Info(fmt.Sprintf("Event %s registered in %s successfully!", event, stageName), map[string]interface{}{
		"context": "GoLife",
		"event":   event,
		"stage":   stageName,
	})

	return nil
}

func (lm *LifeCycle) RemoveEvent(event, stageName string) error {
	l.Info(fmt.Sprintf("Removing event %s from %s...", event, stageName), map[string]interface{}{"context": "GoLife", "event": event, "stage": stageName, "showData": false})
	for i, e := range lm.events {
		if e.Event() == event {
			lm.events = append(lm.events[:i], lm.events[i+1:]...)
			l.Info(fmt.Sprintf("Event %s removed from %s successfully!", event, stageName), map[string]interface{}{"context": "GoLife", "event": event, "stage": stageName, "showData": false})
			return nil
		}
	}
	l.Error(fmt.Sprintf("Event %s not found in %s", event, stageName), map[string]interface{}{"context": "GoLife", "event": event, "stage": stageName, "showData": false})
	return fmt.Errorf("event %s not found in %s", event, stageName)
}
func (lm *LifeCycle) StopEvents() error {
	l.Info("Stopping events...", map[string]interface{}{"context": "GoLife", "showData": false})
	for _, event := range lm.events {
		stopErr := event.StopAll()
		if stopErr != nil {
			return stopErr
		}
	}
	l.Info("Events stopped successfully!", map[string]interface{}{"context": "GoLife", "showData": false})
	return nil
}
func (lm *LifeCycle) StartAll() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	for name, proc := range lm.processes {
		l.Info(fmt.Sprintf("Starting %s...", name), map[string]interface{}{"context": "GoLife", "process": name, "showData": false})
		if err := lm.StartProcess(proc); err != nil {
			l.Error(fmt.Sprintf("Error starting %s: %v", name, err), map[string]interface{}{"context": "GoLife", "process": name, "showData": true})
			return err
		}
	}
	l.Info(fmt.Sprintf("%b Processes started successfully!", len(lm.processes)), map[string]interface{}{"context": "GoLife", "processes": len(lm.processes), "showData": false})
	return nil
}
func (lm *LifeCycle) StartProcess(proc IManagedProcess) error {
	if err := proc.Start(); err != nil {
		l.Error(fmt.Sprintf("Error starting %s: %v", proc.String(), err), map[string]interface{}{"context": "GoLife", "process": proc.String(), "showData": true})
		return err
	}
	l.Info(fmt.Sprintf("%s started successfully!", proc.String()), map[string]interface{}{"context": "GoLife", "process": proc.String(), "showData": false})
	return nil
}
func (lm *LifeCycle) StopAll() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	for name, proc := range lm.processes {
		if err := proc.Stop(); err != nil {
			l.Error(fmt.Sprintf("Error stopping %s: %v", name, err), map[string]interface{}{"context": "GoLife", "process": name, "showData": true})
			return err
		} else {
			l.Info(fmt.Sprintf("%s stopped successfully!", name), map[string]interface{}{"context": "GoLife", "process": name, "showData": false})
			delete(lm.processes, name)
		}
	}
	l.Info(fmt.Sprintf("%b Processes stopped successfully!", len(lm.processes)), map[string]interface{}{"context": "GoLife", "processes": len(lm.processes), "showData": false})
	return nil
}
func (lm *LifeCycle) ListenForSignals() error {
	select {
	case ev := <-lm.eventsCh:
		if ev != nil {
			for _, event := range lm.events {
				event.RegisterEvent(event.Event(), func(data interface{}) {
					event.Trigger(data.(string), event.Event(), data)
				})
			}
		}
	case <-lm.doneChan:
		if len(lm.processes) > 0 {
			startAllErr := lm.StartAll()
			if startAllErr != nil {
				return startAllErr
			}
		}
	case <-lm.sigChan:
		stopAllErr := lm.StopAll()
		if stopAllErr != nil {
			return stopAllErr
		}
		close(lm.doneChan)
	}

	return nil
}
func (lm *LifeCycle) getStageIDByName(name string) string {
	for id, stage := range lm.stages {
		if stage.Name() == name {
			return id
		}
	}
	return ""
}

func NewLifecycleManager(processes map[string]IManagedProcess, stages map[string]IStage, sigChan chan os.Signal, doneChan chan struct{}, events []IManagedProcessEvents, eventsCh chan IManagedProcessEvents) LifeCycleManager {
	stg := make(map[string]IStage)
	if stages != nil {
		stg = stages
	}
	if sigChan == nil {
		sigChan = make(chan os.Signal, 2)
	}
	if doneChan == nil {
		doneChan = make(chan struct{}, 2)
	}
	if events == nil {
		events = make([]IManagedProcessEvents, 0)
	}
	if eventsCh == nil {
		eventsCh = make(chan IManagedProcessEvents, 100)
	}

	mgr := LifeCycle{
		processes: processes,
		stages:    stg,
		sigChan:   sigChan,
		doneChan:  doneChan,
		events:    events,
		eventsMu:  sync.Mutex{},
		eventsCh:  eventsCh,
		mu:        sync.Mutex{},
	}

	signal.Notify(mgr.sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := mgr.ListenForSignals()
		if err != nil {
			l.Error(fmt.Sprintf("Error listening for signals: %v", err), map[string]interface{}{"context": "GoLife", "showData": true})
		}
	}()

	return &mgr
}
func NewLifecycleMgrManual(processes map[string]IManagedProcess, stages map[string]IStage, sigChan chan os.Signal, doneChan chan struct{}, events []IManagedProcessEvents, eventsCh chan IManagedProcessEvents) (LifeCycleManager, error) {
	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}
func NewLifecycleMgrSig() (LifeCycleManager, error) {
	processes := make(map[string]IManagedProcess)
	stages := make(map[string]IStage)
	sigChan := make(chan os.Signal, 2)
	doneChan := make(chan struct{}, 2)
	events := make([]IManagedProcessEvents, 0)
	eventsCh := make(chan IManagedProcessEvents, 100)

	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}
func NewLifecycleMgrDec() (LifeCycleManager, error) {
	processes := make(map[string]IManagedProcess)
	stages := make(map[string]IStage)
	sigChan := make(chan os.Signal, 2)
	doneChan := make(chan struct{}, 2)
	events := make([]IManagedProcessEvents, 0)
	eventsCh := make(chan IManagedProcessEvents, 100)

	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}
