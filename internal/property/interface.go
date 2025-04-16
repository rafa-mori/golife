package property

import (
	"github.com/faelmori/golife/services"
	"reflect"
)

// VoValue is a generic interface for getting and setting a value of type T.
type VoValue[T any] interface {
	// GetType retrieves the type of the value.
	GetType() reflect.Type
	// GetValueWithType retrieves the value of type T.
	GetValueWithType() (any, reflect.Type)
	// GetValue retrieves the value of type T.
	GetValue() T
	// SetValue sets the value of type T.
	SetValue(T, func(any) error) error
	// SetDefaultValue sets the default value of type T.
	SetDefaultValue(T) error
	// AddValidator adds a validation function for the value.
	AddValidator(string, func(any) error) error
	// AddListener adds a change listener for the value.
	AddListener(string, ChangeListener[T]) error
	// AddChainedListener adds a chained listener for the value.
	AddChainedListener(string, string, ChangeListener[T]) error
	// RemoveListener removes a listener for the value.
	RemoveListener(string) error
	// RemoveAllListeners removes all listeners for the value.
	RemoveAllListeners()
}

// PropertyBase defines the base interface for a property.
type PropertyBase interface {
	// GetName retrieves the name of the property.
	GetName() string
	// SetMetadata sets a metadata key-value pair for the property.
	SetMetadata(string, any)
	// GetMetadata retrieves the value of a metadata key. Returns the value and a boolean indicating if the key exists.
	GetMetadata(string) (any, bool)
	// GetType retrieves the type of the property as a string.
	GetType() reflect.Type
	// GetStringType retrieves the type of the property as a string.
	GetStringType() string
}

type PropertyChanCtl[T any] interface {
	// GetChannel retrieves the channel for receiving updates to the property value.
	GetChannel() services.IChannel[T, int]
	// SetChannel sets the channel for receiving updates to the property value.
	SetChannel(services.IChannel[T, int])
	// GetChannelValue retrieves the value from the channel.
	GetChannelValue() any
	// SetChannelValue sets the value in the channel.
	SetChannelValue(any)

	// GetChannelType retrieves the type of the channel.
	GetChannelType() string
	// SetChannelType sets the type of the channel.
	SetChannelType(string, reflect.Type)
	// GetChannelName retrieves the name of the channel.
	GetChannelName() string
	// BroadcastChange broadcasts a change event to the channel.
	BroadcastChange(*T, T)
}

// Property is a generic interface that combines PropertyBase and VoValue.
type Property[T any] interface {
	PropertyBase
	PropertyChanCtl[T]
	VoValue[T]
}
