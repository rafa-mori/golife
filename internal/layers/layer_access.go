package layers

import (
	"fmt"
	p "github.com/faelmori/golife/components/types"
	"reflect"
	"sort"
	"sync"
)

// Layer is a struct that holds the layer information.
type Layer[T any] struct {
	// Mu is the mutexes for this layer.
	Mu *p.Mutexes
	// Scope is the scope of this layer.
	Scope *p.Reference
	// Events is a map of access events.
	Events map[string]map[int]p.ValidationFunc[any]
	// Listeners is a map of access listeners.
	Listeners map[int]p.ValidationFunc[any]
	// State is the state of this layer.
	State *p.ValidationResult
	// Object is the object of this layer.
	Object *T
}

// NewLayer creates a new Layer instance with the provided scope.
func NewLayer[T any](scope string) *Layer[T] {
	// Create a new Reference instance based on the provided scope
	ref := p.NewReference(scope)
	// Create a new Mutexes instance
	mu := p.NewMutexes()
	// Create a new Layer instance
	return &Layer[T]{
		Mu:     mu,
		Scope:  ref,
		Events: make(map[string]map[int]p.ValidationFunc[any]),
		State: &p.ValidationResult{
			IsValid: true,
			Message: "Access layer initialized successfully",
			Error:   nil,
		},
	}
}

// LayerType is a function that returns the type of the Layer.
func (a *Layer[T]) LayerType() reflect.Type { return reflect.TypeFor[T]() }

// LayerObject returns the object of the Layer.
func (a *Layer[T]) LayerObject() *T {
	a.Mu.MuRLock()
	defer a.Mu.MuRUnlock()

	// Return the object of the Layer
	return a.Object
}

// AddLayerEvent adds an event to the Layer.
func (a *Layer[T]) AddLayerEvent(name string, event *p.ValidationFunc[any]) {
	if event == nil {
		return
	}

	a.Mu.MuLock()
	defer a.Mu.MuUnlock()

	// Check if the event already exists
	if _, exists := a.Events[name]; !exists {
		a.Events[name] = make(map[int]p.ValidationFunc[any])
	}
	// Check if the event already exists
	if _, exists := a.Events[name][event.Priority]; exists {
		return
	}
	// MuAdd the event to the map
	a.Events[name][event.Priority] = *event

	fmt.Printf("Access event %s added with priority %d\n", name, event.Priority)
}

// GetLayerEvents returns the events for the given name.
func (a *Layer[T]) GetLayerEvents(name string) map[int]p.ValidationFunc[any] {
	a.Mu.MuRLock()
	defer a.Mu.MuRUnlock()

	// Check if the event exists
	if _, exists := a.Events[name]; !exists {
		return nil
	}

	// Return the events for the given name
	return a.Events[name]
}

// RemoveLayerEvent removes the event with the given name and priority.
func (a *Layer[T]) RemoveLayerEvent(name string, priority int) {
	a.Mu.MuLock()
	defer a.Mu.MuUnlock()

	// Check if the event exists
	if _, exists := a.Events[name]; !exists {
		return
	}

	// Remove the event with the given priority
	delete(a.Events[name], priority)

	// If there are no events left for the given name, remove the name from the map
	if len(a.Events[name]) == 0 {
		delete(a.Events, name)
	}
}

// ExecuteLayerEvent executes the event with the given name and arguments.
func (a *Layer[T]) ExecuteLayerEvent(name string, args ...any) (bool, *p.ValidationResult) {
	// Sort the events by priority
	sortedEvents := sortByPriority(a.Events[name], a.Mu.MuCtxM)
	if sortedEvents == nil {
		return false, &p.ValidationResult{
			Error:   fmt.Errorf("error sorting events for %s", name),
			IsValid: false,
			Message: fmt.Sprintf("Error sorting events for %s", name),
		}
	}

	accessEvents := a.GetLayerEvents(name)
	if accessEvents == nil {
		return false, &p.ValidationResult{
			Error:   fmt.Errorf("access event %s not found", name),
			IsValid: false,
			Message: fmt.Sprintf("Access events for %s not found", name),
		}
	}

	a.Mu.MuLock()
	defer a.Mu.MuUnlock()

	// Execute the events in order of priority
	for _, event := range accessEvents {
		if event.Func != nil {
			var input interface{} = nil
			if len(args) > 0 {
				input = args[0:]
			}
			eventResult := event.Func(input)
			return eventResult != nil, eventResult
		}
	}

	return false, nil
}

