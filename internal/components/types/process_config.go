package types

import (
	ci "github.com/faelmori/golife/internal/components/interfaces"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
)

// ProcessInputConfig is a struct that holds the configuration for the process.
type ProcessInputConfig struct {
	// Logger is the logger for the process
	Logger l.Logger
	// Reference is the reference for the process with ID and Name
	*Reference
	// Mutexes is the mutex for the process
	*Mutexes
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

// NewProcessConfig is the constructor for ProcessInputConfig[T any]
func newProcessConfig(name string, wait, restart bool, typ string, metadata map[string]any, logger l.Logger, debug bool) *ProcessInputConfig {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	if debug {
		gl.SetDebug(debug)
	}
	mu := NewMutexesType()
	ref := NewReference(name)
	pc := &ProcessInputConfig{
		Reference:   ref.GetReference(),
		Mutexes:     mu,
		WaitFor:     wait,
		Restart:     restart,
		ProcessType: typ,
		Metadata:    metadata,
	}
	return pc
}

// NewProcessConfig is the constructor for ProcessInputConfig[T any]
func NewProcessConfig(name string, wait, restart bool, typ string, metadata map[string]any, logger l.Logger, debug bool) ci.IProcessInputConfig {
	return newProcessConfig(name, wait, restart, typ, metadata, logger, debug)
}

// GetWaitFor returns the boolean that indicates if the process should wait for the command to finish.
func (pi *ProcessInputConfig) GetWaitFor() bool {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.WaitFor
}

// GetRestart returns the boolean that indicates if the process should be restarted.
func (pi *ProcessInputConfig) GetRestart() bool {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.Restart
}

// GetProcessType returns the type of the process.
func (pi *ProcessInputConfig) GetProcessType() string {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.ProcessType
}

// GetMetadata returns the metadata for the process.
func (pi *ProcessInputConfig) GetMetadata(key string) (any, bool) {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	if pi.Metadata == nil {
		return nil, false
	}
	value, exists := pi.Metadata[key]
	return value, exists
}

// SetMetadata sets the metadata for the process.
func (pi *ProcessInputConfig) SetMetadata(key string, value any) {
	pi.Mutexes.MuLock()
	defer pi.Mutexes.MuUnlock()

	if pi.Metadata == nil {
		pi.Metadata = make(map[string]any)
	}
	pi.Metadata[key] = value
}
