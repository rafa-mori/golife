package types

import (
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
	EventProperties map[string]Property[any]
	// Event Functions
	EventFuncList func(func(...any) error)
}

func NewManagedEventsConfig() *EventsConfig {
	return &EventsConfig{
		Telemetry:       *NewTelemetry(),
		ThreadingConfig: *NewThreadingConfig(),
		ID:              uuid.New(),
		EventProperties: make(map[string]Property[any]),
	}
}
