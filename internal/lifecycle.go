package internal

import (
	"github.com/faelmori/golife/components/process_input"
	pr "github.com/faelmori/golife/internal/process"
	ev "github.com/faelmori/golife/internal/routines/taskz/events"
	"os/signal"
	"syscall"

	p "github.com/faelmori/golife/components/types"
	st "github.com/faelmori/golife/internal/routines/taskz/stage"
	gl "github.com/faelmori/golife/logger"
	"os/exec"

	l "github.com/faelmori/logz"

	"fmt"
	"os"
)

type ILifeCycle[T any] interface {
	Initialize(pi *process_input.ProcessInput[T]) error
	Snapshot() string
	RegisterBatch(category string, items map[string]any) error
	GetController(category, name string) (any, error)
	EnsureController(category, name string, initializer func() any) error
	StartProcess(proc pr.IManagedProcess[any]) error
	GetProcess(name string) (pr.IManagedProcess[any], error)
	AddProcess(name string, process pr.IManagedProcess[any]) error
	CurrentProcess() string
	RemoveProcess(name string) error
	RegisterProcess(name string, command string, args []string, restart bool, customFn func(obj any) *p.ValidationResult) error
	GetStage(name string) (st.IStage[any], error)
	AddStage(name string, stage st.IStage[any]) error
	CurrentStage() string
	SetCurrentStage(stageName string)
	RemoveStage(name string) error
	RegisterStage(stage st.IStage[any]) error
	GetEvent(name string) (ev.IManagedProcessEvents[any], error)
	AddEvent(name string, event ev.IManagedProcessEvents[any]) error
	RemoveEvent(name string) error
	RegisterEvents(eventMap map[string]ev.IManagedProcessEvents[any]) error
	RegisterEvent(event, stageName string, callback func(interface{})) error
	StopEvents() error
	StartLifecycle() error
	StopLifecycle()
	Send(name string, msg any) error
	Receive(name string) (any, error)
	Start() error
	Stop() error
	Restart() error
	Status() string
	StartAll() error
	StopAll() error
	ListenForSignals() error
	StopListeningForSignals()
}
type LifeCycle[T any] struct {
	// Logger is the Logger instance for the lifecycle manager
	Logger l.Logger

	// Mutexes is the mutex for the lifecycle manager
	*p.Mutexes
	// Reference is the reference ID and name
	*p.Reference

	// initialized is a boolean that indicates if the lifecycle manager is initialized
	initialized bool

	// controllers is a map of controllers
	controllers map[string]any
}

// NewLifeCycle creates a new LifeCycle instance with the provided Logger.
func NewLifeCycle[T any](input *process_input.ProcessInput[T]) ILifeCycle[T] {
	if input == nil {
		gl.Log("error", "Process input is nil")
		return nil
	}
	logger := input.Logger
	lcm := &LifeCycle[T]{
		Logger:    logger,
		Mutexes:   p.NewMutexes(),
		Reference: p.NewReference("LifeCycle"),
	}

	if err := lcm.Initialize(input); err != nil {
		gl.LogObjLogger(lcm, "error", fmt.Sprintf("Error initializing lifecycle manager: %v", err))
		return nil
	}
	if err := lcm.initializeProcess(); err != nil {
		gl.LogObjLogger(lcm, "error", fmt.Sprintf("Error initializing process: %v", err))
		return nil
	}
	if err := lcm.initializeStages(); err != nil {
		gl.LogObjLogger(lcm, "error", fmt.Sprintf("Error initializing stages: %v", err))
		return nil
	}
	if err := lcm.initializeChannels(); err != nil {
		gl.LogObjLogger(lcm, "error", fmt.Sprintf("Error initializing channels: %v", err))
		return nil
	}

	return lcm
}

