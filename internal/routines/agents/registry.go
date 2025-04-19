package agents

import "reflect"

// IChannel is an interface that extends IDynChan and adds methods for monitoring and controlling the channel.
type IChannel[T any, N int] interface {
	// GetChan returns the channel instance.
	GetChan() (chan any, reflect.Type)

	// Listen listens for messages on the channel and processes them. It will block until a message is received or the channel is closed.
	Listen() (<-chan any, reflect.Type, error)

	// Monitor returns the system channel for monitoring.
	Monitor() (chan any, reflect.Type, error)

	// GetLast retrieves the last message sent to the channel.
	GetLast() (any, reflect.Type, error)

	// Send sends a message to the channel.
	Send(v any) error

	// SetLast sets the last message sent to the channel.
	SetLast(v any) error

	// StartSysMonitor starts the system monitoring for the channel.
	StartSysMonitor()

	// StopSysMonitor stops the system monitoring for the channel.
	StopSysMonitor()

	// IsSysEnabled checks if the system monitoring is enabled.
	IsSysEnabled() bool

	// SetSysEnabled sets the system monitoring to enabled or disabled.
	SetSysEnabled(enabled bool)

	// GetType returns the type of the channel.
	GetType() reflect.Type

	// Close closes the channel.
	Close() error

	// IsClosed checks if the channel is closed.
	IsClosed() bool
}
