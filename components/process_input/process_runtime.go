package process_input

import (
	t "github.com/faelmori/golife/components/types"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"reflect"
)

// ProcessRuntimeBase is a struct that holds the runtime information for the process.
type ProcessRuntimeBase[T any] struct {
	// Logger is the logger for the process
	Logger l.Logger
	// Reference is the reference for the process with ID and Name
	*t.Reference
	// Mutexes is the mutex for the process
	*t.Mutexes
	// ObjectType is the type of the object
	ObjectType reflect.Type
	// Object is the object to pass to the command
	Object T
	// Function is a custom function to wrap the command
	Function *t.ValidationFunc[T]
}

// NewProcessRuntimeBase creates a new ProcessRuntimeBase instance.
func NewProcessRuntimeBase[T any](name string, object T, function *t.ValidationFunc[T], logger l.Logger, debug bool) *ProcessRuntimeBase[T] {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	if debug {
		gl.SetDebug(debug)
	}
	ref := t.NewReference(name)
	mu := t.NewMutexes()
	npi := &ProcessRuntimeBase[T]{
		Reference: ref,
		Logger:    logger,
		Mutexes:   mu,
		Object:    object,
		Function:  function,
	}
	return npi
}

// GetObjectType returns the type of the object.
func (pi *ProcessRuntimeBase[T]) GetObjectType() reflect.Type {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	return pi.ObjectType
}

// GetObject returns the object to pass to the command.
func (pi *ProcessRuntimeBase[T]) GetObject() T {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	return pi.Object
}

// GetFunction returns the custom function to wrap the command.
func (pi *ProcessRuntimeBase[T]) GetFunction() *t.ValidationFunc[T] {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	return pi.Function
}
