package workers

import (
	"fmt"
	"github.com/faelmori/golife/internal/routines/agents"
	t "github.com/faelmori/golife/internal/types"
	"github.com/faelmori/golife/services"
	l "github.com/faelmori/logz"
	"sync"
)

type Worker struct {
	t.IWorker
	ID int

	logger l.Logger

	mu   sync.RWMutex
	muL  sync.RWMutex
	wg   sync.WaitGroup
	cond *sync.Cond

	properties map[string]t.Property[any]
	// The size is implicitly defined with the new instance of the interface IChannel.
	//
	// Definition of the channel:
	//
	//		IChannel[T - chan Type, N - buffer size]
	//
	jobChannel    services.IChannel[t.IJob[any], int]    // Canal de trabalho do worker,
	resultChannel services.IChannel[t.IResult, int]      // Canal de resultados do worker
	jobQueue      services.IChannel[t.IAction[any], int] // Canal de trabalho do worker

	stopChannel chan struct{} // Canal de parada do worker
}

// NewWorker cria um novo Worker com propriedades genéricas
func NewWorker(workerID int, logger l.Logger) t.IWorker {
	// Create a new logger if none is provided
	if logger == nil {
		logger = l.GetLogger("Kubex")
	}

	// Create a new worker with the provided ID and logger
	w := &Worker{
		ID: workerID,

		logger: logger,

		mu: sync.RWMutex{},
		wg: sync.WaitGroup{},

		properties:    make(map[string]t.Property[any]),
		jobChannel:    agents.NewChannel[t.IJob[any], int]("jobChannel", nil, 100),
		resultChannel: agents.NewChannel[t.IResult, int]("resultChannel", nil, 100),
		jobQueue:      agents.NewChannel[t.IAction[any], int]("jobQueue", nil, 100),
		stopChannel:   make(chan struct{}, 2),
	}

	w.cond = sync.NewCond(func(wb *Worker) *sync.RWMutex {
		wb.muL = sync.RWMutex{}
		return &wb.muL
	}(w))

	w.properties["status"] = t.NewProperty[string]("status", nil)
	_ = w.properties["status"].SetValue("Stopped", nil)

	return w
}

func (w *Worker) GetWorkerID() int {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.ID
}
func (w *Worker) GetStatus() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if status, ok := w.properties["status"]; ok {
		return status.GetValue().(string)
	} else {
		return "Unknown"
	}
}
func (w *Worker) StartWorkers() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.GetStatus() == "Running" {
		return
	}

	_ = w.properties["status"].SetValue("Running", nil)
	go func(wk t.IWorker) {

		w.wg.Add(1)

		defer func(wkk t.IWorker) {
			if r := recover(); r != nil {
				wkk.Logger().ErrorCtx(fmt.Sprintf("Recovered from panic: %v", r), nil)
				_ = w.properties["status"].SetValue("Stopped", nil)
			} else {
				w.logger.InfoCtx("Worker stopped", nil)

				w.jobChannel.StopSysMonitor()
				w.resultChannel.StopSysMonitor()
				w.jobQueue.StopSysMonitor()

				_ = w.jobChannel.Close()
				_ = w.resultChannel.Close()
				_ = w.jobQueue.Close()

				w.wg.Done()
			}
		}(wk)

		defer w.wg.Done()
		defer close(w.stopChannel)

		for {
			iJob, _ := w.jobChannel.GetChan()
			iRes, _ := w.resultChannel.GetChan()
			select {
			case job := <-iJob:
				if job == nil {
					w.logger.ErrorCtx("Job channel closed", nil)
					continue
				} else {
					jj := job.(t.IJob[any])
					if err := w.HandleJob(jj); err != nil {
						w.logger.ErrorCtx(fmt.Sprintf("Error handling job: %v", err), nil)
					}
					if jj.CanExecute() {
						if err := jj.Execute(); err != nil {
							w.logger.ErrorCtx(fmt.Sprintf("Error executing job: %v", err), nil)
						}
					} else {
						w.logger.ErrorCtx("Job cannot be executed", nil)
					}

				}
			case result := <-iRes:
				if result == nil {
					w.logger.ErrorCtx("Result channel closed", nil)
					continue
				}
				res := result.(t.IResult)
				if err := w.HandleResult(res); err != nil {
					w.logger.ErrorCtx(fmt.Sprintf("Error handling result: %v", err), nil)
				}
				//if res
			case <-w.stopChannel:
				w.logger.InfoCtx("Worker stopped", nil)
				_ = w.properties["status"].SetValue("Stopped", nil)
				w.wg.Done()

				return
			}
		}
	}(w)
}
func (w *Worker) StopWorkers() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.GetStatus() != "Running" {
		return
	}

	_ = w.properties["status"].SetValue("Stopped", nil)
	close(w.stopChannel)
	w.wg.Wait()
}
func (w *Worker) HandleJob(job t.IJob[any]) error {
	// Handle the job here
	return nil
}
func (w *Worker) HandleResult(result t.IResult) error {
	// Handle the result here
	return nil
}
func (w *Worker) GetStopChannel() chan struct{}                       { return w.stopChannel }
func (w *Worker) GetJobChannel() services.IChannel[t.IJob[any], int]  { return w.jobChannel }
func (w *Worker) GetJobQueue() services.IChannel[t.IAction[any], int] { return w.jobQueue }
func (w *Worker) GetResultChannel() services.IChannel[t.IResult, int] { return w.resultChannel }
