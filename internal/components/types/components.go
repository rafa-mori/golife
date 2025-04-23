package types

import (
	ci "github.com/faelmori/golife/internal/components/interfaces"
	l "github.com/faelmori/logz"
)

type Components[T ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]] struct {
	// Logger is the logger for this GoLife instance.
	l.Logger `json:"logger" yaml:"logger" xml:"logger" toml:"logger" gorm:"logger"`

	// Reference is the reference ID and name.
	*Reference

	// ProcessManager is the process manager for this GoLife instance.
	*ProcessManager[T] `json:"process_manager" yaml:"process_manager" xml:"process_manager" toml:"process_manager" gorm:"process_manager"`
	// IProcessManager is the process manager interface for this GoLife instance.
	//processManager ci.IProcessManager[T]

	// StageManager is the stage manager for this GoLife instance.
	*StageManager `json:"stage_manager" yaml:"stage_manager" xml:"stage_manager" toml:"stage_manager" gorm:"stage_manager"`
	// IStageManager is the stage manager interface for this GoLife instance.
	//stageManager ci.IStageManager

	// EventManager is the event manager for this GoLife instance.
	*EventManager `json:"event_manager" yaml:"event_manager" xml:"event_manager" toml:"event_manager" gorm:"event_manager"`
	// IEventManager is the event manager interface for this GoLife instance.
	//eventManager ci.IEventManager

	// SignalManager is the signal manager for this GoLife instance.
	*SignalManager `json:"signal_manager" yaml:"signal_manager" xml:"signal_manager" toml:"signal_manager" gorm:"signal_manager"`
	// ISignalManager is the signal manager interface for this GoLife instance.
	//signalManager ci.ISignalManager
}

func newComponents[T ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]](logger l.Logger) *Components[T] {
	return &Components[T]{
		Logger: logger,

		Reference: newReference("Components"),

		ProcessManager: newProcessManager[T](),
		//processManager: NewProcessManager[T](),

		StageManager: newStageManager(),
		//stageManager: NewStageManager(),

		EventManager: newEventManager(),
		//eventManager: NewEventManager(),
	}
}

func NewComponents[T ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]](logger l.Logger) ci.IComponents[T] {
	return newComponents[T](logger)
}

func (c *Components[T]) GetComponent(name string) (any, bool) {
	switch name {
	case "process_manager":
		return c.ProcessManager, true
	case "stage_manager":
		return c.StageManager, true
	case "event_manager":
		return c.EventManager, true
	case "signal_manager":
		return c.SignalManager, true
	default:
		return nil, false
	}
}
