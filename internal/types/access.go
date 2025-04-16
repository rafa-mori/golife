package types

import (
	"github.com/faelmori/golife/internal/property"
	"github.com/faelmori/golife/internal/utils"
	c "github.com/faelmori/golife/services"
	"github.com/google/uuid"
)

type AccessConfig struct {
	// Telemetry configuration
	Telemetry

	// Threading configuration
	utils.ThreadingConfig

	// ID and Reference
	ID uuid.UUID
	// Access Properties
	AccessProperties map[string]property.Property[any]
	// Access Agents
	AccessAgents map[string]c.IChannel[any, int]
}

func NewAccessConfig() *AccessConfig {
	return &AccessConfig{
		Telemetry:        *NewTelemetry(),
		ThreadingConfig:  *utils.NewThreadingConfig(),
		ID:               uuid.New(),
		AccessProperties: make(map[string]property.Property[any]),
		AccessAgents:     make(map[string]c.IChannel[any, int]),
	}
}
