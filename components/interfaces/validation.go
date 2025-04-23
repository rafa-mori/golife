package interfaces

type IValidation[T any] interface {
	CheckIfWillValidate() bool
	Validate(value *T, args ...any) IValidationResult
	AddValidator(validator IValidationFunc[T]) error
	RemoveValidator(priority int) error
	GetValidator(priority int) (any, error)
	GetValidators() map[int]IValidationFunc[T]
	GetResults() map[int]IValidationResult
	ClearResults()
	IsValid() bool
}
