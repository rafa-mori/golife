package interfaces

type IValidationResult interface {
	GetIsValid() bool
	GetMessage() string
	GetError() error
}
type IValidationFunc[T any] interface {
	GetPriority() int
	SetPriority(priority int)
	GetFunction() func(value *T, args ...any) IValidationResult
	SetFunction(function func(value *T, args ...any) IValidationResult)
	GetResult() IValidationResult
	SetResult(result IValidationResult)
}
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
