package types

import (
	f "github.com/faelmori/golife/internal/property"
	c "github.com/faelmori/golife/services"
	l "github.com/faelmori/logz"
)

type ValidatorFunc[T any] func(value T) error

type MonitorCommand string

const (
	Start   MonitorCommand = "start"
	Stop    MonitorCommand = "stop"
	Restart MonitorCommand = "restart"
)

type IWorker interface {
	Logger() l.Logger
	SetLogger(l.Logger)

	GetWorkerID() int
	GetStatus() string

	StartWorkers()
	StopWorkers()

	HandleJob(job IJob[any]) error
	HandleResult(result IResult) error

	GetStopChannel() chan struct{}

	GetJobChannel() c.IChannel[IJob[any], int]
	GetJobQueue() c.IChannel[IAction[any], int]
	GetResultChannel() c.IChannel[IResult, int]
}

type IWorkerPool interface {
	Logger() l.Logger
	SetLogger(l.Logger)

	GetWorkerCount() int

	GetPoolJobChannel() (c.IChannel[IJob[any], int], error)
	GetPoolResultChannel() (c.IChannel[IResult, int], error)

	GetWorkerLimit() int
	GetWorker(workerID int) (IWorker, error)

	GetWorkerChannel(workerID int) (c.IChannel[IJob[any], int], error)
	GetResultChannel(workerID int) (c.IChannel[IResult, int], error)
	GetJobQueue(workerID int) (c.IChannel[IAction[any], int], error)
	GetResultQueue(workerID int) (c.IChannel[IResult, int], error)
	GetDoneChannel() (chan struct{}, error)

	GetWorkerPool() []IWorker
	SetWorkerPool([]IWorker) error

	RemoveListener(string) error
	AddWorker(int, IWorker) error
	SetWorkerLimit(int) error

	Report() string
	Debug()
	SendToWorker(int, any) error
	AddListener(string, f.ChangeListener[any]) error
}

type IWorkerManager[T any] interface {
	Logger() l.Logger
	SetLogger(l.Logger)

	GetID() string
	GetProperties() map[string]f.Property[any]
	GetWorker(int) (IWorker, error)
	GetWorkerChannel(int) (chan IJob[any], error)
	GetWorkerPool() []IWorker

	SetWorkerPool([]IWorker)
	SetWorkerCount(int) error

	SetWorker(int, IWorker) error
	SetWorkerPoolChannel(int, c.IChannel[IJob[any], int]) error
	SetWorkerChannel(int, c.IChannel[IJob[any], int]) error
	SetWorkerResultChannel(int, c.IChannel[IResult, int]) error
	SetWorkerJobQueue(int, c.IChannel[IAction[any], int]) error
	SetWorkerResultQueue(int, c.IChannel[IResult, int]) error
	SetWorkerStatus(int, string) error
	SetWorkerJobQueueCount(int, int) error

	GetWorkerLimit() int
	GetWorkerCount() int
	GetWorkerStatus() string
	GetWorkerStatusByID(int) string

	GetWorkerPoolInstance() IWorkerPool
	GetWorkerPoolChannel() (c.IChannel[IJob[any], int], error)
	GetWorkerPoolResultChannel() (c.IChannel[IResult, int], error)
	GetWorkerPoolJobQueue() (c.IChannel[IAction[any], int], error)
	GetWorkerPoolResultQueue() (c.IChannel[IResult, int], error)

	AddWorker(worker IWorker) error
	RemoveWorker(workerID int) error
	AddValidator(name string, validator ValidatorFunc[any]) error
	SetWorkerLimit(workerLimit int) error
	MonitorWorkers()
	MonitorPool() chan interface{}
	ValidatePool() error
}
