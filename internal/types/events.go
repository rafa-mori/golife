package types

import (
	f "github.com/faelmori/golife/internal/property"
	"github.com/faelmori/golife/internal/utils"
	"github.com/google/uuid"
)

type EventsConfig struct {
	// Telemetry configuration
	Telemetry
	// Threading configuration
	utils.ThreadingConfig
	// ID and Reference
	ID uuid.UUID
	// Event Properties
	EventProperties map[string]f.Property[any]
	// Event Functions
	EventFuncList func(func(...any) error)
}

func NewManagedEventsConfig() *EventsConfig {
	return &EventsConfig{
		Telemetry:       *NewTelemetry(),
		ThreadingConfig: *utils.NewThreadingConfig(),
		ID:              uuid.New(),
		EventProperties: make(map[string]f.Property[any]),
	}
}
