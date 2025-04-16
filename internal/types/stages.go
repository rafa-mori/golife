package types

import (
	"github.com/faelmori/golife/internal/property"
	"github.com/faelmori/golife/internal/utils"
	c "github.com/faelmori/golife/services"
	"github.com/google/uuid"
)

type StageConfig struct {
	// Telemetry configuration
	Telemetry
	// Threading configuration
	utils.ThreadingConfig
	// ID and Reference
	ID uuid.UUID
	// Stage Properties
	StageProperties map[string]property.Property[any]
	// Stage Agents
	StageAgents map[string]c.IChannel[any, int]
	// Event Map
	EventConfigMap map[string]EventsConfig
}

func NewStageConfig() *StageConfig {
	return &StageConfig{
		Telemetry:       *NewTelemetry(),
		ThreadingConfig: *utils.NewThreadingConfig(),
		ID:              uuid.New(),
		StageProperties: make(map[string]property.Property[any]),
		StageAgents:     make(map[string]c.IChannel[any, int]),
		EventConfigMap:  make(map[string]EventsConfig),
	}
}
