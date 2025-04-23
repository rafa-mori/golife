package interfaces

import l "github.com/faelmori/logz"

type IWorker interface {
	Logger() l.Logger
	SetLogger(l.Logger)

	GetWorkerID() int
	GetStatus() string

	StartWorkers()
	StopWorkers()

	//HandleJob(job IJob[any]) error
	//HandleResult(result IResult) error
	//
	//GetStopChannel() chan struct{}
	//
	//GetJobChannel() c.IChannel[IJob[any], int]
	//GetJobQueue() c.IChannel[IAction[any], int]
	//GetResultChannel() c.IChannel[IResult, int]
}
