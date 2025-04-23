package interfaces

import "reflect"

type IProcessInputRuntimeBase[T any, P IProcessInput[T]] interface {
	GetObjectType() reflect.Type
	GetObject() *P
	GetFunction() IValidationFunc[T]
}
