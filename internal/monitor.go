package internal

import "time"

type IManagedMonit interface {
	Stdout() string
	Stderr() string
	Properties() map[string]interface{}
	Monitor() error

	// Monitoramento
	SetMonitoring(monitoring bool)
	SetInterval(interval time.Duration)
	SetTimeout(timeout time.Duration)
	SetDelay(delay time.Duration)
	SetTimeouts(timeouts []time.Duration)
	SetDelays(delays []time.Duration)

	// Status
	SetRunning(running bool)
	SetStopped(stopped bool)
	SetFailed(failed bool)
	SetSuccess(success bool)

	// Comandos
	Start() error
	Stop() error
	Restart() error
	Reload() error
	IsRunning() bool
	Pid() int
	Wait() error
	Status() string
	String() string
}

type ManagedMonit struct {
	// Monitoramento
	Monitoring bool
	Interval   time.Duration
	Timeout    time.Duration
	Delay      time.Duration
	Timeouts   []time.Duration
	Delays     []time.Duration

	// Status
	Running bool
	Stopped bool
	Failed  bool
	Success bool
}

func (m *ManagedMonit) SetMonitoring(monitoring bool) {
	m.Monitoring = monitoring
}
func (m *ManagedMonit) SetInterval(interval time.Duration) {
	m.Interval = interval
}
func (m *ManagedMonit) SetTimeout(timeout time.Duration) {
	m.Timeout = timeout
}
func (m *ManagedMonit) SetDelay(delay time.Duration) {
	m.Delay = delay
}
func (m *ManagedMonit) SetTimeouts(timeouts []time.Duration) {
	m.Timeouts = timeouts
}
func (m *ManagedMonit) SetDelays(delays []time.Duration) {
	m.Delays = delays
}

func (m *ManagedMonit) SetRunning(running bool) {
	m.Running = running
}
func (m *ManagedMonit) SetStopped(stopped bool) {
	m.Stopped = stopped
}
func (m *ManagedMonit) SetFailed(failed bool) {
	m.Failed = failed
}
func (m *ManagedMonit) SetSuccess(success bool) {
	m.Success = success
}

func (m *ManagedMonit) Start() error {
	return nil
}
func (m *ManagedMonit) Stop() error {
	return nil
}
func (m *ManagedMonit) Restart() error {
	return nil
}
func (m *ManagedMonit) Reload() error {
	return nil
}
func (m *ManagedMonit) IsRunning() bool {
	return false
}
func (m *ManagedMonit) Pid() int {
	return 0
}
func (m *ManagedMonit) Wait() error {
	return nil
}
func (m *ManagedMonit) Status() string {
	return ""
}
func (m *ManagedMonit) String() string {
	return ""
}
func (m *ManagedMonit) Stdout() string {
	return ""
}
func (m *ManagedMonit) Stderr() string {
	return ""
}
func (m *ManagedMonit) Properties() map[string]interface{} {
	return nil
}
func (m *ManagedMonit) Monitor() error {
	return nil
}

func NewManagedMonit() IManagedMonit {
	monit := ManagedMonit{}
	return &monit
}
