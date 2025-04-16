package types

import (
	"github.com/faelmori/golife/internal/property"
	"github.com/faelmori/golife/internal/utils"
	"github.com/google/uuid"
)

type CommandConfig struct {
	// Telemetry configuration
	Telemetry
	// Threading configuration
	utils.ThreadingConfig
	// ID and Reference
	ID uuid.UUID
	// Command Properties
	CmdProperties map[string]property.Property[any]
}

func NewCommandConfig() *CommandConfig {
	return &CommandConfig{
		Telemetry:       *NewTelemetry(),
		ThreadingConfig: *utils.NewThreadingConfig(),
		ID:              uuid.New(),
		CmdProperties:   make(map[string]property.Property[any]),
	}
}
