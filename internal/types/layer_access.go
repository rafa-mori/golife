package types

import "fmt"

type AccessLayer struct {
	ThreadingConfig
	Telemetry

	Scope           string
	AccessEvents    map[string]func(...any) error
	AccessListeners []GenericChannelCallback[any]
}

func NewAccessLayer(scope string) *AccessLayer {
	return &AccessLayer{
		ThreadingConfig: *NewThreadingConfig(),
		Telemetry:       *NewTelemetry(),
		Scope:           scope,
		AccessEvents:    make(map[string]func(...any) error),
	}
}

func (a *AccessLayer) AddAccessEvent(name string, event func(...any) error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.AccessEvents[name] = event
}

func (a *AccessLayer) ExecuteAccessEvent(name string, args ...any) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	event, exists := a.AccessEvents[name]
	if !exists {
		return fmt.Errorf("access event %s not found", name)
	}
	return event(args...)
}
