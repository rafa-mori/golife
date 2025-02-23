package internal

import (
	"fmt"
	"github.com/google/uuid"
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
	RegisterProcess(name string, command string, args []string, restart bool) error
	RegisterEvent(event, stage string) error
	StartProcess(proc IManagedProcess) error
	StartAll() error
	StopAll() error
	Trigger(stage, event string, data interface{})
	DefineStage(name string) IStage
	Send(stage string, msg interface{})
	Receive(stage string) interface{}
	ListenForSignals() error
}
type gWebLifeCycle struct {
	processes map[string]IManagedProcess
	stages    map[string]IStage
	sigChan   chan os.Signal
	doneChan  chan struct{}
	events    []IManagedProcessEvents
	eventsMu  sync.Mutex
	eventsCh  chan IManagedProcessEvents
	triggerCh chan interface{}
	mu        sync.Mutex
}

func (lm *gWebLifeCycle) Trigger(stage, event string, data interface{}) {
	fmt.Printf("Disparando evento %s em %s...\n", event, stage)
	if s, ok := lm.stages[stage]; ok {
		fmt.Printf("Evento %s disparado em %s!\n", event, stage)
		sB := s.(*Stage)
		if fn, ok := sB.EventFns[event]; ok {
			fmt.Printf("Executando evento %s em %s...\n", event, stage)
			fn(data)
		}
	}
}
func (lm *gWebLifeCycle) DefineStage(name string) IStage {
	fmt.Printf("Definindo estágio %s...\n", name)
	return NewStage(uuid.New().String(), name, name, "stage")
}
func (lm *gWebLifeCycle) Send(stage string, msg interface{}) {
	if s, ok := lm.stages[stage]; ok {
		s.Dispatch(func() {
			sB := s.(*Stage)
			if sB.OnEnterFn != nil {
				s.OnEnter(sB.OnEnterFn)
			}
			if sB.OnExitFn != nil {
				s.OnExit(sB.OnExitFn)
			}
		})
	}
}
func (lm *gWebLifeCycle) Receive(stage string) interface{} {
	fmt.Printf("Recebendo mensagem de %s...\n", stage)
	if s, ok := lm.stages[stage]; ok {
		fmt.Printf("Mensagem recebida de %s!\n", stage)
		sB := s.(*Stage)
		return sB.Data
	} else {
		fmt.Printf("Nenhuma mensagem recebida de %s!\n", stage)
		return nil
	}
}
func (lm *gWebLifeCycle) Start() error {
	fmt.Println("Iniciando processos...")
	for _, proc := range lm.processes {
		fmt.Printf("Iniciando %s...\n", proc.String())
		if err := proc.Start(); err != nil {
			fmt.Printf("Erro ao iniciar processo %s: %v\n", proc.String(), err)
			return err
		}
	}
	if len(lm.processes) > 0 {
		fmt.Println("Processos iniciados com sucesso!")
	} else {
		fmt.Println("Nenhum processo registrado para iniciar.")
	}
	return nil
}
func (lm *gWebLifeCycle) Stop() error {
	fmt.Println("Parando processos...")
	for _, proc := range lm.processes {
		fmt.Printf("Parando %s...\n", proc.String())
		if err := proc.Stop(); err != nil {
			fmt.Println("Erro ao parar processo:", err)
			return err
		}
	}
	if len(lm.processes) > 0 {
		fmt.Println("Processos parados com sucesso!")
	} else {
		fmt.Println("Nenhum processo registrado para parar.")
	}
	return nil
}
func (lm *gWebLifeCycle) Restart() error {
	for _, proc := range lm.processes {
		if err := proc.Restart(); err != nil {
			return err
		}
	}
	return nil
}
func (lm *gWebLifeCycle) Status() string {
	//lm.mu.Lock()
	//defer lm.mu.Unlock()
	fmt.Printf("Verificando status dos processos...\n")
	var status string
	for name, proc := range lm.processes {
		fmt.Printf("Processo %s (PID %d) está rodando: %t\n", name, proc.Pid(), proc.IsRunning())
		status += fmt.Sprintf("Processo %s (PID %d) está rodando: %t\n", name, proc.Pid(), proc.IsRunning())
	}
	fmt.Printf("Status dos processos verificado!\n")
	return status
}
func (lm *gWebLifeCycle) RegisterProcess(name string, command string, args []string, restart bool) error {
	//lm.mu.Lock()
	//defer lm.mu.Unlock()

	fmt.Printf("Registrando processo %s...\n", name)
	lm.processes[name] = NewManagedProcess(name, command, args, restart)

	fmt.Printf("Processo %s registrado com sucesso!\n", name)
	return nil
}
func (lm *gWebLifeCycle) RegisterEvent(event, stage string) error {
	//lm.eventsMu.Lock()
	//defer lm.eventsMu.Unlock()
	fmt.Printf("Registrando evento %s em %s...\n", event, stage)
	lm.events = append(lm.events, NewManagedProcessEvents(map[string]func(interface{}){event: func(data interface{}) {
		fmt.Printf("Evento %s disparado em %s!\n", event, stage)
		lm.eventsCh <- NewManagedProcessEvents(map[string]func(interface{}){event: func(data interface{}) {
			fmt.Printf("Executando evento %s em %s...\n", event, stage)
			lm.Trigger(stage, event, data)
		}}, lm.triggerCh)
	}}, lm.triggerCh))
	fmt.Printf("Evento %s registrado em %s com sucesso!\n", event, stage)
	return nil
}
func (lm *gWebLifeCycle) StartAll() error {
	//lm.mu.Lock()
	//defer lm.mu.Unlock()
	fmt.Println("Iniciando processos...")
	for name, proc := range lm.processes {
		fmt.Printf("Iniciando %s...\n", name)
		if err := lm.StartProcess(proc); err != nil {
			fmt.Printf("Erro ao iniciar %s: %v\n", name, err)
			return err
		}
	}
	if len(lm.processes) > 0 {
		fmt.Println("Processos iniciados com sucesso!")
	} else {
		fmt.Println("Nenhum processo registrado para iniciar.")
	}
	return nil
}
func (lm *gWebLifeCycle) StartProcess(proc IManagedProcess) error {
	fmt.Printf(fmt.Sprintf("Iniciando processo %s...\n", proc.String()))
	if err := proc.Start(); err != nil {
		fmt.Printf("Erro ao iniciar processo %s: %v\n", proc.String(), err)
		return err
	}
	fmt.Printf("Processo %s iniciado com sucesso!\n", proc.String())
	return nil
}
func (lm *gWebLifeCycle) StopAll() error {
	//lm.mu.Lock()
	//defer lm.mu.Unlock()
	fmt.Println("Parando processos...")
	for name, proc := range lm.processes {
		if err := proc.Stop(); err != nil {
			fmt.Printf("Erro ao parar %s: %v\n", name, err)
			return err
		} else {
			fmt.Printf("Processo %s parado com sucesso!\n", name)
			delete(lm.processes, name)
		}
	}
	if len(lm.processes) > 0 {
		fmt.Println("Processos parados com sucesso!")
	} else {
		fmt.Println("Nenhum processo registrado para parar.")
	}
	return nil
}
func (lm *gWebLifeCycle) ListenForSignals() error {
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

func NewLifecycleManager(processes map[string]IManagedProcess, stages map[string]IStage, sigChan chan os.Signal, doneChan chan struct{}, events []IManagedProcessEvents, eventsCh chan IManagedProcessEvents) GWebLifeCycleManager {
	mgr := gWebLifeCycle{
		processes: processes,
		stages:    stages,
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
			fmt.Println("Erro ao ouvir sinais:", err)
		}
	}()

	return &mgr
}
