package types

import (
	t "github.com/faelmori/kubex-interfaces/types"
	"github.com/google/uuid"
)

type ManagedEventsConfig struct {
	// Telemetry configuration
	Telemetry
	// Threading configuration
	ThreadingConfig
	// ID and Reference
	ID uuid.UUID
	// Event Properties
	EventProperties map[string]t.Property[any]
	// Event Functions
	EventFuncList func(func(...any) error)
}

func NewManagedEventsConfig() *ManagedEventsConfig {
	return &ManagedEventsConfig{
		Telemetry:       *NewTelemetry(),
		ThreadingConfig: *NewThreadingConfig(),
		ID:              uuid.New(),
		EventProperties: make(map[string]t.Property[any]),
	}
}
