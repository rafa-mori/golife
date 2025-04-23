package interfaces

import (
	l "github.com/faelmori/logz"
)

type IProcessInput[T any] interface {
	Serialize(format string) ([]byte, error)
	Deserialize(data []byte, format string) error
	GetLogger() l.Logger
	GetReference() IReference
	GetName() string
	Validate() IValidationResult
}
