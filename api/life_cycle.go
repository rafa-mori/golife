package api

import (
	i "github.com/faelmori/golife/internal"
	"os"
)

type LifeCycleManager = i.GWebLifeCycleManager

func NewLifecycleManager(processes map[string]i.IManagedProcess, stages map[string]i.IStage, sigChan chan os.Signal, doneChan chan struct{}, events []i.IManagedProcessEvents, eventsCh chan i.IManagedProcessEvents) LifeCycleManager {
	return i.NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh)
}
func NewLifecycleMgrSig() (LifeCycleManager, error) {
	processes := make(map[string]i.IManagedProcess)
	stages := make(map[string]i.IStage)
	sigChan := make(chan os.Signal, 1)
	doneChan := make(chan struct{}, 1)
	eventsCh := make(chan i.IManagedProcessEvents, 1)
	events := make([]i.IManagedProcessEvents, 0)

	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}
func NewLifecycleMgrManual(processes map[string]i.IManagedProcess, stages map[string]i.IStage, sigChan chan os.Signal, doneChan chan struct{}, events []i.IManagedProcessEvents, eventsCh chan i.IManagedProcessEvents) (LifeCycleManager, error) {
	return NewLifecycleManager(processes, stages, sigChan, doneChan, events, eventsCh), nil
}
func NewLifecycleMgrDec() (LifeCycleManager, error) {
	processes := make(map[string]i.IManagedProcess)
	stages := make(map[string]i.IStage)

	return NewLifecycleManager(processes, stages, nil, nil, nil, nil), nil
}
