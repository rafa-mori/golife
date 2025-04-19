package agents

import (
	"github.com/faelmori/golife/services"
	"reflect"
	"sync"
	"sync/atomic"
)

// channel is a struct that implements the IChannel and IDynChan interfaces.
// It provides a more complex implementation of a channel with monitoring capabilities.
type channel[T any, N int] struct {
	services.IChannel[T, N]

	mu           sync.Mutex        // Mutex for thread-safe operations.
	wg           sync.WaitGroup    // WaitGroup for goroutines.
	buffers      N                 // Buffer size of the channel.
	name         string            // Name of the channel.
	chanAny      chan any          // Main channel for communication.
	chanSys      chan any          // System channel for monitoring.
	chanStop     chan struct{}     // Channel to signal stopping of monitoring.
	chanT        chan T            // Main channel for communication.
	last         atomic.Pointer[T] // Last value sent through the channel.
	isSysEnabled bool              // Flag indicating if system monitoring is enabled.
}

// NewLoaderChanInterface creates a new channel with a name, type, and buffer size.
// It returns an instance of IChannel.
func NewLoaderChanInterface[T any, N int](name string, tp *T, buffers N) IChannel[T, N] {
	ch := &channel[T, N]{
		name:     name,
		buffers:  buffers,
		last:     atomic.Pointer[T]{},
		chanSys:  make(chan any, 10),
		chanStop: make(chan struct{}, 1),
	}
	if buffers > 0 {
		ch.chanT = make(chan T, buffers)
	} else {
		ch.chanT = make(chan T, 2)
	}
	if tp != nil {
		ch.last.Store(tp)
	}
	return ch
}

// NewChannel creates a new channel with a name, type, and buffer size.
// It returns an instance of IChannel.
func NewChannel[T any, N int](name string, tp *T, buffers N) IChannel[T, N] {
	ch := &channel[T, N]{
		name:     name,
		buffers:  buffers,
		last:     atomic.Pointer[T]{},
		chanSys:  make(chan any, (buffers+1)/2),
		chanStop: make(chan struct{}, 1),
	}

	if buffers > 0 {
		ch.chanT = make(chan T, buffers)
	} else {
		ch.chanT = make(chan T, 2)
	}
	if tp != nil {
		ch.last.Store(tp)
	}
	return ch
}

// GetType returns the type of the channel.
func (c *channel[T, N]) GetType() reflect.Type { return reflect.TypeFor[T]() }

// GetChan returns the main channel instance.
func (c *channel[T, N]) GetChan() (chan any, reflect.Type) {
	if c.chanAny == nil {
		c.chanAny = make(chan any, 2)
	}
	if c.chanT == nil {
		c.chanT = make(chan T, 2)
	}
	if c.chanSys == nil {
		c.chanSys = make(chan any, 10)
	}
	if c.isSysEnabled {
		go c.StartSysMonitor()
	}
	return c.chanAny, reflect.TypeFor[T]()
}

// Listen listens for messages on the channel and processes them.
func (c *channel[T, N]) Listen() (<-chan any, reflect.Type, error) {
	if c.chanT == nil {
		c.chanT = make(chan T, 2)
	}
	if c.chanSys == nil {
		c.chanSys = make(chan any, 10)
	}
	if c.isSysEnabled {
		go c.StartSysMonitor()
	}
	return c.chanAny, reflect.TypeFor[T](), nil
}

// Monitor returns the system channel for monitoring.
func (c *channel[T, N]) Monitor() (chan any, reflect.Type, error) {
	if c.chanSys == nil {
		c.chanSys = make(chan any, 10)
	}
	if c.isSysEnabled {
		go c.StartSysMonitor()
	}
	return c.chanSys, reflect.TypeFor[T](), nil
}

// GetLast retrieves the last value sent to the channel.
func (c *channel[T, N]) GetLast() (any, reflect.Type, error) {
	if c.last.Load() == nil {
		return nil, reflect.TypeFor[T](), nil
	}
	v := c.last.Load()
	return v, reflect.TypeFor[T](), nil
}

// SetLast sets the last value sent to the channel.
func (c *channel[T, N]) SetLast(v any) error {
	if v == nil {
		return nil
	}
	c.last.Store(v.(*T))
	if c.chanSys != nil {
		c.chanSys <- v
	}
	if c.chanT != nil {
		c.chanT <- v.(T)
	}
	return nil
}

