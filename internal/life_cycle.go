package internal

import (
	"fmt"
	"github.com/faelmori/golife/internal/log"
	l "github.com/faelmori/logz"
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
	RegisterProcess(name string, command string, args []string, restart bool, customFn func() error) error
	RegisterEvent(event, stage string) error
	RemoveEvent(event, stage string) error
	StopEvents() error
	StartProcess(proc IManagedProcess) error
	StartAll() error
	StopAll() error
	Trigger(stage, event string, data interface{})
	DefineStage(name string) IStage
	Send(stage string, msg string)
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
func (lm *gWebLifeCycle) Send(stage string, msg string) {
	if s, ok := lm.stages[stage]; ok {
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

		log.Info(msg, map[string]interface{}{
			"context": "GoLife",
			"stage":   stage,
			"message": msg,
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
	l.Info("Iniciando processos...", map[string]interface{}{
		"context":  "GoLife",
		"showData": false,
	})
	for _, proc := range lm.processes {
		l.Info(fmt.Sprintf("Iniciando processo %s...", proc.String()), map[string]interface{}{
			"context":  "GoLife",
			"process":  proc.String(),
			"showData": false,
		})
		if err := proc.Start(); err != nil {
			l.Error(fmt.Sprintf("Erro ao iniciar processo %s: %v", proc.String(), err), map[string]interface{}{
				"context":  "GoLife",
				"process":  proc.String(),
				"showData": true,
			})
			return err
		}
	}
	l.Info(fmt.Sprintf("%b Processos iniciados com sucesso!", len(lm.processes)), map[string]interface{}{
		"context":   "GoLife",
		"processes": len(lm.processes),
		"showData":  false,
	})
	return nil
}
func (lm *gWebLifeCycle) Stop() error {
	l.Info("Parando processos...", map[string]interface{}{
		"context":  "GoLife",
		"showData": false,
	})
	for _, proc := range lm.processes {
		l.Info(fmt.Sprintf("Parando processo %s...", proc.String()), map[string]interface{}{
			"context":  "GoLife",
			"process":  proc.String(),
			"showData": false,
		})
		if err := proc.Stop(); err != nil {
			l.Error(fmt.Sprintf("Erro ao parar processo %s: %v", proc.String(), err), map[string]interface{}{
				"context":  "GoLife",
				"process":  proc.String(),
				"showData": true,
			})
			return err
		}
	}
	l.Info(fmt.Sprintf("%b Processos parados com sucesso!", len(lm.processes)), map[string]interface{}{
		"context":   "GoLife",
		"processes": len(lm.processes),
		"showData":  false,
	})
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
	lm.mu.Lock()
	defer lm.mu.Unlock()
	l.Info("Verificando status dos processos...", map[string]interface{}{"context": "GoLife", "showData": false})
	var status string
	for name, proc := range lm.processes {
		//fmt.Printf("Processo %s (PID %d) está rodando: %t\n", name, proc.Pid(), proc.IsRunning())
		l.Info(fmt.Sprintf("Processo %s (PID %d) está rodando: %t", name, proc.Pid(), proc.IsRunning()), map[string]interface{}{
			"context":  "GoLife",
			"process":  name,
			"pid":      proc.Pid(),
			"running":  proc.IsRunning(),
			"showData": false,
		})
		status += fmt.Sprintf("Processo %s (PID %d) está rodando: %t\n", name, proc.Pid(), proc.IsRunning())
	}
	l.Info("Status dos processos verificado com sucesso!", map[string]interface{}{"context": "GoLife", "showData": false})
	return status
}
func (lm *gWebLifeCycle) RegisterProcess(name string, command string, args []string, restart bool, customFn func() error) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	l.Info(fmt.Sprintf("Registrando processo %s...", name), map[string]interface{}{"context": "GoLife", "process": name, "showData": false})
	lm.processes[name] = NewManagedProcess(name, command, args, restart, customFn)

	l.Info(fmt.Sprintf("Processo %s registrado com sucesso!", name), map[string]interface{}{"context": "GoLife", "process": name, "showData": false})
	return nil
}
func (lm *gWebLifeCycle) RegisterEvent(event, stage string) error {
	lm.eventsMu.Lock()
	defer lm.eventsMu.Unlock()
	l.Info(fmt.Sprintf("Registrando evento %s em %s...", event, stage), map[string]interface{}{"context": "GoLife", "event": event, "stage": stage, "showData": false})
	lm.events = append(lm.events, NewManagedProcessEvents(map[string]func(interface{}){event: func(data interface{}) {
		l.Info(fmt.Sprintf("Executando evento %s em %s...", event, stage), map[string]interface{}{"context": "GoLife", "event": event, "stage": stage, "showData": false})
		lm.eventsCh <- NewManagedProcessEvents(map[string]func(interface{}){event: func(data interface{}) {
			l.Info(fmt.Sprintf("Disparando evento %s em %s...", event, stage), map[string]interface{}{"context": "GoLife", "event": event, "stage": stage, "showData": false})
			lm.Trigger(stage, event, data)
		}}, lm.triggerCh)
	}}, lm.triggerCh))
	l.Info(fmt.Sprintf("Evento %s registrado em %s com sucesso!", event, stage), map[string]interface{}{"context": "GoLife", "event": event, "stage": stage, "showData": false})
	return nil
}
func (lm *gWebLifeCycle) RemoveEvent(event, stage string) error {
	l.Info(fmt.Sprintf("Removendo evento %s de %s...", event, stage), map[string]interface{}{"context": "GoLife", "event": event, "stage": stage, "showData": false})
	for i, e := range lm.events {
		if e.Event() == event {
			lm.events = append(lm.events[:i], lm.events[i+1:]...)
			l.Info(fmt.Sprintf("Evento %s removido de %s com sucesso!", event, stage), map[string]interface{}{"context": "GoLife", "event": event, "stage": stage, "showData": false})
			return nil
		}
	}
	l.Error(fmt.Sprintf("Evento %s não encontrado em %s", event, stage), map[string]interface{}{"context": "GoLife", "event": event, "stage": stage, "showData": false})
	return fmt.Errorf("evento %s não encontrado em %s", event, stage)
}
func (lm *gWebLifeCycle) StopEvents() error {
	l.Info("Parando eventos...", map[string]interface{}{"context": "GoLife", "showData": false})
	for _, event := range lm.events {
		event.StopAll()
	}
	l.Info("Eventos parados com sucesso!", map[string]interface{}{"context": "GoLife", "showData": false})
	return nil
}
func (lm *gWebLifeCycle) StartAll() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	for name, proc := range lm.processes {
		l.Info(fmt.Sprintf("Iniciando %s...", name), map[string]interface{}{"context": "GoLife", "process": name, "showData": false})
		if err := lm.StartProcess(proc); err != nil {
			l.Error(fmt.Sprintf("Erro ao iniciar %s: %v", name, err), map[string]interface{}{"context": "GoLife", "process": name, "showData": true})
			return err
		}
	}
	l.Info(fmt.Sprintf("%b Processos iniciados com sucesso!", len(lm.processes)), map[string]interface{}{"context": "GoLife", "processes": len(lm.processes), "showData": false})
	return nil
}
func (lm *gWebLifeCycle) StartProcess(proc IManagedProcess) error {
	if err := proc.Start(); err != nil {
		l.Error(fmt.Sprintf("Erro ao iniciar %s: %v", proc.String(), err), map[string]interface{}{"context": "GoLife", "process": proc.String(), "showData": true})
		return err
	}
	l.Info(fmt.Sprintf("%s iniciado com sucesso!", proc.String()), map[string]interface{}{"context": "GoLife", "process": proc.String(), "showData": false})
	return nil
}
func (lm *gWebLifeCycle) StopAll() error {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	for name, proc := range lm.processes {
		if err := proc.Stop(); err != nil {
			l.Error(fmt.Sprintf("Erro ao parar %s: %v", name, err), map[string]interface{}{"context": "GoLife", "process": name, "showData": true})
			return err
		} else {
			l.Info(fmt.Sprintf("%s parado com sucesso!", name), map[string]interface{}{"context": "GoLife", "process": name, "showData": false})
			delete(lm.processes, name)
		}
	}
	l.Info(fmt.Sprintf("%b Processos parados com sucesso!", len(lm.processes)), map[string]interface{}{"context": "GoLife", "processes": len(lm.processes), "showData": false})
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

func NewLifecycleManager(
	processes map[string]IManagedProcess,
	stages map[string]IStage,
	sigChan chan os.Signal,
	doneChan chan struct{},
	events []IManagedProcessEvents,
	eventsCh chan IManagedProcessEvents,
) GWebLifeCycleManager {
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
			l.Error(fmt.Sprintf("Erro ao ouvir sinais: %v", err), map[string]interface{}{"context": "GoLife", "showData": true})
		}
	}()

	return &mgr
}

func NewLifecycleMgrSig() (GWebLifeCycleManager, error) {
	processes := make(map[string]IManagedProcess)
	stages := make(map[string]IStage)
	sigChan := make(chan os.Signal, 1)
	doneChan := make(chan struct{}, 1)
	events := make([]IManagedProcessEvents, 0)
	eventsCh := make(chan IManagedProcessEvents, 1)

	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}

func NewLifecycleMgrManual(
	processes map[string]IManagedProcess,
	stages map[string]IStage,
	sigChan chan os.Signal,
	doneChan chan struct{},
	events []IManagedProcessEvents,
	eventsCh chan IManagedProcessEvents,
) (GWebLifeCycleManager, error) {
	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}

func NewLifecycleMgrDec() (GWebLifeCycleManager, error) {
	processes := make(map[string]IManagedProcess)
	stages := make(map[string]IStage)
	sigChan := make(chan os.Signal, 1)
	doneChan := make(chan struct{}, 1)
	events := []IManagedProcessEvents{}
	eventsCh := make(chan IManagedProcessEvents, 1)

	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}
