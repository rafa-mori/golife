package workers

import (
	"fmt"
	ci "github.com/faelmori/golife/components/interfaces"
	t "github.com/faelmori/golife/components/types"
	"github.com/faelmori/golife/internal/property"
	"github.com/faelmori/golife/services"
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
	"reflect"
	"sync"
)

type WorkerPool struct {
	ci.IWorkerPool
	mu         sync.RWMutex
	wg         sync.WaitGroup
	logger     l.Logger
	ID         string
	Properties map[string]any // Propriedades do WorkerPool
	workers    []ci.IWorker   // Referência aos workers gerenciados pelo pool

	// Channels

	//jobChannel  services.IChannel[ci.IJob[any], int]    // Canal de trabalho do pool
	//jobQueue    services.IChannel[ci.IAction[any], int] // Canal de trabalho do pool
	//resultQueue services.IChannel[ci.IResult, int]      // Canal de resultados do pool
	doneChannel chan struct{} // Canal de resultados do pool
}

// NewWorkerPool cria um novo WorkerPool com propriedades genéricas
func NewWorkerPool(workerLimit int, logger l.Logger) ci.IWorkerPool {
	if logger == nil {
		logger = l.GetLogger("Kubex")
	}
	//var iJob ci.IJob[any]
	//var iResult ci.IResult
	//var iAction ci.IAction[any]
	wp := &WorkerPool{
		mu:         sync.RWMutex{},
		wg:         sync.WaitGroup{},
		logger:     logger,
		ID:         uuid.NewString(),
		Properties: make(map[string]any),
		workers:    make([]ci.IWorker, workerLimit),
		//jobQueue:    agents.NewChannel[ci.IAction[any], int]("jobQueue", &iAction, 100),
		//jobChannel:  agents.NewChannel[ci.IJob[any], int]("jobChannel", &iJob, 100),
		//resultQueue: agents.NewChannel[ci.IResult, int]("resultQueue", &iResult, 100),
		doneChannel: make(chan struct{}, 5),
	}

	// Control
	wkrLimit := t.NewProperty[int]("workerLimit", nil, false, nil)
	// Validator
	//if addValidatorErr := wkrLimit.Prop().AddValidator(p.ValidationFunc[int]{
	//	Priority: 0,
	//	Func: func(value *int, args ...any) *p.ValidationResult {
	//		if *value < 0 {
	//			return &p.ValidationResult{
	//				IsValid: false,
	//				Message: "workerLimit cannot be negative",
	//				Error:   fmt.Errorf("workerLimit cannot be negative"),
	//			}
	//		}
	//		if *value > 50 {
	//			return &p.ValidationResult{
	//				IsValid: false,
	//				Message: "workerLimit cannot be greater than 50",
	//				Error:   fmt.Errorf("workerLimit cannot be greater than 50"),
	//			}
	//		}
	//		return &p.ValidationResult{
	//			IsValid: true,
	//			Message: "workerLimit is valid",
	//			Error:   nil,
	//		}
	//
	//	},
	//	Result: nil,
	//}); addValidatorErr != nil {
	//	wp.logger.Error("Erro ao adicionar validador para workerLimit", map[string]any{
	//		"context":  "WorkerPool",
	//		"action":   "AddValidator",
	//		"error":    addValidatorErr,
	//		"showData": true,
	//	})
	//	workerLimit = 0
	//	wkrLimit.Prop.Set(&workerLimit)
	//	return nil
	//}
	wkrLimit.SetValue(&workerLimit)
	wp.Properties["workerLimit"] = wkrLimit

	zero := 0
	wkrCount := t.NewProperty[int]("workerCount", nil, false, nil)
	wkrCount.SetValue(&zero)
	wp.Properties["workerCount"] = wkrCount

	bufferSize := 10
	wkrBuffer := t.NewProperty[int]("workerBuffer", nil, false, nil)
	wkrBuffer.SetValue(&bufferSize)
	wp.Properties["buffers"] = wkrBuffer

	return wp
}

// Logger retorna o logger do WorkerPool
func (wp *WorkerPool) Logger() l.Logger {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.logger
}

// SetLogger define o logger do WorkerPool
func (wp *WorkerPool) SetLogger(logger l.Logger) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if logger == nil {
		logger = l.GetLogger("Kubex")
	}
	wp.logger = logger
}

// GetWorkerCount retorna o número de workers no pool
func (wp *WorkerPool) GetWorkerCount() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return len(wp.workers)
}

