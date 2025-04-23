package interfaces

import (
	"os/exec"
)

type IManagedProcess[T any] interface {
	GetArgs() []string
	GetCommand() string
	GetFunction() func(T, ...any) (bool, error)
	GetName() string
	GetWaitFor() bool
	GetProcPid() int
	GetProcHandle() uintptr
	GetCmd() *exec.Cmd
	WillRestart() bool

	ExecCmd() error
	Release() error
	Status() string
	Start() error
	Stop() error
	Restart() error
	IsRunning() bool

	Pid() int
	Wait() error
	String() string

	SetArgs(args []string)
	SetCommand(command string)
	SetFunction(func(T, ...any) (bool, error))
	SetName(name string)
	SetWaitFor(wait bool)
	SetProcPid(pid int)
	SetProcHandle(handle uintptr)
	SetCmd(cmd *exec.Cmd)
}
