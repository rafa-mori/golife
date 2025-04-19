package types

import (
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"reflect"
)

type ChannelCtl[T any] struct {
	// Logger is the Logger instance for this Channel instance.
	Logger l.Logger

	// Reference is the reference ID and name.
	*Reference

	// Mutexes is the mutexes for this Channel instance.
	*Mutexes

	// buffers is the number of buffers for the channel.
	buffers int

	// ch is a channel for the value.
	ch chan T

	// chCtl is a channel for the internal control channel.
	chCtl chan string // chCtl is a channel for the control channel. It receives and sends, but only for the control channel.

	// ctl is the exposed control channel that only receives messages.
	ctl <-chan string

	// isCtl is a boolean that indicates if the control channel is set.
	isCtl bool
}

// NewChannelCtl creates a new ChannelCtl instance with the provided name.
func NewChannelCtl[T any](name string, buffers *int, logger l.Logger) *ChannelCtl[T] {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	ref := NewReference(name)
	mu := NewMutexes()
	buf := 3
	if buffers != nil {
		buf = *buffers
	}

	// Create a new ChannelCtl instance
	return &ChannelCtl[T]{
		Logger:    logger,
		Reference: ref,
		Mutexes:   mu,
		ch:        make(chan T, buf),
		chCtl:     make(chan string, 2),
		isCtl:     false,
	}
}

// InitCtl is a method that initializes the control channel.
func (v *ChannelCtl[T]) InitCtl() {
	if v.chCtl == nil {
		gl.Log("info", "Creating control channel for:", v.Name, "ID:", v.ID.String())
		v.chCtl = make(chan string, 2)
	}
	if v.ctl == nil {
		gl.Log("info", "Creating control channel for:", v.Name, "ID:", v.ID.String())
		v.ctl = v.chCtl
	}
}

// StopCtl is a method that stops the control channel.
func (v *ChannelCtl[T]) StopCtl() {
	gl.LogObjLogger(v, "info", "Stopping control channel for:", v.Name, "ID:", v.ID.String())
	if v.chCtl != nil {
		gl.LogObjLogger(v, "info", "Closing control channel for:", v.Name, "ID:", v.ID.String())
		v.chCtl <- "stop"
	}
}

// Channel is a method that returns the channel for the value.
func (v *ChannelCtl[T]) Channel() chan T {
	if v.ch == nil {
		gl.LogObjLogger(v, "info", "Creating channel for:", v.Name, "ID:", v.ID.String())
		v.ch = make(chan T, v.buffers)
	}
	return v.ch
}

// IsCtl is a method that checks if the control channel is set.
func (v *ChannelCtl[T]) IsCtl() bool {
	if v == nil {
		gl.LogObjLogger(v, "error", "Get: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return false
	}
	return v.isCtl
}

// Close is a method that closes the channel.
func (v *ChannelCtl[T]) Close() {
	if v.ch != nil {
		gl.LogObjLogger(v, "info", "Closing channel for:", v.Name, "ID:", v.ID.String())
		close(v.ch)
	}
	if v.chCtl != nil {
		gl.LogObjLogger(v, "info", "Closing control channel for:", v.Name, "ID:", v.ID.String())
		close(v.chCtl)
	}
}

// ChannelType is a method that returns the type of the channel.
func (v *ChannelCtl[T]) ChannelType() reflect.Type { return reflect.TypeFor[T]() }
