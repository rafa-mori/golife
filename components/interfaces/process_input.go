package interfaces

import (
	l "github.com/faelmori/logz"
	"os/exec"
	"reflect"
)

type IProcessInputConfig interface {
	GetWaitFor() bool
	GetRestart() bool
	GetProcessType() string
	GetMetadata(key string) (any, bool)
	SetMetadata(key string, value any)
}

type IProcessInputSystemBase[T any] interface {
	GetCommand() string
	GetArgs() []string
	GetPath() string
	GetCmd() *exec.Cmd
	BuildCmd() *exec.Cmd
	Cmd() *exec.Cmd
	GetProperties() map[string]interface{}
}

type IProcessInputRuntimeBase[T any, P IProcessInput[T]] interface {
	GetObjectType() reflect.Type
	GetObject() *P
	GetFunction() IValidationFunc[T]
}

type IProcessInput[T any] interface {
	Serialize(format string) ([]byte, error)
	Deserialize(data []byte, format string) error
	GetLogger() l.Logger
	GetReference() IReference
	GetName() string
	Validate() IValidationResult
}
