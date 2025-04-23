package interfaces

import (
	"github.com/google/uuid"
	"reflect"
)

type IChannelBase[T any] interface {
	IMutexes

	GetName() string       // The name of the channel.
	GetChannel() chan T    // The channel for the value. Main channel for this struct.
	GetType() reflect.Type // The type of the channel.
	GetBuffers() int       // The number of buffers for the channel.

	SetName(name string) string // Set the name of the channel.
	SetChannel(chan T) chan T   // The channel for the value. Main channel for this struct.
	SetBuffers(buffers int) int // The number of buffers for the channel.

	Close() error // Close the channel.
	Clear() error // Clear the channel.
}

type IChannelCtl[T any] interface {
	IMutexes

	// Structure management

	GetID() uuid.UUID
	GetName() string
	SetName(name string) string

	// Property query

	GetProperty() IProperty[T]

	// SubChannels management

	GetSubChannels() map[string]interface{}
	SetSubChannels(channels map[string]interface{}) map[string]interface{}

	GetSubChannelByName(name string) (IChannelBase[any], reflect.Type, bool)
	SetSubChannelByName(name string, channel IChannelBase[any]) (IChannelBase[any], error)

	GetSubChannelTypeByName(name string) (reflect.Type, bool)

	GetSubChannelBuffersByName(name string) (int, bool)
	SetSubChannelBuffersByName(name string, buffers int) (int, error)

	// Main channel management

	GetMainChannel() chan T
	SetMainChannel(channel chan T) chan T
	GetMainChannelType() reflect.Type

	GetHasMetrics() bool
	SetHasMetrics(hasMetrics bool) bool
	GetBufferSize() int
	SetBufferSize(size int) int

	Close() error

	// Chainable methods

	WithProperty(property IProperty[T]) IChannelCtl[T]
	WithChannel(channel chan T) IChannelCtl[T]
	WithBufferSize(size int) IChannelCtl[T]
	WithMetrics(metrics bool) IChannelCtl[T]
}
