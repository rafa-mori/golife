package types

import (
	"context"
	f "github.com/faelmori/golife/internal/property"
)

// WrapperBasicGenericCallback is a function type that takes two values of type T and returns an error.
type WrapperBasicGenericCallback func(*BasicGenericCallback[any]) (*f.EventMetadata, error)

// BasicGenericCallback is a function type that takes two values of type T and returns an error.
type BasicGenericCallback[T any] func(context.Context, *T, ...any) error

// SwapGenericCallback is a function type that takes two values of type T and returns an error.
type SwapGenericCallback[T any] func(string, *T, T) error

// SwapGenericChannelCallback is a function type that takes two values of type T and returns a channel of type T.
type SwapGenericChannelCallback[T any] func(context string, oldValue *T, newValue T) <-chan T

// SwapGenericChannelCallbackWithError is a function type that takes two values of type T and returns a channel of type T and an error.
type SwapGenericChannelCallbackWithError[T any] func(context string, oldValue *T, newValue T) (<-chan T, error)
