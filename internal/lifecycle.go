package internal

import (
	"fmt"
	ci "github.com/faelmori/golife/components/interfaces"
	pi "github.com/faelmori/golife/components/process_input"
	p "github.com/faelmori/golife/components/types"
	pr "github.com/faelmori/golife/internal/process"
	ev "github.com/faelmori/golife/internal/routines/taskz/events"
	st "github.com/faelmori/golife/internal/routines/taskz/stage"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

type ILifeCycle[T pi.ProcessInput[any]] interface {
	ci.IMutexes

	GetLogger() l.Logger

	Initialize() error

	GetProcess(string) (pr.IManagedProcess[any], error)
	AddProcess(string, *pi.ProcessInput[any]) error
	CurrentProcess() pr.IManagedProcess[any]
	RemoveProcess(string) error

	GetStage(string) (st.IStage[any], error)
	AddStage(string, st.IStage[any]) error
	CurrentStage() st.IStage[any]
	SetCurrentStage(string) (st.IStage[any], error)
	RemoveStage(string) error

	GetEvent(string) (ev.IManagedProcessEvents[any], error)
	AddEvent(string, ev.IManagedProcessEvents[any]) error
	TriggerEvent(string, interface{}) error
	RemoveEvent(string) error

	StartLifecycle() error
	StopLifecycle() error
	RestartLifecycle() error
	StatusLifecycle() string

	GetChannel(name string) (*p.ChannelCtl[any], error)
	GetChannels() map[string]p.ChannelCtl[any]
	AddChannel(name string, channel *p.ChannelCtl[any]) error
	RemoveChannel(name string) error

	GetSupervisor() *pi.ProcessInput[any]
	GetSupervisorProcess() pr.IManagedProcess[pi.ProcessInput[any]]

	ListenForSignals() error

	//ListenForTerminalInput() error
}
type LifeCycle[T pi.ProcessInput[any]] struct {
	// supervisor is the supervisor instance for the lifecycle manager
	supervisor *pi.ProcessInput[any]
	// supervisorProcess is the supervisor process
	supervisorProcess pr.IManagedProcess[pi.ProcessInput[any]]

	// Logger is the Logger instance for the lifecycle manager
	Logger l.Logger

	// Mutexes is the mutex for the lifecycle manager
	*p.Mutexes

	// Reference is the reference ID and name
	*p.Reference

	// ProcessInput is the process input
	ProcessInput *T

	// initialized is a boolean that indicates if the lifecycle manager is initialized
	initialized bool

	// controllers is a map of controllers
	controllers map[string]any
}

func NewLifeCycle[T pi.ProcessInput[any]](processInput *T) ILifeCycle[T] {
	if processInput == nil {
		gl.Log("fatal", "Process input arg is nil")
		return nil
	} else {
		inputT := *processInput
		input := reflect.ValueOf(inputT).Interface().(pi.ProcessInput[any])

		logger := input.Logger
		if logger == nil {
			logger = l.GetLogger("GoLife")
		}

		lcm := &LifeCycle[T]{
			Logger:      logger,
			Mutexes:     p.NewMutexesType(),
			Reference:   p.NewReference("LifeCycle").GetReference(),
			initialized: false,
		}

		lcm.ProcessInput = processInput
		gl.Log("info", "Lifecycle manager initialized successfully")

		if err := lcm.Initialize(); err != nil {
			gl.Log("error", fmt.Sprintf("Error initializing lifecycle manager: %v", err))
			return nil
		}

		return lcm
	}
}

func (lm *LifeCycle[T]) GetLogger() l.Logger {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return nil
	}
	return lm.Logger
}
func (lm *LifeCycle[T]) Initialize() error {
	if lm == nil {
		gl.LogObjLogger(lm, "fatal", "Lifecycle manager is nil")
		return fmt.Errorf("lifecycle manager is nil")
	}

	if lm.ProcessInput == nil {
		gl.LogObjLogger(lm, "fatal", "ProcessInput is nil")
		return fmt.Errorf("LifeCycle process input is nil")
	}
	if lm.Mutexes == nil {
		gl.LogObjLogger(lm, "fatal", "Mutexes is nil")
		return fmt.Errorf("LifeCycle mutexes is nil")
	}
	if lm.Logger == nil {
		gl.LogObjLogger(lm, "fatal", "Logger is nil")
		return fmt.Errorf("LifeCycle logger is nil")
	}

	if lm.controllers == nil {
		gl.LogObjLogger(lm, "debug", "Initializing lifecycle manager controllers")
		lm.controllers = make(map[string]any)
	} else {
		gl.LogObjLogger(lm, "fatal", "Lifecycle manager controllers already initialized")
		return fmt.Errorf("lifecycle manager controllers already initialized")
	}

	// mainInput is the process input
	lm.controllers["mainProcess"] = lm.ProcessInput // main input is the process input
	// processes is a map of processes
	lm.controllers["processes"] = make(map[string]interface{}) //map[string]pProperty[pr.IManagedProcess[any]]
	// stages is a map of stages
	lm.controllers["stages"] = make(map[string]interface{}) //map[string]pProperty[st.IStage[any]]
	// events is a map of events.
	lm.controllers["events"] = make(map[string]interface{}) //map[string]pProperty[ev.IManagedProcessEvents[any]]
	// channels is a map of channels
	lm.controllers["channels"] = make(map[string]interface{}) //map[string]p.ChannelCtl[any]
	// metadata is a map of metadata.
	lm.controllers["metadata"] = make(map[string]interface{})

	if lm.controllers["mainProcess"] == nil {
		gl.LogObjLogger(lm, "fatal", "Main process is nil")
		return fmt.Errorf("main process is nil")
	}

	gl.LogObjLogger(lm, "success", "Lifecycle manager controllers initialized successfully")

	gl.LogObjLogger(lm, "debug", "Initializing lifecycle manager processes")
	if err := initializeProcess(lm); err != nil {
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Error initializing process: %v", err))
		return fmt.Errorf("error initializing process: %v", err)
	}
	gl.LogObjLogger(lm, "success", "Lifecycle manager processes initialized successfully")

	gl.LogObjLogger(lm, "debug", "Initializing lifecycle manager stages")
	if err := initializeStages(lm); err != nil {
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Error initializing stages: %v", err))
		return fmt.Errorf("error initializing stages: %v", err)
	}
	gl.LogObjLogger(lm, "success", "Lifecycle manager stages initialized successfully")

	gl.LogObjLogger(lm, "debug", "Initializing lifecycle manager events")
	if err := initializeChannels(lm); err != nil {
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Error initializing channels: %v", err))
		return fmt.Errorf("error initializing channels: %v", err)
	}
	gl.LogObjLogger(lm, "success", "Lifecycle manager channels initialized successfully")

	lm.initialized = true

	return nil
}
func (lm *LifeCycle[T]) Snapshot() string {
	lm.Mutexes.MuRLock()
	defer lm.Mutexes.MuRUnlock()

	return fmt.Sprintf(
		"Lifecycle [%s]: Processes: %d, Stages: %d, Events: %d, Channels: %d",
		lm.Reference.Name,
		len(lm.controllers["processes"].(map[string]interface{})),
		len(lm.controllers["stages"].(map[string]interface{})),
		len(lm.controllers["events"].(map[string]interface{})),
		len(lm.controllers["channels"].(map[string]interface{})),
	)
}

