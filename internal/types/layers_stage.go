package types

import "fmt"

type StageLayer struct {
	ThreadingConfig
	Telemetry

	Scope          string
	StageEvents    map[string]func(...any) error
	StageListeners []GenericChannelCallback[any]
}

func NewStageLayer(scope string) *StageLayer {
	return &StageLayer{
		ThreadingConfig: *NewThreadingConfig(),
		Telemetry:       *NewTelemetry(),
		Scope:           scope,
		StageEvents:     make(map[string]func(...any) error),
	}
}

func (s *StageLayer) AddStageEvent(name string, event func(...any) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.StageEvents[name] = event
}

func (s *StageLayer) ExecuteStageEvent(name string, args ...any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, exists := s.StageEvents[name]
	if !exists {
		return fmt.Errorf("stage event %s not found", name)
	}
	return event(args...)
}
