package types

type Property[T any] struct {
	// Telemetry is the telemetry for this GoLife instance.
	metrics *Telemetry
	// Prop is the property for this GoLife instance.
	Prop *val[T]
	// Cb is the callback function for this GoLife instance.
	Cb func(any) (bool, error)
}

// NewProperty creates a new property with the given value and Reference.
func NewProperty[T any](name string, v *T, withMetrics bool, cb func(any) (bool, error)) *Property[T] {
	p := &Property[T]{
		Prop: newVal[T](name, v),
		Cb:   cb,
	}
	if withMetrics {
		p.metrics = NewTelemetry()
	}
	return p
}
