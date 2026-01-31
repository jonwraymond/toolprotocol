package task

import (
	"context"
	"time"
)

// State represents task execution state.
type State string

const (
	// StatePending indicates the task is created but not yet started.
	StatePending State = "pending"

	// StateRunning indicates the task is actively executing.
	StateRunning State = "running"

	// StateComplete indicates the task finished successfully.
	StateComplete State = "complete"

	// StateFailed indicates the task finished with an error.
	StateFailed State = "failed"

	// StateCancelled indicates the task was cancelled.
	StateCancelled State = "cancelled"
)

// String returns the string representation of the state.
func (s State) String() string {
	return string(s)
}

// Valid returns true if the state is a known valid state.
func (s State) Valid() bool {
	switch s {
	case StatePending, StateRunning, StateComplete, StateFailed, StateCancelled:
		return true
	default:
		return false
	}
}

// IsTerminal returns true if the state is a terminal state
// (complete, failed, or cancelled) that allows no further transitions.
func (s State) IsTerminal() bool {
	switch s {
	case StateComplete, StateFailed, StateCancelled:
		return true
	default:
		return false
	}
}

// Task represents a long-running operation.
type Task struct {
	// ID is the unique identifier for the task.
	ID string

	// State is the current execution state.
	State State

	// Progress is the completion percentage (0.0 to 1.0).
	Progress float64

	// Message is the current status message.
	Message string

	// Result is the final result (set when complete).
	Result any

	// Error is the error (set when failed).
	Error error

	// CreatedAt is when the task was created.
	CreatedAt time.Time

	// UpdatedAt is when the task was last updated.
	UpdatedAt time.Time

	// CompletedAt is when the task reached a terminal state.
	CompletedAt *time.Time
}

// Clone returns a deep copy of the task.
func (t *Task) Clone() *Task {
	clone := &Task{
		ID:        t.ID,
		State:     t.State,
		Progress:  t.Progress,
		Message:   t.Message,
		Result:    t.Result,
		Error:     t.Error,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	if t.CompletedAt != nil {
		completedAt := *t.CompletedAt
		clone.CompletedAt = &completedAt
	}
	return clone
}

// Manager manages task lifecycle.
type Manager interface {
	// Create creates a new task with the given ID.
	Create(ctx context.Context, id string) (*Task, error)

	// Get retrieves a task by ID.
	Get(ctx context.Context, id string) (*Task, error)

	// List returns all tasks.
	List(ctx context.Context) ([]*Task, error)

	// Update updates task progress and message.
	// Transitions pending → running on first update.
	Update(ctx context.Context, id string, progress float64, message string) error

	// Complete marks the task as complete with the given result.
	Complete(ctx context.Context, id string, result any) error

	// Fail marks the task as failed with the given error.
	Fail(ctx context.Context, id string, err error) error

	// Cancel cancels the task.
	Cancel(ctx context.Context, id string) error

	// Subscribe returns a channel that receives task updates.
	// The channel is closed when the task reaches a terminal state.
	Subscribe(ctx context.Context, id string) (<-chan *Task, error)
}

// Store provides task persistence.
type Store interface {
	// Save saves or updates a task.
	Save(ctx context.Context, task *Task) error

	// Load retrieves a task by ID.
	Load(ctx context.Context, id string) (*Task, error)

	// LoadAll retrieves all tasks.
	LoadAll(ctx context.Context) ([]*Task, error)

	// Delete removes a task by ID.
	Delete(ctx context.Context, id string) error
}

