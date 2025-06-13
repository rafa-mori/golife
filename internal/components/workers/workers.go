package workers

import (
	ci "github.com/rafa-mori/golife/internal/components/interfaces"
	t "github.com/rafa-mori/golife/internal/components/types"
	l "github.com/rafa-mori/logz"

	"reflect"
	"sync"
)

type Worker struct {
	ci.IWorker
	ID int

	logger l.Logger

	mu   sync.RWMutex
	muL  sync.RWMutex
	wg   sync.WaitGroup
	cond *sync.Cond

	properties map[string]ci.IProperty[any]
	// The size is implicitly defined with the new instance of the interface IChannel.
	//
	// Definition of the channel:
	//
	//		IChannel[T - chan Type, N - buffer size]
	//
	//jobChannel    s.IChannel[ci.IJob[any], int]    // Canal de trabalho do worker,
	resultChannel ci.IChannelCtl[ci.IResult] // Canal de resultados do worker
	//jobQueue      s.IChannel[ci.IAction[any], int] // Canal de trabalho do worker

	stopChannel chan struct{} // Canal de parada do worker
}

// NewWorker cria um novo Worker com propriedades genéricas
func NewWorker(workerID int, logger l.Logger) ci.IWorker {
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

		properties: make(map[string]ci.IProperty[any]),
		//jobChannel:    a.NewChannel[ci.IJob[any], int]("jobChannel", nil, 100),
		resultChannel: t.NewChannelCtl[ci.IResult]("resultChannel", nil),
		//jobQueue:      a.NewChannel[ci.IAction[any], int]("jobQueue", nil, 100),
		stopChannel: make(chan struct{}, 2),
	}

	w.cond = sync.NewCond(func(wb *Worker) *sync.RWMutex {
		wb.muL = sync.RWMutex{}
		return &wb.muL
	}(w))

	//w.properties["status"] = property.NewProperty[string]("status", nil)
	//_ = w.properties["status"].SetValue("Stopped", nil)

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
		return reflect.ValueOf(status).Elem().FieldByName("value").String()
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

	//_ = w.properties["status"].SetValue("Running", nil)
	//go func(wk ci.IWorker) {
	//
	//	w.wg.Add(1)
	//
	//	defer func(wkk ci.IWorker) {
	//		if r := recover(); r != nil {
	//			wkk.Logger().Error(fmt.Sprintf("Recovered from panic: %v", r), nil)
	//			_ = w.properties["status"].SetValue("Stopped", nil)
	//		} else {
	//			w.logger.Info("Worker stopped", nil)
	//
	//			w.jobChannel.StopSysMonitor()
	//			w.resultChannel.StopSysMonitor()
	//			w.jobQueue.StopSysMonitor()
	//
	//			_ = w.jobChannel.Close()
	//			_ = w.resultChannel.Close()
	//			_ = w.jobQueue.Close()
	//
	//			w.wg.Done()
	//		}
	//	}(wk)
	//
	//	defer w.wg.Done()
	//	defer close(w.stopChannel)
	//
	//	for {
	//		iJob, _ := w.jobChannel.GetChan()
	//		iRes, _ := w.resultChannel.GetChan()
	//		select {
	//		case job := <-iJob:
	//			if job == nil {
	//				w.logger.Error("Job channel closed", nil)
	//				continue
	//			} else {
	//				jj := job.(ci.IJob[any])
	//				if err := w.HandleJob(jj); err != nil {
	//					w.logger.Error(fmt.Sprintf("Error handling job: %v", err), nil)
	//				}
	//				if jj.CanExecute() {
	//					if err := jj.Execute(); err != nil {
	//						w.logger.Error(fmt.Sprintf("Error executing job: %v", err), nil)
	//					}
	//				} else {
	//					w.logger.Error("Job cannot be executed", nil)
	//				}
	//
	//			}
	//		case result := <-iRes:
	//			if result == nil {
	//				w.logger.Error("Result channel closed", nil)
	//				continue
	//			}
	//			res := result.(ci.IResult)
	//			if err := w.HandleResult(res); err != nil {
	//				w.logger.Error(fmt.Sprintf("Error handling result: %v", err), nil)
	//			}
	//			//if res
	//		case <-w.stopChannel:
	//			w.logger.Info("Worker stopped", nil)
	//			_ = w.properties["status"].SetValue("Stopped", nil)
	//			w.wg.Done()
	//
	//			return
	//		}
	//	}
	//}(w)
}
func (w *Worker) StopWorkers() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.GetStatus() != "Running" {
		return
	}

	/*_ = w.properties["status"].SetValue("Stopped", nil)*/
	close(w.stopChannel)
	w.wg.Wait()
}
func (w *Worker) HandleJob(job any /*ci.IJob[any]*/) error {
	// Handle the job here
	return nil
}
func (w *Worker) HandleResult(result ci.IResult) error {
	// Handle the result here
	return nil
}
func (w *Worker) GetStopChannel() chan struct{} { return w.stopChannel }
func (w *Worker) GetJobChannel() any/*s.IChannel[ci.IJob[any], int]*/ { return nil /*w.jobChannel*/ }
func (w *Worker) GetJobQueue() any/*s.IChannel[ci.IAction[any], int]*/ { return nil /*w.jobQueue*/ }
func (w *Worker) GetResultChannel() any/*s.IChannel[ci.IResult, int]*/ { return w.resultChannel }
