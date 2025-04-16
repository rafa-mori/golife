package types

import (
	c "github.com/faelmori/golife/services"
)

type RoutineConfig struct {
	// Telemetry configuration
	Telemetry
	// Threading configuration
	ThreadingConfig
	// ID and Reference
	ID string
	// Routine Functions
	RoutineFuncList func(func(...any) error)
	// Routine Properties
	RoutineProperties map[string]Property[any]
	// Routine Agents
	RoutineAgents map[string]c.IChannel[any, int]
	// Routine Event Map
	RoutineEventMap map[string]EventsConfig
	// Routine Command Map
	RoutineCommandMap map[string]CommandConfig
	// Routine Access Map
	RoutineAccessMap map[string]AccessConfig
}

func NewRoutineConfig() *RoutineConfig {
	return &RoutineConfig{
		Telemetry:         *NewTelemetry(),
		ThreadingConfig:   *NewThreadingConfig(),
		ID:                "",
		RoutineProperties: make(map[string]Property[any]),
		RoutineAgents:     make(map[string]c.IChannel[any, int]),
		RoutineEventMap:   make(map[string]EventsConfig),
		RoutineCommandMap: make(map[string]CommandConfig),
		RoutineAccessMap:  make(map[string]AccessConfig),
	}
}
