package property

import (
	"fmt"
	a "github.com/faelmori/golife/internal/routines/agents"
	"github.com/faelmori/golife/services"
	"reflect"
	"sync"
	"sync/atomic"
)

type KubexProperty[T any] struct {
	// PropertyBase defines the base interface for a property.
	PropertyBase
	// PropertyChanCtl defines the channel control interface for a property.
	PropertyChanCtl[any]
	// VoValue defines the value object interface for a property.
	VoValue[any]
	// Property defines the property interface for a property.
	//Property[T]

	// MainMutex is a mutex for synchronizing access to the property.
	MainMutex sync.Mutex
	// ChanMutex is a mutex for synchronizing access to the property.
	ChanMutex sync.Mutex
	// ListenerMutex is a mutex for synchronizing access to the property.
	ListenerMutex sync.Mutex
	// LayerMutex is a mutex for synchronizing access to the property.
	LayerMutex sync.Mutex
	// Cond is a condition variable for synchronization.
	Cond atomic.Pointer[sync.Cond]

	// name is the name of the property.
	name string
	// metadata stores additional information about the property.
	metadata Metadata
	// value is the current value of the property.
	value atomic.Pointer[T]
	// IChannel is the channel for receiving updates to the property value.
	chanCtl services.IChannel[T, int]
	// validators is a list of validation functions for the property value.
	validators []func(any) error
	// listeners is a list of change listeners for the property value.
	listeners map[string]ChangeListener[T]
	// cond is a condition variable for synchronization.
	cond *sync.Cond
	// Properties is a map of additional properties associated with this property.
	Properties map[string]any
}

// NewProperty creates a new basic property with the specified name and optional initial value.
func NewProperty[T any](name string, value *T) *KubexProperty[T] {
	var defaultValue T
	kbxProp := KubexProperty[T]{
		name:       name,
		metadata:   make(Metadata),
		validators: make([]func(any) error, 0),
		listeners:  make(map[string]ChangeListener[T]),
		Properties: make(map[string]any),
	}
	if value == nil {
		return &kbxProp
	} else {
		if reflect.TypeFor[T]() == reflect.TypeOf(value) {
			defaultValue = reflect.New(reflect.TypeFor[T]()).Interface().(T)
		} else {
			defaultValue = *value
		}
		_ = kbxProp.SetValue(defaultValue, nil)
		return &kbxProp
	}
}

// NewPropertyWithMethods creates a new full property interface with the specified name and optional initial value.
func NewPropertyWithMethods[T any](name string, value *T) Property[T] {
	var defaultValue T
	kbxProp := KubexProperty[T]{
		name:       name,
		metadata:   make(Metadata),
		validators: make([]func(any) error, 0),
	}
	if value == nil {
		//return reflect.ValueOf(&kbxProp).Interface().(Property[T])
		return &kbxProp
	} else {
		if reflect.TypeFor[T]() == reflect.TypeOf(value) {
			defaultValue = reflect.New(reflect.TypeFor[T]()).Interface().(T)
		} else {
			defaultValue = *value
		}
		_ = kbxProp.SetValue(defaultValue, nil)
		return &kbxProp
	}
}

func (bp *KubexProperty[T]) GetName() string { return bp.name }

func (bp *KubexProperty[T]) GetType() reflect.Type { return reflect.TypeFor[T]() }

func (bp *KubexProperty[T]) GetStringType() string {
	if bp.GetType() == nil {
		return ""
	}
	return bp.GetType().String()
}

func (bp *KubexProperty[T]) GetMetadata(key string) (interface{}, bool) {
	if key == "" {
		if len(bp.metadata) == 0 {
			return nil, false
		}
		return bp.metadata, true
	}
	value, exists := bp.metadata[key]
	return value, exists
}

func (bp *KubexProperty[T]) SetMetadata(key string, value interface{}) {
	if bp.metadata == nil {
		bp.metadata = make(Metadata)
	}
	bp.metadata[key] = value
}

func (bp *KubexProperty[T]) GetValueWithType() (any, reflect.Type) {
	v := bp.value.Load()
	if v == nil {
		return *new(T), nil
	}
	if reflect.TypeFor[T]() == reflect.TypeOf(v) {
		return *v, reflect.TypeFor[T]()
	}
	return *new(T), nil
}

func (bp *KubexProperty[T]) GetValue() *T {
	v := bp.value.Load()
	if v == nil {
		return new(T)
	}
	if reflect.TypeFor[T]() == reflect.TypeOf(v) {
		return v
	}
	return new(T)
}

