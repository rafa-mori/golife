package process

import (
	"fmt"
	ci "github.com/faelmori/golife/internal/components/interfaces"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type ManagedProcess[T any] struct {
	Logger     l.Logger
	Args       []string
	Command    string
	Function   func(T, ...any) (bool, error)
	Cmd        *exec.Cmd
	Name       string
	WaitFor    bool
	ProcPid    int
	ProcHandle uintptr
	mu         sync.Mutex
}

func (p *ManagedProcess[T]) GetArgs() []string                          { return p.Args }
func (p *ManagedProcess[T]) GetCommand() string                         { return p.Command }
func (p *ManagedProcess[T]) GetFunction() func(T, ...any) (bool, error) { return p.Function }
func (p *ManagedProcess[T]) GetName() string                            { return p.Name }
func (p *ManagedProcess[T]) GetWaitFor() bool                           { return p.WaitFor }
func (p *ManagedProcess[T]) GetProcPid() int                            { return p.ProcPid }
func (p *ManagedProcess[T]) GetProcHandle() uintptr                     { return p.ProcHandle }
func (p *ManagedProcess[T]) GetCmd() *exec.Cmd                          { return p.Cmd }
func (p *ManagedProcess[T]) WillRestart() bool                          { return p.Cmd != nil }
func (p *ManagedProcess[T]) Release() error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Cmd != nil && p.Cmd.Process != nil {
		if err := p.Cmd.Process.Release(); err != nil {
			return err
		}
	} else {
		gl.LogObjLogger(p, "warn", "Process is nil or not running")
		return fmt.Errorf("process is nil or not running")
	}
	return nil
}
func (p *ManagedProcess[T]) Status() string {
	if p == nil {
		return "nil"
	}
	if p.Cmd == nil {
		return "nil"
	}
	if p.Cmd.Process == nil {
		return "nil"
	}
	if p.Cmd.ProcessState == nil {
		return "nil"
	}
	return fmt.Sprintf("Process %s (PID %d) is running: %t", p.Name, p.Pid(), p.IsRunning())
}
func (p *ManagedProcess[T]) ExecCmd() error {
	if p == nil {
		gl.LogObjLogger(p, "fatal", "Process is nil")
		return fmt.Errorf("process is nil")
	}
	if p.Command != "" {
		gl.LogObjLogger(p, "info", fmt.Sprintf("Starting process %s with command %s", p.Name, p.Command))
		if p.WaitFor {
			stdout, err := p.Cmd.StdoutPipe()
			if err != nil {
				gl.LogObjLogger(p, "error", fmt.Sprintf("Error creating stdout pipe: %v", err))
				stdout = os.Stdout
			}
			stderr, err := p.Cmd.StderrPipe()
			if err != nil {
				gl.LogObjLogger(p, "error", fmt.Sprintf("Error creating stderr pipe: %v", err))
				stderr = os.Stderr
			}

			defer func() {
				_ = stdout.Close()
				_ = stderr.Close()
			}()
		} else {
			//p.Cmd.SysProcAttr = &syscall.SysProcAttr{
			//	Setsid: true,
			//	//Setpgid: true,
			//	//Pgid: 0,
			//	//Foreground: false,
			//}
		}

		if err := p.Cmd.Start(); err != nil {
			gl.LogObjLogger(p, "error", fmt.Sprintf("Error starting command: %v", err))
			return err
		}

		p.ProcPid = p.Cmd.Process.Pid
		p.ProcHandle = uintptr(p.Cmd.Process.Pid)

		if waitErr := p.Wait(); waitErr != nil {
			gl.LogObjLogger(p, "error", fmt.Sprintf("Error waiting for command to finish: %v", waitErr))
			return waitErr
		}

		if p.Cmd.ProcessState != nil && p.Cmd.ProcessState.Exited() {
			gl.LogObjLogger(p, "info", fmt.Sprintf("Process %s exited with code %d", p.Name, p.Cmd.ProcessState.ExitCode()))
		} else {
			gl.LogObjLogger(p, "warn", fmt.Sprintf("Process %s running with PID %d", p.Name, p.Pid()))
		}

		return nil
	} else {
		gl.LogObjLogger(p, "error", "Command is empty")
		return fmt.Errorf("command is empty")
	}
}
func (p *ManagedProcess[T]) Start() error {
	if p == nil {
		gl.LogObjLogger(p, "fatal", "Process is nil")
		return fmt.Errorf("process is nil")
	}
	if p.IsRunning() {
		gl.LogObjLogger(p, "warn", fmt.Sprintf("Process %s is already running", p.Name))
		return fmt.Errorf("process %s is already running", p.Name)
	}

	if p.Function == nil {
		gl.LogObjLogger(p, "warn", "Custom function is nil, using default execution")
	} else {
		gl.LogObjLogger(p, "info", fmt.Sprintf("Starting process %s with custom function", p.Name))
		go func() {
			if _, err := p.Function(nil, nil); err != nil {
				gl.LogObjLogger(p, "error", fmt.Sprintf("Error executing custom function: %v", err))
			}
		}()
	}
	if err := p.ExecCmd(); err != nil {
		gl.LogObjLogger(p, "error", fmt.Sprintf("Error executing command: %v", err))
		return err
	}
	return nil

}
func (p *ManagedProcess[T]) Stop() error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.IsRunning() {
		return nil
	}

	if err := p.Cmd.Process.Kill(); err != nil {
		return err
	}
	return nil
}
func (p *ManagedProcess[T]) Restart() error {
	if err := p.Stop(); err != nil {
		return err
	}
	return p.Start()
}
func (p *ManagedProcess[T]) IsRunning() bool {
	if p == nil || p.Cmd == nil || p.Cmd.Process == nil {
		return false
	}
	return p.Cmd.ProcessState == nil
}
func (p *ManagedProcess[T]) Pid() int {
	if p == nil || p.ProcPid == 0 {
		return -1
	}
	return p.ProcPid
}
func (p *ManagedProcess[T]) Wait() error {
	if p == nil || p.Cmd == nil || p.Command == "" {
		gl.LogObjLogger(p, "error", "Process is nil or command is empty")
		return nil
	}
	if strings.Contains(p.Command, "_supervisor") {
		defer func(cmd *exec.Cmd) {
			if releaseErr := cmd.Process.Release(); releaseErr != nil {
				gl.LogObjLogger(p, "error", fmt.Sprintf("Error releasing command: %v", releaseErr))
				return
			}
			gl.LogObjLogger(p, "success", fmt.Sprintf("Process %s released with PID %d", p.Name, p.Pid()))
		}(p.Cmd)
		return nil
	}
	if p.WaitFor {
		return p.Cmd.Wait()
	} else {
		return nil
	}
}
func (p *ManagedProcess[T]) String() string {
	return fmt.Sprintf("Process %s (PID %d) is running: %t", p.Name, p.Pid(), p.IsRunning())
}
func (p *ManagedProcess[T]) SetArgs(args []string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Args = args
}
func (p *ManagedProcess[T]) SetCommand(command string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Command = command
}
func (p *ManagedProcess[T]) SetName(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Name = name
}
func (p *ManagedProcess[T]) SetWaitFor(wait bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.WaitFor = wait
}
func (p *ManagedProcess[T]) SetProcPid(pid int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.ProcPid = pid
}
func (p *ManagedProcess[T]) SetProcHandle(handle uintptr) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.ProcHandle = handle
}
func (p *ManagedProcess[T]) SetCmd(cmd *exec.Cmd) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Cmd = cmd
}
func (p *ManagedProcess[T]) SetFunction(customFunc func(T, ...any) (bool, error)) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Function = customFunc
}

func NewManagedProcess[T any](name string, command string, args []string, wait bool, function func(T, ...any) (bool, error)) ci.IManagedProcess[T] {
	envs := os.Environ()
	envPath := os.Getenv("PATH")
	envs = append(envs, fmt.Sprintf("PATH=%s", envPath))

	if args == nil {
		args = make([]string, 0)
	}

	var cmd *exec.Cmd
	if command != "" {
		realCmd, realCmdErr := exec.LookPath(command)
		if realCmdErr == nil {
			cmd = exec.Command(realCmd, args...)
		}
	}

	mgrProc := ManagedProcess[T]{
		Logger:   l.GetLogger("GoLife"),
		Args:     args,
		Cmd:      cmd,
		Command:  command,
		Name:     name,
		WaitFor:  wait,
		Function: function,
		mu:       sync.Mutex{},
	}

	gl.LogObjLogger[ManagedProcess[T]](&mgrProc, "success", fmt.Sprintf("Managed process %s created with command %s", name, command))

	return &mgrProc
}
