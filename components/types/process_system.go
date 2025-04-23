package types

import (
	ci "github.com/faelmori/golife/components/interfaces"
	l "github.com/faelmori/logz"

	"io"
	"os"
	"os/exec"
	"reflect"
)

// ProcessInputSystemBase is a struct that implements the IProcessInputSystemBase interface.
type ProcessInputSystemBase[T any, P any] struct {
	// Logger is the logger for the process
	Logger l.Logger
	// Mutexes is the mutex for the process
	*Mutexes
	// Cmd is the command to run
	cmd *exec.Cmd

	// Reference is the reference for the process with ID and Name
	*Reference `json:"reference" yaml:"reference" xml:"reference" gorm:"reference"`
	// ProcessInputConfig is the configuration for the process
	*ProcessInputConfig `json:"process_config" yaml:"process_config" xml:"process_config" gorm:"process_config"`
	// Command is the command to run
	Command string `json:"command" yaml:"command" xml:"command" gorm:"command"`
	// Args is the arguments to pass to the command
	Args []string `json:"args" yaml:"args" xml:"args" gorm:"args"`
	// Path is the path to the command
	Path string `json:"path" yaml:"path" xml:"path" gorm:"path"`
	// ProcPid is the process ID
	ProcPid int `json:"pid" yaml:"pid" xml:"pid" gorm:"pid"`

	// ProcPointer is the process handle
	ProcPointer uintptr `json:"proc_pointer,omitempty" yaml:"proc_pointer,omitempty" xml:"proc_pointer,omitempty" gorm:"proc_pointer,omitempty"`
	// PropertiesSystemProc is the properties map for the process
	PropertiesSystemProc map[string]interface{} `json:"properties_system_proc,omitempty" yaml:"properties_system_proc,omitempty" xml:"properties_system_proc,omitempty" gorm:"properties_system_proc,omitempty"`
	// Object is the object to pass to the command
	Object *P `json:"object,omitempty" yaml:"object,omitempty" xml:"object,omitempty" gorm:"object,omitempty"`
	// Function is a custom function to wrap the command
	Function ci.IValidationFunc[T] `json:"function,omitempty" yaml:"function,omitempty" xml:"function,omitempty" gorm:"function,omitempty"`
}

// newProcessInputSystemBase creates a new ProcessInput instance with the provided Logger.
func newProcessInputSystemBase[T any, P ci.IProcessInput[T]](name, command string, args []string, waitFor bool, restart bool, function ci.IValidationFunc[T], logger l.Logger, debug bool) *ProcessInputSystemBase[T, P] {
	cfg := newProcessConfig[P](name, waitFor, restart, "system", nil, logger, debug)
	nps := &ProcessInputSystemBase[T, P]{
		Logger:               logger,
		Mutexes:              NewMutexesType(),
		Command:              command,
		Args:                 args,
		Reference:            NewReference(name).GetReference(),
		ProcessInputConfig:   cfg,
		Path:                 "",
		ProcPid:              -1,
		ProcPointer:          0,
		PropertiesSystemProc: make(map[string]any),
		Object:               nil,
		Function:             function,
	}
	return nps
}

// NewProcessInputSystemBase creates a new ProcessInput instance with the provided Logger.
func NewProcessInputSystemBase[T any, P ci.IProcessInput[T]](name, command string, args []string, waitFor bool, restart bool, function ci.IValidationFunc[T], logger l.Logger, debug bool) ci.IProcessInputRuntimeBase[T, P] {
	return newProcessInputSystemBase[T, P](name, command, args, waitFor, restart, function, logger, debug)
}

// GetCommand returns the command to run.
func (pi *ProcessInputSystemBase[T, P]) GetCommand() string {
	if pi == nil {
		return ""
	}
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.Command
}

// GetArgs returns the arguments to pass to the command.
func (pi *ProcessInputSystemBase[T, P]) GetArgs() []string {
	if pi == nil {
		return nil
	}
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.Args
}

// GetPath returns the path to the command.
func (pi *ProcessInputSystemBase[T, P]) GetPath() string {
	if pi == nil {
		return ""
	}
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.Path
}

