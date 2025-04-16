package process

import (
	t "github.com/faelmori/golife/internal/types"
	l "github.com/faelmori/logz"

	"fmt"
	"os"
	"os/exec"
	"sync"
)

type IManagedProcess[T t.ProcessConfig] interface {
	GetArgs() []string
	GetCommand() string
	GetCustomFunc() func() error
	GetName() string
	GetWaitFor() bool
	GetProcPid() int
	GetProcHandle() uintptr
	GetCmd() *exec.Cmd
	WillRestart() bool

	Start() error
	Stop() error
	Restart() error
	IsRunning() bool

	Pid() int
	Wait() error
	String() string

	SetArgs(args []string)
	SetCommand(command string)
	SetCustomFunc(func() error)
	SetName(name string)
	SetWaitFor(wait bool)
	SetProcPid(pid int)
	SetProcHandle(handle uintptr)
	SetCmd(cmd *exec.Cmd)
}

type ManagedProcess[T t.ProcessConfig] struct {
	Args       []string
	Command    string
	CustomFunc func() error
	Cmd        *exec.Cmd
	Name       string
	WaitFor    bool
	ProcPid    int
	ProcHandle uintptr
	mu         sync.Mutex
}

func (p *ManagedProcess[T]) GetArgs() []string           { return p.Args }
func (p *ManagedProcess[T]) GetCommand() string          { return p.Command }
func (p *ManagedProcess[T]) GetCustomFunc() func() error { return p.CustomFunc }
func (p *ManagedProcess[T]) GetName() string             { return p.Name }
func (p *ManagedProcess[T]) GetWaitFor() bool            { return p.WaitFor }
func (p *ManagedProcess[T]) GetProcPid() int             { return p.ProcPid }
func (p *ManagedProcess[T]) GetProcHandle() uintptr      { return p.ProcHandle }
func (p *ManagedProcess[T]) GetCmd() *exec.Cmd           { return p.Cmd }
func (p *ManagedProcess[T]) WillRestart() bool           { return p.Cmd != nil }
func (p *ManagedProcess[T]) Start() error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.IsRunning() {
		return fmt.Errorf("process %s is already running", p.Name)
	}

	if p.CustomFunc != nil {
		l.InfoCtx(fmt.Sprintf("Executing custom function for process %s", p.Name), nil)
		go func() {
			if err := p.CustomFunc(); err != nil {
				l.ErrorCtx(fmt.Sprintf("Error in custom execution of process %s: %v", p.Name, err), nil)
			}
		}()
		return nil
	}
	if p.Command != "" {
		l.InfoCtx(fmt.Sprintf("Starting process %s with command %s", p.Name, p.Command), nil)
		p.Cmd = exec.Command(p.Command, p.Args...)
		if p.WaitFor {
			return p.Cmd.Run()
		} else {
			if err := p.Cmd.Start(); err != nil {
				return err
			}
			p.ProcPid = p.Cmd.Process.Pid
			p.ProcHandle = uintptr(p.Cmd.Process.Pid)
			return p.Cmd.Process.Release()
		}
	} else {
		l.WarnCtx(fmt.Sprintf("No command defined for process %s", p.Name), nil)
		return nil
	}
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
	if p == nil || p.Cmd == nil || p.Cmd.Process == nil {
		return -1
	}
	return p.Cmd.Process.Pid
}
func (p *ManagedProcess[T]) Wait() error {
	if p == nil || p.Cmd == nil {
		return nil
	}
	return p.Cmd.Wait()
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
func (p *ManagedProcess[T]) SetCustomFunc(customFunc func() error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.CustomFunc = customFunc
}

func NewManagedProcess[T t.ProcessConfig](name string, command string, args []string, wait bool, customFunc func() error) IManagedProcess[T] {
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
		Args:       args,
		Cmd:        cmd,
		Command:    command,
		Name:       name,
		WaitFor:    wait,
		CustomFunc: customFunc,
		mu:         sync.Mutex{},
	}
	return &mgrProc
}
