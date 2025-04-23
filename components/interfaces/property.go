package interfaces

import (
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
	"reflect"
)

type IReference interface {
	GetID() uuid.UUID
	GetName() string
	SetName(name string)
	String() string
}

// IPropertyValBase is an interface that defines the methods for a property value.
type IPropertyValBase[T any] interface {
	GetLogger() l.Logger
	GetID() uuid.UUID
	GetName() string
	Value() *T
	StartCtl() <-chan string
	Type() reflect.Type
	Get(async bool) any
	Set(t *T) bool
	Clear() bool
	IsNil() bool
	Serialize(format string) ([]byte, error)
	Deserialize(data []byte, format string) error
}

// IProperty is an interface that defines the methods for a property.
type IProperty[T any] interface {
	GetName() string
	GetValue() T
	SetValue(v *T)
	GetReference() (uuid.UUID, string)
	Prop() IPropertyValBase[T]
	GetLogger() l.Logger
	Serialize(format string) ([]byte, error)
	Deserialize(data []byte, format string) error
	// Telemetry() *ITelemetry
}
