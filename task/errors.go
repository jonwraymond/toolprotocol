package task

import "errors"

// Sentinel errors for task operations.
var (
	// ErrTaskNotFound is returned when a task cannot be found.
	ErrTaskNotFound = errors.New("task not found")

	// ErrTaskExists is returned when creating a task that already exists.
	ErrTaskExists = errors.New("task already exists")

	// ErrInvalidState is returned for invalid task states.
	ErrInvalidState = errors.New("invalid task state")

	// ErrInvalidTransition is returned for invalid state transitions.
	ErrInvalidTransition = errors.New("invalid state transition")

	// ErrEmptyID is returned when a task ID is empty.
	ErrEmptyID = errors.New("task ID cannot be empty")
)

// TaskError wraps an error with task context.
type TaskError struct {
	// TaskID is the ID of the task that caused the error.
	TaskID string

	// Op is the operation that failed.
	Op string

	// Err is the underlying error.
	Err error
}

// Error returns the error message.
func (e *TaskError) Error() string {
	if e.Err == nil {
		return "task " + e.TaskID + ": " + e.Op
	}
	return "task " + e.TaskID + ": " + e.Op + ": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *TaskError) Unwrap() error {
	return e.Err
}
