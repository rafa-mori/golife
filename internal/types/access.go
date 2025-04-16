package types

import (
	c "github.com/faelmori/golife/services"
	"github.com/google/uuid"
)

type AccessConfig struct {
	// Telemetry configuration
	Telemetry

	// Threading configuration
	ThreadingConfig

	// ID and Reference
	ID uuid.UUID
	// Access Properties
	AccessProperties map[string]Property[any]
	// Access Agents
	AccessAgents map[string]c.IChannel[any, int]
}

func NewAccessConfig() *AccessConfig {
	return &AccessConfig{
		Telemetry:        *NewTelemetry(),
		ThreadingConfig:  *NewThreadingConfig(),
		ID:               uuid.New(),
		AccessProperties: make(map[string]Property[any]),
		AccessAgents:     make(map[string]c.IChannel[any, int]),
	}
}
