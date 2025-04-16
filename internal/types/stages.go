package types

import (
	t "github.com/faelmori/golife/internal/types"
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
	StageProperties map[string]t.Property[any]
	// Stage Agents
	StageAgents map[string]t.IChannel[any, int]
	// Event Map
	EventConfigMap map[string]EventsConfig
}

func NewStageConfig() *StageConfig {
	return &StageConfig{
		Telemetry:       *NewTelemetry(),
		ThreadingConfig: *NewThreadingConfig(),
		ID:              uuid.New(),
		StageProperties: make(map[string]t.Property[any]),
		StageAgents:     make(map[string]t.IChannel[any, int]),
		EventConfigMap:  make(map[string]EventsConfig),
	}
}
