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
	GetValue() *T
	// SetValue sets the value of type T.
	SetValue(T, func(T) error) error
	// SetDefaultValue sets the default value of type T.
	SetDefaultValue(*T) error
	// AddValidator adds a validation function for the value.
	AddValidator(string, func(T) error) error
	// AddListener adds a change listener for the value.
	AddListener(string, *ChangeListener[T]) error
	// AddChainedListener adds a chained listener for the value.
	AddChainedListener(string, string, *ChangeListener[T]) error
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
	GetChannelType() reflect.Type
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

	GetName() string
	GetType() reflect.Type
	GetStringType() string
	GetMetadata(string) (interface{}, bool)
	SetMetadata(string, interface{})
	GetValueWithType() (any, reflect.Type)
	//GetValue() T
	SetValue(T, func(T) error) error
	GetChannel() services.IChannel[T, int]
	SetChannel(services.IChannel[T, int])
	GetChannelValue() any
	//AddValidator(name string, validator func(T) error) error
	//AddListener(name string, listener *ChangeListener[T]) error
	BroadcastChange(*T, T)
	SetChannelValue(any)
	//GetChannelType() string
	SetChannelType(string, reflect.Type)
	GetChannelName() string
	//AddChainedListener(s string, s2 string, listener ChangeListener[T]) error
	RemoveListener(string) error
	RemoveAllListeners()
	WaitForCondition(func() bool)
	NotifyCondition()
	notifyListeners(*T, T, *EventMetadata)
	validateAndSet(any) error
}

/*type KubexProperty[T any] struct {
	PropertyBase
	PropertyChanCtl[any]
	VoValue[any]
	MainMutex     sync.Mutex
	ChanMutex     sync.Mutex
	ListenerMutex sync.Mutex
	LayerMutex    sync.Mutex
	Cond          atomic.Pointer[sync.Cond]
	name          string
	metadata      Metadata
	value         atomic.Pointer[T]
	chanCtl       services.IChannel[T, int]
	validators    []func(any) error
	listeners     map[string]ChangeListener[T]
	cond          *sync.Cond
	Properties    map[string]any


	GetName() string
	GetType() reflect.Type
	GetStringType() string
	GetMetadata(key string) (bool)
	SetMetadata(key string, value interface {})
	GetValueWithType() (any, reflect.Type)
	GetValue() T
	SetValue(value T, cb func (T) error) error
	GetChannel() services.IChannel[T, int]
	SetChannel(channel services.IChannel[T, int])
	GetChannelValue() any
	SetDefaultValue(value *T) error
	AddValidator(name string, validator func (T) error) error
	AddListener(name string, listener *property.ChangeListener[T]) error
	BroadcastChange(oldValue *T, newValue T)
	SetChannelValue(a any)
	GetChannelType() string
	SetChannelType(s string, newType reflect.Type)
	GetChannelName() string
	AddChainedListener(s string, s2 string, listener property.ChangeListener[T]) error
	RemoveListener(s string) error
	RemoveAllListeners()
	WaitForCondition(check func () bool)
	NotifyCondition()
	notifyListeners(oldValue *T, newValue T, metadata *property.EventMetadata)
	validateAndSet(value any) error
}*/
