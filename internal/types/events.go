package types

import (
	t "github.com/faelmori/golife/internal/types"
	"github.com/google/uuid"
)

type EventsConfig struct {
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

func NewManagedEventsConfig() *EventsConfig {
	return &EventsConfig{
		Telemetry:       *NewTelemetry(),
		ThreadingConfig: *NewThreadingConfig(),
		ID:              uuid.New(),
		EventProperties: make(map[string]t.Property[any]),
	}
}
