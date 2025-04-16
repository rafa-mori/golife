package types

import (
	"github.com/faelmori/golife/internal/property"
	"github.com/faelmori/golife/internal/utils"
	c "github.com/faelmori/golife/services"
)

type RoutineConfig struct {
	// Telemetry configuration
	Telemetry
	// Threading configuration
	utils.ThreadingConfig
	// ID and Reference
	ID string
	// Routine Functions
	RoutineFuncList func(func(...any) error)
	// Routine Properties
	RoutineProperties map[string]property.Property[any]
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
		ThreadingConfig:   *utils.NewThreadingConfig(),
		ID:                "",
		RoutineProperties: make(map[string]property.Property[any]),
		RoutineAgents:     make(map[string]c.IChannel[any, int]),
		RoutineEventMap:   make(map[string]EventsConfig),
		RoutineCommandMap: make(map[string]CommandConfig),
		RoutineAccessMap:  make(map[string]AccessConfig),
	}
}
