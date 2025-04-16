package types

import (
	t "github.com/faelmori/golife/internal/types"
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
	AccessProperties map[string]t.Property[any]
	// Access Agents
	AccessAgents map[string]t.IChannel[any, int]
}

func NewAccessConfig() *AccessConfig {
	return &AccessConfig{
		Telemetry:        *NewTelemetry(),
		ThreadingConfig:  *NewThreadingConfig(),
		ID:               uuid.New(),
		AccessProperties: make(map[string]t.Property[any]),
		AccessAgents:     make(map[string]t.IChannel[any, int]),
	}
}