func (lm *LifeCycle[T]) GetProcess(name string) (pr.IManagedProcess[any], error) {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return nil, fmt.Errorf("lifecycle manager is nil")
	}
	if process, exists := lm.controllers["processes"].(map[string]any)[name]; exists {
		return process.(pr.IManagedProcess[any]), nil
	}
	return nil, fmt.Errorf("process '%s' does not exist", name)
}
func (lm *LifeCycle[T]) AddProcess(name string, process *pi.ProcessInput[any]) error {

	processes := lm.controllers["processes"].(map[string]interface{})
	if _, exists := processes[name]; exists {
		return fmt.Errorf("process '%s' already exists", name)
	}

	processes[name] = pr.NewManagedProcess(process.Name, process.GetCommand(), process.GetArgs(), process.GetWaitFor(), process.GetFunction())
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Added process '%s'", name))
	return nil
}
func (lm *LifeCycle[T]) CurrentProcess() pr.IManagedProcess[any] {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return nil
	}

	if currentProcess, exists := lm.controllers["currentProcess"]; exists {
		if currentProcess == nil {
			gl.LogObjLogger(lm, "error", "Current process is nil")
			return nil
		}
		if process, ok := currentProcess.(pr.IManagedProcess[any]); ok {
			return process
		}
	}

	return nil
}
func (lm *LifeCycle[T]) RemoveProcess(name string) error {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return fmt.Errorf("lifecycle manager is nil")
	}

	if _, exists := lm.controllers["processes"].(map[string]any)[name]; !exists {
		return fmt.Errorf("process '%s' does not exist", name)
	}

	delete(lm.controllers["processes"].(map[string]any), name)
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Removed process '%s'", name))
	return nil
}

