package types

import (
	ci "github.com/rafa-mori/golife/internal/components/interfaces"
	gl "github.com/rafa-mori/golife/logger"
	l "github.com/rafa-mori/logz"
)

// newProcessInputFromConfig creates a new ProcessInput instance from the provided config data.
func newProcessInputFromConfig[T any](name string, data []byte, format string) (*ProcessInput[T], error) {
	mapper := NewMapper[ProcessInput[T]]()
	logger := l.GetLogger("GoLife")
	npi := newProcessInputSystemBase[T](name, "", nil, false, false, nil, logger, false)
	npOutput := &ProcessInput[T]{
		Reference:               npi.Reference,
		Logger:                  npi.Logger,
		Mutexes:                 npi.Mutexes,
		ProcessInputSystemBase:  npi,
		ProcessInputRuntimeBase: newProcessInputRuntimeBase[T, ci.IProcessInput[T]](name, nil, nil, false, false, logger, false),
		ProcessInputConfig:      npi.ProcessInputConfig,
	}
	err := mapper.Deserialize(data, npOutput, format)
	if err != nil {
		return nil, err
	}
	return npOutput, nil
}

// NewProcessInputFromConfig creates a new ProcessInput instance from the provided config data.
func NewProcessInputFromConfig[T any](name string, data []byte, format string) (ci.IProcessInput[T], error) {
	return newProcessInputFromConfig[T](name, data, format)
}

// newProcessInput creates a new ProcessInput instance with the provided Logger.
func newEmptyProcessInput[T any](name string, logger l.Logger, debug bool) *ProcessInput[T] {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	if debug {
		gl.SetDebug(debug)
	}
	mu := NewMutexesType()
	ref := NewReference(name)
	pc := &ProcessInput[T]{
		Logger:    logger,
		Mutexes:   mu,
		Reference: ref.GetReference(),
		ProcessInputConfig: &ProcessInputConfig{
			Logger:      logger,
			Reference:   ref.GetReference(),
			Mutexes:     mu,
			IsRunning:   false,
			WaitFor:     false,
			Restart:     false,
			ProcessType: "system",
			Metadata:    make(map[string]any),
		},
	}
	return pc
}

// NewEmptyProcessInput creates a new ProcessInput instance with the provided Logger.
func NewEmptyProcessInput[T any](name string, logger l.Logger, debug bool) ci.IProcessInput[T] {
	return newEmptyProcessInput[T](name, logger, debug)
}

// newProcessInput creates a new ProcessInput instance with the provided Logger.
func newProcessInput[T any](name, command string, args []string, waitFor bool, restart bool, function ci.IValidationFunc[T], logger l.Logger, debug bool) *ProcessInput[T] {
	if logger == nil {
		logger = l.GetLogger("GoLife")
	}
	if debug {
		gl.SetDebug(debug)
	}
	npi := newProcessInputSystemBase[T](name, command, args, waitFor, restart, function, logger, debug)
	npOutput := &ProcessInput[T]{
		Reference:               npi.Reference,
		Logger:                  npi.Logger,
		Mutexes:                 npi.Mutexes,
		ProcessInputSystemBase:  npi,
		ProcessInputRuntimeBase: newProcessInputRuntimeBase[T, ci.IProcessInput[T]](name, nil, nil, false, false, logger, false),
		ProcessInputConfig:      npi.ProcessInputConfig,
	}
	return npOutput
}

// NewProcessInput creates a new ProcessInput instance with the provided Logger.
func NewProcessInput[T any](name, command string, args []string, waitFor bool, restart bool, function ci.IValidationFunc[T], logger l.Logger, debug bool) ci.IProcessInput[T] {
	return newProcessInput[T](name, command, args, waitFor, restart, function, logger, debug)
}
