package types

import (
	gl "github.com/rafa-mori/golife/logger"
	"reflect"
)

type ValidationListenerType string

const (
	ValidationListenerTypeBefore  ValidationListenerType = "before"  // Before validation
	ValidationListenerTypeAfter   ValidationListenerType = "after"   // After validation
	ValidationListenerTypeError   ValidationListenerType = "error"   // Error validation
	ValidationListenerTypeSuccess ValidationListenerType = "success" // Success validation
	ValidationListenerTypeDefault ValidationListenerType = "default" // Default validation
)

type ValidationFilterType string

const (
	ValidationFilterTypeEvent    ValidationFilterType = "event"    // Event filter
	ValidationFilterTypeListener ValidationFilterType = "listener" // Listener filter
	ValidationFilterTypeResult   ValidationFilterType = "result"   // Result filter
)

type ValidationListener struct {
	*Mutexes
	Filters   map[ValidationFilterType]func(*ValidationResult) bool
	Handlers  []func(*ValidationResult)
	Listeners map[Reference]map[ValidationListenerType]func(*ValidationResult)
}

func NewValidationListener() *ValidationListener {
	return &ValidationListener{
		Mutexes:   NewMutexesType(),
		Listeners: make(map[Reference]map[ValidationListenerType]func(*ValidationResult)),
		Filters:   make(map[ValidationFilterType]func(*ValidationResult) bool),
		Handlers:  []func(*ValidationResult){},
	}
}

func (vl *ValidationListener) AddFilter(filterType ValidationFilterType, filter func(*ValidationResult) bool) {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return
	}
	vl.Mutexes.MuLock()
	defer vl.Mutexes.MuUnlock()

	if vl.Filters == nil {
		vl.Filters = make(map[ValidationFilterType]func(*ValidationResult) bool)
	}
	if filter == nil {
		gl.Log("error", "RegisterListener: filter is nil")
		return
	}

	vl.Filters[filterType] = filter
}

func (vl *ValidationListener) RemoveFilter(filterType ValidationFilterType) {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return
	}
	vl.Mutexes.MuLock()
	defer vl.Mutexes.MuUnlock()

	delete(vl.Filters, filterType)
}

func (vl *ValidationListener) AddHandler(handler func(*ValidationResult)) {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return
	}
	vl.Mutexes.MuLock()
	defer vl.Mutexes.MuUnlock()

	vl.Handlers = append(vl.Handlers, handler)
}

func (vl *ValidationListener) RemoveHandler(handler func(*ValidationResult)) {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return
	}
	vl.Mutexes.MuLock()
	defer vl.Mutexes.MuUnlock()

	for i, h := range vl.Handlers {
		if reflect.ValueOf(h).Pointer() == reflect.ValueOf(handler).Pointer() {
			vl.Handlers = append(vl.Handlers[:i], vl.Handlers[i+1:]...)
			break
		}
	}
}

func (vl *ValidationListener) AddListener(reference Reference, listenerType ValidationListenerType, handler func(*ValidationResult)) {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return
	}
	vl.Mutexes.MuLock()
	defer vl.Mutexes.MuUnlock()

	if _, exists := vl.Listeners[reference]; !exists {
		vl.Listeners[reference] = make(map[ValidationListenerType]func(*ValidationResult))
	}
	vl.Listeners[reference][listenerType] = handler
}

func (vl *ValidationListener) RemoveListener(reference Reference, listenerType ValidationListenerType) {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return
	}
	vl.Mutexes.MuLock()
	defer vl.Mutexes.MuUnlock()

	if _, exists := vl.Listeners[reference]; exists {
		delete(vl.Listeners[reference], listenerType)
		if len(vl.Listeners[reference]) == 0 {
			delete(vl.Listeners, reference)
		}
	}
}

func (vl *ValidationListener) GetFilters() map[string]func(*ValidationResult) bool {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return nil
	}
	vl.Mutexes.MuLock()
	defer vl.Mutexes.MuUnlock()

	filters := make(map[string]func(*ValidationResult) bool)
	for k, v := range vl.Filters {
		if v == nil {
			gl.Log("error", "RegisterListener: filter is nil")
			continue
		}
		filters[string(k)] = v
	}
	return filters
}

