package interfaces

type IValidationFunc[T any] interface {
	GetPriority() int
	SetPriority(priority int)
	GetFunction() func(value *T, args ...any) IValidationResult
	SetFunction(function func(value *T, args ...any) IValidationResult)
	GetResult() IValidationResult
	SetResult(result IValidationResult)
}
