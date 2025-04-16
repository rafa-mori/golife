package types

import (
	t "github.com/faelmori/kubex-interfaces/types"
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
	CmdProperties map[string]t.Property[any]
	// Command Agents
	CmdAgents map[string]t.IChannel[any, int]
}

func NewCommandConfig() *CommandConfig {
	return &CommandConfig{
		Telemetry:       *NewTelemetry(),
		ThreadingConfig: *NewThreadingConfig(),
		ID:              uuid.New(),
		CmdProperties:   make(map[string]t.Property[any]),
		CmdAgents:       make(map[string]t.IChannel[any, int]),
	}
}
