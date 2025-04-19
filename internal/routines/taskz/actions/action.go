package actions

import (
	"fmt"
	"github.com/faelmori/golife/internal/property"
	f "github.com/faelmori/golife/internal/property"
	c "github.com/faelmori/golife/internal/routines/agents"
	t "github.com/faelmori/golife/internal/types"

	"github.com/google/uuid"
	"reflect"
	"sync"
	"time"
)

type ChangeListener[T any] interface{ f.ChangeListener[T] }

// IAction defines the interface for an action with methods for execution, cancellation, and status retrieval.
type IAction[T any] interface {
	GetRef() uuid.UUID
	GetID() string
	GetType() string
	GetResults() map[string]t.IResult
	GetStatus() string
	GetErrors() []error
	IsRunning() bool
	CanExecute() bool
	Execute() error
	Cancel() error
	GetErrorChannel() chan error
	GetResultChannel() chan T
	GetDoneChannel() chan any
	GetCancelChannel() chan any
	GetProperties() map[string]f.Property[T]
	GetProperty(string) f.Property[T]
	SetProperty(string, any) error
	GetTask() func(T) error
	SetTask(func(T) error)
	SafeSend(string, any) error
	SafeReceive(string) (any, error)
}

// Action is a common struct that implements the IAction interface.
type Action[T any] struct {
	IAction[T]
	// Old Fields - Good logic, without bugs until now
	mu        sync.RWMutex // Mutex to ensure thread-safe access to the struct fields.
	ref       uuid.UUID    // Unique identifier for the action.
	ID        string       // Unique identifier of the action.
	Type      string       // Type of the action.
	isRunning bool         // Indicates whether the action is currently running.

	// New Fields - Let's do it the insane (moderate) way! hehehehe
	Errors     []error                         // List of errors associated with the action.
	Results    map[string]t.IResult            // Map of results associated with the action.
	mapChan    map[string]c.IChannel[any, int] // Map of channels associated with the action.
	Properties map[string]property.Property[T] // Map of properties associated with the action.
	properties map[string]interface{}          // Map of properties associated with the action.
	task       func(T) error                   // Task associated with the action. .
	data       T                               // Data associated with the action.
}

// NewAction creates a new action with the specified type.
// Parameters:
//   - actionType: The type of the action to create.
//
// Returns:
//   - IAction: A new instance of the Action struct.
func NewAction[T any](identifier string, actionType string, data *T, ev func(T) error) IAction[T] {
	uid, uidErr := uuid.NewUUID()
	if uidErr != nil {
		fmt.Printf("Error generating UUID: %v\n", uidErr)
		return nil
	}
	status := "Pending"

	actA := &Action[T]{
		ref:       uid,
		ID:        identifier,
		Type:      actionType,
		isRunning: false,
		Errors:    make([]error, 0),
		Results:   make(map[string]t.IResult),

		mapChan:    make(map[string]c.IChannel[any, int]),
		Properties: make(map[string]property.Property[T]),
		properties: map[string]interface{}{
			"status":     property.NewProperty[string]("status", &status),
			"identifier": property.NewProperty[string]("identifier", &identifier),
			"actionType": property.NewProperty[string]("actionType", &actionType),
			"data":       property.NewProperty[T]("data", data),
		},

		task: ev,
		data: *data,
	}

	// Initialize the map of channels

	actA.mapChan["cancel"] = c.NewChannel[struct{}, int](actA.ref.String(), nil, 10)
	actA.mapChan["result"] = c.NewChannel[t.IResult, int](actA.ref.String(), nil, 10)
	actA.mapChan["error"] = c.NewChannel[any, int](actA.ref.String(), nil, 10)
	actA.mapChan["done"] = c.NewChannel[struct{}, int](actA.ref.String(), nil, 10)

	// Initialize the map of properties

	return actA
}

// GetRef retrieves the unique identifier of the action.
// Returns:
//   - uuid.UUID: The unique identifier of the action.
func (ac *Action[T]) GetRef() uuid.UUID {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.ref
}

// GetID retrieves the unique identifier of the action.
// Returns:
//   - string: The unique identifier of the action.

func (ac *Action[T]) GetID() string {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.ID
}

// GetType retrieves the type of the action.
// Returns:
//   - string: The type of the action.
func (ac *Action[T]) GetType() string {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.Type
}

// GetResults retrieves a map of results associated with the action.
// Returns:
//   - map[string]t.IResult: The map of results.
func (ac *Action[T]) GetResults() map[string]t.IResult {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.Results
}

// GetStatus retrieves the current status of the action.
// Returns:
//   - string: The current status of the action.
func (ac *Action[T]) GetStatus() string {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	if ac.Properties["status"] == nil {
		return "Unknown"
	}
	return fmt.Sprintf("%v", ac.Properties["status"].GetValue())
}

// GetErrors retrieves a list of errors associated with the action.
// Returns:
//   - []error: The list of errors.
func (ac *Action[T]) GetErrors() []error {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.Errors
}

// IsRunning checks if the action is currently running.
// Returns:
//   - bool: True if the action is running, false otherwise.
func (ac *Action[T]) IsRunning() bool {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.isRunning
}

// CanExecute checks if the action can be executed.
// Returns:
//   - bool: True if the action can be executed, false otherwise.
func (ac *Action[T]) CanExecute() bool {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return !ac.isRunning
}

