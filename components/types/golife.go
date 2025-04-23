package types

import (
	ci "github.com/faelmori/golife/components/interfaces"
	l "github.com/faelmori/logz"
)

// GoLife is a generic struct that implements the IGoLife interface.
type GoLife[T ci.ILifeCycle[ci.IProcessInput[any]]] struct {
	// Logger is the Logger instance for this GoLife instance.
	Logger l.Logger
	// Reference is the reference ID and name.
	*Reference
	// Mutexes is the mutexes for this GoLife instance.
	*Mutexes
	// Object is the object to pass to the command.
	Object *T
	// Properties is a map of properties for this GoLife instance.
	properties map[string]interface{}
	// metadata is a map of metadata for this GoLife instance.
	metadata map[string]interface{}
}
