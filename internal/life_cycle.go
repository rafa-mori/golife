package internal

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

type GWebLifeCycleManager interface {
	Start() error
	Stop() error
	Restart() error
	Status() string
	RegisterProcess(name string, command string, args []string, restart bool) error
	RegisterEvent(event, stage string) error
	StartProcess(proc *ManagedProcess) error
	StartAll() error
	StopAll() error
	Trigger(stage, event string, data interface{})
	DefineStage(name string) *Stage
	Send(stage string, msg interface{})
	Receive(stage string) interface{}
	ListenForSignals() error
}

type gWebLifeCycle struct {
	processes map[string]*ManagedProcess
	stages    map[string]*Stage
	sigChan   chan os.Signal
	doneChan  chan struct{}
	events    *ManagedProcessEvents
	eventsMu  sync.Mutex
	eventsCh  chan *ManagedProcessEvents
	mu        sync.Mutex
}

func (lm *gWebLifeCycle) Trigger(stage, event string, data interface{}) {
	if s, ok := lm.stages[stage]; ok {
		if fn, ok := s.EventFns[event]; ok {
			fn(data)
		}
	}
}
func (lm *gWebLifeCycle) DefineStage(name string) *Stage {
	w := NewWorkerPool(1).(*WorkerPool)
	s := &Stage{
		ID:         name,
		Name:       name,
		Type:       "default",
		Desc:       "Default stage",
		Tags:       []string{},
		Meta:       make(map[string]interface{}),
		Data:       nil,
		WorkerPool: w,
		EventFns:   make(map[string]func(interface{})),
	}

	lm.stages[name] = s

	return s
}
func (lm *gWebLifeCycle) Send(stage string, msg interface{}) {
	if s, ok := lm.stages[stage]; ok {
		s.Dispatch(func() {
			if s.OnEnterFn != nil {
				s.OnEnterFn()
			}
			if s.OnExitFn != nil {
				s.OnExitFn()
			}
		})
	}
}
func (lm *gWebLifeCycle) Receive(stage string) interface{} {
	if s, ok := lm.stages[stage]; ok {
		return s.Data
	} else {
		return nil
	}
}
func (lm *gWebLifeCycle) Start() error {
	err := lm.StartAll()
	if err != nil {
		return err
	}
	return nil
}
func (lm *gWebLifeCycle) Stop() error {
	err := lm.StopAll()
	if err != nil {
		return err
	}
	return nil
}
func (lm *gWebLifeCycle) Restart() error {
	stopAllErr := lm.StopAll()
	if stopAllErr != nil {
		return stopAllErr
	}
	startAllErr := lm.StartAll()
	if startAllErr != nil {
		return startAllErr
	}
	return nil
}
func (lm *gWebLifeCycle) Status() string {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	var status string
	for name, proc := range lm.processes {
		status += fmt.Sprintf("Processo %s (PID %d) está rodando: %t\n", name, proc.Pid(), proc.IsRunning())
	}
	return status
}
func (lm *gWebLifeCycle) RegisterProcess(name string, command string, args []string, restart bool) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.processes[name] = &ManagedProcess{
		Cmd:  exec.Command(command, args...),
		Name: name,
	}

	return nil
}
func (lm *gWebLifeCycle) RegisterEvent(event, stage string) error {
	lm.eventsMu.Lock()
	defer lm.eventsMu.Unlock()

	lm.eventsCh <- &ManagedProcessEvents{
		Ev:    event,
		Stage: stage,
		Data:  nil,
	}

	return nil
}
func (lm *gWebLifeCycle) StartAll() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	for name, proc := range lm.processes {
		if err := lm.StartProcess(proc); err != nil {
			fmt.Printf("Erro ao iniciar %s: %v\n", name, err)
		}
	}
	return nil
}
func (lm *gWebLifeCycle) StartProcess(proc *ManagedProcess) error {
	if err := proc.Start(); err != nil {
		return err
	}
	return nil
}
func (lm *gWebLifeCycle) StopAll() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	for name, proc := range lm.processes {
		if err := proc.Stop(); err != nil {
			return err
		} else {
			delete(lm.processes, name)
		}
	}
	return nil
}
func (lm *gWebLifeCycle) ListenForSignals() error {
	select {
	case ev := <-lm.eventsCh:
		lm.eventsMu.Lock()
		defer lm.eventsMu.Unlock()
		lm.Trigger(ev.Stage, ev.Ev, ev.Data)
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

func NewLifecycleManager() GWebLifeCycleManager {
	mgr := &gWebLifeCycle{
		processes: make(map[string]*ManagedProcess),
		sigChan:   make(chan os.Signal, 1),
		doneChan:  make(chan struct{}),
	}

	signal.Notify(mgr.sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := mgr.ListenForSignals()
		if err != nil {
			fmt.Println("Erro ao ouvir sinais:", err)
		}
	}()

	return mgr
}
