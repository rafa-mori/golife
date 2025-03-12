package golife

import (
	"fmt"
	svcs "github.com/faelmori/gkbxsrv/services"
	l "github.com/faelmori/golife/internal"
	"os"
)

var manager l.GWebLifeCycleManager
var brk svcs.BrokerClient
var dbChanData chan map[string]interface{}
var lcErr error

func NewLifecycleMgr(processName, processCmd string, processArgs []string, processWait, restart bool, stages []string, triggers []string, processEvents map[string]func(interface{})) (l.GWebLifeCycleManager, error) {
	return createManager(processName, processCmd, processArgs, processEvents, stages, triggers, processWait, restart)
}

func NewLifecycleMgrDec(processName, processCmd string, processArgs []string, processWait, restart bool, stages []string, triggers []string, processEvents map[string]func(interface{})) (l.GWebLifeCycleManager, error) {
	return createManager(processName, processCmd, processArgs, processEvents, stages, triggers, processWait, restart)
}

func NewLifecycleMgrSig() (l.GWebLifeCycleManager, error) {
	return createManager("", "", nil, nil, nil, nil, false, false)
}

func NewLifecycleMgrManual(processName, processCmd string, processArgs []string, processWait, restart bool, stages []string, triggers []string, processEvents map[string]func(interface{})) (l.GWebLifeCycleManager, error) {
	return createManager(processName, processCmd, processArgs, processEvents, stages, triggers, processWait, restart)
}

func createManager(processName, processCmd string, stages []string, processEvents map[string]func(interface{}), triggers []string, processArgs []string, processWait, restart bool) (l.GWebLifeCycleManager, error) {
	if processName == "" {
		return nil, fmt.Errorf("no process name provided")
	}
	if processCmd == "" {
		return nil, fmt.Errorf("no command provided")
	}
	if len(stages) == 0 {
		stages = []string{"all"}
	}
	if len(processEvents) == 0 {
		processEvents = make(map[string]func(interface{}))
	} else {
		for _, trigger := range triggers {
			for _, stage := range stages {
				processEvents[trigger] = func(data interface{}) {
					manager.Trigger(stage, trigger, data)
				}
			}
		}
	}

	var events []l.IManagedProcessEvents
	var processes = make(map[string]l.IManagedProcess)
	var iStages = make(map[string]l.IStage)
	var sigChan = make(chan os.Signal, 1)
	var doneChan = make(chan struct{}, 1)
	var eventsChan = make(chan interface{}, 1)
	var eventsCh = make(chan l.IManagedProcessEvents, 1)

	for _, stage := range stages {
		iStage := l.NewStage(stage, stage, stage, "stage")
		iStages[stage] = iStage
	}

	for _, trigger := range triggers {
		iStage := l.NewStage(trigger, trigger, trigger, "trigger")
		iStages[trigger] = iStage
	}

	for _, stage := range iStages {
		for _, trigger := range triggers {
			stage.OnEvent(trigger, func(data interface{}) {
				eventsChan <- trigger
			})
		}
	}

	processes[processName] = l.NewManagedProcess(processName, processCmd, processArgs, processWait)
	if processEvents != nil {
		iEvent := l.NewManagedProcessEvents(processEvents, eventsChan)
		events = append(events, iEvent)
	}

	manager = l.NewLifecycleManager(
		processes,
		iStages,
		sigChan,
		doneChan,
		events,
		eventsCh,
	)

	regProcErr := manager.RegisterProcess(processName, processCmd, processArgs, restart)
	if regProcErr != nil {
		return nil, regProcErr
	}

	for _, stage := range iStages {
		manager.DefineStage(stage.Name())
	}

	startAllErr := manager.Start()
	if startAllErr != nil {
		fmt.Println("Erro ao iniciar processos:", startAllErr)
		return nil, startAllErr
	}

	return manager, manager.ListenForSignals()
}
