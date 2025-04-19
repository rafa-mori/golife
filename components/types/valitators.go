package types

import (
	"fmt"
	"sort"
	"sync"
)

type ValidationResult struct {
	IsValid bool
	Message string
	Error   error
}

type ValidationFunc[T any] struct {
	Priority int
	Func     func(value T, args ...any) *ValidationResult
	Result   *ValidationResult
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
	validateFunc func(args ...any) bool
}

// vldtFunc is a function that validates the value.
func vldtFunc[T any](v *Validation[T]) func(args ...any) bool {
	return func(args ...any) bool {
		var valueToValidate T
		if len(args) == 0 {
			return false
		} else {
			if vl, ok := args[0].(T); ok {
				valueToValidate = vl
			} else {
				return false
			}
		}
		if v == nil {
			return false
		}

		v.isValid = true
		v.validatorMap.Range(func(key, value any) bool {
			if validator, ok := value.(ValidationFunc[T]); ok {
				result := validator.Func(valueToValidate, args...)
				if !result.IsValid {
					v.isValid = false
				}
				validator.Result = result
				v.validatorMap.CompareAndSwap(key, value, validator)
			}
			return true
		})
		return v.isValid
	}
}

// NewValidation is a function that creates a new Validation instance.
func NewValidation[T any]() *Validation[T] {
	validation := &Validation[T]{
		isValid:      false,
		validatorMap: sync.Map{},
	}
	validation.validateFunc = vldtFunc(validation)
	return validation
}

func (v *Validation[T]) checkIfWillValidate() bool {
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
func (v *Validation[T]) Validate(args ...any) bool {
	if v == nil || len(args) == 0 {
		return false
	}
	if !v.hasValidation {
		return false
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	validatorMapSnapshot := make(map[int]ValidationFunc[T])
	sortedValByPriority := make([]ValidationFunc[T], 0)

	v.validatorMap.Range(func(key, value any) bool {
		if validator, ok := value.(ValidationFunc[T]); ok {
			if validator.Func != nil {
				// Remove validators without a function and fill the array with the rest.
				sortedValByPriority = append(sortedValByPriority, validator)
			}
		}
		return false
	})

	if len(sortedValByPriority) == 0 {
		return false
	}

	// Sort the validators by priority.
	sort.Slice(sortedValByPriority, func(i, j int) bool {
		return sortedValByPriority[i].Priority < sortedValByPriority[j].Priority
	})

	// Fill the map with the sorted validators, just validators with a function.
	for newPriority, validator := range sortedValByPriority {
		validator.Priority = newPriority
		validatorMapSnapshot[validator.Priority] = validator
	}

	if len(validatorMapSnapshot) == 0 {
		return false
	}

	return v.validateFunc(args...)
}

// AddValidator is a function that adds a validator to the map of validators.
func (v *Validation[T]) AddValidator(validator ValidationFunc[T]) error {
	if v == nil {
		return fmt.Errorf("validation is nil")
	}

	// Will update v.hasValidation always, if this method is called.
	v.checkIfWillValidate()

	if validator.Func == nil {
		return fmt.Errorf("validator function is nil")
	}
	if validator.Priority < 0 {
		return fmt.Errorf("priority must be greater than or equal to 0")
	}
	if _, ok := v.validatorMap.LoadOrStore(validator.Priority, validator); ok {
		return fmt.Errorf("validator with priority %d already exists", validator.Priority)
	}

	// If the validator was added, we need to update v.hasValidation again, just for safety.
	v.checkIfWillValidate()

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
	v.checkIfWillValidate()

	return nil
}

// GetValidator is a function that gets a validator from the map of validators.
func (v *Validation[T]) GetValidator(priority int) (*ValidationFunc[T], error) {
	if v == nil {
		return nil, fmt.Errorf("validation is nil")
	}
	if !v.hasValidation {
		return nil, fmt.Errorf("validation has no validators")
	}
	if validator, ok := v.validatorMap.Load(priority); ok {
		return validator.(*ValidationFunc[T]), nil
	}
	return nil, fmt.Errorf("validator with priority %d does not exist", priority)
}

// GetValidators is a function that gets the map of validators.
func (v *Validation[T]) GetValidators() map[int]ValidationFunc[T] {
	if v == nil {
		return nil
	}
	if !v.hasValidation {
		return nil
	}
	validatorMapSnapshot := make(map[int]ValidationFunc[T])
	v.validatorMap.Range(func(key, value any) bool {
		if validator, ok := value.(ValidationFunc[T]); ok {
			validatorMapSnapshot[validator.Priority] = validator
		}
		return true
	})
	return validatorMapSnapshot
}

// GetResults is a function that gets the map of errors.
func (v *Validation[T]) GetResults() map[int]*ValidationResult {
	if v == nil {
		return nil
	}
	if !v.hasValidation {
		return nil
	}
	results := make(map[int]*ValidationResult)
	v.validatorMap.Range(func(key, value any) bool {
		if validator, ok := value.(*ValidationFunc[T]); ok {
			results[validator.Priority] = validator.Result
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
