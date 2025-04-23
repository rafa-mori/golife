package interfaces

import "github.com/google/uuid"

type IValidationResult interface {
	String() string
	GetID() uuid.UUID
	GetName() string
	GetIsValid() bool
	GetMessage() string
	GetError() error
	GetMetadata(key string) (any, bool)
	SetMetadata(key string, value any)
	GetAllMetadataKeys() []string
}
