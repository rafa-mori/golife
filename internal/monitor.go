package internal

import (
	"time"
	"github.com/rafa-mori/logz"
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
	logz.Info("Monitoring set", map[string]interface{}{"monitoring": monitoring})
}
func (m *ManagedMonit) SetInterval(interval time.Duration) {
	m.Interval = interval
	logz.Info("Interval set", map[string]interface{}{"interval": interval})
}
func (m *ManagedMonit) SetTimeout(timeout time.Duration) {
	m.Timeout = timeout
	logz.Info("Timeout set", map[string]interface{}{"timeout": timeout})
}
func (m *ManagedMonit) SetDelay(delay time.Duration) {
	m.Delay = delay
	logz.Info("Delay set", map[string]interface{}{"delay": delay})
}
func (m *ManagedMonit) SetTimeouts(timeouts []time.Duration) {
	m.Timeouts = timeouts
	logz.Info("Timeouts set", map[string]interface{}{"timeouts": timeouts})
}
func (m *ManagedMonit) SetDelays(delays []time.Duration) {
	m.Delays = delays
	logz.Info("Delays set", map[string]interface{}{"delays": delays})
}
func (m *ManagedMonit) SetRunning(running bool) {
	m.Running = running
	logz.Info("Running set", map[string]interface{}{"running": running})
}
func (m *ManagedMonit) SetStopped(stopped bool) {
	m.Stopped = stopped
	logz.Info("Stopped set", map[string]interface{}{"stopped": stopped})
}
func (m *ManagedMonit) SetFailed(failed bool) {
	m.Failed = failed
	logz.Info("Failed set", map[string]interface{}{"failed": failed})
}
func (m *ManagedMonit) SetSuccess(success bool) {
	m.Success = success
	logz.Info("Success set", map[string]interface{}{"success": success})
}
func (m *ManagedMonit) Start() error {
	logz.Info("Starting monitor", nil)
	return nil
}
func (m *ManagedMonit) Stop() error {
	logz.Info("Stopping monitor", nil)
	return nil
}
func (m *ManagedMonit) Restart() error {
	logz.Info("Restarting monitor", nil)
	return nil
}
func (m *ManagedMonit) Reload() error {
	logz.Info("Reloading monitor", nil)
	return nil
}
func (m *ManagedMonit) IsRunning() bool {
	logz.Info("Checking if monitor is running", nil)
	return false
}
func (m *ManagedMonit) Pid() int {
	logz.Info("Getting monitor PID", nil)
	return 0
}
func (m *ManagedMonit) Wait() error {
	logz.Info("Waiting for monitor", nil)
	return nil
}
func (m *ManagedMonit) Status() string {
	logz.Info("Getting monitor status", nil)
	return ""
}
func (m *ManagedMonit) String() string {
	logz.Info("Getting monitor string representation", nil)
	return ""
}
func (m *ManagedMonit) Stdout() string {
	logz.Info("Getting monitor stdout", nil)
	return ""
}
func (m *ManagedMonit) Stderr() string {
	logz.Info("Getting monitor stderr", nil)
	return ""
}
func (m *ManagedMonit) Properties() map[string]interface{} {
	logz.Info("Getting monitor properties", nil)
	return nil
}
func (m *ManagedMonit) Monitor() error {
	logz.Info("Monitoring", nil)
	return nil
}

func NewManagedMonit() IManagedMonit {
	monit := ManagedMonit{}
	logz.Info("Creating new ManagedMonit", nil)
	return &monit
}
