package internal

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type GWebLifeCycleManager interface {
	Start() error
	Stop() error
	Restart() error
	Status() string
	RegisterProcess(name string, command string, args []string, restart, wait bool) error
	ListenForSignals()
}

type gWebLifeCycle struct {
	processes map[string]*ManagedProcess
	sigChan   chan os.Signal
	doneChan  chan struct{}
	mu        sync.Mutex
}

func (lm *gWebLifeCycle) Start() error {
	return lm.StartAll()
}

func (lm *gWebLifeCycle) Stop() error {
	return lm.StopAll()
}

func (lm *gWebLifeCycle) Restart() error {
	if err := lm.StopAll(); err != nil {
		return err
	}
	return lm.StartAll()
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

func (lm *gWebLifeCycle) RegisterProcess(name, command string, args []string, restart, wait bool) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	proc := NewManagedProcess(name, command, args, wait)
	lm.processes[name] = proc
	return nil
}

func (lm *gWebLifeCycle) StartAll() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	for _, proc := range lm.processes {
		if err := proc.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (lm *gWebLifeCycle) StopAll() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	for name, proc := range lm.processes {
		if err := proc.Stop(); err != nil {
			return err
		}
		delete(lm.processes, name)
	}
	return nil
}

func (lm *gWebLifeCycle) ListenForSignals() {
	signal.Notify(lm.sigChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-lm.sigChan:
			fmt.Println("Recebido sinal de encerramento. Finalizando processos...")
			_ = lm.StopAll()
			os.Exit(0)
		case <-lm.doneChan:
			return
		}
	}
}

func NewLifecycleManager() GWebLifeCycleManager {
	mgr := &gWebLifeCycle{
		processes: make(map[string]*ManagedProcess),
		sigChan:   make(chan os.Signal, 1),
		doneChan:  make(chan struct{}),
	}
	go mgr.ListenForSignals()
	return mgr
}