func (bp *KubexProperty[T]) SetValue(value T, cb func(T) error) error {
	oldValue := bp.value.Load()
	var old *T
	if oldValue == nil {
		old = new(T)
	} else {
		old = reflect.ValueOf(oldValue).Interface().(*T)
	}
	if err := bp.validateAndSet(value); err != nil {
		return err
	}
	bp.notifyListeners(old, value, &EventMetadata{
		Timestamp: "2023-10-01T00:00:00Z",
		Source:    "KubexProperty",
		Details: map[string]interface{}{
			"event": "SetValue",
		},
	})
	if cb != nil {
		return cb(value)
	}
	if bp.chanCtl != nil {
		if ch, _ := bp.chanCtl.GetChan(); ch != nil {
			ch <- value
		}
	}
	return nil
}

func (bp *KubexProperty[T]) GetChannel() services.IChannel[T, int] { return bp.chanCtl }

func (bp *KubexProperty[T]) SetChannel(channel services.IChannel[T, int]) {
	bp.chanCtl = channel
}

func (bp *KubexProperty[T]) GetChannelValue() any {
	if lv, tp, err := bp.chanCtl.GetLast(); err == nil {
		if reflect.TypeFor[T]() == tp {
			return lv.(*T)
		}
	}
	return nil

}

func (bp *KubexProperty[T]) SetDefaultValue(value *T) error {
	if reflect.TypeFor[T]() == reflect.TypeOf(value) {
		bp.value.Store(value)
		return nil
	}
	return fmt.Errorf("type mismatch: expected %s, got %s", reflect.TypeFor[T]().String(), reflect.TypeOf(value).String())
}

func (bp *KubexProperty[T]) AddValidator(name string, validator func(T) error) error {
	if reflect.ValueOf(validator).IsNil() {
		return fmt.Errorf("validator cannot be nil")
	}
	if reflect.TypeFor[T]() != reflect.TypeOf(validator) {
		return fmt.Errorf("type mismatch: expected %s, got %s", reflect.TypeFor[T]().String(), reflect.TypeOf(validator).String())
	}
	innerValidator := func(v any) error {
		vT := reflect.ValueOf(v).Interface().(T)
		return validator(vT)
	}
	bp.validators = append(bp.validators, innerValidator)
	return nil
}

func (bp *KubexProperty[T]) AddListener(name string, listener *ChangeListener[T]) error {
	if _, exists := bp.listeners[name]; exists {
		return fmt.Errorf("listener with name %s already exists", name)
	}
	if reflect.ValueOf(listener).IsNil() {
		var changeListener ChangeListener[T] = func(oldValue *T, newValue T, metadata *EventMetadata) ListenerResponse {
			if reflect.TypeFor[T]() == reflect.TypeOf(oldValue) && reflect.TypeFor[T]() == reflect.TypeOf(newValue) {
				meta := &EventMetadata{
					Timestamp: "2023-10-01T00:00:00Z",
					Source:    "KubexProperty",
					Details: map[string]interface{}{
						"event": "AddListener",
					},
				}
				lll := *listener
				return lll(oldValue, newValue, meta)
			} else {
				return ListenerResponse{
					Success:  false,
					ErrorMsg: "Type mismatch in listener",
					Metadata: &EventMetadata{
						Timestamp: "2023-10-01T00:00:00Z",
						Source:    "KubexProperty",
						Details: map[string]interface{}{
							"event": "AddListener",
							"error": "Type mismatch",
						},
					},
				}
			}
		}
		bp.listeners[name] = changeListener
	}
	if reflect.TypeFor[T]() != reflect.TypeOf(listener) {
		return fmt.Errorf("type mismatch: expected %s, got %s", reflect.TypeFor[T]().String(), reflect.TypeOf(listener).String())

	}
	if reflect.ValueOf(listener).IsNil() {
		return fmt.Errorf("listener cannot be nil")
	}
	bp.listeners[name] = *listener
	return nil
}

func (bp *KubexProperty[T]) BroadcastChange(oldValue *T, newValue T) {
	if bp.chanCtl != nil {
		if ch, _ := bp.chanCtl.GetChan(); ch != nil {
			ch <- newValue
		}
	}
	bp.notifyListeners(oldValue, newValue, &EventMetadata{
		Timestamp: "2023-10-01T00:00:00Z",
		Source:    "KubexProperty",
		Details: map[string]interface{}{
			"event": "BroadcastChange",
		},
	})
}

