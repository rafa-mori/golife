package layers

import (
	"fmt"
	t "github.com/faelmori/golife/internal/types"
	"github.com/faelmori/golife/internal/utils"
	"sync"
)

type StageLayer struct {
	utils.ThreadingConfig
	t.Telemetry

	// Mutex for thread-safe access to events and listeners
	mu          sync.RWMutex
	muListeners sync.RWMutex
	muEvents    sync.RWMutex
	muScope     sync.RWMutex
	// Mutex for thread-safe access to the layer itself
	muLayer sync.RWMutex
	// Mutex Condition for waiting on events
	muLayerCond *sync.Cond

	// Scope for the layer
	// This can be used to identify the layer or its purpose in the application
	// For example, "access", "event", "manager", etc.
	Scope          string
	StageEvents    map[string]func(...any) error
	StageListeners []utils.BasicGenericCallback[any]
}

func NewStageLayer(scope string) *StageLayer {
	return &StageLayer{
		ThreadingConfig: *utils.NewThreadingConfig(),
		Telemetry:       *t.NewTelemetry(),
		Scope:           scope,
		StageEvents:     make(map[string]func(...any) error),
	}
}

func (s *StageLayer) AddStageEvent(name string, event func(...any) error) {
	s.mu.Lock()
	defer s.Unlock()
	s.StageEvents[name] = event
}

func (s *StageLayer) ExecuteStageEvent(name string, args ...any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, exists := s.StageEvents[name]
	if !exists {
		return fmt.Errorf("stage event %s not found", name)
	}
	return event(args...)
}
