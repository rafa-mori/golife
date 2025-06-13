package types

import (
	ci "github.com/rafa-mori/golife/internal/components/interfaces"
	"github.com/google/uuid"
	"reflect"

	"fmt"
	"sort"
	"sync"
)

type ValidationResult struct {
	*Mutexes
	*Reference
	IsValid  bool
	Message  string
	Error    error
	Metadata map[string]any
	Callback func(result *ValidationResult)
}

func newValidationResult(isValid bool, message string, metadata map[string]any, err error) *ValidationResult {
	if metadata == nil {
		metadata = make(map[string]any)
	}
	return &ValidationResult{
		Mutexes:   NewMutexesType(),
		Reference: newReference("ValidationResult"),
		IsValid:   isValid,
		Message:   message,
		Error:     err,
		Metadata:  metadata,
	}
}
func NewValidationResult(isValid bool, message string, metadata map[string]any, err error) ci.IValidationResult {
	return newValidationResult(isValid, message, metadata, err)
}

func (vr *ValidationResult) String() string {
	if vr == nil {
		return ""
	}
	vr.Mutexes.MuRLock()
	defer vr.Mutexes.MuRUnlock()
	if vr.IsValid {
		return "Validation is valid"
	}
	if vr.Error != nil {
		return fmt.Sprintf("Validation is invalid: %s", vr.Error.Error())
	}
	return fmt.Sprintf("Validation is invalid: %s", vr.Message)
}
func (vr *ValidationResult) GetID() uuid.UUID {
	if vr == nil {
		return uuid.Nil
	}
	return vr.Reference.GetID()
}
func (vr *ValidationResult) GetName() string {
	if !reflect.ValueOf(vr).IsValid() {
		return ""
	}
	vr.Mutexes.MuRLock()
	defer vr.Mutexes.MuRUnlock()
	return vr.Reference.GetName()
}
func (vr *ValidationResult) GetIsValid() bool {
	if vr == nil {
		return false
	}
	vr.Mutexes.MuRLock()
	defer vr.Mutexes.MuRUnlock()
	return vr.IsValid
}
func (vr *ValidationResult) GetMessage() string {
	if vr == nil {
		return ""
	}
	vr.Mutexes.MuRLock()
	defer vr.Mutexes.MuRUnlock()
	return vr.Message
}
func (vr *ValidationResult) GetMetadata(key string) (any, bool) {
	if !reflect.ValueOf(vr.Metadata).IsValid() {
		vr.Mutexes.MuLock()
		defer vr.Mutexes.MuUnlock()

		vr.Metadata = make(map[string]any)
		return nil, false
	}

	vr.Mutexes.MuRLock()
	defer vr.Mutexes.MuRUnlock()

	if key == "" {
		return vr.Metadata, true
	}
	value, exists := vr.Metadata[key]

	return value, exists
}
func (vr *ValidationResult) SetMetadata(key string, value any) {
	if vr == nil {
		return
	}
	vr.Mutexes.MuLock()
	defer vr.Mutexes.MuUnlock()

	if vr.Metadata == nil {
		vr.Metadata = make(map[string]any)
	}
	if !reflect.ValueOf(value).IsValid() {
		return
	}
	if key == "" {
		return
	} else if key == "all" {
		if vl, ok := value.(map[string]any); ok {
			vr.Metadata = vl
			return
		} else if vl, ok := value.(ValidationResult); ok {
			vr.Metadata = vl.Metadata
			return
		}
	}

	vr.Metadata[key] = value
}
func (vr *ValidationResult) GetAllMetadataKeys() []string {
	if vr == nil || vr.Metadata == nil {
		return nil
	}

	vr.Mutexes.MuRLock()
	defer vr.Mutexes.MuRUnlock()

	keys := make([]string, 0, len(vr.Metadata))
	for key := range vr.Metadata {
		keys = append(keys, key)
	}
	return keys
}
func (vr *ValidationResult) GetError() error {
	if vr == nil {
		return nil
	}
	return vr.Error
}

type ValidationFunc[T any] struct {
	Priority int
	Func     func(value *T, args ...any) ci.IValidationResult
	Result   ci.IValidationResult
}

func newValidationFunc[T any](priority int, f func(value *T, args ...any) ci.IValidationResult) *ValidationFunc[T] {
	return &ValidationFunc[T]{
		Priority: priority,
		Func:     f,
		Result:   nil,
	}
}
func NewValidationFunc[T any](priority int, f func(value *T, args ...any) ci.IValidationResult) ci.IValidationFunc[T] {
	validFunc := newValidationFunc[T](priority, nil)
	validFunc.Func = f
	return validFunc
}

