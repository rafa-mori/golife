package interfaces

import (
	l "github.com/faelmori/logz"
)

type IProcessManager[T IProcessInput[any]] interface {
	GetProcess(string) (IManagedProcess[any], error)
	AddProcess(string, *T) error
	CurrentProcess() IManagedProcess[any]
	RemoveProcess(string) error
}

type IStageManager interface {
	GetStage(string) (IStage[any], error)
	AddStage(string, IStage[any]) error
	CurrentStage() IStage[any]
	SetCurrentStage(string) (IStage[any], error)
	RemoveStage(string) error
}

type IEventManager interface {
	GetEvent(string) (IManagedProcessEvents[any], error)
	AddEvent(string, IManagedProcessEvents[any]) error
	TriggerEvent(string, interface{}) error
	RemoveEvent(string) error
}

type ISignalManager interface {
	ListenForSignals() error
}

type ILifeCycle[T IProcessInput[any]] interface {
	IMutexes

	Reference() IReference

	GetConfig() *T
	SetConfig(*T)
	ValidateConfig() error

	GetLogger() l.Logger
	SetLogger(l.Logger)

	Initialize() error

	IProcessManager[T] // Gerencia processos no Lifecycle
	IStageManager      // Gerencia stages e suas transições
	IEventManager      // Gerencia eventos
	ISignalManager     // Gerencia sinais do sistema

	StartLifecycle() error
	StopLifecycle() error
	RestartLifecycle() error
	StatusLifecycle() string
}
