package workers

import (
	"fmt"
	c "github.com/faelmori/golife/internal/routines/agents"
	t "github.com/faelmori/golife/internal/types"
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
	"sync"
	"time"
)

type WorkerManager[T any] struct {
	t.IWorkerManager[T]
	ID string

	mu   sync.RWMutex
	wg   sync.WaitGroup
	cond *sync.Cond

	logger l.Logger

	Properties map[string]t.Property[any]

	workerPool t.IWorkerPool //t.IWorkerPool
}

// NewWorkerManager cria um novo WorkerManager que gerencia o WorkerPool
func NewWorkerManager[T any](pool t.IWorkerPool, logger l.Logger) t.IWorkerManager[T] {
	if logger == nil {
		logger = l.GetLogger("Kubex")
	}
	wm := &WorkerManager[T]{
		ID: uuid.NewString(),

		logger: logger,

		mu:   sync.RWMutex{},
		wg:   sync.WaitGroup{},
		cond: sync.NewCond(&sync.Mutex{}),

		Properties: make(map[string]t.Property[any]),

		workerPool: pool,
	}

	// Propriedades de controle
	wm.Properties["status"] = t.NewProperty[string]("status", nil)
	wm.Properties["status"].SetValue("Stopped", nil)
	wm.Properties["workerCount"] = t.NewProperty[int]("workerCount", nil)
	wm.Properties["workerCount"].SetValue(0, nil)
	wm.Properties["monitorInterval"] = t.NewProperty[int]("monitorInterval", nil)
	wm.Properties["monitorInterval"].SetValue(500, nil)

	return wm
}

// AddWorker adiciona um worker ao pool
func (wm *WorkerManager[T]) AddWorker() (t.IWorker, error) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if wm.workerPool.GetWorkerPool() == nil {
		if err := wm.workerPool.SetWorkerPool(make([]t.IWorker, 0)); err != nil {
			return nil, err
		}
	}
	if len(wm.workerPool.GetWorkerPool()) >= wm.Properties["workerLimit"].GetValue().(int) {
		return nil, fmt.Errorf("worker limit reached")
	}

	worker := NewWorker(wm.Properties["workerCount"].GetValue().(int), wm.logger)
	if err := wm.workerPool.AddWorker(wm.Properties["workerCount"].GetValue().(int), worker); err != nil {
		return nil, err
	}
	if setValErr := wm.Properties["workerCount"].SetValue(len(wm.workerPool.(*WorkerPool).workers), nil); setValErr != nil {
		return nil, setValErr
	}
	return worker, nil
}

// AddWorkerObj adiciona um worker ao pool
func (wm *WorkerManager[T]) AddWorkerObj(worker t.IWorker) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if wm.workerPool.GetWorkerPool() == nil {
		if err := wm.workerPool.SetWorkerPool(make([]t.IWorker, 0)); err != nil {
			return err
		}
	}
	if worker == nil {
		return fmt.Errorf("worker cannot be nil")
	}

	if len(wm.workerPool.GetWorkerPool()) >= wm.Properties["workerLimit"].GetValue().(int) {
		return fmt.Errorf("worker limit reached")
	}
	if err := wm.workerPool.SetWorkerPool(append(wm.workerPool.GetWorkerPool(), worker)); err != nil {
		return err
	}
	if setValErr := wm.Properties["workerCount"].SetValue(len(wm.workerPool.(*WorkerPool).workers), nil); setValErr != nil {
		return setValErr
	}
	return nil
}

// RemoveWorker remove um worker do pool
func (wm *WorkerManager[T]) RemoveWorker(workerID int) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if workerID < 0 || workerID >= len(wm.workerPool.(*WorkerPool).workers) {
		return fmt.Errorf("worker ID out of range")
	}
	wm.workerPool.(*WorkerPool).workers = append(wm.workerPool.(*WorkerPool).workers[:workerID], wm.workerPool.(*WorkerPool).workers[workerID+1:]...)
	if setValErr := wm.Properties["workerCount"].SetValue(len(wm.workerPool.(*WorkerPool).workers), nil); setValErr != nil {
		return setValErr
	}
	return nil
}

// AddValidator adiciona um validador para a propriedade
func (wm *WorkerManager[T]) AddValidator(name string, validator t.ValidatorFunc[any]) error {
	if _, exists := wm.Properties[name]; exists {
		if addValidatorErr := wm.Properties[name].AddValidator(name, validator); addValidatorErr != nil {
			return addValidatorErr
		}
	} else {
		return fmt.Errorf("property %s does not exist", name)
	}
	return nil
}