func (lm *LifeCycle[T]) Initialize(pi *process_input.ProcessInput[T]) error {
	if pi == nil {
		gl.LogObjLogger(lm, "error", "Process input is nil")
		return fmt.Errorf("process input is nil")
	}
	if lm.Mutexes == nil {
		lm.Mutexes = p.NewMutexes()
	}

	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	if lm.Logger == nil {
		lm.Logger = l.GetLogger("GoLife")
	}

	if lm.controllers == nil {
		controllers := map[string]any{
			// mainInput is the process input
			"mainProcess": pi, // main input is the process input
			// processes is a map of processes
			"processes": make(map[string]interface{}), //map[string]pProperty[pr.IManagedProcess[any]]
			// stages is a map of stages
			"stages": make(map[string]interface{}), //map[string]pProperty[st.IStage[any]]
			// events is a map of events.
			"events": make(map[string]interface{}), //map[string]pProperty[ev.IManagedProcessEvents[any]]
			// channels is a map of channels
			"channels": make(map[string]interface{}), //map[string]p.ChannelCtl[any]
			// metadata is a map of metadata.
			"metadata": make(map[string]interface{}),
		}
		lm.controllers = controllers
	}

	lm.initialized = true

	return nil
}
func (lm *LifeCycle[T]) initializeChannels() error {
	if lm.controllers["channels"] == nil {
		gl.LogObjLogger(lm, "error", "Channels are nil")
		return fmt.Errorf("channels are nil")
	}

	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	bufSm, bufMd, bufLg := GetDefaultBufferSizes()
	channels := lm.controllers["channels"].(map[string]interface{})
	channels["chanCtl"] = p.NewChannelCtl[any]("chanCtl", &bufSm, lm.Logger)
	channels["chanProcess"] = p.NewChannelCtl[any]("chanProcess", &bufLg, lm.Logger)
	channels["chanStage"] = p.NewChannelCtl[any]("chanStage", &bufMd, lm.Logger)
	channels["chanEvent"] = p.NewChannelCtl[any]("chanEvent", &bufLg, lm.Logger)
	channels["chanSignal"] = p.NewChannelCtl[os.Signal]("chanSignal", &bufSm, lm.Logger)
	channels["chanDone"] = p.NewChannelCtl[bool]("chanDone", &bufSm, lm.Logger)
	channels["chanExit"] = p.NewChannelCtl[bool]("chanExit", &bufSm, lm.Logger)
	channels["chanError"] = p.NewChannelCtl[error]("chanError", &bufSm, lm.Logger)
	channels["chanMessage"] = p.NewChannelCtl[any]("chanMessage", &bufLg, lm.Logger)
	lm.controllers["channels"] = channels
	return nil
}
func (lm *LifeCycle[T]) initializeProcess() error {
	if lm.controllers["processes"] == nil {
		gl.LogObjLogger(lm, "error", "Processes are nil")
		return fmt.Errorf("processes are nil")
	}

	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	processes := lm.controllers["processes"].(map[string]interface{})
	mainProcessInput := lm.controllers["mainProcess"].(*process_input.ProcessInput[T])
	processes["mainProcess"] = pr.NewManagedProcess[any](mainProcessInput.Name, mainProcessInput.Command, mainProcessInput.Args, mainProcessInput.WaitFor, mainProcessInput.CustomFunc)
	mainProcess := processes["mainProcess"].(*pr.ManagedProcess[any])
	mainProcess.SetProcPid(mainProcessInput.ProcPid)
	mainProcess.SetProcHandle(mainProcessInput.ProcPointer)
	mainProcess.SetCmd(exec.Command(mainProcessInput.Command, mainProcessInput.Args...))
	mainProcess.SetName(mainProcessInput.Name)
	mainProcess.SetCustomFunc(mainProcessInput.CustomFunc)

	return nil
}
func (lm *LifeCycle[T]) initializeStages() error {
	if lm.controllers["stages"] == nil {
		gl.LogObjLogger(lm, "error", "Stages are nil")
		return fmt.Errorf("stages are nil")
	}
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()
	stages := lm.controllers["stages"].(map[string]interface{})
	baseStages := getBaseStages()
	for _, stage := range baseStages {
		if _, ok := stages[stage.ID().String()]; !ok {
			stages[stage.ID().String()] = stage
		}
	}
	lm.controllers["stages"] = stages
	return nil
}
func (lm *LifeCycle[T]) initializeEvents(scope string) error {
	if lm.controllers["events"] == nil {
		gl.LogObjLogger(lm, "error", "Events are nil")
		return fmt.Errorf("events are nil")
	}
	err := lm.RegisterEvents(getBaseEvents())
	if err != nil {
		gl.LogObjLogger(lm, "error", err.Error())
	}
	return nil
}

