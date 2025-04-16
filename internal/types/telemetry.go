package types

import "time"

// Telemetry is a struct that holds telemetry data for a process
type Telemetry struct {
	LastUpdated time.Time
	Metrics     map[string]float64
}

// NewTelemetry creates a new Telemetry instance
func NewTelemetry() *Telemetry {
	return &Telemetry{
		LastUpdated: time.Now(),
		Metrics:     make(map[string]float64),
	}
}

// UpdateMetrics updates the telemetry metrics with the given data
func (t *Telemetry) UpdateMetrics(data map[string]float64) {
	t.LastUpdated = time.Now()
	for key, value := range data {
		t.Metrics[key] = value
	}
}

// GetMetrics returns the telemetry metrics
func (t *Telemetry) GetMetrics() map[string]float64 {
	return t.Metrics
}

// GetLastUpdated returns the last updated time of the telemetry
func (t *Telemetry) GetLastUpdated() time.Time { return t.LastUpdated }

// ResetMetrics resets the telemetry metrics
func (t *Telemetry) ResetMetrics() {
	t.Metrics = make(map[string]float64)
	t.LastUpdated = time.Now()
}