func (lm *LifeCycle[T]) GetStage(name string) (st.IStage[any], error) {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return nil, fmt.Errorf("lifecycle manager is nil")
	}

	lm.Mutexes.MuRLock()
	defer lm.Mutexes.MuRUnlock()

	if stage, exists := lm.controllers["stages"].(map[string]any)[name]; exists {
		if stageObj, ok := stage.(st.IStage[any]); ok {
			return stageObj, nil
		}
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Failed to cast stage '%s'", name))
		return nil, fmt.Errorf("failed to cast stage '%s'", name)
	}
	return nil, fmt.Errorf("stage '%s' does not exist", name)
}
func (lm *LifeCycle[T]) AddStage(name string, stage st.IStage[any]) error {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return fmt.Errorf("lifecycle manager is nil")
	}
	if stage == nil {
		gl.LogObjLogger(lm, "error", "Stage is nil")
		return fmt.Errorf("stage is nil")
	}
	if name == "" {
		gl.LogObjLogger(lm, "error", "Stage name is empty")
		return fmt.Errorf("stage name is empty")
	}

	if stageRegistered, err := lm.GetStage(name); err != nil {
		lm.Mutexes.MuLock()
		defer lm.Mutexes.MuUnlock()

		lm.controllers["stages"].(map[string]any)[name] = stage
		gl.LogObjLogger(lm, "info", fmt.Sprintf("Added stage '%s'", name))
	} else {
		gl.LogObjLogger(lm, "warn", fmt.Sprintf("Stage '%s' already registered: %v", name, stageRegistered.GetStageID()))
	}
	return nil
}
func (lm *LifeCycle[T]) CurrentStage() st.IStage[any] {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return nil
	}

	lm.Mutexes.MuRLock()
	defer lm.Mutexes.MuRUnlock()

	if currentStage, exists := lm.controllers["currentStage"]; exists {
		if stage, ok := currentStage.(st.IStage[any]); ok {
			return stage
		}
	} else {
		if currentStage == nil {
			gl.LogObjLogger(lm, "info", "Stages not initialized")
			return nil
		}
	}
	return nil
}
func (lm *LifeCycle[T]) SetCurrentStage(stageName string) (st.IStage[any], error) {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return nil, fmt.Errorf("lifecycle manager is nil")
	}

	lm.Mutexes.MuLock()
	defer lm.Mutexes.MuUnlock()

	if currentStage := lm.CurrentStage(); currentStage != nil {
		if currentStage.Name() == stageName {
			gl.LogObjLogger(lm, "info", fmt.Sprintf("Stage '%s' is already current", stageName))
			return currentStage, nil
		}
		if !currentStage.CanTransitionTo(stageName) {
			gl.LogObjLogger(lm, "error", fmt.Sprintf("Stage '%s' cannot transition to '%s'", currentStage.Name(), stageName))
			return nil, fmt.Errorf("stage '%s' cannot transition to '%s'", currentStage.Name(), stageName)
		}
	}
	if stageObj, exists := lm.controllers["stages"].(map[string]any)[stageName]; !exists {
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Stage '%s' does not exist", stageName))
		return nil, fmt.Errorf("stage '%s' does not exist", stageName)
	} else {
		if stage, ok := stageObj.(st.IStage[any]); ok {
			lm.controllers["currentStage"] = stage

			gl.LogObjLogger(lm, "info", fmt.Sprintf("Transitioned to stage '%s'", stageName))
			if stage.EventExists("start") {
				startEvent := stage.GetEvent("start")
				if startEvent != nil {
					if pp := lm.CurrentProcess(); pp != nil {
						gl.LogObjLogger(lm, "info", fmt.Sprintf("Triggering 'start' event in %s stage with process: '%s'", stage.Name(), pp.GetName()))
						startEvent(pp.GetName())
					} else {
						gl.LogObjLogger(lm, "info", fmt.Sprintf("Triggering 'start' event in stage: '%s'", stage.Name()))
						startEvent(nil)
					}
				}
			} else {
				gl.LogObjLogger(lm, "info", fmt.Sprintf("No 'start' event found in stage: '%s'", stage.Name()))
			}

			return stage, nil
		} else {
			gl.LogObjLogger(lm, "error", fmt.Sprintf("Failed to cast stage '%s'", stageName))
			return nil, fmt.Errorf("failed to cast stage '%s'", stageName)
		}
	}
}
func (lm *LifeCycle[T]) RemoveStage(name string) error {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return fmt.Errorf("lifecycle manager is nil")
	}
	if _, err := lm.GetStage(name); err != nil {
		return fmt.Errorf("stage '%s' does not exist", name)
	}

	delete(lm.controllers["stages"].(map[string]any), name)

	gl.LogObjLogger(lm, "info", fmt.Sprintf("Removed stage '%s'", name))

	return nil
}

