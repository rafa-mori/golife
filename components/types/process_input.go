package types

import (
	ci "github.com/faelmori/golife/components/interfaces"
	gl "github.com/faelmori/golife/logger"
	l "github.com/faelmori/logz"

	"fmt"
	"strings"
)

// ProcessInput is a struct that holds the input for the process.
type ProcessInput[T any] struct {
	// Logger is the logger for the process
	Logger l.Logger
	// Mutexes is the mutex for the process
	*Mutexes
	// chCtl is the channel for controlling the process, this will be passed through the flow into the many functions, frames/layers
	// and will be used to control the process, stages, events, etc.
	*ChannelCtl[string]

	// Reference is the reference for the process with ID and Name
	*Reference `json:"reference" yaml:"reference" xml:"reference" gorm:"reference"`
	// ProcessInputConfig is the configuration for the process
	*ProcessInputConfig `json:"config" yaml:"config" xml:"config" gorm:"config"`
	// ProcessInputRuntimeBase is the system information for the process
	*ProcessInputRuntimeBase[T, ci.IProcessInput[T]] `json:"system,omitempty" yaml:"system,omitempty" xml:"system,omitempty" gorm:"system,omitempty"`
	// ProcessInputRuntimeBase is the runtime information for the process
	*ProcessInputSystemBase[T, ci.IProcessInput[T]] `json:"runtime,omitempty" yaml:"runtime,omitempty" xml:"runtime,omitempty" gorm:"runtime,omitempty"`
}

// newProcessInputFromConfig creates a new ProcessInput instance from the provided config data.
func newProcessInputFromConfig[T any](name string, data []byte, format string) (*ProcessInput[T], error) {
	mapper := NewMapper[ProcessInput[T]]()
	logger := l.GetLogger("GoLife")
	npi := newProcessInputSystemBase[T, ci.IProcessInput[T]](name, "", nil, false, false, nil, logger, false)
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
func newProcessInput[T any](name string, logger l.Logger, debug bool) *ProcessInput[T] {
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

// NewProcessInput creates a new ProcessInput instance with the provided Logger.
func NewProcessInput[T any](name string, logger l.Logger, debug bool) ci.IProcessInput[T] {
	return newProcessInput[T](name, logger, debug)
}

// Serialize serializes the ProcessInput instance to the specified format.
func (pi *ProcessInput[T]) Serialize(format string) ([]byte, error) {
	mapper := NewMapper[ProcessInput[T]]()
	return mapper.Serialize(nil, pi, format)
}

// Deserialize deserializes the data into the ProcessInput instance.
func (pi *ProcessInput[T]) Deserialize(data []byte, format string) error {
	mapper := NewMapper[ProcessInput[T]]()
	return mapper.Deserialize(data, pi, format)
}

// GetLogger returns the Logger instance for the process.
func (pi *ProcessInput[T]) GetLogger() l.Logger {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.Logger
}

// GetReference returns the reference ID and name.
func (pi *ProcessInput[T]) GetReference() ci.IReference {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.Reference
}

// GetName returns the name of the process.
func (pi *ProcessInput[T]) GetName() string {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.Reference.Name
}

// Validate validates the ProcessInput instance.
func (pi *ProcessInput[T]) Validate() ci.IValidationResult {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()
	validations := map[string]func() bool{
		"Command is nil":       func() bool { return pi.BuildCmd() != nil },
		"Command is empty":     func() bool { return pi.Command != "" },
		"Args is nil":          func() bool { return pi.Args != nil },
		"Path is empty":        func() bool { return pi.Path != "" },
		"ProcessType is empty": func() bool { return pi.ProcessType != "" },
		"Reference is nil":     func() bool { return pi.Reference != nil },
		"Mutexes is nil":       func() bool { return pi.Mutexes != nil },
	}

	messageStringBuilder := strings.Builder{}
	var isInvalid bool
	for message, isValid := range validations {
		if !isValid() {
			isInvalid = true
			if messageStringBuilder.Len() > 0 {
				messageStringBuilder.WriteString(message + "\n")
			} else {
				messageStringBuilder.WriteString("ProcessInput is invalid:\n")
				messageStringBuilder.WriteString(message + "\n")
			}
		}
	}

	if isInvalid {
		return &ValidationResult{
			IsValid: false,
			Message: messageStringBuilder.String(),
			Error:   fmt.Errorf("ProcessInput validation failed. Error: %s", messageStringBuilder.String()),
		}
	}

	return &ValidationResult{
		IsValid: true,
		Message: "ProcessInput is valid",
		Error:   nil,
	}
}
