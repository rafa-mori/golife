package types

import (
	l "github.com/faelmori/golife/logger"
	"github.com/google/uuid"
	"reflect"
	"sync/atomic"
)

// Reference is a struct that holds the Reference ID and name.
type Reference struct {
	// refID is the unique identifier for this context.
	ID uuid.UUID
	// refName is the name of the context.
	Name string
}

// NewReference is a function that creates a new Reference instance.
func NewReference(name string) *Reference {
	return &Reference{
		ID:   uuid.New(),
		Name: name,
	}
}

// val is a type for the value.
type val[T any] struct {
	// v is the value.
	*atomic.Pointer[T]

	// Reference is the identifiers for the context.
	*Reference

	//muCtx is the mutexes for the context.
	*muCtx

	// validation is the validation for the value.
	*Validation[T]

	// Channel is the channel for the value.
	*ChannelCtl[T]
}

// NewVal is a function that creates a new val instance.
func newVal[T any](name string, v *T) *val[T] {
	ref := NewReference(name)

	// Create a new val instance
	vv := atomic.Pointer[T]{}
	if v != nil {
		vv.Store(v)
	} else {
		vv.Store(new(T))
	}

	// Create a new mutexes instance
	mu := newMuCtx()

	// Create a new validation instance
	validation := NewValidation[T]()

	l.Log("info", "Created new val instance for:", name, "ID:", ref.ID.String())
	return &val[T]{
		Pointer:    &vv,
		Validation: validation,
		Reference:  ref,
		muCtx:      mu,
	}
}

// StartCtl is a method that starts the control channel.
func (v *val[T]) StartCtl() <-chan string {
	l.Log("info", "Starting control channel for:", v.Name, "ID:", v.ID.String())
	go monitorRoutine[T](v)
	return v.ctl
}

// Type is a method that returns the type of the value.
func (v *val[T]) Type() reflect.Type { return reflect.TypeFor[T]() }

// Get is a method that returns the value.
func (v *val[T]) Get(async bool) any {
	if v == nil {
		l.Log("error", "Get: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return nil
	}
	vl := v.Load()
	if async {
		if v.ch != nil {
			if vl == nil {
				if v.Type().Kind() != reflect.Ptr {
					l.Log("info", "Creating async value for:", v.Name, "ID:", v.ID.String())
					vl = new(T)
					v.ch <- *vl
				}
			} else {
				l.Log("info", "Sending async value for:", v.Name, "ID:", v.ID.String())
				v.ch <- *vl
			}
		} else {
			l.Log("warn", "Get: channel is nil, cannot send async value (", reflect.TypeFor[T]().String(), ")")
		}
	}
	return vl
}

// Set is a method that sets the value.
func (v *val[T]) Set(t *T) bool {
	if v == nil {
		l.Log("error", "Get: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return false
	}
	if t == nil {
		l.Log("error", "Set: nil is not a valid value (", reflect.TypeFor[T]().String(), ")")
		return false
	}
	if v.hasValidation {
		if !v.Validate(*t) {
			l.Log("error", "Set: validation failed (", reflect.TypeFor[T]().String(), ")")
			return false
		}
	}
	if v.CompareAndSwap(v.Load(), t) {
		l.Log("info", "Set: changed value for:", v.Name, "ID:", v.ID.String())
		if v.ch != nil {
			v.ch <- *t
		}
		return true
	}
	l.Log("error", "Set: value not changed (", reflect.TypeFor[T]().String(), ")")
	return false
}

// Clear is a method that clears the value.
func (v *val[T]) Clear() {
	if v == nil {
		l.Log("error", "Get: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return
	}
	if v.Load() != nil {
		l.Log("info", "Clearing value for:", v.Name, "ID:", v.ID.String())
		vl := *new(T)

		v.Store(&vl)

		if v.ch != nil {
			l.Log("info", "Sending clear value for:", v.Name, "ID:", v.ID.String())
			v.ch <- vl
		}
	}
}

// IsNil is a method that checks if the value is nil.
func (v *val[T]) IsNil() bool {
	if v == nil {
		l.Log("error", "Get: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return true
	}
	return v.Load() == nil
}
