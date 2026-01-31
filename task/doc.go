// Package task provides task lifecycle management for long-running
// operations across all protocols.
//
// This package enables tracking of task state, progress, and results
// regardless of the underlying protocol (MCP, A2A, ACP).
//
// # State Machine
//
// Tasks follow a state machine with the following states:
//
//   - pending: Task created but not yet started
//   - running: Task is actively executing
//   - complete: Task finished successfully
//   - failed: Task finished with an error
//   - cancelled: Task was cancelled
//
// Valid transitions:
//
//	pending → running (on first Update or explicit Start)
//	running → complete (on Complete)
//	running → failed (on Fail)
//	pending|running → cancelled (on Cancel)
//
// Terminal states (complete, failed, cancelled) allow no further transitions.
//
// # Manager Interface
//
// The Manager interface provides task lifecycle operations:
//
//	type Manager interface {
//	    Create(ctx context.Context, id string) (*Task, error)
//	    Get(ctx context.Context, id string) (*Task, error)
//	    List(ctx context.Context) ([]*Task, error)
//	    Update(ctx context.Context, id string, progress float64, message string) error
//	    Complete(ctx context.Context, id string, result any) error
//	    Fail(ctx context.Context, id string, err error) error
//	    Cancel(ctx context.Context, id string) error
//	    Subscribe(ctx context.Context, id string) (<-chan *Task, error)
//	}
//
// # Usage
//
//	// Create a manager
//	mgr := task.NewManager()
//
//	// Create a task
//	t, err := mgr.Create(ctx, "task-1")
//
//	// Update progress
//	mgr.Update(ctx, "task-1", 0.5, "Processing...")
//
//	// Complete the task
//	mgr.Complete(ctx, "task-1", result)
//
//	// Subscribe to updates
//	ch, err := mgr.Subscribe(ctx, "task-1")
//	for update := range ch {
//	    fmt.Printf("Progress: %.0f%%\n", update.Progress*100)
//	}
package task