// SetWorkerLimit define o limite de workers do pool
func (wm *WorkerManager[T]) SetWorkerLimit(workerLimit int) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	workerPool := wm.workerPool.(*WorkerPool)
	if workerLimit <= 0 {
		return fmt.Errorf("worker limit must be greater than 0")
	}
	if workerLimit < len(wm.workerPool.(*WorkerPool).workers) {
		return fmt.Errorf("worker limit cannot be less than current worker count")
	}
	if setValueErr := workerPool.Properties["workerLimit"].SetValue(workerLimit, nil); setValueErr != nil {
		return setValueErr
	}
	return nil
}

// MonitorWorkers monitora os workers do pool
func (wm *WorkerManager[T]) MonitorWorkers() {
	interval := wm.Properties["monitorInterval"].GetValue().(int)
	go func() {
		for {
			if wm.Properties["status"].GetValue().(string) != "Running" {
				fmt.Println("Worker monitoring stopped.")
				break
			}
			for _, worker := range wm.workerPool.(*WorkerPool).workers {
				fmt.Printf("Worker ID: %d | Status: %v | Jobs: %d\n",
					worker.GetWorkerID(), worker.GetStatus(), worker.GetWorkerID())
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}()
}

// MonitorPool inicia um monitoramento do pool de workers
func (wm *WorkerManager[T]) MonitorPool() chan interface{} {
	if _, exists := wm.Properties["monitorCtl"]; !exists {
		wm.Properties["monitorCtl"] = t.NewProperty[string]("monitorCtl", nil)
		if setValErr := wm.Properties["monitorCtl"].SetValue("Stopped", nil); setValErr != nil {
			wm.logger.ErrorCtx("Failed to set monitor control value", map[string]any{
				"context":  "WorkerManager",
				"action":   "SetValue",
				"error":    setValErr,
				"showData": true,
			})
			return nil
		}
		wm.Properties["monitorCtl"].SetChannel(c.NewChannel[string]("monitorCtl", nil, 5))
	}

	iChanCtl := wm.Properties["monitorCtl"].GetChannel()
	chanCtl, _ := iChanCtl.GetChan()

	commands := map[t.MonitorCommand]func(){
		t.Start: func() {
			fmt.Println("Starting monitor")
			if setValErr := wm.Properties["monitorCtl"].SetValue("Running", nil); setValErr != nil {
				wm.logger.ErrorCtx("Failed to set monitor control value", map[string]any{
					"context":  "WorkerManager",
					"action":   "SetValue",
					"error":    setValErr,
					"showData": true,
				})
				return
			}
			wm.MonitorWorkers()
		},
		t.Restart: func() {
			fmt.Println("Restarting monitor")
			if setValErr := wm.Properties["monitorCtl"].SetValue("Stopping", nil); setValErr != nil {
				wm.logger.ErrorCtx("Failed to set monitor control value", map[string]any{
					"context":  "WorkerManager",
					"action":   "SetValue",
					"error":    setValErr,
					"showData": true,
				})
				return
			}
			wm.MonitorWorkers()
		},
		t.Stop: func() {
			fmt.Println("Stopping monitor")
			if setValErr := wm.Properties["monitorCtl"].SetValue("Stopped", nil); setValErr != nil {
				wm.logger.ErrorCtx("Failed to set monitor control value", map[string]any{
					"context":  "WorkerManager",
					"action":   "SetValue",
					"error":    setValErr,
					"showData": true,
				})
				return
			}
		},
	}

	go func(chanCtl chan any) {
		interval := wm.Properties["monitorInterval"].GetValue().(int)
		for {
			fmt.Printf("Pool Info | WorkerCount: %d | Limit: %d\n",
				len(wm.workerPool.(*WorkerPool).workers),
				wm.Properties["workerLimit"].GetValue().(int))

			select {
			case <-time.After(time.Duration(interval) * time.Millisecond):
				interval = wm.Properties["monitorInterval"].GetValue().(int)
			case msg := <-chanCtl:
				if cmd, ok := commands[t.MonitorCommand(msg.(string))]; ok {
					cmd()
				} else {
					fmt.Printf("Unknown command: %v\n", msg)
				}
			}
		}
	}(chanCtl)

	return nil
}

func (wm *WorkerManager[T]) ValidatePool() error {
	if len(wm.workerPool.(*WorkerPool).workers) > wm.workerPool.(*WorkerPool).Properties["workerLimit"].GetValue().(int) {
		return fmt.Errorf("worker count exceeds worker limit")
	}
	return nil
}
