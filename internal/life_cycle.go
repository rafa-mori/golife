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
	RegisterProcess(name string, command string, args []string, restart bool)
	StartAll()
	StopAll()
	Trigger(stage, event string, data interface{})
	DefineStage(name string) *Stage
	Send(stage string, msg interface{})
	Receive(stage string) interface{}
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
	s := &Stage{
		ID:         name,
		Name:       name,
		Type:       "default",
		Desc:       "Default stage",
		Tags:       []string{},
		Meta:       make(map[string]interface{}),
		Data:       nil,
		WorkerPool: NewWorkerPool(1),
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
	lm.StartAll()
	return nil
}

func (lm *gWebLifeCycle) Stop() error {
	lm.StopAll()
	return nil
}

func (lm *gWebLifeCycle) Restart() error {
	lm.StopAll()
	lm.StartAll()
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

func (lm *gWebLifeCycle) RegisterProcess(name string, command string, args []string, restart bool) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.processes[name] = &ManagedProcess{
		Cmd:  exec.Command(command, args...),
		Name: name,
	}
}

func (lm *gWebLifeCycle) StartAll() {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	for name, proc := range lm.processes {
		if err := lm.startProcess(proc); err != nil {
			fmt.Printf("Erro ao iniciar %s: %v\n", name, err)
		}
	}
}

func (lm *gWebLifeCycle) startProcess(proc *ManagedProcess) error {
	if err := proc.Start(); err != nil {
		return err
	}
	return nil
}

func (lm *gWebLifeCycle) StopAll() {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	for name, proc := range lm.processes {
		if err := proc.Stop(); err != nil {
			fmt.Printf("Erro ao encerrar %s: %v\n", name, err)
		}
	}
}

func (lm *gWebLifeCycle) listenForSignals() {
	select {
	case ev := <-lm.eventsCh:
		lm.eventsMu.Lock()
		defer lm.eventsMu.Unlock()
		lm.Trigger(ev.Stage, ev.Ev, ev.Data)
	case <-lm.doneChan:
		if len(lm.processes) > 0 {
			lm.StartAll()
		}
	case <-lm.sigChan:
		lm.StopAll()
		close(lm.doneChan)
	}
}

func NewLifecycleManager() GWebLifeCycleManager {
	mgr := &gWebLifeCycle{
		processes: make(map[string]*ManagedProcess),
		sigChan:   make(chan os.Signal, 1),
		doneChan:  make(chan struct{}),
	}

	signal.Notify(mgr.sigChan, syscall.SIGINT, syscall.SIGTERM)
	go mgr.listenForSignals()

	return mgr
}