func (lm *LifeCycle[T]) Snapshot() string {
	lm.Mutexes.RLock()
	defer lm.Mutexes.RUnlock()

	return fmt.Sprintf(
		"Lifecycle [%s]: Processes: %d, Stages: %d, Events: %d, Channels: %d",
		lm.Reference.Name,
		len(lm.controllers["processes"].(map[string]interface{})),
		len(lm.controllers["stages"].(map[string]interface{})),
		len(lm.controllers["events"].(map[string]interface{})),
		len(lm.controllers["channels"].(map[string]interface{})),
	)
}
func (lm *LifeCycle[T]) RegisterBatch(category string, items map[string]any) error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	if lm.controllers[category] == nil {
		lm.controllers[category] = make(map[string]interface{})
	}

	categoryMap := lm.controllers[category].(map[string]interface{})
	for name, item := range items {
		if _, exists := categoryMap[name]; exists {
			continue
		}
		categoryMap[name] = item
		gl.LogObjLogger(lm, "info", fmt.Sprintf("Registered '%s' in category '%s'", name, category))
	}
	return nil
}

func (lm *LifeCycle[T]) GetController(category, name string) (any, error) {
	lm.Mutexes.RLock()
	defer lm.Mutexes.RUnlock()

	if categoryMap, exists := lm.controllers[category]; !exists {
		return nil, fmt.Errorf("category '%s' not found", category)
	} else {
		if controller, exists := categoryMap.(map[string]interface{})[name]; !exists {
			return nil, fmt.Errorf("controller '%s' not found in category '%s'", name, category)
		} else {
			return controller, nil
		}
	}
}
func (lm *LifeCycle[T]) EnsureController(category, name string, initializer func() any) error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	if lm.controllers[category] == nil {
		lm.controllers[category] = make(map[string]interface{})
	}

	categoryMap := lm.controllers[category].(map[string]interface{})
	if _, exists := categoryMap[name]; !exists {
		categoryMap[name] = initializer()
		gl.LogObjLogger(lm, "info", fmt.Sprintf("Initialized controller '%s' in category '%s'", name, category))
	}

	return nil
}

func (lm *LifeCycle[T]) StartProcess(proc pr.IManagedProcess[any]) error {
	if err := proc.Start(); err != nil {
		gl.LogObjLogger(lm, "error", fmt.Sprintf("ErrorCtx starting process %s: %v", proc.String(), err))
		return err
	}
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Process %s started successfully!", proc.String()))
	return nil
}
func (lm *LifeCycle[T]) GetProcess(name string) (pr.IManagedProcess[any], error) {
	lm.Mutexes.RLock()
	defer lm.Mutexes.RUnlock()

	processes := lm.controllers["processes"].(map[string]interface{})
	if process, exists := processes[name]; exists {
		return process.(pr.IManagedProcess[any]), nil
	}
	return nil, fmt.Errorf("process '%s' does not exist", name)
}
func (lm *LifeCycle[T]) AddProcess(name string, process pr.IManagedProcess[any]) error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	processes := lm.controllers["processes"].(map[string]interface{})
	if _, exists := processes[name]; exists {
		return fmt.Errorf("process '%s' already exists", name)
	}

	processes[name] = process
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Added process '%s'", name))
	return nil
}
func (lm *LifeCycle[T]) CurrentProcess() string {
	lm.Mutexes.RLock()
	defer lm.Mutexes.RUnlock()

	if currentProcess, exists := lm.controllers["currentProcess"]; exists {
		if process, ok := currentProcess.(pr.IManagedProcess[any]); ok {
			return process.GetName()
		}
	}

	return ""
}
func (lm *LifeCycle[T]) RemoveProcess(name string) error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	processes := lm.controllers["processes"].(map[string]interface{})
	if _, exists := processes[name]; !exists {
		return fmt.Errorf("process '%s' does not exist", name)
	}

	delete(processes, name)
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Removed process '%s'", name))
	return nil
}
func (lm *LifeCycle[T]) RegisterProcess(name string, command string, args []string, restart bool, customFn func(obj any) *p.ValidationResult) error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	processes := lm.controllers["processes"].(map[string]interface{})
	if _, exists := processes[name]; exists {
		return fmt.Errorf("process '%s' already exists", name)
	}

	processes[name] = pr.NewManagedProcess[any](name, command, args, restart, nil)
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Process '%s' registered successfully!", name))
	return nil
}

