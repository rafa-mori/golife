package types

import (
	c "github.com/faelmori/golife/internal/routines/chan"
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

	// Manager Properties
	ManagerProperties map[string]Property[any]
	// Manager Agents
	ManagerAgentsMap map[string]c.IChannel[any, int]
	// Manager Stages
	ManagerStagesMap map[string]StageConfig
}

func NewManagerConfig() *ManagerConfig {
	return &ManagerConfig{
		Telemetry:         *NewTelemetry(),
		ThreadingConfig:   *NewThreadingConfig(),
		ID:                uuid.New(),
		ManagerProperties: make(map[string]Property[any]),
		ManagerAgentsMap:  make(map[string]c.IChannel[any, int]),

		ManagerStagesMap: make(map[string]StageConfig),

		//ManagedRoutineMap: make(map[string]*r.ManagedGoroutine),
		//ManagedProcessMap: make(map[string]*p.ManagedProcess),
		//ManagedEventsMap:  make(map[string]*e.ManagedProcessEvents[any]),
	}
}

//func RegisterProcess(manager *ManagerConfig, process *p.ManagedProcess) {
//	manager.ManagedProcessMap[process.Name] = process
//}

//
//func (m *ManagerConfig) RegisterLayerEvent(layerScope, eventName string, event func(...any) error) error {
//	layer, exists := m.ManagerStagesMap[layerScope]
//	if !exists {
//		return fmt.Errorf("layer %s not found", layerScope)
//	}
//	evConfig := *NewManagedEventsConfig()
//
//	layer.StageEventMap[eventName] = evConfig.EventFuncList
//	return nil
//}
