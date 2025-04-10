package internal

import (
	"github.com/faelmori/kbxutils/factory"
	l "github.com/faelmori/logz"

	"context"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type IManagedGoroutine interface {
	Start() error
	Stop() error
	Pause() error
	Resume() error
	IsRunning() bool
	Wait() error
	String() string

	Send(msg interface{})
	Receive() interface{}

	SetArgs(args []string)
	SetCommand(command string)
	SetName(name string)
	SetWaitFor(wait bool)

	SetGoroutineFn(fn func())
	SetGoroutineCh(ch chan struct{})
	SetGoroutineErr(err error)
	SetGoroutineDone(done bool)
	SetGoroutineWG(wg sync.WaitGroup)
	SetGoroutineMu(mu sync.Mutex)
	SetGoroutineOnce(once sync.Once)
	SetGoroutineCond(cond *sync.Cond)
	SetGoroutineLock(lock sync.Mutex)
	SetGoroutineDoneCh(ch chan struct{})
	SetGoroutineErrCh(ch chan error)
	SetGoroutineCancel(cancel func())
	SetGoroutineCtx(ctx context.Context)
	SetGoroutineCancelFn(cancelFn context.CancelFunc)
	SetGoroutineTimeout(timeout time.Duration)
	SetGoroutineDeadline(deadline time.Time)
	SetGoroutineDeadlineSet(set bool)

	GetGoroutineFn() func()
	GetGoroutineCh() chan struct{}
	GetGoroutineErr() error
	GetGoroutineDone() bool
	GetGoroutineWG() sync.WaitGroup
	GetGoroutineMu() sync.Mutex
	GetGoroutineOnce() sync.Once
	GetGoroutineCond() *sync.Cond
	GetGoroutineLock() sync.Mutex
	GetGoroutineDoneCh() chan struct{}
	GetGoroutineErrCh() chan error
	GetGoroutineCancel() func()
	GetGoroutineCtx() context.Context
	GetGoroutineCancelFn() context.CancelFunc
	GetGoroutineTimeout() time.Duration
	GetGoroutineDeadline() time.Time
	GetGoroutineDeadlineSet() bool

	Copy() IManagedGoroutine
}

type ManagedGoroutine struct {
	// Processo gerenciado (terceira opção - goroutine)
	goroutineFn          func()
	goroutineCh          chan struct{}
	goroutineErr         error
	goroutineDone        bool
	goroutineWG          sync.WaitGroup
	goroutineMu          sync.Mutex
	goroutineOnce        sync.Once
	goroutineCond        *sync.Cond
	goroutineLock        sync.Mutex
	goroutineDoneCh      chan struct{}
	goroutineErrCh       chan error
	goroutineCancel      func()
	goroutineCtx         context.Context
	goroutineCancelFn    context.CancelFunc
	goroutineTimeout     time.Duration
	goroutineDeadline    time.Time
	goroutineDeadlineSet bool
}

func (m *ManagedGoroutine) Start() error {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	if m.IsRunning() {
		return nil
	}

	m.goroutineWG.Add(1)
	go func() {
		defer m.goroutineWG.Done()
		m.goroutineFn()
		logActivity("Goroutine started")
	}()

	return nil
}
func (m *ManagedGoroutine) Stop() error {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	if !m.IsRunning() {
		return nil
	}

	m.goroutineCancelFn()
	m.goroutineWG.Wait()
	logActivity("Goroutine stopped")

	return nil
}
func (m *ManagedGoroutine) Pause() error {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	if !m.IsRunning() {
		return nil
	}

	m.goroutineCancelFn()
	logActivity("Goroutine paused")

	return nil
}
func (m *ManagedGoroutine) Resume() error {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	if m.IsRunning() {
		return nil
	}

	m.goroutineWG.Add(1)
	go func() {
		defer m.goroutineWG.Done()
		m.goroutineFn()
		logActivity("Goroutine resumed")
	}()

	return nil
}
func (m *ManagedGoroutine) IsRunning() bool {
	if m == nil {
		return false
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineDone
}
func (m *ManagedGoroutine) Wait() error {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineWG.Wait()

	return nil
}
func (m *ManagedGoroutine) String() string {
	if m == nil {
		return ""
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return ""
}
func (m *ManagedGoroutine) Send(msg interface{}) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

}
func (m *ManagedGoroutine) Receive() interface{} {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return nil
}
func (m *ManagedGoroutine) SetArgs(args []string) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()
}
func (m *ManagedGoroutine) SetCommand(command string) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()
}
func (m *ManagedGoroutine) SetName(name string) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()
}
func (m *ManagedGoroutine) SetWaitFor(wait bool) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()
}
func (m *ManagedGoroutine) SetGoroutineFn(fn func()) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineFn = fn
}
func (m *ManagedGoroutine) SetGoroutineCh(ch chan struct{}) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineCh = ch
}
func (m *ManagedGoroutine) SetGoroutineErr(err error) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineErr = err
}
func (m *ManagedGoroutine) SetGoroutineDone(done bool) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineDone = done
}
func (m *ManagedGoroutine) SetGoroutineWG(wg sync.WaitGroup) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineWG = wg
}
func (m *ManagedGoroutine) SetGoroutineMu(mu sync.Mutex) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineMu = mu
}
func (m *ManagedGoroutine) SetGoroutineOnce(once sync.Once) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineOnce = once
}
func (m *ManagedGoroutine) SetGoroutineCond(cond *sync.Cond) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineCond = cond
}
func (m *ManagedGoroutine) SetGoroutineLock(lock sync.Mutex) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineLock = lock
}
func (m *ManagedGoroutine) SetGoroutineDoneCh(ch chan struct{}) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineDoneCh = ch
}
func (m *ManagedGoroutine) SetGoroutineErrCh(ch chan error) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineErrCh = ch
}
func (m *ManagedGoroutine) SetGoroutineCancel(cancel func()) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineCancel = cancel
}
func (m *ManagedGoroutine) SetGoroutineCtx(ctx context.Context) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineCtx = ctx
}
func (m *ManagedGoroutine) SetGoroutineCancelFn(cancelFn context.CancelFunc) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineCancelFn = cancelFn
}
func (m *ManagedGoroutine) SetGoroutineTimeout(timeout time.Duration) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineTimeout = timeout
}
func (m *ManagedGoroutine) SetGoroutineDeadline(deadline time.Time) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineDeadline = deadline
}
func (m *ManagedGoroutine) SetGoroutineDeadlineSet(set bool) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineDeadlineSet = set
}
func (m *ManagedGoroutine) GetGoroutineFn() func() {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineFn
}
func (m *ManagedGoroutine) GetGoroutineCh() chan struct{} {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineCh
}
func (m *ManagedGoroutine) GetGoroutineErr() error {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineErr
}
func (m *ManagedGoroutine) GetGoroutineDone() bool {
	if m == nil {
		return false
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineDone
}
func (m *ManagedGoroutine) GetGoroutineWG() sync.WaitGroup {
	if m == nil {
		return sync.WaitGroup{}
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineWG
}
func (m *ManagedGoroutine) GetGoroutineMu() sync.Mutex {
	if m == nil {
		return sync.Mutex{}
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineMu
}
func (m *ManagedGoroutine) GetGoroutineOnce() sync.Once {
	if m == nil {
		return sync.Once{}
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineOnce
}
func (m *ManagedGoroutine) GetGoroutineCond() *sync.Cond {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineCond
}
func (m *ManagedGoroutine) GetGoroutineLock() sync.Mutex {
	if m == nil {
		return sync.Mutex{}
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineLock
}
func (m *ManagedGoroutine) GetGoroutineDoneCh() chan struct{} {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineDoneCh
}
func (m *ManagedGoroutine) GetGoroutineErrCh() chan error {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineErrCh
}
func (m *ManagedGoroutine) GetGoroutineCancel() func() {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineCancel
}
func (m *ManagedGoroutine) GetGoroutineCtx() context.Context {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineCtx
}
func (m *ManagedGoroutine) GetGoroutineCancelFn() context.CancelFunc {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineCancelFn
}
func (m *ManagedGoroutine) GetGoroutineTimeout() time.Duration {
	if m == nil {
		return 0
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineTimeout
}
func (m *ManagedGoroutine) GetGoroutineDeadline() time.Time {
	if m == nil {
		return time.Time{}
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineDeadline
}
func (m *ManagedGoroutine) GetGoroutineDeadlineSet() bool {
	if m == nil {
		return false
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineDeadlineSet
}
func (m *ManagedGoroutine) Copy() IManagedGoroutine {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return &ManagedGoroutine{
		goroutineFn:          m.goroutineFn,
		goroutineCh:          m.goroutineCh,
		goroutineErr:         m.goroutineErr,
		goroutineDone:        m.goroutineDone,
		goroutineWG:          m.goroutineWG,
		goroutineMu:          m.goroutineMu,
		goroutineOnce:        m.goroutineOnce,
		goroutineCond:        m.goroutineCond,
		goroutineLock:        m.goroutineLock,
		goroutineDoneCh:      m.goroutineDoneCh,
		goroutineErrCh:       m.goroutineErrCh,
		goroutineCancel:      m.goroutineCancel,
		goroutineCtx:         m.goroutineCtx,
		goroutineCancelFn:    m.goroutineCancelFn,
		goroutineTimeout:     m.goroutineTimeout,
		goroutineDeadline:    m.goroutineDeadline,
		goroutineDeadlineSet: m.goroutineDeadlineSet,
	}
}

func NewManagedGoroutine(fn func()) *ManagedGoroutine {
	return &ManagedGoroutine{
		goroutineFn: fn,
	}
}

func logActivity(activity string) {
	fs := *factory.NewFilesystemService("")
	cfgDir := filepath.Dir(fs.GetConfigFilePath())
	logFilePath := filepath.Join(cfgDir, "goroutine.log")

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l.ErrorCtx("Erro ao abrir arquivo de log", map[string]interface{}{"error": err})
		return
	}
	defer func(logFile *os.File) {
		_ = logFile.Close()
	}(logFile)

	logEntry := time.Now().Format(time.RFC3339) + ": " + activity + "\n"
	if _, err := logFile.WriteString(logEntry); err != nil {
		l.ErrorCtx("Erro ao escrever no arquivo de log", map[string]interface{}{"error": err})
	}
}
