package types

import "fmt"

type BaseLayer struct {
	ThreadingConfig
	Telemetry

	Scope     string
	Events    map[string]func(...any) error
	Listeners []GenericChannelCallback[any]
}

func (b *BaseLayer) AddEvent(name string, event func(...any) error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Events[name] = event
}

func (b *BaseLayer) ExecuteEvent(name string, args ...any) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	event, exists := b.Events[name]
	if !exists {
		return fmt.Errorf("event %s not found", name)
	}
	return event(args...)
}
