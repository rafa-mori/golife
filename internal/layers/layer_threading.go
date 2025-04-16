package layers

import (
	t "github.com/faelmori/golife/internal/types"
	"sync"
)

// LayerThreadingConfig is a struct that contains the configuration for threading
type LayerThreadingConfig struct {
	maxThreads int
	maxQueue   int

	// conditionMap is a map of string keys to sync.Cond values.
	conditionMap map[string]*sync.Cond

	// mutexes for synchronization and thread-safety
	muLayer     sync.RWMutex
	muScope     sync.RWMutex
	muEvents    sync.RWMutex
	muListeners sync.RWMutex

	muLayerCond *sync.Cond
}

// LayerThreading is a struct that contains the threading configuration and telemetry
type LayerThreading struct {
	LayerThreadingConfig
	Telemetry
	Scope     string
	Events    map[string]func(...any) error
	Listeners []func(...any) error
}
type Telemetry struct {
	t.Telemetry
}
type BasicGenericCallback[T any] func(T) error

// NewLayerThreading creates a new LayerThreading instance
func NewLayerThreading(scope string) *LayerThreading {
	return &LayerThreading{
		LayerThreadingConfig: LayerThreadingConfig{},
		Telemetry:            Telemetry{Telemetry: *t.NewTelemetry()},
		Scope:                scope,
		Events:               make(map[string]func(...any) error),
	}
}
