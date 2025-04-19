package types

import (
	f "github.com/faelmori/golife/internal/property"
	u "github.com/faelmori/golife/internal/utils"
	c "github.com/faelmori/golife/services"
	"github.com/google/uuid"
	"os"
	"syscall"
)

type ProcessParameters[T f.KubexProperty[any]] struct {
	Name string
	Args []string
	Env  []string
	CWD  string
	Type string
	Host string
	Port int
	User string
}

type ProcessConfig[T f.KubexProperty[any]] struct {
	// Telemetry configuration
	Telemetry

	// Threading configuration
	u.ThreadingConfig

	// ID and Reference
	ID uuid.UUID
	// Process ID
	Pid int
	// Type of process
	Type string
	// Name of the process
	Name string
	// Basic process properties
	ProcessProperties map[string]any
	// Process Agents
	ProcessAgents map[string]c.IChannel[any, int]
	// Process Stages
	ProcessStagesMap map[string]StageConfig
	// Process Events
	ProcessEventsMap map[string]EventsConfig
}

func (pc *ProcessConfig[T]) InitDefaults(args *ProcessParameters[f.KubexProperty[any]]) {
	if args == nil {
		args = &ProcessParameters[f.KubexProperty[any]]{}
	}

	pc.ProcessProperties = make(map[string]any)

	name := f.NewProperty[string]("name", nil)
	pid := f.NewProperty[int]("pid", nil)

	pc.ProcessProperties["name"] = name
	pc.ProcessProperties["pid"] = pid
	pc.ProcessProperties["cwd"] = f.NewProperty[any]("cwd", nil)
	pc.ProcessProperties["args"] = f.NewProperty[any]("args", nil)
	pc.ProcessProperties["env"] = f.NewProperty[any]("env", nil)
	//_ = pc.ProcessProperties["env"].SetValue(args.Env, nil)

	pc.ProcessProperties["host"] = f.NewProperty[any]("host", nil)
	//_ = pc.ProcessProperties["host"].SetValue(args.Host, nil)

	pc.ProcessProperties["port"] = f.NewProperty[any]("port", nil)
	//_ = pc.ProcessProperties["port"].SetValue(args.Port, nil)

	pc.ProcessProperties["user"] = f.NewProperty[any]("user", nil)
	//_ = pc.ProcessProperties["user"].SetValue(args.User, nil)

	pc.ProcessProperties["pid"] = f.NewProperty[any]("pid", nil)
	//_ = pc.ProcessProperties["pid"].SetValue(pc.Pid, nil)

}

func (pc *ProcessConfig[T]) GetProcessSysPid() int {
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

func NewProcessConfig[T any](args ProcessParameters[f.KubexProperty[any]]) *ProcessConfig[f.KubexProperty[any]] {
	mc := ProcessConfig[f.KubexProperty[any]]{
		Telemetry:         *NewTelemetry(),
		ThreadingConfig:   *u.NewThreadingConfig(),
		ID:                uuid.New(),
		ProcessProperties: make(map[string]any),
	}

	mc.Type = args.Type
	mc.Pid = mc.GetProcessSysPid()

	mc.InitDefaults(nil)

	return &mc
}

func (pc *ProcessConfig[T]) GetID() uuid.UUID {
	pc.Telemetry.logger.Info("GetID", nil)
	return pc.ID
}
