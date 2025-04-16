package types

import (
	c "github.com/faelmori/golife/services"
	"github.com/google/uuid"
)

type StageConfig struct {
	// Telemetry configuration
	Telemetry
	// Threading configuration
	ThreadingConfig
	// ID and Reference
	ID uuid.UUID
	// Stage Properties
	StageProperties map[string]Property[any]
	// Stage Agents
	StageAgents map[string]c.IChannel[any, int]
	// Event Map
	EventConfigMap map[string]EventsConfig
}

func NewStageConfig() *StageConfig {
	return &StageConfig{
		Telemetry:       *NewTelemetry(),
		ThreadingConfig: *NewThreadingConfig(),
		ID:              uuid.New(),
		StageProperties: make(map[string]Property[any]),
		StageAgents:     make(map[string]c.IChannel[any, int]),
		EventConfigMap:  make(map[string]EventsConfig),
	}
}
