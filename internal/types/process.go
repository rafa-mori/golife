package types

import (
	c "github.com/faelmori/golife/internal/routines/chan"
	"github.com/google/uuid"
	"os"
	"syscall"
)

type ProcessParameters struct {
	Name string
	Args []string
	Env  []string
	CWD  string
	Type string
	Host string
	Port int
	User string
}

type ProcessConfig struct {
	// Telemetry configuration
	Telemetry

	// Threading configuration
	ThreadingConfig

	// ID and Reference
	ID uuid.UUID
	// Process ID
	Pid int
	// Type of process
	Type string
	// Name of the process
	Name string
	// Basic process properties
	ProcessProperties map[string]Property[any]
	// Process Agents
	ProcessAgents map[string]c.IChannel[any, int]
	// Process Stages
	ProcessStagesMap map[string]StageConfig
	// Process Events
	ProcessEventsMap map[string]ManagedEventsConfig
}

func (pc *ProcessConfig) InitDefaults(args *ProcessParameters) {
	if args == nil {
		args = &ProcessParameters{}
	}

	pc.ProcessProperties["name"] = NewProperty[string]("name", nil)
	_ = pc.ProcessProperties["name"].SetValue(args.Name, nil)

	pc.ProcessProperties["cwd"] = NewProperty[any]("cwd", nil)
	_ = pc.ProcessProperties["cwd"].SetValue(args.CWD, nil)

	pc.ProcessProperties["args"] = NewProperty[any]("args", nil)
	_ = pc.ProcessProperties["args"].SetValue(args.Args, nil)

	pc.ProcessProperties["env"] = NewProperty[any]("env", nil)
	_ = pc.ProcessProperties["env"].SetValue(args.Env, nil)

	pc.ProcessProperties["host"] = NewProperty[any]("host", nil)
	_ = pc.ProcessProperties["host"].SetValue(args.Host, nil)

	pc.ProcessProperties["port"] = NewProperty[any]("port", nil)
	_ = pc.ProcessProperties["port"].SetValue(args.Port, nil)

	pc.ProcessProperties["user"] = NewProperty[any]("user", nil)
	_ = pc.ProcessProperties["user"].SetValue(args.User, nil)

	pc.ProcessProperties["pid"] = NewProperty[any]("pid", nil)
	_ = pc.ProcessProperties["pid"].SetValue(pc.Pid, nil)

}

func (pc *ProcessConfig) GetProcessSysPid() int {
	var pidA, pidB, pidC int
	var err error
	pidA = os.Getpid()
	pidB = os.Getppid()
	if pidC, err = syscall.Getpgid(pidA); err != nil {
		return -1
	} else {
		if pidC == 0 {
			pidC = pidA
		}
		if pidC == pidB {
			pidC = pidA
		}
	}
	return pidC
}

func NewProcessConfig(args ProcessParameters) *ProcessConfig {
	mc := &ProcessConfig{
		Telemetry:         *NewTelemetry(),
		ThreadingConfig:   *NewThreadingConfig(),
		ID:                uuid.New(),
		ProcessProperties: make(map[string]Property[any]),
	}

	mc.Type = args.Type
	mc.Pid = mc.GetProcessSysPid()

	mc.InitDefaults(&args)

	return mc
}