func (lm *LifeCycle[T]) GetStage(name string) (st.IStage[any], error) {
	lm.Mutexes.RLock()
	defer lm.Mutexes.RUnlock()

	stages := lm.controllers["stages"].(map[string]interface{})
	if stage, exists := stages[name]; exists {
		return stage.(st.IStage[any]), nil
	}
	return nil, fmt.Errorf("stage '%s' does not exist", name)
}
func (lm *LifeCycle[T]) AddStage(name string, stage st.IStage[any]) error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	stages := lm.controllers["stages"].(map[string]interface{})
	if _, exists := stages[name]; exists {
		return fmt.Errorf("stage '%s' already exists", name)
	}

	stages[name] = stage
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Added stage '%s'", name))
	return nil
}
func (lm *LifeCycle[T]) CurrentStage() string {
	lm.Mutexes.RLock()
	defer lm.Mutexes.RUnlock()

	if currentStage, exists := lm.controllers["currentStage"]; exists {
		if stage, ok := currentStage.(st.IStage[any]); ok {
			return stage.Name()
		}
	}
	return ""
}
func (lm *LifeCycle[T]) SetCurrentStage(stageName string) {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	stages := lm.controllers["stages"].(map[string]interface{})
	if stageObj, exists := stages[stageName]; !exists {
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Stage '%s' does not exist", stageName))
		return
	} else {
		if stage, ok := stageObj.(st.IStage[any]); ok {
			if stage.CanTransitionTo(stageName) {
				lm.controllers["currentStage"] = stage
				gl.LogObjLogger(lm, "info", fmt.Sprintf("Transitioned to stage '%s'", stageName))
			}
		} else {
			gl.LogObjLogger(lm, "error", fmt.Sprintf("Failed to cast stage '%s'", stageName))
			return
		}

	}

	gl.LogObjLogger(lm, "info", fmt.Sprintf("Current stage set to '%s'", stageName))
}
func (lm *LifeCycle[T]) RemoveStage(name string) error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	stages := lm.controllers["stages"].(map[string]interface{})
	if _, exists := stages[name]; !exists {
		return fmt.Errorf("stage '%s' does not exist", name)
	}

	delete(stages, name)
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Removed stage '%s'", name))
	return nil
}
func (lm *LifeCycle[T]) RegisterStage(stage st.IStage[any]) error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	stages := lm.controllers["stages"].(map[string]interface{})
	if _, exists := stages[stage.Name()]; exists {
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Stage '%s' already registered!", stage.Name()))
		return fmt.Errorf("stage '%s' already registered", stage.Name())
	}

	stages[stage.Name()] = stage
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Stage '%s' registered successfully!", stage.Name()))
	return nil
}

func (lm *LifeCycle[T]) GetEvent(name string) (ev.IManagedProcessEvents[any], error) {
	lm.Mutexes.RLock()
	defer lm.Mutexes.RUnlock()

	events := lm.controllers["events"].(map[string]interface{})
	if event, exists := events[name]; exists {
		return event.(ev.IManagedProcessEvents[any]), nil
	}
	return nil, fmt.Errorf("event '%s' does not exist", name)
}
func (lm *LifeCycle[T]) AddEvent(name string, event ev.IManagedProcessEvents[any]) error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	events := lm.controllers["events"].(map[string]interface{})
	if _, exists := events[name]; exists {
		return fmt.Errorf("event '%s' already exists", name)
	}

	events[name] = event
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Added event '%s'", name))
	return nil
}
func (lm *LifeCycle[T]) RemoveEvent(name string) error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	events := lm.controllers["events"].(map[string]interface{})
	if _, exists := events[name]; !exists {
		return fmt.Errorf("event '%s' does not exist", name)
	}

	delete(events, name)
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Removed event '%s'", name))
	return nil
}
func (lm *LifeCycle[T]) RegisterEvents(eventMap map[string]ev.IManagedProcessEvents[any]) error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	if lm.controllers["events"] == nil {
		lm.controllers["events"] = make(map[string]interface{})
	}

	events := lm.controllers["events"].(map[string]interface{})
	for name, event := range eventMap {
		if _, exists := events[name]; exists {
			return fmt.Errorf("event '%s' already exists", name)
		}
		events[name] = event
		gl.LogObjLogger(lm, "info", fmt.Sprintf("Added event '%s'", name))
	}

	return nil
}
func (lm *LifeCycle[T]) RegisterEvent(event, stageName string, callback func(interface{})) error {
	if stage, stageErr := lm.GetStage(stageName); stageErr != nil {
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Stage %s not found!", stageErr))
		return fmt.Errorf("stage %s not found", stageErr)
	} else {
		lm.Mutexes.Lock()
		defer lm.Mutexes.Unlock()

		gl.LogObjLogger(lm, "info", fmt.Sprintf("Stage %s found!", stage.Name()))
		stage.OnEvent(event, callback)
	}
	gl.LogObjLogger(lm, "info", fmt.Sprintf("Event %s registered successfully in stage %s!", event, stageName))
	return nil
}
func (lm *LifeCycle[T]) StopEvents() error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	events := lm.controllers["events"].(map[string]interface{})
	for name, event := range events {
		if evnt, ok := event.(ev.IManagedProcessEvents[any]); ok {
			if err := evnt.StopAll(); err != nil {
				gl.LogObjLogger(lm, "error", fmt.Sprintf("Error stopping event '%s': %v", name, err))
				return err
			}
			gl.LogObjLogger(lm, "info", fmt.Sprintf("Event '%s' stopped successfully!", name))
		} else {
			gl.LogObjLogger(lm, "error", fmt.Sprintf("Failed to cast event '%s'", name))
			return fmt.Errorf("failed to cast event '%s'", name)
		}
	}
	gl.LogObjLogger(lm, "info", "All events stopped successfully!")
	return nil
}

