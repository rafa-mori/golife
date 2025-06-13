package types

import (
	"context"
	ci "github.com/rafa-mori/golife/internal/components/interfaces"
	gl "github.com/rafa-mori/golife/logger"
	l "github.com/rafa-mori/logz"
	"os"

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

// Send sends a message to the channel.
func (pi *ProcessInput[T]) Send(msg string, cb any) error {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	if pi.ChannelCtl == nil || pi.ChannelCtl.ch == nil {
		gl.LogObjLogger(pi, "notice", "ChannelCtl is nil, creating a new one")
		pi.ChannelCtl = NewChannelCtl[string](pi.Name, pi.Logger).(*ChannelCtl[string])
		pi.ChannelCtl.ch = make(chan string, 1)
		pi.ChannelCtl.ch <- msg
		return nil
	}

	select {
	case pi.ChannelCtl.ch <- msg:
	default:
		return fmt.Errorf("channel buffer full, message not sent")
	}

	if cb != nil {
		if callback, ok := cb.(func(string)); ok {
			callback(msg)
		} else {
			gl.LogObjLogger(pi, "error", "Callback is not of type func(string)")
			return fmt.Errorf("callback is not of type func(string)")
		}
	}
	return nil
}

// Receive receives a message from the channel.
func (pi *ProcessInput[T]) Receive(ctx context.Context, cb any) (any, error) {
	pi.Mutexes.MuRLock()
	defer pi.Mutexes.MuRUnlock()

	if pi.ChannelCtl == nil || pi.ChannelCtl.ch == nil {
		return nil, fmt.Errorf("control channel not initialized")
	}

	interceptChan := make(chan string, 1)

	go func() {
		for {
			select {
			case msg := <-pi.ChannelCtl.ch:
				interceptChan <- msg
				if callback, ok := cb.(func(string)); ok {
					callback(msg)
				} else {
					gl.LogObjLogger(pi, "error", "Callback is not of type func(string)")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return <-interceptChan, nil
}

func (pi *ProcessInput[T]) SaveToFile(filename, format string) error {
	data, err := pi.Serialize(format)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (pi *ProcessInput[T]) LoadFromFile(filename, format string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return pi.Deserialize(data, format)
}
