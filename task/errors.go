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
