package process_input

import (
	t "github.com/faelmori/golife/components/types"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"io"
	"os"
	"os/exec"
)

// ProcessSystemBase is a struct that holds the system process information.
type ProcessSystemBase[T any] struct {
	// Logger is the logger for the process
	Logger l.Logger
	// Reference is the reference for the process with ID and Name
	*t.Reference
	// Mutexes is the mutex for the process
	*t.Mutexes
	// Command is the command to run
	Command string `json:"command" yaml:"command" xml:"command" gorm:"command"`
	// Args is the arguments to pass to the command
	Args []string `json:"args" yaml:"args" xml:"args" gorm:"args"`
	// Path is the env path to pass to the command
	Path string
	// Cmd is the command to run
	cmd *exec.Cmd

	PropertiesSystemProc map[string]interface{}

	//// ProcPid is the process ID
	//ProcPid int `json:"pid" yaml:"pid" xml:"pid" gorm:"pid"`
	//// ProcPidFile is the process ID file
	//ProcPidFile string `json:"pid_file" yaml:"pid_file" xml:"pid_file" gorm:"pid_file"`
	//// ProcHandle is the process handle
	//ProcPointer uintptr `json:"proc_pointer" yaml:"proc_pointer" xml:"proc_pointer" gorm:"proc_pointer"`
	//// Path is the path to the command
	//Path string `json:"path" yaml:"path" xml:"path" gorm:"path"`

}

// NewSystemProcessInput creates a new ProcessInput instance with the provided Logger.
func NewSystemProcessInput[T any](name, command string, args []string, waitFor bool, restart bool, customFn func(obj T, args ...any) *t.ValidationResult, logger l.Logger) *ProcessSystemBase[T] {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	ref := t.NewReference(name)
	mu := t.NewMutexes()
	npi := &ProcessSystemBase[T]{
		Reference:            ref,
		Logger:               logger,
		Mutexes:              mu,
		Command:              command,
		Args:                 args,
		PropertiesSystemProc: make(map[string]interface{}),
	}
	// This method is used to initialize the cmd field with right and more robust way
	if npi.BuildCmd() == nil {
		gl.LogObjLogger[ProcessSystemBase[T]](npi, "error", "Command is nil")
		return nil
	}

	return npi
}

// GetCommand returns the command to run.
func (pi *ProcessSystemBase[T]) GetCommand() string {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	return pi.Command
}

// GetArgs returns the arguments to pass to the command.
func (pi *ProcessSystemBase[T]) GetArgs() []string {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	return pi.Args
}

// GetPath returns the path to the command.
func (pi *ProcessSystemBase[T]) GetPath() string {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	return pi.Path
}

// GetCmd returns the command to run.
func (pi *ProcessSystemBase[T]) GetCmd() *exec.Cmd {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	if pi.cmd == nil {
		return pi.BuildCmd()
	}

	return pi.cmd
}

// BuildCmd builds the command to run.
func (pi *ProcessSystemBase[T]) BuildCmd() *exec.Cmd {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

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
func (pi *ProcessSystemBase[T]) GetProperties() map[string]interface{} {
	if pi.PropertiesSystemProc == nil {
		pi.PropertiesSystemProc = make(map[string]interface{})
	}
	return pi.PropertiesSystemProc
}

// Cmd returns the command to run.
func (pi *ProcessSystemBase[T]) Cmd() *exec.Cmd {
	pi.Mutexes.RLock()
	defer pi.Mutexes.RUnlock()

	if pi.cmd == nil {
		return pi.BuildCmd()
	}

	return pi.cmd
}