// AddLayerListener adds a listener to the Layer.
func (a *Layer[T]) AddLayerListener(listener *p.ValidationFunc[any]) {
	if listener == nil {
		return
	}
	a.Mu.MuLock()
	defer a.Mu.MuUnlock()
	// Check if the listener already exists
	if _, exists := a.Listeners[listener.Priority]; exists {
		return
	}
	a.Listeners[listener.Priority] = *listener
}

// GetLayerListeners returns the listeners for the Layer.
func (a *Layer[T]) GetLayerListeners() map[int]p.ValidationFunc[any] {
	a.Mu.MuRLock()
	defer a.Mu.MuRUnlock()
	return a.Listeners
}

// RemoveLayerListener removes the listener with the given priority.
func (a *Layer[T]) RemoveLayerListener(priority int) {
	a.Mu.MuLock()
	defer a.Mu.MuUnlock()

	// Check if the listener exists
	if _, exists := a.Listeners[priority]; !exists {
		return
	}

	// Remove the listener with the given priority
	delete(a.Listeners, priority)
}

// ExecuteLayerListener executes the listener with the given priority and arguments.
func (a *Layer[T]) ExecuteLayerListener(priority int, args ...any) (bool, *p.ValidationResult) {
	a.Mu.MuLock()
	defer a.Mu.MuUnlock()

	// Check if the listener exists
	if _, exists := a.Listeners[priority]; !exists {
		return false, &p.ValidationResult{
			Error:   fmt.Errorf("access listener with priority %d not found", priority),
			IsValid: false,
			Message: fmt.Sprintf("Access listener with priority %d not found", priority),
		}
	}

	// Execute the listener
	listener := a.Listeners[priority]
	if listener.Func != nil {
		var input interface{} = nil
		if len(args) > 0 {
			input = args[0:]
		}
		listenerResult := listener.Func(input)
		return listenerResult != nil, listenerResult
	}

	return false, nil
}

// NotifyLayerListeners notifies all listeners with the given arguments.
func (a *Layer[T]) NotifyLayerListeners(args ...any) {
	// Sort the listeners by priority
	sortedListeners := sortByPriority(a.Listeners, a.Mu.MuCtxM)
	if sortedListeners == nil {
		fmt.Println("Error sorting access listeners")
		return
	}

	// Notify all access listeners
	for priority, _ := range sortedListeners {
		if ok, err := a.ExecuteLayerListener(priority, args...); err != nil {
			fmt.Printf("Error executing access listener with priority %d: %v\n", priority, err)
		} else if !ok {
			fmt.Printf("Access listener with priority %d returned false\n", priority)
		} else {
			fmt.Printf("Access listener with priority %d executed successfully\n", priority)
		}
	}
}

// GetLayerState returns the state of the Layer.
func (a *Layer[T]) GetLayerState() *p.ValidationResult {
	a.Mu.MuRLock()
	defer a.Mu.MuRUnlock()

	// Return the state of the Layer
	return a.State
}

// Snapshot returns a string representation of the Layer.
func (a *Layer[T]) Snapshot() string {
	a.Mu.MuRLock()
	defer a.Mu.MuRUnlock()

	valid := "valid"
	if !a.State.IsValid {
		valid = "invalid"
	}
	return fmt.Sprintf("Layer %s (%s): %d events, %d listeners. ", a.Scope.Name, valid, len(a.Events), len(a.Listeners))
}

// sortByPriority sorts the given slice of ValidationFunc by priority.
func sortByPriority(funcMap map[int]p.ValidationFunc[any], mutex *sync.RWMutex) map[int]p.ValidationFunc[any] {
	mutex.RLock()
	defer mutex.RUnlock()

	// Sort the events by priority
	sortedEvents := make([]p.ValidationFunc[any], 0)
	for _, v := range funcMap {
		sortedEvents = append(sortedEvents, v)
	}

	// Sort the events by priority
	sort.Slice(sortedEvents, func(i, j int) bool {
		return sortedEvents[i].Priority < sortedEvents[j].Priority
	})

	// Create a new map with the sorted events
	sortedEventsMap := make(map[int]p.ValidationFunc[any])
	for _, v := range sortedEvents {
		sortedEventsMap[v.Priority] = v
	}

	// Update the Events map with the sorted events
	funcMap = sortedEventsMap

	// Return the sorted events
	return funcMap
}
