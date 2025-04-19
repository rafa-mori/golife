package process_input

import (
	p "github.com/faelmori/golife/components/types"
	l "github.com/faelmori/logz"
)

// ProcessConfig is a struct that holds the configuration for the process.
type ProcessConfig[T any] struct {
	// Logger is the logger for the process
	Logger l.Logger
	// Reference is the reference for the process with ID and Name
	*p.Reference
	// Mutexes is the mutex for the process
	*p.Mutexes
	// IsRunning is a boolean that indicates if the process is running
	IsRunning bool `json:"is_running" yaml:"is_running" xml:"is_running" gorm:"is_running"`
	// WaitFor is a boolean that indicates if the process should wait for the command to finish
	WaitFor bool `json:"wait_for" yaml:"wait_for" xml:"wait_for" gorm:"wait_for"`
	// Restart is a boolean that indicates if the process should be restarted
	Restart bool `json:"restart" yaml:"restart" xml:"restart" gorm:"restart"`
	// ProcessType is the type of the process
	ProcessType string `json:"process_type" yaml:"process_type" xml:"process_type" gorm:"process_type"`
	// metadata is a map of metadata for the process
	Metadata map[string]any `json:"metadata" yaml:"metadata" xml:"metadata" gorm:"metadata"`
}

// NewProcessConfig is the constructor for ProcessConfig[T any]
func NewProcessConfig[T any](name string, wait, restart bool, typ string, metadata map[string]any, logger l.Logger, debug bool) *ProcessConfig[T] {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	mu := p.NewMutexes()
	ref := p.NewReference(name)
	pc := &ProcessConfig[T]{
		Reference:   ref,
		Mutexes:     mu,
		WaitFor:     wait,
		Restart:     restart,
		ProcessType: typ,
		Metadata:    metadata,
	}
	return pc
}

// GetWaitFor returns the boolean that indicates if the process should wait for the command to finish.
func (pi *ProcessConfig[T]) GetWaitFor() bool {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	return pi.WaitFor
}

// GetRestart returns the boolean that indicates if the process should be restarted.
func (pi *ProcessConfig[T]) GetRestart() bool {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	return pi.Restart
}

// GetProcessType returns the type of the process.
func (pi *ProcessConfig[T]) GetProcessType() string {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	return pi.ProcessType
}

// GetMetadata returns the metadata for the process.
func (pi *ProcessConfig[T]) GetMetadata(key string) (any, bool) {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	if pi.Metadata == nil {
		return nil, false
	}
	value, exists := pi.Metadata[key]
	return value, exists
}

// SetMetadata sets the metadata for the process.
func (pi *ProcessConfig[T]) SetMetadata(key string, value any) {
	pi.Mutexes.Lock()
	defer pi.Mutexes.Unlock()

	if pi.Metadata == nil {
		pi.Metadata = make(map[string]any)
	}
	pi.Metadata[key] = value
}