func (lm *LifeCycle[T]) GetEvent(name string) (ev.IManagedProcessEvents[any], error) {
	lm.Mutexes.MuRLock()
	defer lm.Mutexes.MuRUnlock()

	events := lm.controllers["events"].(map[string]interface{})
	if event, exists := events[name]; exists {
		return event.(ev.IManagedProcessEvents[any]), nil
	}
	return nil, fmt.Errorf("event '%s' does not exist", name)
}
func (lm *LifeCycle[T]) AddEvent(name string, event ev.IManagedProcessEvents[any]) error {

	events := lm.controllers["events"].(map[string]interface{})
	if _, exists := events[name]; exists {
		return fmt.Errorf("event '%s' already exists", name)
	}

	events[name] = event
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Added event '%s'", name))
	return nil
}
func (lm *LifeCycle[T]) RemoveEvent(name string) error {

	events := lm.controllers["events"].(map[string]interface{})
	if _, exists := events[name]; !exists {
		return fmt.Errorf("event '%s' does not exist", name)
	}

	delete(events, name)
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Removed event '%s'", name))
	return nil
}
func (lm *LifeCycle[T]) TriggerEvent(name string, data interface{}) error {
	events := lm.controllers["events"].(map[string]interface{})
	if event, exists := events[name]; exists {
		if evnt, ok := event.(ev.IManagedProcessEvents[any]); ok {
			evnt.Trigger("", name, data)
			gl.LogObjLogger(lm, "info", fmt.Sprintf("Event '%s' triggered successfully!", name))
			return nil
		}
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Failed to cast event '%s'", name))
		return fmt.Errorf("failed to cast event '%s'", name)
	}
	return fmt.Errorf("event '%s' does not exist", name)
}

func (lm *LifeCycle[T]) GetChannel(name string) (*p.ChannelCtl[any], error) {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return nil, fmt.Errorf("lifecycle manager is nil")
	}

	lm.Mutexes.MuRLock()
	defer lm.Mutexes.MuRUnlock()

	if channel, exists := lm.controllers["channels"].(map[string]any)[name]; exists {
		if channelObj, ok := channel.(p.ChannelCtl[any]); ok {
			return &channelObj, nil
		}
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Failed to cast channel '%s'", name))
		return nil, fmt.Errorf("failed to cast channel '%s'", name)
	}
	return nil, fmt.Errorf("channel '%s' does not exist", name)
}
func (lm *LifeCycle[T]) GetChannels() map[string]p.ChannelCtl[any] {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return nil
	}

	lm.Mutexes.MuRLock()
	defer lm.Mutexes.MuRUnlock()

	channels := make(map[string]p.ChannelCtl[any])
	for name, channel := range lm.controllers["channels"].(map[string]any) {
		if channelObj, ok := channel.(p.ChannelCtl[any]); ok {
			channels[name] = channelObj
		}
	}
	return channels
}
func (lm *LifeCycle[T]) AddChannel(name string, channel *p.ChannelCtl[any]) error {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return fmt.Errorf("lifecycle manager is nil")
	}
	if channel == nil {
		gl.LogObjLogger(lm, "error", "Channel is nil")
		return fmt.Errorf("channel is nil")
	}
	if name == "" {
		gl.LogObjLogger(lm, "error", "Channel name is empty")
		return fmt.Errorf("channel name is empty")
	}

	if channelRegistered, err := lm.GetChannel(name); err != nil {
		lm.Mutexes.MuLock()
		defer lm.Mutexes.MuUnlock()

		lm.controllers["channels"].(map[string]any)[name] = channel
		gl.LogObjLogger(lm, "info", fmt.Sprintf("Added channel '%s'", name))
	} else {
		gl.LogObjLogger(lm, "warn", fmt.Sprintf("Channel '%s' already registered: %v", name, channelRegistered))
	}
	return nil
}
func (lm *LifeCycle[T]) RemoveChannel(name string) error {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return fmt.Errorf("lifecycle manager is nil")
	}
	if _, err := lm.GetChannel(name); err != nil {
		return fmt.Errorf("channel '%s' does not exist", name)
	}

	delete(lm.controllers["channels"].(map[string]any), name)

	gl.LogObjLogger(lm, "info", fmt.Sprintf("Removed channel '%s'", name))

	return nil
}

