package jobs

import (
	"fmt"
	//a "github.com/faelmori/gastype/internal/routinces/taskz/actions"
	t "github.com/faelmori/golife/internal/types"
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
	"time"
)

type IJob interface {
	t.IAction
	GetID() string
}

// Job represents a job that executes an action and tracks its state and results.
type Job struct {
	logger       l.Logger             // Logger instance for logging job-related information.
	ID           string               // Unique identifier for the job.
	CreateTime   string               // Creation time of the job in RFC3339 format.
	CancelTime   time.Time            // Time when the job was canceled.
	FinishTime   time.Time            // Time when the job was finished.
	Results      map[string]t.IResult // Map of results associated with the job.
	Errors       []error              // List of errors encountered during the job execution.
	Running      bool                 // Indicates whether the job is currently running.
	CancelChanel chan struct{}        // Channel used to signal job cancellation.
	DoneChanel   chan struct{}        // Channel used to signal job completion.
	Action       t.IAction            // Action associated with the job.
}

// NewJob creates a new job.
// Parameters:
//   - action: The action to be executed by the job.
//   - cancelChanel: Channel to signal job cancellation.
//   - doneChanel: Channel to signal job completion.
//   - logger: Logger instance for the job.
//
// Returns:
//   - *Job: A new instance of the Job struct.
func NewJob(action t.IAction, cancelChanel chan struct{}, doneChanel chan struct{}, logger l.Logger) *Job {
	if logger == nil {
		logger = l.GetLogger("GasType")
	}
	if cancelChanel == nil {
		cancelChanel = make(chan struct{}, 2)
	}
	if doneChanel == nil {
		doneChanel = make(chan struct{}, 2)
	}
	return &Job{
		logger:       logger,
		ID:           uuid.NewString(),
		Results:      make(map[string]t.IResult),
		Errors:       make([]error, 0),
		Running:      false,
		CancelChanel: cancelChanel,
		DoneChanel:   doneChanel,
		FinishTime:   time.Time{},
		CancelTime:   time.Time{},
		CreateTime:   time.Now().Format(time.RFC3339),
		Action:       action,
	}
}

// GetID returns the job ID.
// Returns:
//   - string: The unique identifier of the job.
func (jb *Job) GetID() string {
	return jb.ID
}

// GetAction returns the action associated with the job.
// Returns:
//   - a.IAction: The action instance.
func (jb *Job) GetAction() t.IAction {
	return jb.Action
}

// GetType returns the type of the action associated with the job.
// Returns:
//   - string: The type of the action.
func (jb *Job) GetType() string {
	return jb.Action.GetType()
}

// GetResults returns the results of the action associated with the job.
// Returns:
//   - map[string]t.IResult: A map of results.
func (jb *Job) GetResults() map[string]t.IResult {
	act := jb.GetAction()
	mp := make(map[string]t.IResult)
	for k, v := range act.GetResults() {
		mp[k] = v
	}
	return mp
}

// GetStatus returns the current status of the action associated with the job.
// Returns:
//   - string: The status of the action.
func (jb *Job) GetStatus() string {
	return jb.Action.GetStatus()
}

// GetErrors returns the errors encountered during the action execution.
// Returns:
//   - []error: A list of errors.
func (jb *Job) GetErrors() []error {
	return jb.Action.GetErrors()
}

// IsRunning checks if the action associated with the job is currently running.
// Returns:
//   - bool: True if the action is running, false otherwise.
func (jb *Job) IsRunning() bool {
	return jb.Action.IsRunning()
}

// CanExecute checks if the action associated with the job can be executed.
// Returns:
//   - bool: True if the action can be executed, false otherwise.
func (jb *Job) CanExecute() bool {
	return jb.Action.CanExecute()
}

// Execute starts the execution of the action associated with the job.
// Returns:
//   - error: An error if the job is already running or cannot be executed, nil otherwise.
func (jb *Job) Execute() error {
	if jb.Running {
		jb.logger.ErrorCtx("Job is already running", map[string]interface{}{"job_id": jb.ID})
		return nil
	}
	if jb.Action.CanExecute() {
		if err := jb.Action.Execute(); err != nil {
			return err
		}
	} else {
		jb.logger.ErrorCtx("Job cannot be executed", map[string]interface{}{"job_id": jb.ID})
		return nil
	}
	return nil
}

// GetErrorChannel returns the error channel of the action associated with the job.
// Returns:
//   - chan error: The error channel.
func (jb *Job) GetErrorChannel() chan error {
	return jb.Action.GetErrorChannel()
}

// GetDoneChannel returns the done channel of the action associated with the job.
// Returns:
//   - chan struct{}: The done channel.
func (jb *Job) GetDoneChannel() chan struct{} {
	return jb.Action.GetDoneChannel()
}

// GetCancelChannel returns the cancel channel of the action associated with the job.
// Returns:
//   - chan struct{}: The cancel channel.
func (jb *Job) GetCancelChannel() chan struct{} {
	return jb.Action.GetCancelChannel()
}

// GetResultsChannel returns the results channel of the action associated with the job.
// Returns:
//   - chan t.IResult: The results channel.
func (jb *Job) GetResultsChannel() chan t.IResult {
	return jb.Action.GetResultsChannel()
}

// GetCreateTime returns the creation time of the job.
// Returns:
//   - time.Time: The creation time.
func (jb *Job) GetCreateTime() time.Time {
	createTime, err := time.Parse(time.RFC3339, jb.CreateTime)
	if err != nil {
		jb.logger.ErrorCtx("Error parsing create time", map[string]interface{}{"job_id": jb.ID, "error": err})
		return time.Time{}
	}
	return createTime
}

// GetFinishTime returns the finish time of the job.
// Returns:
//   - time.Time: The finish time.
func (jb *Job) GetFinishTime() time.Time {
	return jb.FinishTime
}

// GetCancelTime returns the cancel time of the job.
// Returns:
//   - time.Time: The cancel time.
func (jb *Job) GetCancelTime() time.Time {
	return jb.CancelTime
}

// SetFinishTime sets the finish time of the job.
// Parameters:
//   - t: The finish time to set.
func (jb *Job) SetFinishTime(t time.Time) {
	jb.FinishTime = t
}

// SetCancelTime sets the cancel time of the job.
// Parameters:
//   - t: The cancel time to set.
func (jb *Job) SetCancelTime(t time.Time) {
	jb.CancelTime = t
}

// Cancel cancels the job and updates its state.
// Returns:
//   - error: An error if the job is not running, nil otherwise.
func (jb *Job) Cancel() error {
	if jb.Running {
		ch := jb.GetCancelChannel()
		if ch != nil {
			defer close(ch)
			ch <- struct{}{}
		}
		jb.Running = false
		jb.CancelTime = time.Now()
		jb.logger.InfoCtx("Job cancelled", map[string]interface{}{"job_id": jb.ID})
		return nil
	} else {
		jb.logger.ErrorCtx("Job is not running", map[string]interface{}{"job_id": jb.ID})
		return fmt.Errorf("job is not running")
	}
}
