package process_input

import (
	t "github.com/faelmori/golife/components/types"
	l "github.com/faelmori/logz"
)

// ProcessInput is a struct that holds the input for the process.
type ProcessInput[T any] struct {
	// Logger is the logger for the process
	Logger l.Logger

	// Reference is the reference for the process with ID and Name
	*t.Reference

	// Mutexes is the mutex for the process
	*t.Mutexes

	// ProcessSystemBase is the system information for the process
	*ProcessSystemBase[ProcessInput[T], T]

	// ProcessRuntimeBase is the runtime information for the process
	*ProcessRuntimeBase[ProcessInput[T], T]

	// ProcessConfig is the configuration for the process
	*ProcessConfig[T]
}

// NewProcessRuntimeBase creates a new ProcessRuntimeBase instance.
func NewProcessRuntimeBase[T any](name string, object *T, function *t.ValidationFunc[ProcessInput[T]], waitFor, restart bool, logger l.Logger, debug bool) *ProcessInput[T] {
	npi := newProcessRuntimeBase[ProcessInput[T], T](name, object, function, waitFor, restart, logger, debug)
	npr := newSystemProcessInput[ProcessInput[T], T](name, "", nil, false, false, nil, logger, debug)
	return &ProcessInput[T]{
		Reference:          npi.Reference,
		Logger:             npi.Logger,
		Mutexes:            npi.Mutexes,
		ProcessRuntimeBase: npi,
		ProcessSystemBase:  npr,
		ProcessConfig:      npi.ProcessConfig,
	}
}

// NewSystemProcessInput creates a new ProcessInput instance with the provided Logger.
func NewSystemProcessInput[T any](name, command string, args []string, waitFor bool, restart bool, function *t.ValidationFunc[ProcessInput[T]], logger l.Logger, debug bool) *ProcessInput[T] {
	npi := newSystemProcessInput[ProcessInput[T], T](name, command, args, waitFor, restart, function, logger, debug)
	npr := newProcessRuntimeBase[ProcessInput[T], T](name, new(T), function, waitFor, restart, logger, debug)
	return &ProcessInput[T]{
		Reference:          npi.Reference,
		Logger:             npi.Logger,
		Mutexes:            npi.Mutexes,
		ProcessSystemBase:  npi,
		ProcessRuntimeBase: npr,
		ProcessConfig:      npi.ProcessConfig,
	}
}

// Serialize serializes the ProcessInput instance to the specified format.
func (pi *ProcessInput[T]) Serialize(format string) ([]byte, error) {
	mapper := t.NewMapper[ProcessInput[T]]()
	return mapper.Serialize(nil, pi, format)
}

// Deserialize deserializes the data into the ProcessInput instance.
func (pi *ProcessInput[T]) Deserialize(data []byte, format string) error {
	mapper := t.NewMapper[ProcessInput[T]]()
	return mapper.Deserialize(data, pi, format)
}

// GetLogger returns the Logger instance for the process.
func (pi *ProcessInput[T]) GetLogger() l.Logger {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	return pi.Logger
}

// GetReference returns the reference ID and name.
func (pi *ProcessInput[T]) GetReference() *t.Reference {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	return pi.Reference
}

// GetName returns the name of the process.
func (pi *ProcessInput[T]) GetName() string {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	return pi.Reference.Name
}

// Validate validates the ProcessInput instance.
func (pi *ProcessInput[T]) Validate() *t.ValidationResult {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	cmd := pi.BuildCmd()
	if cmd == nil {
		return &t.ValidationResult{
			IsValid: false,
			Message: "Command is nil",
			Error:   nil,
		}
	}
	if pi.Command == "" {
		return &t.ValidationResult{
			IsValid: false,
			Message: "Command is empty",
			Error:   nil,
		}
	}
	if pi.Args == nil {
		return &t.ValidationResult{
			IsValid: false,
			Message: "Args is nil",
			Error:   nil,
		}
	}
	if pi.Path == "" {
		return &t.ValidationResult{
			IsValid: false,
			Message: "Path is empty",
			Error:   nil,
		}
	}
	if pi.ProcessType == "" {
		return &t.ValidationResult{
			IsValid: false,
			Message: "ProcessType is empty",
			Error:   nil,
		}
	}
	if pi.ObjectType == nil {
		return &t.ValidationResult{
			IsValid: false,
			Message: "ObjectType is nil",
			Error:   nil,
		}
	}
	if pi.Reference == nil {
		return &t.ValidationResult{
			IsValid: false,
			Message: "Reference is nil",
			Error:   nil,
		}
	}
	if pi.Mutexes == nil {
		return &t.ValidationResult{
			IsValid: false,
			Message: "Mutexes is nil",
			Error:   nil,
		}
	}

	return &t.ValidationResult{
		IsValid: true,
		Message: "ProcessInput is valid",
		Error:   nil,
	}
}

// NewProcessInputFromConfig creates a new ProcessInput instance from the provided config data.
func NewProcessInputFromConfig[T any](name string, data []byte, format string) (*ProcessInput[T], error) {
	mapper := t.NewMapper[ProcessInput[T]]()
	logger := l.GetLogger("GoLife")
	npi := newSystemProcessInput[ProcessInput[T], T](name, "", nil, false, false, nil, logger, false)
	npr := newProcessRuntimeBase[ProcessInput[T], T](name, new(T), nil, false, false, logger, false)
	np := &ProcessInput[T]{
		Reference:          npr.Reference,
		Logger:             npr.Logger,
		Mutexes:            npr.Mutexes,
		ProcessSystemBase:  npi,
		ProcessRuntimeBase: npr,
		ProcessConfig:      npr.ProcessConfig,
	}
	err := mapper.Deserialize(data, np, format)
	if err != nil {
		return nil, err
	}
	return np, nil
}
