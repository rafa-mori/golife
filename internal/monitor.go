package internal

import (
	l "github.com/faelmori/logz"

	"time"
)

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
	l.InfoCtx("Monitoring set", map[string]interface{}{"monitoring": monitoring})
}
func (m *ManagedMonit) SetInterval(interval time.Duration) {
	m.Interval = interval
	l.InfoCtx("Interval set", map[string]interface{}{"interval": interval})
}
func (m *ManagedMonit) SetTimeout(timeout time.Duration) {
	m.Timeout = timeout
	l.InfoCtx("Timeout set", map[string]interface{}{"timeout": timeout})
}
func (m *ManagedMonit) SetDelay(delay time.Duration) {
	m.Delay = delay
	l.InfoCtx("Delay set", map[string]interface{}{"delay": delay})
}
func (m *ManagedMonit) SetTimeouts(timeouts []time.Duration) {
	m.Timeouts = timeouts
	l.InfoCtx("Timeouts set", map[string]interface{}{"timeouts": timeouts})
}
func (m *ManagedMonit) SetDelays(delays []time.Duration) {
	m.Delays = delays
	l.InfoCtx("Delays set", map[string]interface{}{"delays": delays})
}
func (m *ManagedMonit) SetRunning(running bool) {
	m.Running = running
	l.InfoCtx("Running set", map[string]interface{}{"running": running})
}
func (m *ManagedMonit) SetStopped(stopped bool) {
	m.Stopped = stopped
	l.InfoCtx("Stopped set", map[string]interface{}{"stopped": stopped})
}
func (m *ManagedMonit) SetFailed(failed bool) {
	m.Failed = failed
	l.InfoCtx("Failed set", map[string]interface{}{"failed": failed})
}
func (m *ManagedMonit) SetSuccess(success bool) {
	m.Success = success
	l.InfoCtx("SuccessCtx set", map[string]interface{}{"success": success})
}
func (m *ManagedMonit) Start() error {
	l.InfoCtx("Starting monitor", nil)
	return nil
}
func (m *ManagedMonit) Stop() error {
	l.InfoCtx("Stopping monitor", nil)
	return nil
}
func (m *ManagedMonit) Restart() error {
	l.InfoCtx("Restarting monitor", nil)
	return nil
}
func (m *ManagedMonit) Reload() error {
	l.InfoCtx("Reloading monitor", nil)
	return nil
}
func (m *ManagedMonit) IsRunning() bool {
	l.InfoCtx("Checking if monitor is running", nil)
	return false
}
func (m *ManagedMonit) Pid() int {
	l.InfoCtx("Getting monitor PID", nil)
	return 0
}
func (m *ManagedMonit) Wait() error {
	l.InfoCtx("Waiting for monitor", nil)
	return nil
}
func (m *ManagedMonit) Status() string {
	l.InfoCtx("Getting monitor status", nil)
	return ""
}
func (m *ManagedMonit) String() string {
	l.InfoCtx("Getting monitor string representation", nil)
	return ""
}
func (m *ManagedMonit) Stdout() string {
	l.InfoCtx("Getting monitor stdout", nil)
	return ""
}
func (m *ManagedMonit) Stderr() string {
	l.InfoCtx("Getting monitor stderr", nil)
	return ""
}
func (m *ManagedMonit) Properties() map[string]interface{} {
	l.InfoCtx("Getting monitor properties", nil)
	return nil
}
func (m *ManagedMonit) Monitor() error {
	l.InfoCtx("Monitoring", nil)
	return nil
}

func NewManagedMonit() IManagedMonit {
	monit := ManagedMonit{}
	l.InfoCtx("Creating new ManagedMonit", nil)
	return &monit
}
