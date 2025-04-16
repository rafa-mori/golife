package events

import (
	"github.com/faelmori/golife/internal/process"
	t "github.com/faelmori/golife/internal/types"
	c "github.com/faelmori/golife/services"
	"github.com/google/uuid"
	"sync"
)

type IManagedProcessEvents[T any] interface {
	Event() string
	RegisterEvent(event string, fn func(interface{}))
	Trigger(stage, event string, data interface{})
	Send(stage string, msg interface{})
	Receive(stage string) interface{}
	ListenForSignals() error
	RegisterProcess(name string, command string, args []string, restart bool) error
	StartProcess(proc *process.ManagedProcess[any]) error
	StartAll() error
	StopAll() error
	StopProcess(proc *process.ManagedProcess[any]) error
	RestartProcess(proc *process.ManagedProcess[any]) error
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

//type ManagedProcessEvents[T any] struct {
//	EventFns  map[string]func(interface{})
//	TriggerCh chan interface{}
//	Name      string
//
//	// Internals
//	mu    sync.Mutex
//	Data  interface{}
//	Ev    string
//	Fn    func(interface{}) error
//	Stage string
//	User  string
//	Pass  string
//	Args  []string
//	Env   []string
//	Dir   string
//	Port  int
//	Host  string
//	Cert  string
//	Key   string
//	CA    string
//	SSL   bool
//	TLS   bool
//	Auth  bool
//}

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

func (m *ManagedProcessEvents[T]) Event() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.Name
}
func (m *ManagedProcessEvents[T]) RegisterEvent(event string, fn func(interface{})) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.EventFns[event] = fn
}
func (m *ManagedProcessEvents[T]) Trigger(stage, event string, data interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if fn, ok := m.EventFns[event]; ok {
		fn(data)
	}
}
func (m *ManagedProcessEvents[T]) Send(stage string, msg interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TriggerCh <- msg
}
func (m *ManagedProcessEvents[T]) Receive(stage string) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	return <-m.TriggerCh
}
func (m *ManagedProcessEvents[T]) ListenForSignals() error {
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
func (m *ManagedProcessEvents[T]) RegisterProcess(name string, command string, args []string, restart bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents[T]) StartProcess(proc *process.ManagedProcess[any]) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents[T]) StartAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents[T]) StopAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents[T]) StopProcess(proc *process.ManagedProcess[any]) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents[T]) RestartProcess(proc *process.ManagedProcess[any]) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents[T]) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return false
}
func (m *ManagedProcessEvents[T]) Pid() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return 0
}
func (m *ManagedProcessEvents[T]) Wait() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents[T]) Status() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return ""
}
func (m *ManagedProcessEvents[T]) String() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return ""
}
func (m *ManagedProcessEvents[T]) SetArgs(args []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Args = args
}
func (m *ManagedProcessEvents[T]) SetEnv(env []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Env = env
}
func (m *ManagedProcessEvents[T]) SetDir(dir string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Dir = dir
}
func (m *ManagedProcessEvents[T]) SetPort(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Port = port
}
func (m *ManagedProcessEvents[T]) SetHost(host string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Host = host
}
func (m *ManagedProcessEvents[T]) SetUser(user string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.User = user
}
func (m *ManagedProcessEvents[T]) SetPass(pass string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Pass = pass
}
func (m *ManagedProcessEvents[T]) SetCert(cert string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Cert = cert
}
func (m *ManagedProcessEvents[T]) SetKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Key = key
}
func (m *ManagedProcessEvents[T]) SetCA(ca string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CA = ca
}
func (m *ManagedProcessEvents[T]) SetSSL(ssl bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SSL = ssl
}
func (m *ManagedProcessEvents[T]) SetTLS(tls bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TLS = tls
}
func (m *ManagedProcessEvents[T]) SetAuth(auth bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Auth = auth
}

// </editor-fold>

func NewManagedProcessEvents[T any]() IManagedProcessEvents[T] {
	events := ManagedProcessEvents[T]{
		ID: uuid.New().String(),
		mu: sync.Mutex{},
		wg: sync.WaitGroup{},
	}
	return &events
}
