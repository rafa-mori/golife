package process_input

import (
	"fmt"
	t "github.com/faelmori/golife/components/types"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"reflect"
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
	*ProcessSystemBase[T]

	// ProcessRuntimeBase is the runtime information for the process
	*ProcessRuntimeBase[T]

	// ProcessConfig is the configuration for the process
	*ProcessConfig[T]
}

// NewProcessInput creates a new ProcessInput instance from the provided config data and Logger.
func NewProcessInput[T any](name string, logger l.Logger, debug bool) *ProcessInput[T] {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	if debug {
		gl.SetDebug(debug)
	}
	ref := t.NewReference(name)
	mu := t.NewMutexes()
	npi := &ProcessInput[T]{
		Logger:    logger,
		Mutexes:   mu,
		Reference: ref,
		ProcessSystemBase: &ProcessSystemBase[T]{
			Command:     "",
			Args:        nil,
			Path:        "",
			ProcPid:     0,
			ProcPidFile: "",
			ProcPointer: 0,
		},
		ProcessRuntimeBase: &ProcessRuntimeBase[T]{
			Mutexes: mu,
			Object:  *new(T),
			Function: &t.ValidationFunc[T]{
				Priority: 0,
				Func: func(obj T, args ...any) *t.ValidationResult {
					objT := reflect.ValueOf(obj).Interface().(T)
					if len(args) > 0 {
						for _, arg := range args {
							if reflect.TypeOf(arg) == reflect.TypeFor[func(T, ...any) *t.ValidationResult]() {
								return arg.(func(T, ...any) *t.ValidationResult)(objT, args...)
							}
						}
						return nil
					}
					return nil
				},
				Result: nil,
			},
		},
		ProcessConfig: &ProcessConfig[T]{
			IsRunning:   false,
			WaitFor:     false,
			Restart:     false,
			ProcessType: "",
		},
	}
	return npi
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
	mu := t.NewMutexes()
	ref := t.NewReference(name)
	npi := &ProcessInput[T]{
		Logger:    logger,
		Mutexes:   mu,
		Reference: ref,
		ProcessSystemBase: &ProcessSystemBase[T]{
			Mutexes: mu,
			Command: "",
			Args:    nil,
		},
		ProcessRuntimeBase: &ProcessRuntimeBase[T]{
			Mutexes:    mu,
			ObjectType: reflect.TypeOf((*T)(nil)).Elem(),
			Object:     *new(T),
			Function: &t.ValidationFunc[T]{
				Priority: 0,
				Func: func(obj T, args ...any) *t.ValidationResult {
					objT := reflect.ValueOf(obj).Interface().(T)
					if len(args) > 0 {
						for _, arg := range args {
							if reflect.TypeOf(arg) == reflect.TypeFor[func(T, ...any) *t.ValidationResult]() {
								return arg.(func(T, ...any) *t.ValidationResult)(objT, args...)
							}
						}
						return nil
					}
					return nil
				},
			},
		},
	}
	err := mapper.Deserialize(data, npi, format)
	if err != nil {
		return nil, err
	}
	// This method is used to initialize the cmd field with right and more robust way
	if npi.BuildCmd() == nil {
		gl.LogObjLogger[ProcessInput[T]](npi, "error", "Command is nil")
		return nil, fmt.Errorf("command is nil")
	}
	return npi, nil
}
