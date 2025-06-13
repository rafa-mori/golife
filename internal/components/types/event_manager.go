package types

import (
	"fmt"
	ci "github.com/rafa-mori/golife/internal/components/interfaces"
	gl "github.com/rafa-mori/golife/logger"
)

type EventManager struct {
	events map[string]ci.IManagedProcessEvents[any]
}

func newEventManager() *EventManager {
	return &EventManager{events: make(map[string]ci.IManagedProcessEvents[any])}
}
func NewEventManager() ci.IEventManager { return newEventManager() }

func (em *EventManager) GetEvent(name string) (ci.IManagedProcessEvents[any], error) {
	if event, exists := em.events[name]; exists {
		return event, nil
	}
	gl.LogObjLogger(&em, "error", "Event not found", "event", name)
	return nil, fmt.Errorf("event %s not found", name)
}
func (em *EventManager) AddEvent(name string, eventObj ci.IManagedProcessEvents[any]) error {
	if _, exists := em.events[name]; exists {
		gl.LogObjLogger(&em, "error", "Event already exists", "event", name)
		return fmt.Errorf("event %s already exists", name)
	}
	em.events[name] = eventObj
	return nil
}
func (em *EventManager) TriggerEvent(name string, data interface{}) error {
	// TRIGGER EVENT WITH ALL NECESSARY DATA
	// FOR NOW, JUST LOGGING
	if _, exists := em.events[name]; exists {
		//event.Trigger(data)
		return nil
	}
	gl.LogObjLogger(&em, "error", "Event not found", "event", name)
	return fmt.Errorf("event %s not found", name)
}
func (em *EventManager) RemoveEvent(name string) error {
	if _, exists := em.events[name]; exists {
		delete(em.events, name)
		return nil
	}
	gl.LogObjLogger(&em, "error", "Event not found", "event", name)
	return fmt.Errorf("event %s not found", name)
}
