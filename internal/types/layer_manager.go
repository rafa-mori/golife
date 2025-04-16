package types

import (
	"fmt"
	"sync"
)

type ManagerLayer struct {
	ThreadingConfig
	Telemetry

	Scope            string
	ManagerEvents    map[string]func(...any) error
	ManagerListeners []GenericChannelCallback[any]
}

func NewManagerLayer(scope string) *ManagerLayer {
	return &ManagerLayer{
		ThreadingConfig: *NewThreadingConfig(),
		Telemetry:       *NewTelemetry(),
		Scope:           scope,
		ManagerEvents:   make(map[string]func(...any) error),
	}
}

func (m *ManagerLayer) AddManagerEvent(name string, event func(...any) error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ManagerEvents[name] = event
}

func (m *ManagerLayer) ExecuteManagerEvent(name string, args ...any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	event, exists := m.ManagerEvents[name]
	if !exists {
		return fmt.Errorf("manager event %s not found", name)
	}
	return event(args...)
}

type LayerManager struct {
	layers map[string]any // Um mapa para armazenar diferentes tipos de layers
	mu     sync.RWMutex
}

// NewLayerManager cria uma nova instância do LayerManager
func NewLayerManager() *LayerManager {
	return &LayerManager{
		layers: make(map[string]any),
	}
}

// AddLayer registra uma camada no LayerManager
func (lm *LayerManager) AddLayer(name string, layer any) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	lm.layers[name] = layer
}

// GetLayer recupera uma camada pelo nome
func (lm *LayerManager) GetLayer(name string) (any, error) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	layer, exists := lm.layers[name]
	if !exists {
		return nil, fmt.Errorf("layer %s não encontrada", name)
	}
	return layer, nil
}

// ExecuteLayerEvent executa um evento em uma camada específica
func (lm *LayerManager) ExecuteLayerEvent(layerName, eventName string, args ...any) error {
	layer, err := lm.GetLayer(layerName)
	if err != nil {
		return err
	}

	switch l := layer.(type) {
	case *EventLayer:
		return l.ExecuteEvent(eventName, args...)
	case *StageLayer:
		return l.ExecuteStageEvent(eventName, args...)
	case *AccessLayer:
		return l.ExecuteAccessEvent(eventName, args...)
	default:
		return fmt.Errorf("tipo de camada desconhecido para %s", layerName)
	}
}
