package internal

import (
	"os/exec"
)

type ManagedProcess struct {
	Cmd  *exec.Cmd
	Name string
}

func (p *ManagedProcess) Start() error {
	return p.Cmd.Start()
}

func (p *ManagedProcess) Stop() error {
	if p.Cmd != nil {
		return p.Cmd.Process.Kill()
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
	if p.Cmd != nil {
		return p.Cmd.ProcessState == nil
	}
	return false
}

func (p *ManagedProcess) Pid() int {
	if p.Cmd != nil {
		return p.Cmd.Process.Pid
	}
	return 0
}