func (lm *LifeCycle[T]) GetSupervisor() *pi.ProcessInput[any] {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return nil
	}

	lm.Mutexes.MuRLock()
	defer lm.Mutexes.MuRUnlock()

	return lm.supervisor
}
func (lm *LifeCycle[T]) GetSupervisorProcess() pr.IManagedProcess[pi.ProcessInput[any]] {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return nil
	}

	lm.Mutexes.MuRLock()
	defer lm.Mutexes.MuRUnlock()

	return lm.supervisorProcess
}

func (lm *LifeCycle[T]) StartLifecycle() error {

	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return fmt.Errorf("lifecycle manager is nil")
	}
	gl.LogObjLogger(lm, "info", "Starting lifecycle via stages")

	if len(lm.controllers["stages"].(map[string]interface{})) == 0 {
		return fmt.Errorf("no stages available to start lifecycle")
	}

	if stageInit, stageInitErr := lm.SetCurrentStage("init"); stageInitErr != nil {
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Error setting current stage to 'init': %v", stageInitErr))
		return fmt.Errorf("error setting current stage to 'init': %v", stageInitErr)
	} else {
		if stageInit == nil {
			gl.LogObjLogger(lm, "error", "Stage 'init' is nil")
			return fmt.Errorf("stage 'init' is nil")
		}
		if stageExecute, stageExecuteErr := lm.SetCurrentStage("execute"); stageExecuteErr != nil {
			gl.LogObjLogger(lm, "error", fmt.Sprintf("Error setting current stage to 'execute': %v", stageExecuteErr))
			return fmt.Errorf("error setting current stage to 'execute': %v", stageExecuteErr)
		} else {
			if stageExecute == nil {
				gl.LogObjLogger(lm, "error", "Stage 'execute' is nil")
				return fmt.Errorf("stage 'execute' is nil")
			}
		}
	}

	gl.LogObjLogger(lm, "success", "Lifecycle started via stages successfully!")
	return nil
}
func (lm *LifeCycle[T]) StopLifecycle() error {

	gl.LogObjLogger(lm, "info", "Stopping lifecycle")
	if processes, ok := lm.controllers["processes"].(map[string]interface{}); ok {
		for _, proc := range processes {
			if prc, ok := proc.(*pr.ManagedProcess[any]); ok {
				if err := prc.Stop(); err != nil {
					gl.LogObjLogger(lm, "error", fmt.Sprintf("Error stopping process '%s': %v", prc.Name, err))
					return err
				} // Presumes that the method Stop() exists.
				gl.LogObjLogger(lm, "info", fmt.Sprintf("Stopped process '%s'", prc.Name))
			}
		}
	}

	if channels, ok := lm.controllers["channels"].(map[string]interface{}); ok {
		for name, ch := range channels {
			if ctl, ok := ch.(p.ChannelCtl[any]); ok {
				_ = ctl.Close() // Presumes that the method Close() exists.
				gl.LogObjLogger(lm, "info", fmt.Sprintf("Closed channel '%s'", name))
			}
		}
	}

	return nil
}

