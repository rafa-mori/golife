package golife

import (
	"fmt"
	l "github.com/faelmori/golife/api"
	lg "github.com/faelmori/logz"
	"os"
)

type LifecycleManager = l.LifeCycleManager
type ManagedProcess = l.ManagedProcess
type ManagedProcessEvents = l.ManagedProcessEvent
type Stage = l.Stage
type WorkerPool = l.WorkerPool

func NewManagedProcess(name, cmd string, args []string, wait bool, customFn func() error) ManagedProcess {
	return l.NewManagedProcess(name, cmd, args, wait, customFn)
}

func NewLifecycleMgr(
	processName, processCmd string,
	processArgs []string,
	processWait, restart bool,
	stages, triggers []string,
	processEvents map[string]func(interface{}),
	customFn func() error,
	sigChan chan os.Signal,
	doneChan chan struct{},
	eventsChan chan interface{},
	eventsCh chan l.ManagedProcessEvent,
) (LifecycleManager, error) {
	return createManager(processName, processCmd, processArgs, processEvents, stages, triggers, processWait, restart, customFn, sigChan, doneChan, eventsChan, eventsCh)
}

func NewLifecycleMgrDec(
	processName, processCmd string,
	processArgs []string,
	processWait, restart bool,
	stages, triggers []string,
	processEvents map[string]func(interface{}),
	customFn func() error,
) (LifecycleManager, error) {
	return createManager(processName, processCmd, processArgs, processEvents, stages, triggers, processWait, restart, customFn, nil, nil, nil, nil)
}

func NewLifecycleMgrSig(
	sigChan chan os.Signal,
	doneChan chan struct{},
	eventsChan chan interface{},
	eventsCh chan l.ManagedProcessEvent,
) (LifecycleManager, error) {
	return createManager("", "", nil, nil, nil, nil, false, false, nil, sigChan, doneChan, eventsChan, eventsCh)
}

func NewLifecycleMgrManual(
	processName, processCmd string,
	processArgs []string,
	processWait,
	restart bool,
	stages, triggers []string,
	processEvents map[string]func(interface{}),
	customFn func() error,
	sigChan chan os.Signal,
	doneChan chan struct{},
	eventsChan chan interface{},
	eventsCh chan l.ManagedProcessEvent,
) (LifecycleManager, error) {
	return createManager(processName, processCmd, processArgs, processEvents, stages, triggers, processWait, restart, customFn, sigChan, doneChan, eventsChan, eventsCh)
}
func NewStage(name, desc, stageType string) Stage {
	return l.NewStage(name, desc, stageType)
}

func createManager(
	processName, processCmd string,
	stages []string,
	processEvents map[string]func(interface{}),
	triggers, processArgs []string,
	processWait, restart bool,
	customFn func() error,
	sigChan chan os.Signal,
	doneChan chan struct{},
	eventsChan chan interface{},
	eventsCh chan l.ManagedProcessEvent,
) (LifecycleManager, error) {
	if processName == "" {
		lg.Error("No process name provided", nil)
		return nil, fmt.Errorf("no process name provided")
	}
	var events []l.ManagedProcessEvent
	var processes = make(map[string]l.ManagedProcess)
	var iStages = make(map[string]l.Stage)
	if stages == nil {
		sigChan = make(chan os.Signal, 1)
	}
	if doneChan == nil {
		doneChan = make(chan struct{}, 1)
	}
	if eventsChan == nil {
		eventsChan = make(chan interface{}, 1)
	}
	if eventsCh == nil {
		eventsCh = make(chan l.ManagedProcessEvent, 1)
	}

	lcm, err := l.NewLifecycleMgrChan(sigChan, doneChan, eventsCh)
	if err != nil {
		return nil, err
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
					l.Trigger(lcm, stage, trigger, data)
				}
			}
		}
	}

	for _, stage := range stages {
		iStage := l.NewStage(stage, stage, "stage")
		iStages[stage] = iStage
	}

	for _, trigger := range triggers {
		iStage := l.NewStage(trigger, trigger, "trigger")
		iStages[trigger] = iStage
	}

	for _, stage := range iStages {
		for _, trigger := range triggers {
			stage.OnEvent(trigger, func(data interface{}) {
				eventsChan <- trigger
			})
		}
	}

	processes[processName] = l.NewManagedProcess(processName, processCmd, processArgs, processWait, customFn)
	if processEvents != nil {
		iEvent := l.NewEvent(processEvents, eventsChan)
		events = append(events, iEvent)
	}

	regProcErr := lcm.RegisterProcess(processName, processCmd, processArgs, restart, customFn)
	if regProcErr != nil {
		return nil, regProcErr
	}

	for _, stage := range iStages {
		if regStageErr := lcm.RegisterStage(stage); regStageErr != nil {
			return nil, regStageErr
		}
	}

	startAllErr := lcm.Start()
	if startAllErr != nil {
		fmt.Println("Erro ao iniciar processos:", startAllErr)
		return nil, startAllErr
	}

	return lcm, lcm.ListenForSignals()
}
