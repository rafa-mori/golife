package interfaces

import (
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
)

type IComponents[T IProperty[IProcessInput[IManagedProcess[any]]]] interface {
	IProcessManager[T]          // Gerencia processos no Lifecycle
	IStageManager               // Gerencia stages e suas transições
	IEventManager               // Gerencia eventos
	ISignalManager[chan string] // Gerencia sinais do sistema
	GetComponent(name string) (any, bool)
}

type IProcessManager[T IProperty[IProcessInput[IManagedProcess[any]]]] interface {
	GetProcess(string) (IProcessInput[IManagedProcess[any]], error)
	AddProcess(string, T) error
	CurrentProcess() IProcessInput[IManagedProcess[any]]
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

type ISignalManager[T chan string] interface {
	ListenForSignals() error
	StopListening()
}

type ILifeCycle[T any, P IProperty[IProcessInput[T]]] interface {
	IMutexes

	GetID() uuid.UUID
	GetName() string

	GetConfig() *P
	SetConfig(P)
	ValidateConfig() error

	GetLogger() l.Logger
	SetLogger(l.Logger)

	// AddComponent(name string, component any)
	GetComponent(name string) (any, bool)
	// RemoveComponent(name string) error

	Initialize() error
	ValidateLifecycle() IValidationResult
	StartLifecycle() error
	StopLifecycle() error
	RestartLifecycle() error
	StatusLifecycle() string

	Shutdown() error
}