func (lm *LifeCycle[T]) RestartLifecycle() error {
	if lm == nil {
		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
		return fmt.Errorf("lifecycle manager is nil")
	}

	processes := lm.controllers["processes"].(map[string]interface{})
	for name, proc := range processes {
		if prc, ok := proc.(pr.IManagedProcess[any]); ok {
			if err := prc.Restart(); err != nil {
				gl.LogObjLogger(lm, "error", fmt.Sprintf("Error restarting process '%s': %v", name, err))
				return err
			}
			gl.LogObjLogger(lm, "info", fmt.Sprintf("Process '%s' restarted successfully!", name))
		}
	}
	return nil
}
func (lm *LifeCycle[T]) StatusLifecycle() string {
	//lm.Mutexes.MuRLock()
	//defer lm.Mutexes.MuRUnlock()

	processes := lm.controllers["processes"].(map[string]interface{})
	status := "Processes Status:\n"
	for name, proc := range processes {
		if prc, ok := proc.(pr.IManagedProcess[any]); ok {
			status += fmt.Sprintf("Process '%s': %s\n", name, prc.Status())
		}
	}
	return status
}

func (lm *LifeCycle[T]) ListenForSignals() error {
	if lm.controllers["channels"] == nil {
		gl.LogObjLogger(lm, "error", "Channels are nil")
		return fmt.Errorf("channels are nil")
	}
	channels := lm.controllers["channels"].(map[string]interface{})
	if sigChan, ok := channels["chanSignal"].(ci.IChannelCtl[os.Signal]); ok {
		cSignal, cSignalType, cSignalExists := sigChan.GetSubChannelByName("signal")
		if !cSignalExists {
			gl.LogObjLogger(lm, "error", "Signal channel is nil")
			return fmt.Errorf("signal channel is nil")
		}
		if cSignalType != reflect.TypeFor[os.Signal]() {
			gl.LogObjLogger(lm, "error", "Signal channel type is not os.Signal")
			return fmt.Errorf("signal channel type is not os.Signal")
		}
		chanSignal := reflect.ValueOf(cSignal.GetChannel()).Interface().(chan os.Signal)
		signal.Notify(chanSignal, os.Interrupt, os.Kill, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

		go func(chanSignal chan os.Signal) {
			for sig := range chanSignal {
				gl.LogObjLogger(lm, "info", fmt.Sprintf("Received signal: %s", sig))
				chanSignal <- sig
			}
		}(chanSignal)

		gl.LogObjLogger(lm, "info", fmt.Sprintf("Listening for signals on channel '%s'", sigChan.GetName()))
	} else {
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Failed to cast signal channel"))
		return fmt.Errorf("failed to cast signal channel")
	}
	return nil
}

//func (lm *LifeCycle[T]) ListenForTerminalInput() error {
//	if lm == nil {
//		gl.LogObjLogger(lm, "error", "Lifecycle manager is nil")
//		return fmt.Errorf("lifecycle manager is nil")
//	}
//
//	objChErr := p.NewChannelCtl[error]("chanError", func(b int) *int { return &b }(1), lm.Logger)
//	objChDone := p.NewChannelCtl[bool]("chanDone", func(b int) *int { return &b }(1), lm.Logger)
//	objChSignal := p.NewChannelCtl[os.Signal]("chanSignal", func(b int) *int { return &b }(1), lm.Logger)
//
//	chErr := objChErr.Channel()
//	chDone := objChDone.Channel()
//	chSignal := objChSignal.Channel()
//
//	go listenStdin(lm, chErr, chDone, chSignal)
//
//	gl.LogObjLogger(lm, "info", fmt.Sprintf("Listening for terminal input on channel '%s'", objChErr.Name))
//
//	for {
//		select {
//		case errorMsg := <-chErr:
//			if errorMsg != nil {
//				gl.LogObjLogger(lm, "error", fmt.Sprintf("Error from channel 'chanError': %v", errorMsg))
//			} else {
//				continue
//			}
//		case <-time.After(2 * time.Second):
//			if cst := lm.CurrentStage(); cst != nil {
//				if cst.Name() == "end" {
//					gl.LogObjLogger(lm, "info", fmt.Sprintf("%s: End stage reached, stopping listening for terminal input", "Process monitor"))
//					chDone <- true
//					return nil
//				} else {
//					gl.LogObjLogger(lm, "info", fmt.Sprintf("%s: Current stage: %s", "Process monitor", cst.Name()))
//				}
//			}
//		default:
//			continue
//		}
//	}
//}
