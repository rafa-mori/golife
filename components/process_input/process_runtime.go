package process_input

import (
	t "github.com/faelmori/golife/components/types"
	l "github.com/faelmori/logz"
	"reflect"
)

// ProcessRuntimeBase is a struct that holds the runtime information for the process.
type ProcessRuntimeBase[T any, P any] struct {
	// Logger is the logger for the process
	Logger l.Logger
	// Mutexes is the mutex for the process
	*t.Mutexes
	// Reference is the reference for the process with ID and Name
	*t.Reference `json:"reference" yaml:"reference" xml:"reference" gorm:"reference"`
	// ProcessConfig is the configuration for the process
	*ProcessConfig[P] `json:"process_config" yaml:"process_config" xml:"process_config" gorm:"process_config"`
	// Object is the object to pass to the command
	Object *P `json:"object,omitempty" yaml:"object,omitempty" xml:"object,omitempty" gorm:"object,omitempty"`
	// Function is a custom function to wrap the command
	Function *t.ValidationFunc[ProcessInput[P]] `json:"function,omitempty" yaml:"function,omitempty" xml:"function,omitempty" gorm:"function,omitempty"`
}

// newProcessRuntimeBase creates a new ProcessRuntimeBase instance.
func newProcessRuntimeBase[T ProcessInput[P], P any](name string, object *P, function *t.ValidationFunc[ProcessInput[P]], waitFor, restart bool, logger l.Logger, debug bool) *ProcessRuntimeBase[T, P] {
	cfg := NewProcessConfig[P](name, waitFor, restart, "runtime", nil, logger, debug)
	npi := ProcessRuntimeBase[T, P]{
		Reference:     cfg.Reference,
		Logger:        cfg.Logger,
		Mutexes:       cfg.Mutexes,
		Object:        object,
		Function:      function,
		ProcessConfig: cfg,
	}
	return &npi
}

// GetObjectType returns the type of the object.
func (pi *ProcessRuntimeBase[T, P]) GetObjectType() reflect.Type { return reflect.TypeFor[T]() }

// GetObject returns the object to pass to the command.
func (pi *ProcessRuntimeBase[T, P]) GetObject() *P {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.Object
}

// GetFunction returns the custom function to wrap the command.
func (pi *ProcessRuntimeBase[T, P]) GetFunction() *t.ValidationFunc[ProcessInput[P]] {
	if pi == nil {
		return nil
	}
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.Function
}
