package internal

import (
	"testing"
)

func TestOnEvent(t *testing.T) {
	stage := NewStage("1", "testStage", "Test Stage", "stage")
	triggered := false
	stage.OnEvent("testEvent", func(data interface{}) {
		triggered = true
	})

	stage.Trigger("testStage", "testEvent", nil)
	if !triggered {
		t.Errorf("Expected event to be triggered")
	}
}

func TestAutoScale(t *testing.T) {
	stage := NewStage("1", "testStage", "Test Stage", "stage")
	stage.AutoScale(5)

	if stage.WorkerPool == nil {
		t.Errorf("Expected worker pool to be initialized")
	}

	if cap(stage.WorkerPool.Tasks) != 5 {
		t.Errorf("Expected worker pool size to be 5, got %d", cap(stage.WorkerPool.Tasks))
	}
}

func TestOnEnter(t *testing.T) {
	stage := NewStage("1", "testStage", "Test Stage", "stage")
	enterCalled := false
	stage.OnEnter(func() {
		enterCalled = true
	})

	stage.OnEnterFn()
	if !enterCalled {
		t.Errorf("Expected OnEnter function to be called")
	}
}

func TestOnExit(t *testing.T) {
	stage := NewStage("1", "testStage", "Test Stage", "stage")
	exitCalled := false
	stage.OnExit(func() {
		exitCalled = true
	})

	stage.OnExitFn()
	if !exitCalled {
		t.Errorf("Expected OnExit function to be called")
	}
}

func TestDispatch(t *testing.T) {
	stage := NewStage("1", "testStage", "Test Stage", "stage")
	stage.AutoScale(1)

	taskCompleted := false
	stage.Dispatch(func() {
		taskCompleted = true
	})

	stage.WorkerPool.Wg.Wait()
	if !taskCompleted {
		t.Errorf("Expected dispatched task to be completed")
	}
}
