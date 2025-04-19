package routines

import (
	"fmt"
	"log"

	"context"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type IManagedGoroutine[T any] interface {
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

	Copy() IManagedGoroutine[T]
}

type ManagedGoroutine[T any] struct {
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

func (m *ManagedGoroutine[T]) Start() error {
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
func (m *ManagedGoroutine[T]) Stop() error {
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
func (m *ManagedGoroutine[T]) Pause() error {
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
func (m *ManagedGoroutine[T]) Resume() error {
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
func (m *ManagedGoroutine[T]) IsRunning() bool {
	if m == nil {
		return false
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineDone
}
func (m *ManagedGoroutine[T]) Wait() error {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineWG.Wait()

	return nil
}
func (m *ManagedGoroutine[T]) String() string {
	if m == nil {
		return ""
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return ""
}
func (m *ManagedGoroutine[T]) Send(msg interface{}) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

}
func (m *ManagedGoroutine[T]) Receive() interface{} {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return nil
}
func (m *ManagedGoroutine[T]) SetArgs(args []string) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()
}
func (m *ManagedGoroutine[T]) SetCommand(command string) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()
}
func (m *ManagedGoroutine[T]) SetName(name string) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()
}
func (m *ManagedGoroutine[T]) SetWaitFor(wait bool) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()
}
func (m *ManagedGoroutine[T]) SetGoroutineFn(fn func()) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineFn = fn
}
func (m *ManagedGoroutine[T]) SetGoroutineCh(ch chan struct{}) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineCh = ch
}
func (m *ManagedGoroutine[T]) SetGoroutineErr(err error) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineErr = err
}
func (m *ManagedGoroutine[T]) SetGoroutineDone(done bool) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineDone = done
}
func (m *ManagedGoroutine[T]) SetGoroutineWG(wg sync.WaitGroup) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineWG = wg
}
func (m *ManagedGoroutine[T]) SetGoroutineMu(mu sync.Mutex) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineMu = mu
}
func (m *ManagedGoroutine[T]) SetGoroutineOnce(once sync.Once) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineOnce = once
}
func (m *ManagedGoroutine[T]) SetGoroutineCond(cond *sync.Cond) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineCond = cond
}
func (m *ManagedGoroutine[T]) SetGoroutineLock(lock sync.Mutex) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineLock = lock
}
func (m *ManagedGoroutine[T]) SetGoroutineDoneCh(ch chan struct{}) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineDoneCh = ch
}
func (m *ManagedGoroutine[T]) SetGoroutineErrCh(ch chan error) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineErrCh = ch
}
func (m *ManagedGoroutine[T]) SetGoroutineCancel(cancel func()) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineCancel = cancel
}
func (m *ManagedGoroutine[T]) SetGoroutineCtx(ctx context.Context) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineCtx = ctx
}
func (m *ManagedGoroutine[T]) SetGoroutineCancelFn(cancelFn context.CancelFunc) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineCancelFn = cancelFn
}
func (m *ManagedGoroutine[T]) SetGoroutineTimeout(timeout time.Duration) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineTimeout = timeout
}
func (m *ManagedGoroutine[T]) SetGoroutineDeadline(deadline time.Time) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineDeadline = deadline
}
func (m *ManagedGoroutine[T]) SetGoroutineDeadlineSet(set bool) {
	if m == nil {
		return
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	m.goroutineDeadlineSet = set
}
func (m *ManagedGoroutine[T]) GetGoroutineFn() func() {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineFn
}
func (m *ManagedGoroutine[T]) GetGoroutineCh() chan struct{} {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineCh
}
func (m *ManagedGoroutine[T]) GetGoroutineErr() error {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineErr
}
func (m *ManagedGoroutine[T]) GetGoroutineDone() bool {
	if m == nil {
		return false
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineDone
}
func (m *ManagedGoroutine[T]) GetGoroutineWG() sync.WaitGroup {
	if m == nil {
		return sync.WaitGroup{}
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineWG
}
func (m *ManagedGoroutine[T]) GetGoroutineMu() sync.Mutex {
	if m == nil {
		return sync.Mutex{}
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineMu
}
func (m *ManagedGoroutine[T]) GetGoroutineOnce() sync.Once {
	if m == nil {
		return sync.Once{}
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineOnce
}
func (m *ManagedGoroutine[T]) GetGoroutineCond() *sync.Cond {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineCond
}
func (m *ManagedGoroutine[T]) GetGoroutineLock() sync.Mutex {
	if m == nil {
		return sync.Mutex{}
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineLock
}
func (m *ManagedGoroutine[T]) GetGoroutineDoneCh() chan struct{} {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineDoneCh
}
func (m *ManagedGoroutine[T]) GetGoroutineErrCh() chan error {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineErrCh
}
func (m *ManagedGoroutine[T]) GetGoroutineCancel() func() {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineCancel
}
func (m *ManagedGoroutine[T]) GetGoroutineCtx() context.Context {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineCtx
}
func (m *ManagedGoroutine[T]) GetGoroutineCancelFn() context.CancelFunc {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineCancelFn
}
func (m *ManagedGoroutine[T]) GetGoroutineTimeout() time.Duration {
	if m == nil {
		return 0
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineTimeout
}
func (m *ManagedGoroutine[T]) GetGoroutineDeadline() time.Time {
	if m == nil {
		return time.Time{}
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineDeadline
}
func (m *ManagedGoroutine[T]) GetGoroutineDeadlineSet() bool {
	if m == nil {
		return false
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return m.goroutineDeadlineSet
}
func (m *ManagedGoroutine[T]) Copy() IManagedGoroutine[T] {
	if m == nil {
		return nil
	}
	m.goroutineMu.Lock()
	defer m.goroutineMu.Unlock()

	return &ManagedGoroutine[T]{
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

func NewManagedGoroutine[T any](fn func()) *ManagedGoroutine[T] {
	return &ManagedGoroutine[T]{
		goroutineFn: fn,
	}
}

func logActivity(activity string) {
	cfgDir := filepath.Dir("")
	logFilePath := filepath.Join(cfgDir, "goroutine.log")

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(fmt.Sprintf("Erro ao abrir arquivo de log", err.Error()), activity)
		return
	}
	defer func(logFile *os.File) {
		_ = logFile.Close()
	}(logFile)

	logEntry := time.Now().Format(time.RFC3339) + ": " + activity + "\n"
	if _, err := logFile.WriteString(logEntry); err != nil {
		fmt.Println("Erro ao escrever no arquivo de log", map[string]interface{}{"error": err})
	}
}
