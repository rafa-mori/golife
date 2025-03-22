package internal

import "sync"

type IManagedProcessEvents interface {
	Event() string
	RegisterEvent(event string, fn func(interface{}))
	Trigger(stage, event string, data interface{})
	Send(stage string, msg interface{})
	Receive(stage string) interface{}
	ListenForSignals() error
	RegisterProcess(name string, command string, args []string, restart bool) error
	StartProcess(proc *ManagedProcess) error
	StartAll() error
	StopAll() error
	StopProcess(proc *ManagedProcess) error
	RestartProcess(proc *ManagedProcess) error
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
type ManagedProcessEvents struct {
	EventFns  map[string]func(interface{})
	TriggerCh chan interface{}
	Name      string

	// Internals
	mu    sync.Mutex
	Data  interface{}
	Ev    string
	Fn    func(interface{}) error
	Stage string
	User  string
	Pass  string
	Args  []string
	Env   []string
	Dir   string
	Port  int
	Host  string
	Cert  string
	Key   string
	CA    string
	SSL   bool
	TLS   bool
	Auth  bool
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
func (m *ManagedProcessEvents) StartProcess(proc *ManagedProcess) error {
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
func (m *ManagedProcessEvents) StopProcess(proc *ManagedProcess) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}
func (m *ManagedProcessEvents) RestartProcess(proc *ManagedProcess) error {
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
