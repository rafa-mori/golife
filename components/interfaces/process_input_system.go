package interfaces

import "os/exec"

type IProcessInputSystemBase[T any] interface {
	GetCommand() string
	GetArgs() []string
	GetPath() string
	GetCmd() *exec.Cmd
	BuildCmd() *exec.Cmd
	Cmd() *exec.Cmd
	GetProperties() map[string]interface{}
}
