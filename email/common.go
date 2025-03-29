package email

import (
	"fmt"
)

// ExtraData is a map that holds extra headers for an email, where the key is a string and the value is of type T.
type ExtraData[T any] map[string]T

// ExtraFields is an interface that defines methods for managing extra headers in an email.
type ExtraFields[T any] interface {
	GetData(key string) (T, error)
	SetData(key string, value T) error
	GetAllData() ExtraData[T]
	SetAllData(data ExtraData[T]) error
	DeleteData(key string) error
}

// EmailExtraFields is a struct that implements the ExtraFields interface for email headers.
type EmailExtraFields[T any] struct{ data ExtraData[T] }

func (e *EmailExtraFields[T]) GetData(key string) (T, error) {
	value, exists := e.data[key]
	if !exists {
		return value, fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}

func (e *EmailExtraFields[T]) SetData(key string, value T) error {
	e.data[key] = value
	return nil
}

func (e *EmailExtraFields[T]) GetAllData() ExtraData[T] {
	return e.data
}

func (e *EmailExtraFields[T]) SetAllData(data ExtraData[T]) error {
	e.data = data
	return nil
}

func (e *EmailExtraFields[T]) DeleteData(key string) error {
	delete(e.data, key)
	return nil
}
