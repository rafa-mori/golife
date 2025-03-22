package api

import (
	i "github.com/faelmori/golife/internal"
)

type Event = i.IManagedProcessEvents

func NewEvent(eventFns map[string]func(interface{}), triggerCh chan interface{}) Event {
	return i.NewManagedProcessEvents(eventFns, triggerCh)
}

func Trigger(lc LifeCycleManager, stage string, event string, data interface{}) {
	lc.Trigger(stage, event, data)
}

func RegisterEvent(lc LifeCycleManager, stage string, event string, fn func(interface{})) error {
	return lc.RegisterEvent(stage, event, fn)
}

func RegisterStage(lc LifeCycleManager, stage Stage) error {
	return lc.RegisterStage(stage)
}

func RegisterProcess(lc LifeCycleManager, process ManagedProcess) error {
	return lc.RegisterProcess(process.GetName(), process.GetCommand(), process.GetArgs(), process.WillRestart(), process.GetCustomFunc())
}
