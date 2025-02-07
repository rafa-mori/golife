package golife

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

type GWebLifeCycle struct {
	processes map[string]*ManagedProcess
	stages    map[string]*Stage
	sigChan   chan os.Signal
	doneChan  chan struct{}
	mu        sync.Mutex
}

func (lm *GWebLifeCycle) Start() error {
	lm.StartAll()
	return nil
}

func (lm *GWebLifeCycle) Stop() error {
	lm.StopAll()
	return nil
}

func (lm *GWebLifeCycle) Restart() error {
	lm.StopAll()
	lm.StartAll()
	return nil
}

func (lm *GWebLifeCycle) Status() string {
	return "Running"
}

func (lm *GWebLifeCycle) RegisterProcess(name string, command string, args []string, restart bool) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.processes[name] = &ManagedProcess{
		Cmd:  exec.Command(command, args...),
		Name: name,
	}
}

func (lm *GWebLifeCycle) StartAll() {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	for name, proc := range lm.processes {
		if err := lm.startProcess(proc); err != nil {
			fmt.Printf("Erro ao iniciar %s: %v\n", name, err)
		}
	}
}

func (lm *GWebLifeCycle) startProcess(proc *ManagedProcess) error {
	if err := proc.Start(); err != nil {
		return err
	}
	fmt.Printf("Processo %s iniciado com PID %d\n", proc.Name, proc.Pid())
	return nil
}

func (lm *GWebLifeCycle) StopAll() {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	for name, proc := range lm.processes {
		if err := proc.Stop(); err != nil {
			fmt.Printf("Erro ao encerrar %s: %v\n", name, err)
		}
	}
}

func (lm *GWebLifeCycle) listenForSignals() {
	<-lm.sigChan
	fmt.Println("Recebido sinal de encerramento, finalizando processos...")
	lm.StopAll()
	os.Exit(0)
}

func NewLifecycleManager() *GWebLifeCycle {
	mgr := &GWebLifeCycle{
		processes: make(map[string]*ManagedProcess),
		sigChan:   make(chan os.Signal, 1),
		doneChan:  make(chan struct{}),
	}

	signal.Notify(mgr.sigChan, syscall.SIGINT, syscall.SIGTERM)
	go mgr.listenForSignals()

	return mgr
}
