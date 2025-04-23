package types

import (
	ci "github.com/faelmori/golife/internal/components/interfaces"
	l "github.com/faelmori/logz"

	"reflect"
)

// ProcessInputRuntimeBase is a struct that holds the runtime information for the process.
type ProcessInputRuntimeBase[T any, P any] struct {
	// Logger is the logger for the process
	Logger l.Logger
	// Mutexes is the mutex for the process
	*Mutexes
	// Reference is the reference for the process with ID and Name
	*Reference `json:"reference" yaml:"reference" xml:"reference" gorm:"reference"`
	// ProcessInputConfig is the configuration for the process
	*ProcessInputConfig `json:"process_config" yaml:"process_config" xml:"process_config" gorm:"process_config"`
	// Object is the object to pass to the command
	Object *P `json:"object,omitempty" yaml:"object,omitempty" xml:"object,omitempty" gorm:"object,omitempty"`
	// Function is a custom function to wrap the command
	Function ci.IValidationFunc[T] `json:"function,omitempty" yaml:"function,omitempty" xml:"function,omitempty" gorm:"function,omitempty"`
}

// newProcessInputRuntimeBase creates a new ProcessInputRuntimeBase instance.
func newProcessInputRuntimeBase[T any, P ci.IProcessInput[T]](name string, object *P, function ci.IValidationFunc[T], waitFor, restart bool, logger l.Logger, debug bool) *ProcessInputRuntimeBase[T, P] {
	cfg := newProcessConfig(name, waitFor, restart, "runtime", nil, logger, debug)
	return &ProcessInputRuntimeBase[T, P]{
		Reference:          cfg.Reference,
		Logger:             cfg.Logger,
		Mutexes:            cfg.Mutexes,
		Object:             object,
		Function:           function,
		ProcessInputConfig: cfg,
	}
}

// NewProcessInputRuntimeBase creates a new ProcessInputRuntimeBase instance.
func NewProcessInputRuntimeBase[T any, P ci.IProcessInput[T]](name string, object *P, function ci.IValidationFunc[T], waitFor, restart bool, logger l.Logger, debug bool) ci.IProcessInputRuntimeBase[T, P] {
	cfg := newProcessConfig(name, waitFor, restart, "runtime", nil, logger, debug)
	return &ProcessInputRuntimeBase[T, P]{
		Reference:          cfg.Reference,
		Logger:             cfg.Logger,
		Mutexes:            cfg.Mutexes,
		Object:             object,
		Function:           function,
		ProcessInputConfig: cfg,
	}
}

// GetObjectType returns the type of the object.
func (pi *ProcessInputRuntimeBase[T, P]) GetObjectType() reflect.Type { return reflect.TypeFor[T]() }

// GetObject returns the object to pass to the command.
func (pi *ProcessInputRuntimeBase[T, P]) GetObject() *P {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.Object
}

// GetFunction returns the custom function to wrap the command.
func (pi *ProcessInputRuntimeBase[T, P]) GetFunction() ci.IValidationFunc[T] {
	if pi == nil {
		return nil
	}
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.Function
}
