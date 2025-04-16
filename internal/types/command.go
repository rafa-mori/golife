package types

import (
	"github.com/google/uuid"
)

type CommandConfig struct {
	// Telemetry configuration
	Telemetry
	// Threading configuration
	ThreadingConfig
	// ID and Reference
	ID uuid.UUID
	// Command Properties
	CmdProperties map[string]Property[any]
}

func NewCommandConfig() *CommandConfig {
	return &CommandConfig{
		Telemetry:       *NewTelemetry(),
		ThreadingConfig: *NewThreadingConfig(),
		ID:              uuid.New(),
		CmdProperties:   make(map[string]Property[any]),
	}
}