func (vf *ValidationFunc[T]) GetPriority() int {
	if vf == nil {
		return -1
	}
	return vf.Priority
}
func (vf *ValidationFunc[T]) SetPriority(priority int) {
	if vf == nil {
		return
	}
	vf.Priority = priority
}
func (vf *ValidationFunc[T]) GetFunction() func(value *T, args ...any) ci.IValidationResult {
	if vf == nil {
		return nil
	}
	return vf.Func
}
func (vf *ValidationFunc[T]) SetFunction(f func(value *T, args ...any) ci.IValidationResult) {
	if vf == nil {
		return
	}
	vf.Func = f
}
func (vf *ValidationFunc[T]) GetResult() ci.IValidationResult {
	if vf == nil {
		return nil
	}
	return vf.Result
}
func (vf *ValidationFunc[T]) SetResult(result ci.IValidationResult) {
	if vf == nil {
		return
	}
	vf.Result = result
}

// Validation is a struct that holds the validation function and the errors.
type Validation[T any] struct {
	mu sync.RWMutex
	// isValid is a boolean that indicates if the value is valid.
	isValid bool
	// hasValidate is a boolean that indicates if the value will be validated.
	hasValidation bool
	// validatorMap is the map of validators.
	validatorMap sync.Map
	// validateFunc is the function that validates the value.
	validateFunc func(value *T, args ...any) ci.IValidationResult
}

// vldtFunc is a function that validates the value.
func vldtFunc[T any](v *Validation[T]) func(value *T, args ...any) ci.IValidationResult {
	return func(value *T, args ...any) ci.IValidationResult {
		if v == nil {
			return newValidationResult(false, "validation is nil", nil, fmt.Errorf("validation is nil"))
		}
		if !v.IsValid() {
			return newValidationResult(false, "validation is invalid", nil, fmt.Errorf("validation is invalid"))
		}

		for _, arg := range args {
			if validator, ok := arg.(*ValidationFunc[T]); ok {
				if validator.Func != nil {
					result := validator.Func(value, args...)
					if result != nil && !result.GetIsValid() {
						return result
					}
				}
			}
		}

		return newValidationResult(true, "validation is valid", nil, nil)
	}
}

func VldtFunc[T any](v ci.IValidation[T]) func(value *T, args ...any) ci.IValidationResult {
	return func(value *T, args ...any) ci.IValidationResult {
		if v == nil {
			return NewValidationResult(false, "validation is nil", nil, fmt.Errorf("validation is nil"))
		}
		if !v.IsValid() {
			return NewValidationResult(false, "validation is invalid", nil, fmt.Errorf("validation is invalid"))
		}

		//v.mu.Lock()
		//defer v.mu.Unlock()

		for _, arg := range args {
			if validator, ok := arg.(ci.IValidationFunc[T]); ok {
				if validator.GetFunction() != nil {
					result := validator.GetFunction()(value, args...)
					if result != nil && !result.GetIsValid() {
						return result
					}
				}
			}
		}

		return NewValidationResult(true, "validation is valid", nil, nil)
	}
}

func newValidation[T any]() *Validation[T] {
	validation := &Validation[T]{
		isValid:      false,
		validatorMap: sync.Map{},
	}
	validation.validateFunc = vldtFunc(validation)
	return validation
}
func NewValidation[T any]() ci.IValidation[T] {
	validation := &Validation[T]{
		isValid:      false,
		validatorMap: sync.Map{},
	}
	validation.validateFunc = vldtFunc(validation)
	return validation
}

func (v *Validation[T]) CheckIfWillValidate() bool {
	if v == nil {
		return false
	}

	v.mu.RLock()
	defer v.mu.RUnlock()

	hasValidator := false
	v.validatorMap.Range(func(key, value any) bool {
		if _, vld := key.(int); vld {
			if _, ok := value.(ValidationFunc[T]); ok {
				hasValidator = true
				return false
			}
		}
		return true
	})
	v.hasValidation = hasValidator
	return hasValidator
}

