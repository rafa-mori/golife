package workers

import (
	"fmt"
	"github.com/faelmori/golife/internal/routines/chan"
	t "github.com/faelmori/golife/internal/types"
	l "github.com/faelmori/logz"
	"sync"
)

type Worker struct {
	t.IWorker
	mu             sync.RWMutex
	wg             sync.WaitGroup
	logger         l.Logger
	ID             int
	Properties     map[string]t.Property[any]
	JobChannel     _chan.IChannel[t.IJob, int]    // Canal de trabalho do worker
	ResultChannel  _chan.IChannel[t.IResult, int] // Canal de resultados do worker
	JobQueue       _chan.IChannel[t.IAction, int] // Canal de trabalho do worker
	StopChannel    chan struct{}                  // Canal de parada do worker
	JobQueueSize   int                            // Tamanho do buffer para os canais (Max 100)
	JobChannelSize int                            // Tamanho do buffer para os canais (Max 100)
}

// NewWorker cria um novo Worker com propriedades genéricas
func NewWorker(workerID int, logger l.Logger) t.IWorker {
	if logger == nil {
		logger = l.GetLogger("Kubex")
	}
	w := &Worker{
		mu:             sync.RWMutex{},
		wg:             sync.WaitGroup{},
		logger:         logger,
		ID:             workerID,
		Properties:     make(map[string]t.Property[any]),
		JobChannel:     _chan.NewChannel[t.IJob, int]("jobChannel", nil, 100),
		ResultChannel:  _chan.NewChannel[t.IResult, int]("resultChannel", nil, 100),
		JobQueue:       _chan.NewChannel[t.IAction, int]("jobQueue", nil, 100),
		StopChannel:    make(chan struct{}, 5),
		JobQueueSize:   100,
		JobChannelSize: 100,
	}

	w.Properties["status"] = t.NewProperty[string]("status", nil)
	_ = w.Properties["status"].SetValue("Stopped", nil)

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
	if status, ok := w.Properties["status"]; ok {
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

	_ = w.Properties["status"].SetValue("Running", nil)
	go func(wk t.IWorker) {

		w.wg.Add(1)

		defer func(wkk t.IWorker) {
			if r := recover(); r != nil {
				wkk.Logger().ErrorCtx(fmt.Sprintf("Recovered from panic: %v", r), nil)
				_ = w.Properties["status"].SetValue("Stopped", nil)
			} else {
				w.logger.InfoCtx("Worker stopped", nil)

				w.JobChannel.StopSysMonitor()
				w.ResultChannel.StopSysMonitor()
				w.JobQueue.StopSysMonitor()

				_ = w.JobChannel.Close()
				_ = w.ResultChannel.Close()
				_ = w.JobQueue.Close()

				w.wg.Done()
			}
		}(wk)

		defer w.wg.Done()
		defer close(w.StopChannel)

		for {
			iJob, _ := w.JobChannel.GetChan()
			iRes, _ := w.ResultChannel.GetChan()
			select {
			case job := <-iJob:
				if job == nil {
					w.logger.ErrorCtx("Job channel closed", nil)
					continue
				} else {
					jj := job.(t.IJob)
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
			case <-w.StopChannel:
				w.logger.InfoCtx("Worker stopped", nil)
				_ = w.Properties["status"].SetValue("Stopped", nil)
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

	_ = w.Properties["status"].SetValue("Stopped", nil)
	close(w.StopChannel)
	w.wg.Wait()
}
func (w *Worker) HandleJob(job t.IJob) error {
	// Handle the job here
	return nil
}
func (w *Worker) HandleResult(result t.IResult) error {
	// Handle the result here
	return nil
}
func (w *Worker) GetStopChannel() chan struct{}                    { return w.StopChannel }
func (w *Worker) GetJobChannel() _chan.IChannel[t.IJob, int]       { return w.JobChannel }
func (w *Worker) GetJobQueue() _chan.IChannel[t.IAction, int]      { return w.JobQueue }
func (w *Worker) GetResultChannel() _chan.IChannel[t.IResult, int] { return w.ResultChannel }