// GetPoolJobChannel retorna o canal de trabalho do pool
func (wp *WorkerPool) GetPoolJobChannel() (services.IChannel[any, int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	//if wp.jobChannel != nil {
	//	return wp.jobChannel, nil
	//}
	return nil, fmt.Errorf("failed to get job channel")
}

// GetPoolResultChannel retorna o canal de resultados do pool
func (wp *WorkerPool) GetPoolResultChannel() (services.IChannel[ci.IResult, int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	//if wp.resultQueue != nil {
	//	return wp.resultQueue, nil
	//}
	return nil, fmt.Errorf("failed to get result channel")
}

// GetJobQueue retorna o canal de trabalho do pool
func (wp *WorkerPool) GetJobQueue(workerID int) (services.IChannel[any /*ci.IAction[any]*/, int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if workerID < 0 || workerID >= len(wp.workers) {
		return nil, fmt.Errorf("worker ID out of range")
	}
	if wp.workers[workerID] == nil {
		return nil, fmt.Errorf("worker not found")
	}
	//if wp.workers[workerID].GetJobQueue() != nil {
	//	return wp.workers[workerID].GetJobQueue(), nil
	//}
	return nil, fmt.Errorf("failed to get job queue")
}

// GetDoneChannel retorna o canal de resultados do pool
func (wp *WorkerPool) GetDoneChannel() (chan struct{}, error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if wp.doneChannel != nil {
		return wp.doneChannel, nil
	}
	return nil, fmt.Errorf("failed to get done channel")
}

// GetWorkerLimit retorna o limite de workers do pool
func (wp *WorkerPool) GetWorkerLimit() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	valInt := 0
	if wkrLimit, ok := wp.Properties["workerLimit"]; !ok {
		if wkrLimitValT, ok := reflect.ValueOf(wkrLimit).Interface().(t.Property[int]); ok {
			valInt = wkrLimitValT.GetValue()
		}
	}
	return valInt
}

// GetWorker retorna um worker específico do pool
func (wp *WorkerPool) GetWorker(workerID int) (ci.IWorker, error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if workerID < 0 || workerID >= len(wp.workers) {
		return nil, fmt.Errorf("worker ID out of range")
	}
	return wp.workers[workerID], nil
}

// GetWorkerChannel retorna o canal de trabalho de um worker específico
func (wp *WorkerPool) GetWorkerChannel(workerID int) (services.IChannel[any /*ci.IJob[any]*/, int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if workerID < 0 || workerID >= len(wp.workers) {
		return nil, fmt.Errorf("worker ID out of range")
	}
	//return wp.workers[workerID].GetJobChannel(), nil
	return nil, nil
}

// GetResultChannel retorna o canal de resultados de um worker específico
func (wp *WorkerPool) GetResultChannel(workerID int) (services.IChannel[ci.IResult, int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if workerID < 0 || workerID >= len(wp.workers) {
		return nil, fmt.Errorf("worker ID out of range")
	}
	return nil, nil
}

// GetWorkerPool retorna o pool de workers
func (wp *WorkerPool) GetWorkerPool() []ci.IWorker {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.workers
}

// Debug imprime informações de depuração sobre o WorkerPool
func (wp *WorkerPool) Debug() {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	valInt := 0
	if wkrLimit, ok := wp.Properties["workerLimit"]; !ok {
		if wkrLimitValT, ok := reflect.ValueOf(wkrLimit).Interface().(t.Property[int]); ok {
			valInt = wkrLimitValT.GetValue()
		}
	}

	fmt.Printf("WorkerPool ID: %s\n", wp.ID)
	fmt.Printf("WorkerCount: %d | WorkerLimit: %d\n", len(wp.workers), valInt)
	for i, worker := range wp.workers {
		fmt.Printf("Worker %d | Status: %v\n", i, worker.GetStatus())
	}
}

// SendToWorker envia um trabalho para um worker específico
func (wp *WorkerPool) SendToWorker(workerID int, job any) error {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	if workerID < 0 || workerID >= len(wp.workers) {
		return fmt.Errorf("worker ID out of range")
	}

	//jobCh := wp.workers[workerID].GetJobChannel()

	//return jobCh.Send(job)
	return nil
}

// Report gera um relatório do estado do Wo-rkerPool
func (wp *WorkerPool) Report() string {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	valInt := 0
	if wkrLimit, ok := wp.Properties["workerLimit"]; !ok {
		if wkrLimitValT, ok := reflect.ValueOf(wkrLimit).Interface().(t.Property[int]); ok {
			valInt = wkrLimitValT.GetValue()
		}
	}

	report := fmt.Sprintf("WorkerPool Report\nWorkerCount: %d | WorkerLimit: %d\n", len(wp.workers), valInt)
	for i, worker := range wp.workers {
		report += fmt.Sprintf("Worker %d | Status: %v\n", i, worker.GetStatus())
	}
	return report
}

// AddListener adiciona um listener a um evento específico
func (wp *WorkerPool) AddListener(event string, listener property.ChangeListener[any]) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	//if property, ok := wp.Properties[event]; ok {
	//var ltn property.ChangeListener[any] = func(oldValue, newValue any, metadata property.EventMetadata) property.ListenerResponse {
	//	return listener(oldValue, newValue, metadata)
	//}
	//if err := property.AddListener(event, ltn); err != nil {
	//	return err
	//}
	//} else {
	//	wp.logger.ErrorCtx("Event not found", map[string]any{
	//		"context": "WorkerPool",
	//		"event":   event,
	//	})
	//}
	return fmt.Errorf("event %s not found", event)
}

// RemoveListener remove um listener de um evento específico
func (wp *WorkerPool) RemoveListener(event string) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	//if property, ok := wp.Properties[event]; ok {
	//	if err := property.RemoveListener(event); err != nil {
	//		return err
	//	}
	//} else {
	//	wp.logger.ErrorCtx("Event not found", map[string]any{
	//		"context": "WorkerPool",
	//		"event":   event,
	//	})
	//}
	return fmt.Errorf("event %s not found", event)
}

// AddWorker adiciona um novo worker ao pool
func (wp *WorkerPool) AddWorker(workerID int, worker ci.IWorker) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if workerID < 0 || workerID >= len(wp.workers) {
		return fmt.Errorf("worker ID out of range")
	}
	wp.workers[workerID] = worker
	return nil
}

// SetWorkerLimit define o limite de workers do pool
func (wp *WorkerPool) SetWorkerLimit(limit int) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if limit < 0 {
		return fmt.Errorf("worker limit cannot be negative")
	}
	if wkrLimit, ok := wp.Properties["workerLimit"]; !ok {
		if wkrLimitValT, ok := reflect.ValueOf(wkrLimit).Interface().(t.Property[int]); ok {
			wkrLimitValT.SetValue(&limit)
			wp.Properties["workerLimit"] = wkrLimitValT
		}
	}
	return nil
}

