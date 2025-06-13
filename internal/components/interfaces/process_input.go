package interfaces

import (
	"context"
	l "github.com/rafa-mori/logz"
)

type IProcessInput[T any] interface {
	Serialize(format string) ([]byte, error)
	Deserialize(data []byte, format string) error
	GetLogger() l.Logger
	GetReference() IReference
	GetName() string
	Validate() IValidationResult
	Send(string, any) error
	Receive(ctx context.Context, cb any) (any, error)
	SaveToFile(filename, format string) error
	LoadFromFile(filename, format string) error
}
