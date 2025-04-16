package layers

import (
	"fmt"
	t "github.com/faelmori/golife/internal/types"
	u "github.com/faelmori/golife/internal/utils"
)

type AccessLayer struct {
	u.ThreadingConfig
	t.Telemetry

	Scope           string
	AccessEvents    map[string]func(...any) error
	AccessListeners []BasicGenericCallback[any]
}

func NewAccessLayer(scope string) *AccessLayer {
	return &AccessLayer{
		ThreadingConfig: *u.NewThreadingConfig(),
		Telemetry:       *t.NewTelemetry(),
		Scope:           scope,
		AccessEvents:    make(map[string]func(...any) error),
	}
}

func (a *AccessLayer) AddAccessEvent(name string, event func(...any) error) {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	a.AccessEvents[name] = event
}

func (a *AccessLayer) ExecuteAccessEvent(name string, args ...any) error {
	a.Mu.RLock()
	defer a.Mu.RUnlock()
	event, exists := a.AccessEvents[name]
	if !exists {
		return fmt.Errorf("access event %s not found", name)
	}
	return event(args...)
}
