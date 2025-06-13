package types

import (
	l "github.com/rafa-mori/logz"
	"sync"
	"time"
)

// TelemetryIdentifier is a struct that holds the identifier for telemetry data
type TelemetryIdentifier struct {
	// ID is the unique identifier for this telemetry instance
	ID string
	// Name is the name of the telemetry instance
	Name string
	// Logger is the Logger instance for this telemetry
	Logger l.Logger
	// Type is the type of telemetry (e.g., CPU, Memory, etc.)
	Type string
}

// TelemetryMutex is a struct that holds mutexes for synchronizing access to telemetry data
type TelemetryMutex struct {
	// mutex is a mutex for synchronizing access to the telemetry data
	mutex *sync.RWMutex
	// mutexL is a mutex for synchronizing access to the Logger
	mutexL *sync.RWMutex
	// mutexC is a mutex for synchronizing access to the channels
	mutexC *sync.RWMutex
	// mutexW is a mutex for synchronizing access to the wait group
	mutexW *sync.RWMutex
}

// TelemetryChannel is a struct that holds channels for telemetry data
type TelemetryChannel struct {
	// channelOut is a channel for sending telemetry data
	channelOut chan map[string]float64
	// channelIn is a channel for receiving telemetry data
	channelIn chan map[string]float64
	// channelErr is a channel for sending error messages
	channelErr chan error
	// channelDone is a channel for signaling when the telemetry is done
	channelDone chan struct{}
	// channelExit is a channel for signaling when the telemetry should exit
	channelExit chan struct{}
}

// TelemetryData is a struct that holds telemetry data
type TelemetryData struct {
	// LastUpdated is the last time the telemetry data was updated
	LastUpdated time.Time
	// Metrics is a map of metric names to their values
	Metrics map[string]float64
}

// TelemetryLogger is a struct that holds a Logger for telemetry data
type TelemetryLogger struct {
	// logger is the logger instance
	logger l.Logger
}

// TelemetryProperty is a struct that holds a property for telemetry data
type TelemetryProperty struct {
	// property is the property instance
	property any
}

// TelemetryConfig is a struct that holds the configuration for telemetry data
type TelemetryConfig struct {
	// config is the configuration instance
	config any
}

// Telemetry is a struct that holds telemetry data for a process
type Telemetry struct {
	// TelemetryIdentifier is the identifier for telemetry data
	TelemetryIdentifier
	// TelemetryLogger is the Logger for telemetry data
	TelemetryLogger
	// TelemetryData is the telemetry data
	TelemetryData
	// TelemetryMutex is a mutex for synchronizing access to the telemetry data
	TelemetryMutex
	// TelemetryChannel is a struct that holds channels for telemetry data
	TelemetryChannel
	// TelemetryProperty is a struct that holds a property for telemetry data
	TelemetryProperty
	// TelemetryConfig is a struct that holds the configuration for telemetry data
	TelemetryConfig
}

// NewTelemetry creates a new Telemetry instance
func NewTelemetry() *Telemetry {
	return &Telemetry{
		TelemetryIdentifier: TelemetryIdentifier{
			ID:     "default",
			Name:   "default",
			Logger: l.GetLogger("Telemetry"),
			Type:   "default",
		},
		TelemetryLogger: TelemetryLogger{
			logger: l.GetLogger("Telemetry"),
		},
		TelemetryData: TelemetryData{
			LastUpdated: time.Now(),
			Metrics:     make(map[string]float64),
		},
		TelemetryMutex: TelemetryMutex{
			mutex:  &sync.RWMutex{},
			mutexL: &sync.RWMutex{},
			mutexC: &sync.RWMutex{},
			mutexW: &sync.RWMutex{},
		},
		TelemetryChannel: TelemetryChannel{
			channelOut:  make(chan map[string]float64, 1),
			channelIn:   make(chan map[string]float64, 1),
			channelErr:  make(chan error, 1),
			channelDone: make(chan struct{}, 1),
			channelExit: make(chan struct{}, 1),
		},
	}
}

// UpdateMetrics updates the telemetry metrics with the given data
func (t *Telemetry) UpdateMetrics(data map[string]float64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.LastUpdated = time.Now()
	for key, value := range data {
		t.Metrics[key] = value
	}
}

// GetMetrics returns the telemetry metrics
func (t *Telemetry) GetMetrics() map[string]float64 {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.Metrics
}

// GetLastUpdated returns the last updated time of the telemetry
func (t *Telemetry) GetLastUpdated() time.Time {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.LastUpdated
}

// ResetMetrics resets the telemetry metrics
func (t *Telemetry) ResetMetrics() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.Metrics = make(map[string]float64)
	t.LastUpdated = time.Now()
}
