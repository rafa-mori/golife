package internal

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type ManagedProcess struct {
	Args       []string
	Command    string
	Cmd        *exec.Cmd
	Name       string
	WaitFor    bool
	ProcPid    int
	ProcHandle uintptr
	mu         sync.Mutex
}

func (p *ManagedProcess) Start() error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.IsRunning() {
		return fmt.Errorf("processo %s já está rodando", p.Name)
	}

	p.Cmd = exec.Command(p.Command, p.Args...)
	if p.WaitFor {
		return p.Cmd.Run()
	} else {

		if procErr := p.Cmd.Start(); procErr != nil {
			return procErr
		} else {
			switch p.Cmd.ProcessState {
			case nil:
				p.ProcPid = p.Cmd.Process.Pid
				p.ProcHandle = uintptr(p.Cmd.Process.Pid)
				if releaseErr := p.Cmd.Process.Release(); releaseErr != nil {
					return releaseErr
				}
				return nil
			default:
				return fmt.Errorf("processo %s não está rodando", p.Name)
			}
		}

		//p.ProcPid, p.ProcHandle, procErr = syscall.StartProcess(p.Command, p.Args, &syscall.ProcAttr{
		//	Dir:   p.Cmd.Dir,
		//	Env:   p.Cmd.Env,
		//	Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		//})
		//if procErr != nil {
		//	return procErr
		//}
		//releaseProcErr := p.Cmd.Process.Release()
		//if releaseProcErr != nil {
		//	return releaseProcErr
		//}
	}

	if err := p.Cmd.Start(); err != nil {
		return err
	}
	return nil
}

func (p *ManagedProcess) Stop() error {
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

func (p *ManagedProcess) Restart() error {
	if err := p.Stop(); err != nil {
		return err
	}
	return p.Start()
}

func (p *ManagedProcess) IsRunning() bool {
	if p == nil || p.Cmd == nil || p.Cmd.Process == nil {
		return false
	}
	return p.Cmd.ProcessState == nil
}

func (p *ManagedProcess) Pid() int {
	if p == nil || p.Cmd == nil || p.Cmd.Process == nil {
		return -1
	}
	return p.Cmd.Process.Pid
}

func (p *ManagedProcess) Wait() error {
	if p == nil || p.Cmd == nil {
		return nil
	}
	return p.Cmd.Wait()
}

func NewManagedProcess(name string, command string, args []string, wait bool) *ManagedProcess {
	envs := os.Environ()
	envPath := os.Getenv("PATH")
	envs = append(envs, fmt.Sprintf("PATH=%s", envPath))
	realCmd, realCmdErr := exec.LookPath(command)
	if realCmdErr != nil {
		return nil
	}
	return &ManagedProcess{
		Args:    []string{},
		Cmd:     exec.Command(realCmd, args...),
		Command: command,
		Name:    name,
		WaitFor: wait,
		mu:      sync.Mutex{},
	}
}
