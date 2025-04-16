package layers

import (
	"fmt"
	"github.com/faelmori/golife/internal/types"
	"github.com/faelmori/golife/internal/utils"
)

type EventLayer struct {
	utils.ThreadingConfig
	types.Telemetry

	Scope     string
	Events    map[string]func(...any) error
	Listeners []utils.BasicGenericCallback[any]
}

func NewEventLayer(scope string) *EventLayer {
	return &EventLayer{
		ThreadingConfig: *utils.NewThreadingConfig(),
		Telemetry:       *types.NewTelemetry(),
		Scope:           scope,
		Events:          make(map[string]func(...any) error),
	}
}

func (e *EventLayer) AddEvent(name string, event func(...any) error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Events[name] = event
}
func (e *EventLayer) RemoveEvent(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.Events, name)
}

func (e *EventLayer) GetEvent(name string) (func(...any) error, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	event, exists := e.Events[name]
	return event, exists
}
func (e *EventLayer) ExecuteEvent(name string, args ...any) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	event, exists := e.Events[name]
	if !exists {
		return fmt.Errorf("event %s not found", name)
	}
	return event(args...)
}

func (e *EventLayer) SetScope(scope string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Scope = scope
}
func (e *EventLayer) GetScope() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Scope
}

func (e *EventLayer) GetEvents() map[string]func(...any) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Events
}
func (e *EventLayer) SetEvents(events map[string]func(...any) error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Events = events
}

func (e *EventLayer) GetListeners() map[string]func(...any) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Events
}