func (bp *KubexProperty[T]) SetChannelValue(a any) {
	if reflect.TypeFor[T]() == reflect.TypeOf(a) {
		if err := bp.chanCtl.SetLast(a); err != nil {
			fmt.Printf("Error setting channel value: %s\n", err.Error())
			return
		}
	} else {
		fmt.Printf("Type mismatch: expected %s, got %s\n", reflect.TypeFor[T]().String(), reflect.TypeOf(a).String())
	}
}

func (bp *KubexProperty[T]) GetChannelType() reflect.Type {
	if bp.chanCtl != nil {
		return bp.chanCtl.GetType()
	}
	return nil
}

func (bp *KubexProperty[T]) SetChannelType(s string, newType reflect.Type) {
	if bp.chanCtl != nil {
		ch, _ := bp.chanCtl.GetChan()
		size := len(ch)
		typ := reflect.TypeOf(newType.Elem())
		val := reflect.New(typ).Interface()
		bp.chanCtl = a.NewChannel[any, int](s, &val, size)
	}
}

func (bp *KubexProperty[T]) GetChannelName() string {
	//TODO implement me
	panic("implement me")
}

func (bp *KubexProperty[T]) AddChainedListener(s string, s2 string, listener *ChangeListener[T]) error {
	if _, exists := bp.listeners[s]; exists {
		return fmt.Errorf("listener with name %s already exists", s)
	}
	if _, exists := bp.listeners[s2]; !exists {
		return fmt.Errorf("listener with name %s does not exist", s2)
	}
	var ltn ChangeListener[T] = func(oldValue *T, newValue T, metadata *EventMetadata) ListenerResponse {
		if reflect.TypeFor[T]() == reflect.TypeOf(oldValue) && reflect.TypeFor[T]() == reflect.TypeOf(newValue) {
			meta := &EventMetadata{
				Timestamp: "2023-10-01T00:00:00Z",
				Source:    "KubexProperty",
				Details: map[string]interface{}{
					"event": "AddChainedListener",
				},
			}

			lll := *listener

			return lll(oldValue, newValue, meta)
		} else {
			return ListenerResponse{
				Success:  false,
				ErrorMsg: "Type mismatch in listener",
				Metadata: &EventMetadata{
					Timestamp: "2023-10-01T00:00:00Z",
					Source:    "KubexProperty",
					Details: map[string]interface{}{
						"event": "AddChainedListener",
						"error": "Type mismatch",
					},
				},
			}
		}
	}
	bp.listeners[s] = ltn
	return nil
}

func (bp *KubexProperty[T]) RemoveListener(s string) error {
	if _, exists := bp.listeners[s]; !exists {
		return fmt.Errorf("listener with name %s does not exist", s)
	}
	delete(bp.listeners, s)
	return nil
}

func (bp *KubexProperty[T]) RemoveAllListeners() {
	if len(bp.listeners) > 0 {
		for k := range bp.listeners {
			delete(bp.listeners, k)
		}
	}
}

func (bp *KubexProperty[T]) WaitForCondition(check func() bool) {
	cond := bp.Cond.Load()
	if cond == nil {
		cond = sync.NewCond(&sync.Mutex{})
		bp.Cond.Store(cond)
	}
	cond.L.Lock()
	defer cond.L.Unlock()

	for !check() {
		cond.Wait()
	}
}

func (bp *KubexProperty[T]) NotifyCondition() {
	cond := bp.Cond.Load()
	if cond == nil {
		cond = sync.NewCond(&sync.Mutex{})
		bp.Cond.Store(cond)
	}

	cond.L.Lock()
	defer cond.L.Unlock()
	cond.Broadcast() // Sinaliza para todas as goroutines esperando
}

func (bp *KubexProperty[T]) notifyListeners(oldValue *T, newValue T, metadata *EventMetadata) {
	for _, listener := range bp.listeners {
		listenerResponse := listener(oldValue, newValue, metadata)
		if listenerResponse.ErrorMsg != "" {
			fmt.Printf("Error notifying listener: %s\n", listenerResponse.ErrorMsg)
			return
		}
	}
}

func (bp *KubexProperty[T]) validateAndSet(value any) error {
	if reflect.TypeFor[T]() != reflect.TypeOf(value) {
		return fmt.Errorf("type mismatch: expected %s, got %s", reflect.TypeFor[T]().String(), reflect.TypeOf(value).String())
	}
	for _, validator := range bp.validators {
		if err := validator(value.(T)); err != nil {
			return err
		}
	}
	bp.value.Store(value.(*T))
	return nil
}
