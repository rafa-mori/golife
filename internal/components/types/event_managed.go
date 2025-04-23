package types

import (
	ci "github.com/faelmori/golife/internal/components/interfaces"
	"github.com/google/uuid"
	"sync"
)

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
	EventProperties map[string]ci.IProperty[any]
	EventFuncList   map[string]func(interface{}) error
	EventAgents     map[string]ci.IChannelCtl[any]

	// Original Payload/Data/State
	Data T

	// Command properties
	//CmdName string
	//CmdArgs []string
	//CmdEnv  []string
	//CmdDir  string
	CmdProperties map[string]ci.IProperty[any]

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
	AccessProperties map[string]ci.IProperty[any]
}

// <editor-fold defaultstate="collapsed" desc="ManagedProcessEvents">

func (m *ManagedProcessEvents[T]) Event() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.ID
}
func (m *ManagedProcessEvents[T]) RegisterEvent(event string, fn func(interface{})) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.EventFuncList[event] = fn
}
func (m *ManagedProcessEvents[T]) Trigger(stage, event string, data interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//if fn, ok := m.EventFns[event]; ok {
	//	fn(data)
	//}
}
func (m *ManagedProcessEvents[T]) Send(stage string, msg interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.TriggerCh <- msg
}
func (m *ManagedProcessEvents[T]) Receive(stage string) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	//return <-m.TriggerCh
	return nil
}
func (m *ManagedProcessEvents[T]) ListenForSignals() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	//go func() {
	//	for {
	//		//select {
	//		//case <-m.TriggerCh:
	//		//	// Do something
	//		//}
	//	}
	//}()

	return nil
}
func (m *ManagedProcessEvents[T]) RegisterProcess(name string, command string, args []string, restart bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents[T]) StartProcess(proc ci.IManagedProcess[any]) error {
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
func (m *ManagedProcessEvents[T]) StopProcess(proc ci.IManagedProcess[any]) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents[T]) RestartProcess(proc ci.IManagedProcess[any]) error {
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

	//m.Args = args
}
func (m *ManagedProcessEvents[T]) SetEnv(env []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.Env = env
}
func (m *ManagedProcessEvents[T]) SetDir(dir string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.Dir = dir
}
func (m *ManagedProcessEvents[T]) SetPort(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.Port = port
}
func (m *ManagedProcessEvents[T]) SetHost(host string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.Host = host
}
func (m *ManagedProcessEvents[T]) SetUser(user string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.User = user
}
func (m *ManagedProcessEvents[T]) SetPass(pass string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.Pass = pass
}
func (m *ManagedProcessEvents[T]) SetCert(cert string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.Cert = cert
}
func (m *ManagedProcessEvents[T]) SetKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.Key = key
}
func (m *ManagedProcessEvents[T]) SetCA(ca string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.CA = ca
}
func (m *ManagedProcessEvents[T]) SetSSL(ssl bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.SSL = ssl
}
func (m *ManagedProcessEvents[T]) SetTLS(tls bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.TLS = tls
}
func (m *ManagedProcessEvents[T]) SetAuth(auth bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	//m.Auth = auth
}

// </editor-fold>

func NewManagedProcessEvents[T any]() ci.IManagedProcessEvents[T] {
	events := ManagedProcessEvents[T]{
		ID: uuid.New().String(),
		mu: sync.Mutex{},
		wg: sync.WaitGroup{},
	}
	return &events
}