// Execute performs the action and sets its running state to true.
// Returns:
//   - error: An error if the action is already running, nil otherwise.
func (ac *Action[T]) Execute() error {
	start := time.Now()
	ac.mu.Lock()
	if ac.isRunning {
		ac.mu.Unlock()
		return fmt.Errorf("action already running")
	}
	ac.isRunning = true
	ac.mu.Unlock()

	val := ac.Properties["data"].GetValue()
	valT := reflect.ValueOf(val).Interface().(T)
	err := ac.task(valT)
	duration := time.Since(start)

	if err != nil {
		ac.Errors = append(ac.Errors, err)
		fmt.Printf("Action %s failed in %v\n", ac.ID, duration)
		return err
	}
	fmt.Printf("Action %s succeeded in %v\n", ac.ID, duration)
	return nil
}

// Cancel cancels the action and sets its running state to false.
// Returns:
//   - error: An error if the action is not running, nil otherwise.
func (ac *Action[T]) Cancel() error {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	if !ac.isRunning {
		return nil
	}
	ac.isRunning = false
	return nil
}

// GetResultChannel retrieves the channel used to communicate action results.
// Returns:
//   - chan t.IResult: The result channel (currently nil).
func (ac *Action[T]) GetResultChannel() chan T {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	if ac.mapChan["result"] == nil {
		return nil
	}
	if ch, _ := ac.mapChan["result"].GetChan(); ch != nil {
		if reflect.TypeOf(ch) == reflect.TypeOf(make(chan T)) {
			return reflect.ValueOf(ch).Interface().(chan T)
		}
	}
	return nil
}

// GetErrorChannel retrieves the channel used to communicate errors.
// Returns:
//   - chan error: The error channel (currently nil).
func (ac *Action[T]) GetErrorChannel() chan error {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	if ac.mapChan["error"] == nil {
		return nil
	}
	if ch, _ := ac.mapChan["error"].GetChan(); ch != nil {
		if reflect.TypeOf(ch) == reflect.TypeOf(make(chan error)) {
			return reflect.ValueOf(ch).Interface().(chan error)
		}
	}
	return nil
}

// GetDoneChannel retrieves the channel used to signal action completion.
// Returns:
//   - chan struct{}: The done channel (currently nil).
func (ac *Action[T]) GetDoneChannel() chan any {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return nil
}

// GetCancelChannel retrieves the channel used to signal action cancellation.
func (ac *Action[T]) GetCancelChannel() chan any {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	if ac.mapChan["cancel"] == nil {
		return nil
	}
	if ch, _ := ac.mapChan["cancel"].GetChan(); ch != nil {
		if reflect.TypeOf(ch) == reflect.TypeOf(make(chan struct{})) {
			return ch
		}
	}
	return nil
}

// GetProperties retrieves a map of properties associated with the action.
// Returns:
//   - map[string]types.Property[any]: The map of properties.
func (ac *Action[T]) GetProperties() map[string]f.Property[T] {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.Properties
}

// GetProperty retrieves a specific property by its name.
// Parameters:
//   - name: The name of the property to retrieve.
//
// Returns:
//   - types.Property[any]: The property associated with the specified name.
func (ac *Action[T]) GetProperty(name string) f.Property[T] {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	if prop, ok := ac.Properties[name]; ok {
		var fn f.ChangeListener[T] = property.NewListener[T]("status", func(oldValue *T, newValue T, metadata *property.EventMetadata) property.ListenerResponse {
			fmt.Printf("Action %s changed status: %v -> %v\n", ac.ID, oldValue, newValue)
			return property.ListenerResponse{
				Success:  true,
				Metadata: metadata,
				ErrorMsg: "",
			}
		}).N
		if prop != nil {
			err := prop.AddListener("status", &fn)
			if err != nil {
				return nil
			}
			return prop
		}
	}
	return nil
}

// SetProperty sets a specific property by its name.
// Parameters:
//   - name: The name of the property to set.
//   - value: The value to set for the property.
//
// Returns:
//   - error: An error if the property is not found, nil otherwise.
func (ac *Action[T]) SetProperty(name string, value any) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	if prop, ok := ac.Properties[name]; ok {
		valueT := reflect.ValueOf(value).Interface().(T)
		return prop.SetValue(valueT, nil)
	}
	return nil
}

// GetTask retrieves the task associated with the action.
// Returns:
//   - func(T) error: The task function.
func (ac *Action[T]) GetTask() func(T) error {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.task
}

// SetTask sets the task associated with the action.
// Parameters:
//   - task: The task function to set.
func (ac *Action[T]) SetTask(task func(T) error) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.task = task
}

// SafeSend a value to the specified channel.
// Parameters:
//   - channel: The name of the channel to send the value to.
//   - value: The value to send.
//
// Returns:
//   - error: An error if the channel is not found or closed, nil otherwise.
//
// The function uses a read lock to ensure thread-safe access to the channel map.
func (ac *Action[T]) SafeSend(channel string, value any) error {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	if ch, ok := ac.mapChan[channel]; ok {
		if ch == nil {
			return fmt.Errorf("channel %s is closed", channel)
		}
		return ch.Send(value)
	}
	return fmt.Errorf("channel %s not found", channel)
}

// SafeReceive retrieves a value from the specified channel.
// Parameters:
//   - channel: The name of the channel to receive the value from.
//
// Returns:
//   - any: The received value.
//   - error: An error if the channel is not found or closed, nil otherwise.
//
// The function uses a read lock to ensure thread-safe access to the channel map.
func (ac *Action[T]) SafeReceive(channel string) (any, error) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	if ch, ok := ac.mapChan[channel]; ok {
		if ch == nil {
			return nil, fmt.Errorf("channel %s is closed", channel)
		}
		if cha, _ := ch.GetChan(); cha != nil {
			if reflect.TypeOf(cha) == reflect.TypeOf(make(chan any)) {
				return <-cha, nil
			}
			return nil, fmt.Errorf("channel %s is not of type chan any", channel)
		} else {
			return nil, fmt.Errorf("channel %s is closed", channel)
		}
	}
	return nil, fmt.Errorf("channel %s not found", channel)
}
