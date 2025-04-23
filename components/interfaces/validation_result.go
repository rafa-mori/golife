package interfaces

type IValidationResult interface {
	GetIsValid() bool
	GetMessage() string
	GetError() error
}