func (vl *ValidationListener) GetHandlersByName(name string) []func(*ValidationResult) {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return nil
	}
	vl.Mutexes.MuLock()
	defer vl.Mutexes.MuUnlock()

	for _, handler := range vl.Handlers {
		if handler == nil {
			gl.Log("error", "RegisterListener: handler is nil")
			continue
		}
		if name == "" {
			gl.Log("error", "RegisterListener: name is empty")
			continue
		}
		if handler == nil {
			gl.Log("error", "RegisterListener: handler is nil")
			continue
		}
		return []func(*ValidationResult){handler}
	}
	return nil
}

func (vl *ValidationListener) GetHandlers() []func(*ValidationResult) {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return nil
	}
	vl.Mutexes.MuLock()
	defer vl.Mutexes.MuUnlock()

	handlers := make([]func(*ValidationResult), len(vl.Handlers))
	copy(handlers, vl.Handlers)
	return handlers
}

func (vl *ValidationListener) GetListeners() map[Reference]map[ValidationListenerType]func(*ValidationResult) {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return nil
	}
	vl.Mutexes.MuLock()
	defer vl.Mutexes.MuUnlock()

	listeners := make(map[Reference]map[ValidationListenerType]func(*ValidationResult))
	for k, v := range vl.Listeners {
		listeners[k] = v
	}
	return listeners
}

func (vl *ValidationListener) GetListenersByName(name string) map[ValidationListenerType]func(*ValidationResult) {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return nil
	}
	vl.Mutexes.MuLock()
	defer vl.Mutexes.MuUnlock()

	for k, v := range vl.Listeners {
		if k.GetName() == name {
			return v
		}
	}
	return nil
}

func (vl *ValidationListener) GetListenersKeys() map[string]Reference {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return nil
	}
	vl.Mutexes.MuLock()
	defer vl.Mutexes.MuUnlock()

	keys := make(map[string]Reference)
	for k := range vl.Listeners {
		keys[k.GetName()] = k
	}
	return keys
}

func (vl *ValidationListener) RegisterListener(reference Reference, handler func(*ValidationResult)) {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return
	}
	vl.Mutexes.MuLock()
	defer vl.Mutexes.MuUnlock()

	if handler == nil {
		gl.Log("error", "RegisterListener: handler is nil")
		return
	}

	if _, exists := vl.Listeners[reference]; !exists {
		vl.Listeners[reference] = make(map[ValidationListenerType]func(*ValidationResult))
	}
	vl.Listeners[reference][ValidationListenerTypeDefault] = handler
}

func (vl *ValidationListener) Trigger(event string, result *ValidationResult) {
	if vl == nil {
		gl.Log("error", "RegisterListener: ValidationListener is nil")
		return
	}

	vl.Mutexes.MuRLock()
	defer vl.Mutexes.MuRUnlock()

	if result == nil {
		gl.Log("error", "RegisterListener: result is nil")
		return
	}
	if event == "" {
		gl.Log("error", "RegisterListener: event is empty")
		return
	}

	if listenerZ := vl.GetListenersByName(event); listenerZ != nil {
		// Check event filters
		for _, filter := range vl.Filters {
			if filter == nil {
				gl.Log("error", "RegisterListener: filter is nil")
				continue
			}
			if !filter(result) {
				gl.Log("info", "RegisterListener: filter failed")
				return
			}
		}
		for _, listener := range listenerZ {
			if listener == nil {
				gl.Log("error", "RegisterListener: listener is nil")
				continue
			}
			// Check listener filters
			for _, filter := range vl.Filters {
				if filter == nil {
					gl.Log("error", "RegisterListener: filter is nil")
					continue
				}
				if !filter(result) {
					gl.Log("info", "RegisterListener: filter failed")
					return
				}
			}
			// Async dispatch
			go listener(result)
		}
	}
}
