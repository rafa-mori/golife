package internal

import (
	"testing"
)

func TestTrigger(t *testing.T) {
	manager := NewLifecycleManager(nil, nil, nil, nil, nil, nil)
	manager.RegisterEvent("testEvent", "testStage")

	triggered := false
	manager.Trigger("testStage", "testEvent", nil)
	if !triggered {
		t.Errorf("Expected event to be triggered")
	}
}

func TestRegisterEvent(t *testing.T) {
	manager := NewLifecycleManager(nil, nil, nil, nil, nil, nil)
	err := manager.RegisterEvent("testEvent", "testStage")
	if err != nil {
		t.Errorf("Expected event to be registered, got error: %v", err)
	}
}

func TestRemoveEvent(t *testing.T) {
	manager := NewLifecycleManager(nil, nil, nil, nil, nil, nil)
	manager.RegisterEvent("testEvent", "testStage")
	err := manager.RemoveEvent("testEvent", "testStage")
	if err != nil {
		t.Errorf("Expected event to be removed, got error: %v", err)
	}
}

func TestStopEvents(t *testing.T) {
	manager := NewLifecycleManager(nil, nil, nil, nil, nil, nil)
	manager.RegisterEvent("testEvent", "testStage")
	err := manager.StopEvents()
	if err != nil {
		t.Errorf("Expected all events to be stopped, got error: %v", err)
	}
}