// GetCmd returns the command to run.
func (pi *ProcessInputSystemBase[T, P]) GetCmd() *exec.Cmd {
	if pi == nil {
		return nil
	}
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	if pi.cmd == nil {
		return pi.BuildCmd()
	}

	return pi.cmd
}

// BuildCmd builds the command to run.
func (pi *ProcessInputSystemBase[T, P]) BuildCmd() *exec.Cmd {
	if pi == nil {
		return nil
	}
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	cmd := exec.Command(pi.Command, pi.Args...)
	cmd.Dir = pi.GetPath()
	cmd.Env = append(cmd.Env, "PATH="+pi.GetPath())
	pi.cmd = cmd

	if pi.PropertiesSystemProc == nil {
		pi.PropertiesSystemProc = make(map[string]interface{})
	}

	pi.PropertiesSystemProc["ProcessState"] = func() *os.ProcessState {
		if pi.cmd == nil {
			return nil
		}
		return pi.cmd.ProcessState
	}
	pi.PropertiesSystemProc["Path"] = func() string {
		if pi.cmd == nil {
			return ""
		}
		return pi.cmd.Path
	}
	pi.PropertiesSystemProc["ProcPid"] = func() int {
		if pi.cmd == nil {
			return -1
		}
		if pi.cmd.Process == nil {
			return -1
		}
		return pi.cmd.Process.Pid
	}
	pi.PropertiesSystemProc["Args"] = func() []string {
		if pi.cmd == nil {
			return nil
		}
		return pi.cmd.Args
	}
	pi.PropertiesSystemProc["Env"] = func() []string {
		if pi.cmd == nil {
			return nil
		}
		return pi.cmd.Env
	}
	pi.PropertiesSystemProc["Dir"] = func() string {
		if pi.cmd == nil {
			return ""
		}
		return pi.cmd.Dir
	}
	pi.PropertiesSystemProc["Stdin"] = func() io.Reader {
		if pi.cmd == nil {
			return nil
		}
		return pi.cmd.Stdin
	}
	pi.PropertiesSystemProc["Stdout"] = func() io.Writer {
		if pi.cmd == nil {
			return nil
		}
		return pi.cmd.Stdout
	}
	pi.PropertiesSystemProc["Stderr"] = func() io.Writer {
		if pi.cmd == nil {
			return nil
		}
		return pi.cmd.Stderr
	}
	pi.PropertiesSystemProc["ExtraFiles"] = func() []*os.File {
		if pi.cmd == nil {
			return nil
		}
		return pi.cmd.ExtraFiles
	}
	pi.PropertiesSystemProc["Cancel"] = func() func() error {
		if pi.cmd == nil {
			return nil
		}
		return pi.cmd.Cancel
	}

	return pi.cmd
}

// GetProperties returns the properties map for the process
func (pi *ProcessInputSystemBase[T, P]) GetProperties() map[string]interface{} {
	if pi.PropertiesSystemProc == nil {
		pi.PropertiesSystemProc = make(map[string]interface{})
	}
	return pi.PropertiesSystemProc
}

// Cmd returns the command to run.
func (pi *ProcessInputSystemBase[T, P]) Cmd() *exec.Cmd {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	if pi.cmd == nil {
		return pi.BuildCmd()
	}

	return pi.cmd
}

// GetObjectType returns the process type.
func (pi *ProcessInputSystemBase[T, P]) GetObjectType() reflect.Type {
	if pi == nil {
		return nil
	}
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	if pi.Object == nil {
		return nil
	}
	return reflect.TypeFor[T]()
}

// GetObject returns the object to pass to the command.
func (pi *ProcessInputSystemBase[T, P]) GetObject() *P {
	if pi == nil {
		return nil
	}
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	if pi.Object == nil {
		return nil
	}
	return pi.Object
}

// GetFunction returns the custom function to wrap the command.
func (pi *ProcessInputSystemBase[T, P]) GetFunction() ci.IValidationFunc[T] {
	if pi == nil {
		return nil
	}
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	if pi.Function == nil {
		return nil
	}
	return pi.Function
}
