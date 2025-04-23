package types

import (
	"fmt"
	gl "github.com/faelmori/golife/logger"
	"github.com/google/uuid"
	"reflect"
	"runtime"
)

type IReference interface {
	GetID() uuid.UUID
	GetName() string
	SetName(name string)
	String() string
	GetReference() *Reference
}

// Reference is a struct that holds the Reference ID and name.
type Reference struct {
	// refID is the unique identifier for this context.
	ID uuid.UUID
	// refName is the name of the context.
	Name string
}

// newReference is a function that creates a new Reference instance.
func newReference(name string) *Reference {
	if name == "" {
		pc, _, line, ok := runtime.Caller(1)
		if ok {
			fn := runtime.FuncForPC(pc)
			name = fmt.Sprintf("%s:%d", fn.Name(), line)
		} else {
			name = "unknown"
		}
	}
	return &Reference{
		ID:   uuid.New(),
		Name: name,
	}
}

// NewReference is a function that creates a new IReference instance.
func NewReference(name string) IReference {
	return newReference(name)
}

// String is a method that returns the string representation of the reference.
func (r *Reference) String() string {
	return fmt.Sprintf("ID: %s, Name: %s", r.ID.String(), r.Name)
}

// GetID is a method that returns the ID of the reference.
func (r *Reference) GetID() uuid.UUID {
	if r == nil {
		gl.Log("error", "GetID: reference does not exist (", reflect.TypeFor[Reference]().String(), ")")
		return uuid.Nil
	}
	return r.ID
}

// GetName is a method that returns the name of the reference.
func (r *Reference) GetName() string {
	if r == nil {
		gl.Log("error", "GetName: reference does not exist (", reflect.TypeFor[Reference]().String(), ")")
		return ""
	}
	return r.Name
}

// SetName is a method that sets the name of the reference.
func (r *Reference) SetName(name string) {
	if r == nil {
		gl.Log("error", "SetName: reference does not exist (", reflect.TypeFor[Reference]().String(), ")")
		return
	}
	r.Name = name
}

// GetReference is a method that returns the reference struct (non-interface).
func (r *Reference) GetReference() *Reference {
	if r == nil {
		gl.Log("error", "GetReference: reference does not exist (", reflect.TypeFor[Reference]().String(), ")")
		return nil
	}
	return r
}