// StartSysMonitor starts a goroutine to monitor the system channel for logging or auditing purposes.
func (c *channel[T, N]) StartSysMonitor() {
	// When the channel is created, it will be in a closed state.
	// The channel will be opened when the first message is sent and keep it open until the channel is closed.
	// The channel will be closed when the last message is sent and the channel is closed.
	// Until then, the channel will be in a closed state.
	if c.chanSys == nil {
		if c.buffers > 0 {
			c.chanSys = make(chan any, (c.buffers+1)/2)
		} else {
			c.buffers = 10
			c.chanSys = make(chan any, c.buffers)
		}
	}
	if c.chanT == nil {
		if c.buffers > 0 {
			c.chanT = make(chan T, c.buffers)
		} else {
			c.buffers = 10
			c.chanT = make(chan T, c.buffers)
		}
	}
	defer c.StopSysMonitor()
	if c.isSysEnabled {
		return
	}
	c.wg.Add(1)
	go func() {
		c.isSysEnabled = true
		defer c.wg.Done()
		defer c.StopSysMonitor()
		for {
			select {
			case v := <-c.chanSys:
				c.last.Store(v.(*T))
			case currMsg := <-c.chanT:
				vl := reflect.ValueOf(currMsg)
				if !vl.IsValid() || vl.IsNil() || vl.IsZero() {
					continue
				}
				typeOfCurrMsg := vl.Type()
				lastMsg := c.last.Load()
				if typeOfCurrMsg != reflect.TypeOf(lastMsg) {
					continue
				}
				if typeOfCurrMsg.Kind() == reflect.Ptr || typeOfCurrMsg.Kind() == reflect.Interface {
					if reflect.DeepEqual(currMsg, lastMsg) {
						continue
					}
					if c.last.CompareAndSwap(c.last.Load(), &currMsg) {
						c.chanSys <- currMsg
					} else {
						continue
					}
				}
			case <-c.chanStop:
				return
			}
		}
	}()
}

// StopSysMonitor stops the system monitoring goroutine and closes all associated channels.
func (c *channel[T, N]) StopSysMonitor() {
	defer c.wg.Done()

	if c.chanSys != nil {
		c.mu.Lock()
		close(c.chanSys)
		c.chanSys = nil
		c.mu.Unlock()
	}
	if c.chanT != nil {
		c.mu.Lock()
		close(c.chanT)
		c.chanT = nil
		c.mu.Unlock()
	}
	if c.chanStop != nil {
		c.mu.Lock()
		close(c.chanStop)
		c.chanStop = nil
	}
	if c.chanT == nil && c.chanSys == nil && c.chanStop == nil {
		c.isSysEnabled = false
	} else {
		c.isSysEnabled = true
	}
}

// IsSysEnabled checks if the system monitoring is enabled.
func (c *channel[T, N]) IsSysEnabled() bool { return c.isSysEnabled }

// SetSysEnabled sets the system monitoring to enabled or disabled.
func (c *channel[T, N]) SetSysEnabled(enabled bool) {
	c.isSysEnabled = enabled
	if c.isSysEnabled {
		c.StartSysMonitor()
	} else {
		c.StopSysMonitor()
	}
}

// Send sends a message to the channel.
func (c *channel[T, N]) Send(v any) error {
	if v == nil {
		return nil
	}
	if c.chanT == nil {
		c.chanT = make(chan T, 2)
	}
	if c.chanSys == nil {
		c.chanSys = make(chan any, 10)
	}
	c.chanT <- v.(T)
	if c.chanSys != nil {
		c.chanSys <- v
	}
	if c.last.Load() == nil {
		c.last.Store(v.(*T))
	}
	return nil
}

// Close closes the channel and stops the monitoring goroutine.
func (c *channel[T, N]) Close() error {
	if c.chanT != nil {
		close(c.chanT)
		c.chanT = nil
	}
	if c.chanSys != nil {
		close(c.chanSys)
		c.chanSys = nil
	}
	if c.chanStop != nil {
		close(c.chanStop)
		c.chanStop = nil
	}
	return nil
}

// IsClosed checks if the channel is closed.
func (c *channel[T, N]) IsClosed() bool {
	if c.chanT == nil && c.chanSys == nil && c.chanStop == nil {
		return true
	}
	return false
}
