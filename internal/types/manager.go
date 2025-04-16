package types

import (
	"fmt"
	c "github.com/faelmori/golife/services"
	"github.com/google/uuid"
)

type ManagerConfig struct {
	// Telemetry configuration
	Telemetry
	// Threading configuration
	ThreadingConfig
	// ID and Reference
	ID uuid.UUID

	// MANAGED

	//// Managed Routines
	//ManagedRoutineMap map[string]*r.ManagedGoroutine
	//// Managed Process
	//ManagedProcessMap map[string]*p.ManagedProcess
	//// Managed Events
	//ManagedEventsMap map[string]*e.ManagedProcessEvents[any]

	// MANAGER

	// Manager Routines
	RoutineConfigMap map[string]*RoutineConfig
	// Manager Properties
	ManagerProperties map[string]Property[any]
	// Manager Agents
	ManagerAgentsMap map[string]c.IChannel[any, int]
	// Manager Stages
	StagesConfigMap map[string]*StageConfig
	// Manager Processes
	ProcessConfigMap map[string]*ProcessConfig
	// Manager Events
	EventsConfigMap map[string]*EventsConfig
	// Manager Stages
	StageConfigMap map[string]*StageConfig
}

func NewManagerConfig() *ManagerConfig {
	return &ManagerConfig{
		Telemetry:         *NewTelemetry(),
		ThreadingConfig:   *NewThreadingConfig(),
		ID:                uuid.New(),
		ManagerProperties: make(map[string]Property[any]),
		ManagerAgentsMap:  make(map[string]c.IChannel[any, int]),

		RoutineConfigMap: make(map[string]*RoutineConfig),
		StagesConfigMap:  make(map[string]*StageConfig),
		ProcessConfigMap: make(map[string]*ProcessConfig),
		EventsConfigMap:  make(map[string]*EventsConfig),
		StageConfigMap:   make(map[string]*StageConfig),
	}
}

func RegisterProcess(manager *ManagerConfig, process *ProcessConfig) {
	manager.ProcessConfigMap[process.Name] = process
}

func (m *ManagerConfig) RegisterLayerEvent(layerScope, eventName string, event func(...any) error) error {
	layer, exists := m.StageConfigMap[layerScope]
	if !exists {
		return fmt.Errorf("layer %s not found", layerScope)
	}
	layer.EventConfigMap[eventName] = *NewManagedEventsConfig()
	return nil
}