func (lm *LifeCycle[T]) StartLifecycle() error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	gl.LogObjLogger(lm, "info", "Starting lifecycle")
	if channels, ok := lm.controllers["channels"].(map[string]interface{}); ok {
		for name, ch := range channels {
			if ctl, ok := ch.(p.ChannelCtl[any]); ok {
				ctl.InitCtl() // Presumes that the method InitCtl() exists.
				gl.LogObjLogger(lm, "info", fmt.Sprintf("Initialized channel '%s'", name))
			}
		}
	}

	if processes, ok := lm.controllers["processes"].(map[string]interface{}); ok {
		for _, proc := range processes {
			if prc, ok := proc.(*pr.ManagedProcess[any]); ok {
				if err := prc.Start(); err != nil {
					return fmt.Errorf("error starting process '%s': %v", prc.Name, err)
				} // Presumes that the method Start() exists.
				gl.LogObjLogger(lm, "info", fmt.Sprintf("Started process '%s'", prc.Name))
			}
		}
	}

	return nil
}
func (lm *LifeCycle[T]) StopLifecycle() {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	gl.LogObjLogger(lm, "info", "Stopping lifecycle")
	if processes, ok := lm.controllers["processes"].(map[string]interface{}); ok {
		for _, proc := range processes {
			if prc, ok := proc.(*pr.ManagedProcess[any]); ok {
				if err := prc.Stop(); err != nil {
					gl.LogObjLogger(lm, "error", fmt.Sprintf("Error stopping process '%s': %v", prc.Name, err))
					return
				} // Presumes that the method Stop() exists.
				gl.LogObjLogger(lm, "info", fmt.Sprintf("Stopped process '%s'", prc.Name))
			}
		}
	}

	if channels, ok := lm.controllers["channels"].(map[string]interface{}); ok {
		for name, ch := range channels {
			if ctl, ok := ch.(p.ChannelCtl[any]); ok {
				ctl.Close() // Presume que o método Close() existe.
				gl.LogObjLogger(lm, "info", fmt.Sprintf("Closed channel '%s'", name))
			}
		}
	}
}

func (lm *LifeCycle[T]) Send(name string, msg any) error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	if channels, ok := lm.controllers["channels"].(map[string]interface{}); ok {
		if ch, exists := channels[name]; exists {
			if ctl, ok := ch.(p.ChannelCtl[any]); ok {
				ctl.Channel() <- msg // Presumes that the method Channel() exists.
				gl.LogObjLogger(lm, "info", fmt.Sprintf("Sent message to channel '%s'", name))
				return nil
			}
		}
	}
	return fmt.Errorf("channel '%s' does not exist", name)
}
func (lm *LifeCycle[T]) Receive(name string) (any, error) {
	lm.Mutexes.RLock()
	defer lm.Mutexes.RUnlock()

	if channels, ok := lm.controllers["channels"].(map[string]interface{}); ok {
		if ch, exists := channels[name]; exists {
			if ctl, ok := ch.(p.ChannelCtl[any]); ok {
				msg := <-ctl.Channel() // Presumes that the method Channel() exists.
				gl.LogObjLogger(lm, "info", fmt.Sprintf("Received message from channel '%s'", name))
				return msg, nil
			}
		}
	}
	return nil, fmt.Errorf("channel '%s' does not exist", name)
}