// SetWorkerPool define o pool de workers
func (wp *WorkerPool) SetWorkerPool(workerPool []ci.IWorker) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if workerPool == nil {
		return fmt.Errorf("worker pool cannot be nil")

	}
	valInt := 0
	if wkrLimit, ok := wp.Properties["workerLimit"]; !ok {
		if wkrLimitValT, ok := reflect.ValueOf(wkrLimit).Interface().(t.Property[int]); ok {
			valInt = wkrLimitValT.GetValue()
		}
	}
	if len(workerPool) > valInt {
		return fmt.Errorf("worker pool exceeds worker limit")
	}
	wp.workers = workerPool
	return nil
}

// getChannel retorna um canal específico do WorkerPool
func (wp *WorkerPool) getChannel(key string) (any, error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if _, ok := wp.Properties[key]; !ok {
		return nil, fmt.Errorf("channel %s not found", key)
	}
	if wkrProp, ok := wp.Properties[key]; !ok {
		//if wkrPropT, ok := reflect.ValueOf(wkrProp).Interface().(t.Property[any]); ok {
		if _, ok := reflect.ValueOf(wkrProp).Interface().(t.Property[any]); ok {
			//return wkrPropT.Channel(), nil
			return nil, nil
		} else {
			//return reflect.ValueOf(wkrProp).Interface().(*t.Property[any]).Prop.Channel(), nil
		}
	}
	return nil, fmt.Errorf("failed to get channel %s", key)
}

// validateWorkerID valida o ID do worker
func (wp *WorkerPool) validateWorkerID(workerID int) error {
	if workerID < 0 || workerID >= len(wp.workers) {
		return fmt.Errorf("worker ID %d out of range", workerID)
	}
	return nil
}

// getWorkerChannel retorna o canal de um worker específico
func (wp *WorkerPool) getWorkerChannel(workerID int, channelFunc func(ci.IWorker) services.IChannel[any /*ci.IJob[any]*/, int]) (services.IChannel[any /*ci.IJob[any]*/, int], error) {
	if err := wp.validateWorkerID(workerID); err != nil {
		return nil, err
	}
	return channelFunc(wp.workers[workerID]), nil
}
