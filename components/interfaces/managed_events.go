package interfaces

type IManagedProcessEvents[T any] interface {
	Event() string
	RegisterEvent(event string, fn func(interface{}))
	Trigger(stage, event string, data interface{})
	Send(stage string, msg interface{})
	Receive(stage string) interface{}
	ListenForSignals() error
	RegisterProcess(name string, command string, args []string, restart bool) error
	StartProcess(proc IManagedProcess[any]) error
	StartAll() error
	StopAll() error
	StopProcess(proc IManagedProcess[any]) error
	RestartProcess(proc IManagedProcess[any]) error
	IsRunning() bool
	Pid() int
	Wait() error
	Status() string
	String() string
	SetArgs(args []string)
	SetEnv(env []string)
	SetDir(dir string)
	SetPort(port int)
	SetHost(host string)
	SetUser(user string)
	SetPass(pass string)
	SetCert(cert string)
	SetKey(key string)
	SetCA(ca string)
	SetSSL(ssl bool)
	SetTLS(tls bool)
	SetAuth(auth bool)
}
