package events

import (
	"github.com/faelmori/golife/internal/process"
	c "github.com/faelmori/golife/internal/routines/chan"
	t "github.com/faelmori/golife/internal/types"
	"sync"
)

type IManagedProcessEvents interface {
	Event() string
	RegisterEvent(event string, fn func(interface{}))
	Trigger(stage, event string, data interface{})
	Send(stage string, msg interface{})
	Receive(stage string) interface{}
	ListenForSignals() error
	RegisterProcess(name string, command string, args []string, restart bool) error
	StartProcess(proc *process.ManagedProcess) error
	StartAll() error
	StopAll() error
	StopProcess(proc *process.ManagedProcess) error
	RestartProcess(proc *process.ManagedProcess) error
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
type ManagedProcessEvents[T any] struct {
	// ID and Reference
	ID string

	// Thread-safe channel for event functions
	mu sync.Mutex
	wg sync.WaitGroup

	// Event Dynamic Secure Properties
	// Ex:
	//		EventName  string
	//		EventFunc  func(...any) error
	//		EventStage string
	EventProperties map[string]t.Property[any]
	EventFuncList   map[string]t.GenericChannelCallback[any]
	EventAgents     map[string]c.IChannel[any, int]

	// Original Payload/Data/State
	Data T

	// Command properties
	//CmdName string
	//CmdArgs []string
	//CmdEnv  []string
	//CmdDir  string
	CmdProperties map[string]t.Property[any]

	// Access properties
	// Ex:
	//User string
	//Pass string
	//Port   int
	//Host   string
	//Cert   string
	//Key    string
	//CA     string
	//IsSSL  bool
	//IsTLS  bool
	//IsAuth bool
	AccessProperties map[string]t.Property[any]
}

// <editor-fold defaultstate="collapsed" desc="ManagedProcessEvents">

func (m *ManagedProcessEvents) Event() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.Name
}
func (m *ManagedProcessEvents) RegisterEvent(event string, fn func(interface{})) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.EventFns[event] = fn
}
func (m *ManagedProcessEvents) Trigger(stage, event string, data interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if fn, ok := m.EventFns[event]; ok {
		fn(data)
	}
}
func (m *ManagedProcessEvents) Send(stage string, msg interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TriggerCh <- msg
}
func (m *ManagedProcessEvents) Receive(stage string) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	return <-m.TriggerCh
}
func (m *ManagedProcessEvents) ListenForSignals() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	go func() {
		for {
			select {
			case <-m.TriggerCh:
				// Do something
			}
		}
	}()

	return nil
}
func (m *ManagedProcessEvents) RegisterProcess(name string, command string, args []string, restart bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents) StartProcess(proc *process.ManagedProcess) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents) StartAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents) StopAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents) StopProcess(proc *process.ManagedProcess) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents) RestartProcess(proc *process.ManagedProcess) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return false
}
func (m *ManagedProcessEvents) Pid() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return 0
}
func (m *ManagedProcessEvents) Wait() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents) Status() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return ""
}
func (m *ManagedProcessEvents) String() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return ""
}
func (m *ManagedProcessEvents) SetArgs(args []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Args = args
}
func (m *ManagedProcessEvents) SetEnv(env []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Env = env
}
func (m *ManagedProcessEvents) SetDir(dir string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Dir = dir
}
func (m *ManagedProcessEvents) SetPort(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Port = port
}
func (m *ManagedProcessEvents) SetHost(host string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Host = host
}
func (m *ManagedProcessEvents) SetUser(user string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.User = user
}
func (m *ManagedProcessEvents) SetPass(pass string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Pass = pass
}
func (m *ManagedProcessEvents) SetCert(cert string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Cert = cert
}
func (m *ManagedProcessEvents) SetKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Key = key
}
func (m *ManagedProcessEvents) SetCA(ca string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CA = ca
}
func (m *ManagedProcessEvents) SetSSL(ssl bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SSL = ssl
}
func (m *ManagedProcessEvents) SetTLS(tls bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TLS = tls
}
func (m *ManagedProcessEvents) SetAuth(auth bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Auth = auth
}

// </editor-fold>

func NewManagedProcessEvents(eventFns map[string]func(interface{}), triggerCh chan interface{}) IManagedProcessEvents {
	if eventFns == nil {
		eventFns = make(map[string]func(interface{}))
	}
	if triggerCh == nil {
		triggerCh = make(chan interface{}, 100)
	}

	events := ManagedProcessEvents{
		EventFns:  eventFns,
		TriggerCh: triggerCh,
	}
	return &events
}
