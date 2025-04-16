package utils

import "sync"

// IDeferFunc is a function that takes a function and an interface{} and returns a function
type IDeferFunc func(func(), interface{}) func() error

// IDefer is a struct that holds the defer function
type IDefer interface {
	Defer() IDeferFunc
}

// IMutex is a struct that holds the mutexes
type IMutex interface {
	TryLock() bool
	Lock()
	Unlock()
	RLock()
	RUnlock()
	LockFunc(func())
	UnlockFunc(func())
	LockFuncWithArgs(func(interface{}), interface{})
	UnlockFuncWithArgs(func(interface{}), interface{})
}

// ISync is a struct that holds the sync.WaitGroup
type ISync interface {
	Add(delta int)
	Wait()
	Done()
	WaitGroup() *sync.WaitGroup
	WaitGroupAdd(delta int)
	WaitGroupDone()
}

// IThreading is a struct that holds the defer function
type IThreading interface {
	IMutex
	ISync
	IDefer
}

type ThreadingConfig struct {
	// Thread-safe channel for event functions
	Mu sync.RWMutex
	// Thread-safe channel for process functions
	wg sync.WaitGroup
}

func NewThreadingConfig() *ThreadingConfig {
	return &ThreadingConfig{
		Mu: sync.RWMutex{},
		wg: sync.WaitGroup{},
	}
}
