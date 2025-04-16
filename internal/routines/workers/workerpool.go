package workers

import (
	"fmt"
	"github.com/faelmori/golife/internal/routines/agents"
	t "github.com/faelmori/golife/internal/types"
	u "github.com/faelmori/golife/internal/utils"
	"github.com/faelmori/golife/services"
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
	"sync"
)

type WorkerPool struct {
	t.IWorkerPool
	mu         sync.RWMutex
	wg         sync.WaitGroup
	logger     l.Logger
	ID         string
	Properties map[string]t.Property[any]
	workers    []t.IWorker // Referência aos workers gerenciados pelo pool

	// Channels

	jobChannel  services.IChannel[t.IJob[any], int]    // Canal de trabalho do pool
	jobQueue    services.IChannel[t.IAction[any], int] // Canal de trabalho do pool
	resultQueue services.IChannel[t.IResult, int]      // Canal de resultados do pool
	doneChannel chan struct{}                          // Canal de resultados do pool
}

// NewWorkerPool cria um novo WorkerPool com propriedades genéricas
func NewWorkerPool(workerLimit int, logger l.Logger) t.IWorkerPool {
	if logger == nil {
		logger = l.GetLogger("Kubex")
	}
	var iJob t.IJob[any]
	var iResult t.IResult
	var iAction t.IAction[any]
	wp := &WorkerPool{
		mu:          sync.RWMutex{},
		wg:          sync.WaitGroup{},
		logger:      logger,
		ID:          uuid.NewString(),
		Properties:  make(map[string]t.Property[any]),
		workers:     make([]t.IWorker, workerLimit),
		jobQueue:    agents.NewChannel[t.IAction[any], int]("jobQueue", &iAction, 100),
		jobChannel:  agents.NewChannel[t.IJob[any], int]("jobChannel", &iJob, 100),
		resultQueue: agents.NewChannel[t.IResult, int]("resultQueue", &iResult, 100),
		doneChannel: make(chan struct{}, 5),
	}

	// Control
	wp.Properties["workerLimit"] = t.NewProperty[int]("workerLimit", nil)
	wp.Properties["workerLimit"].SetValue(workerLimit, nil)

	wp.Properties["workerCount"] = t.NewProperty[int]("workerCount", nil)
	wp.Properties["workerCount"].SetValue(0, nil)

	wp.Properties["buffers"] = t.NewProperty[int]("buffers", nil) // Tamanho do buffer para os canais (Max 100)
	_ = wp.Properties["buffers"].SetValue(100, nil)

	// Validator
	if addValidatorErr := wp.Properties["workerLimit"].AddValidator("workerLimit", u.ValidateWorkerLimit); addValidatorErr != nil {
		wp.logger.ErrorCtx("Erro ao adicionar validador para workerLimit", map[string]any{
			"context":  "WorkerPool",
			"action":   "AddValidator",
			"error":    addValidatorErr,
			"showData": true,
		})
		if setValErr := wp.Properties["workerLimit"].SetValue(0, nil); setValErr != nil {
			wp.logger.ErrorCtx("Erro ao definir o valor padrão para workerLimit", map[string]any{
				"context":  "WorkerPool",
				"action":   "SetValue",
				"error":    setValErr,
				"showData": true,
			})
		}
		return nil
	}

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
func (wp *WorkerPool) GetPoolJobChannel() (services.IChannel[t.IJob[any], int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if wp.jobChannel != nil {
		return wp.jobChannel, nil
	}
	return nil, fmt.Errorf("failed to get job channel")
}

// GetPoolResultChannel retorna o canal de resultados do pool
func (wp *WorkerPool) GetPoolResultChannel() (services.IChannel[t.IResult, int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if wp.resultQueue != nil {
		return wp.resultQueue, nil
	}
	return nil, fmt.Errorf("failed to get result channel")
}

// GetJobQueue retorna o canal de trabalho do pool
func (wp *WorkerPool) GetJobQueue(workerID int) (services.IChannel[t.IAction[any], int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if workerID < 0 || workerID >= len(wp.workers) {
		return nil, fmt.Errorf("worker ID out of range")
	}
	if wp.workers[workerID] == nil {
		return nil, fmt.Errorf("worker not found")
	}
	if wp.workers[workerID].GetJobQueue() != nil {
		return wp.workers[workerID].GetJobQueue(), nil
	}
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
	return wp.Properties["workerLimit"].GetValue().(int)
}

// GetWorker retorna um worker específico do pool
func (wp *WorkerPool) GetWorker(workerID int) (t.IWorker, error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if workerID < 0 || workerID >= len(wp.workers) {
		return nil, fmt.Errorf("worker ID out of range")
	}
	return wp.workers[workerID], nil
}

// GetWorkerChannel retorna o canal de trabalho de um worker específico
func (wp *WorkerPool) GetWorkerChannel(workerID int) (services.IChannel[t.IJob[any], int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if workerID < 0 || workerID >= len(wp.workers) {
		return nil, fmt.Errorf("worker ID out of range")
	}
	return wp.workers[workerID].GetJobChannel(), nil
}

// GetResultChannel retorna o canal de resultados de um worker específico
func (wp *WorkerPool) GetResultChannel(workerID int) (services.IChannel[t.IResult, int], error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if workerID < 0 || workerID >= len(wp.workers) {
		return nil, fmt.Errorf("worker ID out of range")
	}
	return wp.workers[workerID].GetResultChannel(), nil
}

// GetWorkerPool retorna o pool de workers
func (wp *WorkerPool) GetWorkerPool() []t.IWorker {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.workers
}

// Debug imprime informações de depuração sobre o WorkerPool
func (wp *WorkerPool) Debug() {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	fmt.Printf("WorkerPool ID: %s\n", wp.ID)
	fmt.Printf("WorkerCount: %d | WorkerLimit: %d\n",
		len(wp.workers), wp.Properties["workerLimit"].GetValue())
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

	jobCh := wp.workers[workerID].GetJobChannel()

	return jobCh.Send(job)
}

// Report gera um relatório do estado do Wo-rkerPool
func (wp *WorkerPool) Report() string {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	report := fmt.Sprintf("WorkerPool Report\nWorkerCount: %d | WorkerLimit: %d\n",
		len(wp.workers), wp.Properties["workerLimit"].GetValue())
	for i, worker := range wp.workers {
		report += fmt.Sprintf("Worker %d | Status: %v\n", i, worker.GetStatus())
	}
	return report
}

// AddListener adiciona um listener a um evento específico
func (wp *WorkerPool) AddListener(event string, listener t.ChangeListener[any]) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if property, ok := wp.Properties[event]; ok {
		var ltn t.ChangeListener[any] = func(oldValue, newValue any, metadata t.EventMetadata) t.ListenerResponse {
			return listener(oldValue, newValue, metadata)
		}
		if err := property.AddListener(event, ltn); err != nil {
			return err
		}
	} else {
		wp.logger.ErrorCtx("Event not found", map[string]any{
			"context": "WorkerPool",
			"event":   event,
		})
	}
	return fmt.Errorf("event %s not found", event)
}

// RemoveListener remove um listener de um evento específico
func (wp *WorkerPool) RemoveListener(event string) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if property, ok := wp.Properties[event]; ok {
		if err := property.RemoveListener(event); err != nil {
			return err
		}
	} else {
		wp.logger.ErrorCtx("Event not found", map[string]any{
			"context": "WorkerPool",
			"event":   event,
		})
	}
	return fmt.Errorf("event %s not found", event)
}

// AddWorker adiciona um novo worker ao pool
func (wp *WorkerPool) AddWorker(workerID int, worker t.IWorker) error {
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
	if err := wp.Properties["workerLimit"].SetValue(limit, nil); err != nil {
		return err
	}
	return nil
}

// SetWorkerPool define o pool de workers
func (wp *WorkerPool) SetWorkerPool(workerPool []t.IWorker) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if workerPool == nil {
		return fmt.Errorf("worker pool cannot be nil")
	}
	if len(workerPool) > wp.Properties["workerLimit"].GetValue().(int) {
		return fmt.Errorf("worker pool exceeds worker limit")
	}
	wp.workers = workerPool
	return nil
}

// getChannel retorna um canal específico do WorkerPool
func (wp *WorkerPool) getChannel(key string) (any, error) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	if ch, ok := wp.Properties[key].GetValue().(chan any); ok {
		return ch, nil
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
func (wp *WorkerPool) getWorkerChannel(workerID int, channelFunc func(t.IWorker) services.IChannel[t.IJob[any], int]) (services.IChannel[t.IJob[any], int], error) {
	if err := wp.validateWorkerID(workerID); err != nil {
		return nil, err
	}
	return channelFunc(wp.workers[workerID]), nil
}
