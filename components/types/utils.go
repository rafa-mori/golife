package types

import (
	"fmt"
	l "github.com/faelmori/golife/logger"
	"reflect"
	"sync"
)

func monitor[T any](v *val[T]) {
	select {
	case msg := <-v.chCtl:
		switch msg {
		case "stop":
			if v.ch != nil {
				l.Log("debug", "Stopping channel for:", v.Name, "ID:", v.ID.String())
				v.ch <- *v.Load()
				return
			}
		case "get":
			if v.ch != nil {
				l.Log("debug", "Getting value from channel for:", v.Name, "ID:", v.ID.String())
				v.ch <- *v.Load()
			}
		case "set":
			if v.ch != nil {
				l.Log("debug", "Waiting for value from channel for:", v.Name, "ID:", v.ID.String())
				nVal := <-v.ch
				if reflect.ValueOf(nVal).IsValid() {
					if reflect.ValueOf(nVal).CanConvert(reflect.TypeFor[T]()) {
						l.Log("debug", "Setting value from channel for:", v.Name, "ID:", v.ID.String())
						v.Set(&nVal)
					} else {
						l.Log("error", "Set: invalid type for channel value (", reflect.TypeFor[T]().String(), ")")
					}
				}
			}
		case "clear":
			if v.ch != nil {
				l.Log("debug", "Clearing channel for:", v.Name, "ID:", v.ID.String())
				v.Clear()
			}
		}
	}
}
func monitorDefer[T any](v *val[T]) {
	if r := recover(); r != nil {
		l.Log("error", "Recovering from panic in monitor routine for:", v.Name, "ID:", v.ID.String(), "Error:", fmt.Sprintf("%v", r))
		// In recovering from panic, we need to check if the channel is nil.
		// If it is nil, we need to create a new channel.
		if v.ch == nil {
			l.Log("debug", "Creating new channel for:", v.Name, "ID:", v.ID.String())
			// If the channel is nil, create a new channel.
			v.ch = make(chan T, 3)
		}
		if v.chCtl == nil {
			l.Log("debug", "Creating new control channel for:", v.Name, "ID:", v.ID.String())
			// If the control channel is nil, create a new control channel.
			v.chCtl = make(chan string, 2)
		}
		if !v.isCtl {
			l.Log("debug", "Setting control channel for:", v.Name, "ID:", v.ID.String())
			// If the control channel is not set, set it to the control channel.
			v.isCtl = true
		}
	} else {
		l.Log("debug", "Exiting monitor routine for:", v.Name, "ID:", v.ID.String())
		// When the monitor routine is done, we need to close the channels.
		// If the channel is not nil, close it.
		if v.ch != nil {
			close(v.ch)
		}
		if v.chCtl != nil {
			close(v.chCtl)
		}
		v.ch = nil
		v.chCtl = nil
		v.isCtl = false
	}
	// Always set the control channel to the control channel that only receives messages.
	v.ctl = v.chCtl
	// Always check the v mutexes to see if someone is locking the mutex.
	if v.muCtx != nil {
		if v.muCtx.MuCtxM != nil {
			if vMuCtxM := (*sync.RWMutex)(v.muCtx.MuCtxM); vMuCtxM != nil {
				if vMuCtxM.TryLock() {
					// If the mutex was locked, unlock it.
					vMuCtxM.Unlock()
				} else if vMuCtxM.TryRLock() {
					// If the mutex was locked, unlock it.
					vMuCtxM.RUnlock()
				}
				l.Log("debug", "Unlocked mutex for:", v.Name, "ID:", v.ID.String())
			}
		}
		if v.muCtx.MuCtxL != nil {
			if vMuCtxL := (*sync.RWMutex)(v.muCtx.MuCtxL); vMuCtxL != nil {
				if vMuCtxL.TryRLock() {
					vMuCtxL.RUnlock()
				} else if vMuCtxL.TryLock() {
					vMuCtxL.Unlock()
				}
				l.Log("debug", "Unlocked mutex for:", v.Name, "ID:", v.ID.String())
			}
		}
		if v.muCtx.MuCtxCond != nil {
			if vMuCtxCond := v.muCtx.MuCtxCond; vMuCtxCond != nil {
				l.Log("debug", "Broadcasting condition variable for:", v.Name, "ID:", v.ID.String())
				vMuCtxCond.Broadcast()
			}
		}
	}
}
func monitorRoutine[T any](v *val[T]) {
	if v.chCtl == nil {
		l.Log("debug", "Creating new control channel for:", v.Name, "ID:", v.ID.String())
		v.chCtl = make(chan string, 2)
	}
	if v.ch == nil {
		l.Log("debug", "Creating new channel for:", v.Name, "ID:", v.ID.String())
		v.ch = make(chan T, 3)
	}
	if v.ctl == nil {
		l.Log("debug", "Setting control channel for:", v.Name, "ID:", v.ID.String())
		v.ctl = v.chCtl
	}
	if !v.isCtl {
		l.Log("debug", "Setting monitor routine for:", v.Name, "ID:", v.ID.String())
		v.isCtl = true
		defer monitorDefer[T](v)
		for {
			monitor[T](v)
			if v.ch == nil {
				l.Log("debug", "Channel is nil for:", v.Name, "ID:", v.ID.String(), "Exiting monitor routine")
				break
			}
			if v.chCtl == nil {
				l.Log("debug", "Control channel is nil for:", v.Name, "ID:", v.ID.String(), "Exiting monitor routine")
				break
			}
			if v.ctl == nil {
				l.Log("debug", "Control channel is nil for:", v.Name, "ID:", v.ID.String(), "Exiting monitor routine")
				break
			}
		}
	}
}
