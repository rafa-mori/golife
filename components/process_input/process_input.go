package process_input

import (
	"fmt"
	ci "github.com/faelmori/golife/components/interfaces"
	t "github.com/faelmori/golife/components/types"
	l "github.com/faelmori/logz"
	"strings"
)

// ProcessInput is a struct that holds the input for the process.
type ProcessInput[T any] struct {
	// Logger is the logger for the process
	Logger l.Logger
	// Mutexes is the mutex for the process
	*t.Mutexes
	// chCtl is the channel for controlling the process, this will be passed through the flow into the many functions, frames/layers
	// and will be used to control the process, stages, events, etc.
	*t.ChannelCtl[string]

	// Reference is the reference for the process with ID and Name
	*t.Reference `json:"reference" yaml:"reference" xml:"reference" gorm:"reference"`
	// ProcessConfig is the configuration for the process
	*ProcessConfig[T] `json:"config" yaml:"config" xml:"config" gorm:"config"`
	// ProcessSystemBase is the system information for the process
	*ProcessSystemBase[ProcessInput[T], T] `json:"system,omitempty" yaml:"system,omitempty" xml:"system,omitempty" gorm:"system,omitempty"`
	// ProcessRuntimeBase is the runtime information for the process
	*ProcessRuntimeBase[ProcessInput[T], T] `json:"runtime,omitempty" yaml:"runtime,omitempty" xml:"runtime,omitempty" gorm:"runtime,omitempty"`
}

// NewProcessInputFromConfig creates a new ProcessInput instance from the provided config data.
func NewProcessInputFromConfig[T any](name string, data []byte, format string) (*ProcessInput[T], error) {
	mapper := t.NewMapper[ProcessInput[T]]()
	logger := l.GetLogger("GoLife")
	npi := newSystemProcessInput[ProcessInput[T], T](name, "", nil, false, false, nil, logger, false)
	npr := newProcessRuntimeBase[ProcessInput[T], T](name, new(T), nil, false, false, logger, false)
	np := &ProcessInput[T]{
		Reference:          npr.Reference,
		Logger:             npr.Logger,
		Mutexes:            npr.Mutexes,
		ProcessSystemBase:  npi,
		ProcessRuntimeBase: npr,
		ProcessConfig:      npr.ProcessConfig,
	}
	err := mapper.Deserialize(data, np, format)
	if err != nil {
		return nil, err
	}
	return np, nil
}

// NewSystemProcessInput creates a new ProcessInput instance with the provided Logger.
func NewSystemProcessInput[T any](name, command string, args []string, waitFor bool, restart bool, function *t.ValidationFunc[ProcessInput[T]], logger l.Logger, debug bool) *ProcessInput[T] {
	npi := newSystemProcessInput[ProcessInput[T], T](name, command, args, waitFor, restart, function, logger, debug)
	npr := newProcessRuntimeBase[ProcessInput[T], T](name, new(T), function, waitFor, restart, logger, debug)

	pi := &ProcessInput[T]{
		Reference:          npi.Reference,
		Logger:             npi.Logger,
		Mutexes:            npi.Mutexes,
		ProcessSystemBase:  npi,
		ProcessRuntimeBase: npr,
		ProcessConfig:      npi.ProcessConfig,
	}

	prp := t.NewProperty[ProcessInput[T]]("ctl", pi, false, nil)

	chCtl := t.NewChannelCtlWithProperty[ProcessInput[T], ci.IProperty[ProcessInput[T]]](
		name,
		func(b int) *int { return &b }(10),
		prp,
		true,
		logger,
	)

	chs := func(cm map[string]any) map[string]any {
		pi.Mutexes.MuLock()
		defer pi.Mutexes.MuUnlock()
		cm["ctl"] = chCtl
		return cm
	}(pi.GetSubChannels())

	pi.SetSubChannels(chs)

	return pi
}

// Serialize serializes the ProcessInput instance to the specified format.
func (pi *ProcessInput[T]) Serialize(format string) ([]byte, error) {
	mapper := t.NewMapper[ProcessInput[T]]()
	return mapper.Serialize(nil, pi, format)
}

// Deserialize deserializes the data into the ProcessInput instance.
func (pi *ProcessInput[T]) Deserialize(data []byte, format string) error {
	mapper := t.NewMapper[ProcessInput[T]]()
	return mapper.Deserialize(data, pi, format)
}

// GetLogger returns the Logger instance for the process.
func (pi *ProcessInput[T]) GetLogger() l.Logger {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	return pi.Logger
}

// GetReference returns the reference ID and name.
func (pi *ProcessInput[T]) GetReference() *t.Reference {
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
func (pi *ProcessInput[T]) Validate() *t.ValidationResult {
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
		return &t.ValidationResult{
			IsValid: false,
			Message: messageStringBuilder.String(),
			Error:   fmt.Errorf("ProcessInput validation failed. Error: %s", messageStringBuilder.String()),
		}
	}

	return &t.ValidationResult{
		IsValid: true,
		Message: "ProcessInput is valid",
		Error:   nil,
	}
}
