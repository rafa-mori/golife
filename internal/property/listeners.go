package property

import (
	"fmt"
	"reflect"
	"time"
)

// ListenerResponse represents the response from a listener.
type ListenerResponse struct {
	Success  bool
	ErrorMsg string
	Metadata *EventMetadata
}

// ChangeListener is a function type that takes two values of type T and returns a ListenerResponse.
type ChangeListener[T any] func(*T, T, *EventMetadata) ListenerResponse

// NewListener creates a new ChangeListener with the specified name and listener function.
func NewListener[T any](name string, listener ChangeListener[T]) ChangeListener[T] {
	if listener == nil {
		return nil
	}
	var nLtn ChangeListener[T] = func(oldValue *T, newValue T, metadata *EventMetadata) ListenerResponse {
		if reflect.TypeFor[T]() == reflect.TypeOf(oldValue) && reflect.TypeFor[T]() == reflect.TypeOf(newValue) {
			res := listener(oldValue, newValue, metadata)
			if res.Success {
				return res
			} else {
				return ListenerResponse{
					Success:  false,
					ErrorMsg: res.ErrorMsg,
					Metadata: &EventMetadata{
						Timestamp: time.Now().String(),
						Source:    "WorkerPool",
						Details: map[string]interface{}{
							"event": "NewListener",
						},
					},
				}
			}
		}
		return ListenerResponse{
			Success:  false,
			ErrorMsg: fmt.Sprintf("type mismatch: expected %s, got %s", reflect.TypeFor[T]().String(), reflect.TypeOf(oldValue).String()),
			Metadata: &EventMetadata{
				Timestamp: time.Now().String(),
				Source:    "WorkerPool",
				Details: map[string]interface{}{
					"event": "NewListener",
				},
			},
		}
	}
	return nLtn
}

// N is a no-op function that returns a ListenerResponse.
func (cl ChangeListener[T]) N(oldValue *T, newValue T, metadata *EventMetadata) ListenerResponse {
	if cl == nil {
		return ListenerResponse{
			Success:  false,
			ErrorMsg: "ChangeListener is nil",
			Metadata: &EventMetadata{
				Timestamp: metadata.Timestamp,
				Source:    metadata.Source,
				Details: map[string]interface{}{
					"event": "ChangeListener",
				},
			},
		}
	}
	return cl(oldValue, newValue, metadata)
}

// M is a method that takes two values of type any and returns a ListenerResponse.
func (cl ChangeListener[T]) M(oldValue, newValue any, metadata *EventMetadata) ListenerResponse {
	if cl == nil {
		return ListenerResponse{Success: false, ErrorMsg: "Listener cannot be nil", Metadata: metadata}
	}
	listenerResponse := cl(oldValue.(*T), newValue.(T), metadata)
	if listenerResponse.ErrorMsg != "" {
		return ListenerResponse{Success: false, ErrorMsg: listenerResponse.ErrorMsg, Metadata: metadata}
	}
	return ListenerResponse{Success: true, Metadata: metadata}
}

func (cl ChangeListener[T]) Broadcast(oldValue *T, newValue T) {
	if cl == nil {
		return
	}
	listenerResponse := cl(oldValue, newValue, &EventMetadata{
		Timestamp: time.Now().String(),
		Source:    "WorkerPool",
		Details: map[string]interface{}{
			"event": "Broadcast",
		},
	})
	if listenerResponse.ErrorMsg != "" {
		fmt.Printf("Error broadcasting change: %s\n", listenerResponse.ErrorMsg)
	}
}