func (lm *LifeCycle[T]) Start() error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	processes := lm.controllers["processes"].(map[string]interface{})
	for name, proc := range processes {
		if prc, ok := proc.(pr.IManagedProcess[any]); ok {
			if err := prc.Start(); err != nil {
				gl.LogObjLogger(lm, "error", fmt.Sprintf("Error starting process '%s': %v", name, err))
				return err
			}
			gl.LogObjLogger(lm, "info", fmt.Sprintf("Process '%s' started successfully!", name))
		}
	}
	return nil
}
func (lm *LifeCycle[T]) Stop() error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	processes := lm.controllers["processes"].(map[string]interface{})
	for name, proc := range processes {
		if prc, ok := proc.(pr.IManagedProcess[any]); ok {
			if err := prc.Stop(); err != nil {
				gl.LogObjLogger(lm, "error", fmt.Sprintf("Error starting process '%s': %v", name, err))
				return err
			}
			gl.LogObjLogger(lm, "info", fmt.Sprintf("Process '%s' started successfully!", name))
		}
	}
	return nil
}
func (lm *LifeCycle[T]) Restart() error {
	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

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
func (lm *LifeCycle[T]) Status() string {
	lm.Mutexes.RLock()
	defer lm.Mutexes.RUnlock()

	processes := lm.controllers["processes"].(map[string]interface{})
	status := "Processes Status:\n"
	for name, proc := range processes {
		if prc, ok := proc.(pr.IManagedProcess[any]); ok {
			status += fmt.Sprintf("Process '%s': %s\n", name, prc.Status())
		}
	}
	return status
}

func (lm *LifeCycle[T]) StartAll() error {
	if err := lm.Start(); err != nil {
		return err
	}

	if err := lm.StartLifecycle(); err != nil {
		return err
	}

	gl.LogObjLogger(lm, "info", "Lifecycle started successfully!")
	return nil
}
func (lm *LifeCycle[T]) StopAll() error {
	//for name, proc := range lm.processes {
	//	if err := proc.Stop(); err != nil {
	//		gl.LogObjLogger(lm, "error", fmt.Sprintf("ErrorCtx stopping process %s: %v", name, err))
	//		return err
	//	} else {
	//		gl.LogObjLogger(lm, "info", fmt.Sprintf("Process %s stopped successfully!", name))
	//		delete(lm.processes, name)
	//	}
	//}
	gl.LogObjLogger(lm, "info", fmt.Sprintf("All processes stopped successfully!"))
	return nil
}

func (lm *LifeCycle[T]) ListenForSignals() error {
	if lm.controllers["channels"] == nil {
		gl.LogObjLogger(lm, "error", "Channels are nil")
		return fmt.Errorf("channels are nil")
	}

	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	channels := lm.controllers["channels"].(map[string]interface{})
	if sigChan, ok := channels["chanSignal"].(p.ChannelCtl[os.Signal]); ok {
		chanSignal := sigChan.Channel()
		signal.Notify(chanSignal, os.Interrupt, os.Kill, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

		go func(chanSignal chan os.Signal) {
			for sig := range chanSignal {
				gl.LogObjLogger(lm, "info", fmt.Sprintf("Received signal: %s", sig))
				sigChan.Channel() <- sig
			}
		}(chanSignal)

		gl.LogObjLogger(lm, "info", fmt.Sprintf("Listening for signals on channel '%s'", sigChan.Name))
	} else {
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Failed to cast signal channel"))
		return fmt.Errorf("failed to cast signal channel")
	}
	return nil
}
func (lm *LifeCycle[T]) StopListeningForSignals() {
	if lm.controllers["channels"] == nil {
		gl.LogObjLogger(lm, "error", "Channels are nil")
		return
	}

	lm.Mutexes.Lock()
	defer lm.Mutexes.Unlock()

	channels := lm.controllers["channels"].(map[string]interface{})
	if sigChan, ok := channels["chanSignal"].(p.ChannelCtl[os.Signal]); ok {
		sigChan.Close()
		gl.LogObjLogger(lm, "info", fmt.Sprintf("Signal channel closed successfully!"))
	} else {
		gl.LogObjLogger(lm, "error", fmt.Sprintf("Failed to cast signal channel"))
	}
}

func GetDefaultBufferSizes() (sm, md, lg int) {
	return 2, 5, 10
}
func getBaseStages() map[string]st.IStage[any] {
	return map[string]st.IStage[any]{
		"init":  st.NewStage[any]("init", "Initialization stage", "base", nil),
		"start": st.NewStage[any]("start", "Start stage", "base", nil),
		"end":   st.NewStage[any]("end", "End stage", "base", nil),
	}
}
func getBaseEvents() map[string]ev.IManagedProcessEvents[any] {
	return map[string]ev.IManagedProcessEvents[any]{
		"start": ev.NewManagedProcessEvents[any](),
		"stop":  ev.NewManagedProcessEvents[any](),
	}
}