// Validate is the function that validates the value.
func (v *Validation[T]) Validate(value *T, args ...any) ci.IValidationResult {
	if v == nil {
		return NewValidationResult(false, "validation is nil", nil, fmt.Errorf("validation is nil"))
	}
	if value == nil {
		return NewValidationResult(false, "value is nil", nil, fmt.Errorf("value is nil"))
	}
	if !v.hasValidation {
		return NewValidationResult(false, "validation has no validators", nil, fmt.Errorf("validation has no validators"))
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	results := make([]ci.IValidationResult, 0)
	v.validatorMap.Range(func(key, val any) bool {
		if validator, ok := val.(ValidationFunc[T]); ok {
			result := validator.Func(value, args...)
			results = append(results, result)
			if result != nil && !result.GetIsValid() {
				v.isValid = false
				return false
			}
		}
		return true
	})

	if len(results) > 0 {
		sort.Slice(results, func(i, j int) bool {
			return results[i].GetMessage() < results[j].GetMessage()
		})
	}

	v.isValid = true
	for _, result := range results {
		if result != nil && !result.GetIsValid() {
			v.isValid = false
			break
		}
	}

	return NewValidationResult(v.isValid, "validation is valid", nil, nil)
}

// AddValidator is a function that adds a validator to the map of validators.
func (v *Validation[T]) AddValidator(validator ci.IValidationFunc[T]) error {
	if v == nil {
		return fmt.Errorf("validation is nil")
	}

	// Will update v.hasValidation always, if this method is called.
	v.CheckIfWillValidate()

	if validator.GetFunction() == nil {
		return fmt.Errorf("validator function is nil")
	}
	if validator.GetPriority() < 0 {
		return fmt.Errorf("priority must be greater than or equal to 0")
	}
	if _, ok := v.validatorMap.LoadOrStore(validator.GetPriority(), validator); ok {
		return fmt.Errorf("validator with priority %d already exists", validator.GetPriority())
	}

	// If the validator was added, we need to update v.hasValidation again, just for safety.
	v.CheckIfWillValidate()

	return nil
}

// RemoveValidator is a function that removes a validator from the map of validators.
func (v *Validation[T]) RemoveValidator(priority int) error {
	if v == nil {
		return fmt.Errorf("validation is nil")
	}
	if _, ok := v.validatorMap.LoadAndDelete(priority); !ok {
		return fmt.Errorf("validator with priority %d does not exist", priority)
	}

	// If the validator was removed, we need to update v.hasValidation.
	v.CheckIfWillValidate()

	return nil
}

// GetValidator is a function that gets a validator from the map of validators.
func (v *Validation[T]) GetValidator(priority int) (any, error) {
	if v == nil {
		return nil, fmt.Errorf("validation is nil")
	}
	if !v.hasValidation {
		return nil, fmt.Errorf("validation has no validators")
	}
	if validator, ok := v.validatorMap.Load(priority); ok {
		return validator, nil
	}
	return nil, fmt.Errorf("validator with priority %d does not exist", priority)
}

// GetValidators is a function that gets the map of validators.
func (v *Validation[T]) GetValidators() map[int]ci.IValidationFunc[T] {
	if v == nil {
		return nil
	}
	if !v.hasValidation {
		return nil
	}
	validatorMapSnapshot := make(map[int]ci.IValidationFunc[T])
	v.validatorMap.Range(func(key, value any) bool {
		if validator, ok := value.(ci.IValidationFunc[T]); ok {
			validatorMapSnapshot[validator.GetPriority()] = validator
		}
		return true
	})
	return validatorMapSnapshot
}

// GetResults is a function that gets the map of errors.
func (v *Validation[T]) GetResults() map[int]ci.IValidationResult {
	if v == nil {
		return nil
	}
	if !v.hasValidation {
		return nil
	}
	results := make(map[int]ci.IValidationResult)
	v.validatorMap.Range(func(key, value any) bool {
		if validator, ok := value.(ci.IValidationFunc[T]); ok {
			results[validator.GetPriority()] = validator.GetResult()
		}
		return true
	})
	return results
}

// ClearResults is a function that clears the map of errors.
func (v *Validation[T]) ClearResults() {
	if v == nil {
		return
	}
	if !v.hasValidation {
		return
	}
	v.validatorMap.Range(func(key, value any) bool {
		if validator, ok := value.(ValidationFunc[T]); ok {
			validator.Result = nil
			v.validatorMap.Store(key, validator)
		}
		return true
	})
}

// IsValid is a function that gets the boolean that indicates if the value is valid.
func (v *Validation[T]) IsValid() bool {
	if v == nil {
		// If the validation is nil, we need to return false.
		// But we will Log that the validation is nil.
		return false
	}
	if !v.hasValidation {
		// If the validation has no validators, we need to return false.
		// But we will Log that the validation has no validators.
		return false
	}
	return v.isValid
}
