package interfaces

import (
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
)

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
