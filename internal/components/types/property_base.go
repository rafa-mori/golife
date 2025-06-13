package types

import (
	ci "github.com/rafa-mori/golife/internal/components/interfaces"
	gl "github.com/rafa-mori/golife/logger"
	l "github.com/rafa-mori/logz"

	"fmt"
	"github.com/google/uuid"
	"reflect"
	"sync/atomic"
)

// PropertyValBase is a type for the value.
type PropertyValBase[T any] struct {
	// Logger is the logger for this context.
	Logger l.Logger

	// v is the value.
	*atomic.Pointer[T]

	// Reference is the identifiers for the context.
	// IReference
	*Reference

	//muCtx is the mutexes for the context.
	*Mutexes

	// validation is the validation for the value.
	*Validation[T]

	// Channel is the channel for the value.
	ci.IChannelCtl[T]
	channelCtl *ChannelCtl[T]
}

// NewVal is a function that creates a new PropertyValBase instance.
func newVal[T any](name string, v *T) *PropertyValBase[T] {
	ref := NewReference(name)

	// Create a new PropertyValBase instance
	vv := atomic.Pointer[T]{}
	if v != nil {
		vv.Store(v)
	} else {
		vv.Store(new(T))
	}

	// Create a new mutexes instance
	mu := NewMutexesType()

	// Create a new validation instance
	validation := newValidation[T]()

	gl.Log("debug", "Created new PropertyValBase instance for:", name, "ID:", ref.GetID().String())

	return &PropertyValBase[T]{
		Pointer:    &vv,
		Validation: validation,
		Reference:  ref.GetReference(),
		channelCtl: NewChannelCtl[T](name, nil).(*ChannelCtl[T]),
		Mutexes:    mu,
	}
}

func NewVal[T any](name string, v *T) ci.IPropertyValBase[T] {
	ref := NewReference(name)

	// Create a new PropertyValBase instance
	vv := atomic.Pointer[T]{}
	if v != nil {
		vv.Store(v)
	} else {
		vv.Store(new(T))
	}

	// Create a new mutexes instance
	mu := NewMutexesType()

	// Create a new validation instance
	validation := newValidation[T]()

	gl.Log("debug", "Created new PropertyValBase instance for:", name, "ID:", ref.GetID().String())

	return &PropertyValBase[T]{
		Pointer:    &vv,
		Validation: validation,
		Reference:  ref.GetReference(),
		channelCtl: NewChannelCtl[T](name, nil).(*ChannelCtl[T]),
		Mutexes:    mu,
	}
}

// GetLogger is a method that returns the logger for the value.
func (v *PropertyValBase[T]) GetLogger() l.Logger {
	if v == nil {
		gl.Log("error", "GetLogger: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return nil
	}
	return v.Logger
}

// GetName is a method that returns the name of the value.
func (v *PropertyValBase[T]) GetName() string {
	if v == nil {
		gl.Log("error", "GetName: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return ""
	}
	return v.Name
}

// GetID is a method that returns the ID of the value.
func (v *PropertyValBase[T]) GetID() uuid.UUID {
	if v == nil {
		gl.Log("error", "GetID: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return uuid.Nil
	}
	return v.ID
}

// Value is a method that returns the value.
func (v *PropertyValBase[T]) Value() *T {
	if v == nil {
		gl.Log("error", "Value: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return nil
	}
	return v.Load()
}

// StartCtl is a method that starts the control channel.
func (v *PropertyValBase[T]) StartCtl() <-chan string {
	if v == nil {
		gl.Log("error", "StartCtl: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return nil
	}
	if v.channelCtl == nil {
		gl.Log("error", "StartCtl: channel control is nil (", reflect.TypeFor[T]().String(), ")")
		return nil
	}
	return v.channelCtl.Channels["ctl"].(<-chan string)
}

// Type is a method that returns the type of the value.
func (v *PropertyValBase[T]) Type() reflect.Type { return reflect.TypeFor[T]() }

// Get is a method that returns the value.
func (v *PropertyValBase[T]) Get(async bool) any {
	if v == nil {
		gl.Log("error", "Get: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return nil
	}
	if async {
		if v.channelCtl != nil {
			gl.Log("debug", "Getting value from channel for:", v.Name, "ID:", v.ID.String())
			v.channelCtl.Channels["get"].(chan T) <- *v.Load()
		}
	} else {
		gl.Log("debug", "Getting value for:", v.Name, "ID:", v.ID.String())
		return v.Load()
	}
	return nil
}

// Set is a method that sets the value.
func (v *PropertyValBase[T]) Set(t *T) bool {
	if v == nil {
		gl.Log("error", "Set: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return false
	}
	if t == nil {
		gl.Log("error", "Set: value is nil (", reflect.TypeFor[T]().String(), ")")
		return false
	}
	if v.Validation != nil {
		if ok := v.Validation.Validate(t); !ok.GetIsValid() {
			gl.Log("error", fmt.Sprintf("Set: validation error (%s): %v", reflect.TypeFor[T]().String(), v.Validation.GetResults()))
			return false
		}
	}
	v.Store(t)
	if v.channelCtl != nil {
		gl.Log("debug", "Setting value for:", v.Name, "ID:", v.ID.String())
		v.channelCtl.Channels["set"].(chan T) <- *t
	}
	return true
}

// Clear is a method that clears the value.
func (v *PropertyValBase[T]) Clear() bool {
	if v == nil {
		gl.Log("error", "Clear: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return false
	}
	if v.channelCtl != nil {
		gl.Log("debug", "Clearing value for:", v.Name, "ID:", v.ID.String())
		v.channelCtl.Channels["clear"].(chan string) <- "clear"
	}
	return true
}

// IsNil is a method that checks if the value is nil.
func (v *PropertyValBase[T]) IsNil() bool {
	if v == nil {
		gl.Log("error", "Get: property does not exist (", reflect.TypeFor[T]().String(), ")")
		return true
	}
	return v.Load() == nil
}

// Serialize is a method that serializes the value.
func (v *PropertyValBase[T]) Serialize(format string) ([]byte, error) {
	mapper := NewMapper[T]()
	if value := v.Value(); value == nil {
		return nil, fmt.Errorf("value is nil")
	} else {
		return mapper.Serialize(nil, value, format)
	}
}

// Deserialize is a method that deserializes the data into the value.
func (v *PropertyValBase[T]) Deserialize(data []byte, format string) error {
	mapper := NewMapper[T]()
	if value := v.Value(); value == nil {
		return fmt.Errorf("value is nil")
	} else {
		return mapper.Deserialize(data, value, format)
	}
}
