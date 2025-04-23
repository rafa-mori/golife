package utils

import (
	"fmt"
	ci "github.com/faelmori/golife/components/interfaces"
	gl "github.com/faelmori/golife/logger"
	"reflect"
)

func chanRoutineCtl[T any](v ci.IChannelCtl[T], chCtl chan string, ch chan T) {
	select {
	case msg := <-chCtl:
		switch msg {
		case "stop":
			if ch != nil {
				gl.Log("debug", "Stopping channel for:", v.GetName(), "ID:", v.GetID().String())
				ch <- v.GetProperty().GetValue()
				return
			}
		case "get":
			if ch != nil {
				gl.Log("debug", "Getting value from channel for:", v.GetName(), "ID:", v.GetID().String())
				ch <- v.GetProperty().GetValue()
			}
		case "set":
			if ch != nil {
				gl.Log("debug", "Waiting for value from channel for:", v.GetName(), "ID:", v.GetID().String())
				nVal := <-ch
				if reflect.ValueOf(nVal).IsValid() {
					if reflect.ValueOf(nVal).CanConvert(reflect.TypeFor[T]()) {
						gl.Log("debug", "Setting value from channel for:", v.GetName(), "ID:", v.GetID().String())
						v.GetProperty().SetValue(&nVal)
					} else {
						gl.Log("error", "Set: invalid type for channel value (", reflect.TypeFor[T]().String(), ")")
					}
				}
			}
		case "save":
			if ch != nil {
				gl.Log("debug", "Saving value from channel for:", v.GetName(), "ID:", v.GetID().String())
				nVal := <-ch
				if reflect.ValueOf(nVal).IsValid() {
					if reflect.ValueOf(nVal).CanConvert(reflect.TypeFor[T]()) {
						gl.Log("debug", "Saving value from channel for:", v.GetName(), "ID:", v.GetID().String())
						v.GetProperty().SetValue(&nVal)
					} else {
						gl.Log("error", "Save: invalid type for channel value (", reflect.TypeFor[T]().String(), ")")
					}
				}
			}
		case "clear":
			if ch != nil {
				gl.Log("debug", "Clearing channel for:", v.GetName(), "ID:", v.GetID().String())
				v.GetProperty().SetValue(nil)
			}
		}
	}
}
func chanRoutineDefer[T any](v ci.IChannelCtl[T], chCtl chan string, ch chan T) {
	if r := recover(); r != nil {
		gl.Log("error", "Recovering from panic in monitor routine for:", v.GetName(), "ID:", v.GetID().String(), "Error:", fmt.Sprintf("%v", r))
		// In recovering from panic, we need to check if the channel is nil.
		// If it is nil, we need to create a new channel.
		if ch == nil {
			gl.Log("debug", "Creating new channel for:", v.GetName(), "ID:", v.GetID().String())
			// If the channel is nil, create a new channel.
			ch = make(chan T, 3)
		}
		if chCtl == nil {
			gl.Log("debug", "Creating new control channel for:", v.GetName(), "ID:", v.GetID().String())
			// If the control channel is nil, create a new control channel.
			chCtl = make(chan string, 2)
		}
	} else {
		gl.Log("debug", "Exiting monitor routine for:", v.GetName(), "ID:", v.GetID().String())
		// When the monitor routine is done, we need to close the channels.
		// If the channel is not nil, close it.
		if ch != nil {
			close(ch)
		}
		if chCtl != nil {
			close(chCtl)
		}
		ch = nil
		chCtl = nil
	}
	// Always check the v mutexes to see if someone is locking the mutex or not on exit (defer).
	if v.MuTryLock() {
		// If the mutex was locked, unlock it.
		v.MuUnlock()
	} else if v.MuTryRLock() {
		// If the mutex was locked, unlock it.
		v.MuRUnlock()
	}
}
func chanRoutineWrapper[T any](v ci.IChannelCtl[T]) {
	gl.Log("debug", "Setting monitor routine for:", v.GetName(), "ID:", v.GetID().String())
	if rawChCtl, chCtlType, chCtlOk := v.GetSubChannelByName("ctl"); !chCtlOk {
		gl.LogObjLogger(&v, "error", "ChannelCtl: no control channel found")
		return
	} else {
		if chCtlType != reflect.TypeOf("string") {
			gl.LogObjLogger(&v, "error", "ChannelCtl: control channel is not a string channel")
			return
		}
		chCtl := reflect.ValueOf(rawChCtl).Interface().(chan string)
		ch := v.GetMainChannel()

		defer chanRoutineDefer[T](v, chCtl, ch)
		for {
			chanRoutineCtl[T](v, chCtl, ch)
			if ch == nil {
				gl.Log("debug", "Channel is nil for:", v.GetName(), "ID:", v.GetID().String(), "Exiting monitor routine")
				break
			}
			if chCtl == nil {
				gl.Log("debug", "Control channel is nil for:", v.GetName(), "ID:", v.GetID().String(), "Exiting monitor routine")
				break
			}
		}
	}
}

func GetDefaultBufferSizes() (sm, md, lg int) { return 2, 5, 10 }
