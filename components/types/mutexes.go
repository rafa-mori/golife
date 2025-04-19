package types

import "sync"

// muCtx is the mutex context map
type muCtx struct {
	// MuCtxM is a mutex for the ctx map.
	MuCtxM *sync.RWMutex
	// MuCtxL is a mutex for sync.Cond in the ctx map.
	MuCtxL *sync.RWMutex
	// MuCtxCond is a condition variable for the ctx map.
	MuCtxCond *sync.Cond
	// MuCtxWg is a wait group for the ctx map.
	MuCtxWg *sync.WaitGroup
}

// newMuCtx creates a new mutex context map
func newMuCtx() *muCtx {
	mu := &muCtx{
		MuCtxM:    &sync.RWMutex{},
		MuCtxCond: &sync.Cond{},
		MuCtxWg:   &sync.WaitGroup{},
	}
	mu.MuCtxCond = sync.NewCond(func(gll *muCtx) *sync.RWMutex {
		gll.MuCtxL = &sync.RWMutex{}
		return gll.MuCtxL
	}(mu))
	return mu
}

// Mutexes is a struct that holds the mutex context map
type Mutexes struct {
	// muCtx is the mutex context map
	*muCtx
}

// NewMutexes creates a new mutex context map
func NewMutexes() *Mutexes {
	return &Mutexes{newMuCtx()}
}

// Lock locks the mutex
func (m *Mutexes) Lock() { m.MuCtxM.Lock() }

// Unlock unlocks the mutex
func (m *Mutexes) Unlock() { m.MuCtxM.Unlock() }

// RLock locks the mutex for reading
func (m *Mutexes) RLock() { m.MuCtxL.RLock() }

// RUnlock unlocks the mutex for reading
func (m *Mutexes) RUnlock() { m.MuCtxL.RUnlock() }

// WaitCond waits for the condition variable to be signaled
func (m *Mutexes) WaitCond() { m.MuCtxCond.Wait() }

// SignalCond signals the condition variable
func (m *Mutexes) SignalCond() { m.MuCtxCond.Signal() }

// BroadcastCond broadcasts the condition variable
func (m *Mutexes) BroadcastCond() { m.MuCtxCond.Broadcast() }

// Add adds a delta to the wait group counter
func (m *Mutexes) Add(delta int) { m.MuCtxWg.Add(delta) }

// Done signals that the wait group is done
func (m *Mutexes) Done() { m.MuCtxWg.Done() }

// Wait waits for the wait group counter to reach zero
func (m *Mutexes) Wait() { m.MuCtxWg.Wait() }
