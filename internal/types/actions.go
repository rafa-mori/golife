package types

import (
	"time"
)

// IActionManager defines the interface for managing actions and their associated channels.
type IActionManager[T any] interface {
	// GetResults retrieves a map of results associated with the action.
	GetResults() map[string]IResult

	// GetErrorChannel returns the channel used to communicate errors.
	GetErrorChannel() chan error

	// GetDoneChannel returns the channel used to signal the completion of the action.
	GetDoneChannel() chan any

	// GetCancelChannel returns the channel used to signal the cancellation of the action.
	GetCancelChannel() chan any

	// GetResultsChannel returns the channel used to communicate results.
	GetResultsChannel() chan T
}

// IActionBase defines the base interface for an action, including its metadata and status.
type IActionBase[T any] interface {
	// GetID retrieves the unique identifier of the action.
	GetID() string

	// GetType retrieves the type of the action.
	GetType() string

	// GetStatus retrieves the current status of the action.
	GetStatus() string

	// GetErrors retrieves a list of errors associated with the action.
	GetErrors() []error
}

// IAction defines the interface for an executable action, combining base and manager functionalities.
type IAction[T any] interface {
	IActionBase[T]
	IActionManager[T]

	// IsRunning checks if the action is currently running.
	IsRunning() bool

	// CanExecute checks if the action can be executed.
	CanExecute() bool

	// Execute performs the action and returns an error if it fails.
	Execute() error

	// Cancel cancels the action and returns an error if it fails.
	Cancel() error
}

// IJob defines the interface for a job, which is a specialized action with additional time tracking.
type IJob[T any] interface {
	IAction[T]

	// GetAction retrieves the associated action of the job.
	GetAction() IAction[T]

	// GetResults retrieves a map of results associated with the job.
	GetResults() map[string]IResult

	// GetErrorChannel returns the channel used to communicate errors for the job.
	GetErrorChannel() chan error

	// GetDoneChannel returns the channel used to signal the completion of the job.
	GetDoneChannel() chan any

	// GetCancelChannel returns the channel used to signal the cancellation of the job.
	GetCancelChannel() chan any

	// GetResultsChannel returns the channel used to communicate results for the job.
	GetResultsChannel() chan T

	// GetCreateTime retrieves the creation time of the job.
	GetCreateTime() time.Time

	// GetFinishTime retrieves the finish time of the job.
	GetFinishTime() time.Time

	// GetCancelTime retrieves the cancellation time of the job.
	GetCancelTime() time.Time

	// SetFinishTime sets the finish time of the job.
	SetFinishTime(t time.Time)

	// SetCancelTime sets the cancellation time of the job.
	SetCancelTime(t time.Time)
}
