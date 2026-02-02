// Package task provides task lifecycle management for long-running
// operations across all protocols.
//
// This package enables tracking of task state, progress, and results
// regardless of the underlying protocol (MCP, A2A, ACP).
//
// # Ecosystem Position
//
// task manages long-running operations with progress tracking:
//
//	┌─────────────────────────────────────────────────────────────────┐
//	│                      Task Management Flow                       │
//	├─────────────────────────────────────────────────────────────────┤
//	│                                                                 │
//	│   Protocol              task                  Subscriber        │
//	│   ┌─────────┐        ┌───────────┐         ┌─────────┐         │
//	│   │ Request │────────│  Create   │─────────│ Handler │         │
//	│   │         │        │           │         │         │         │
//	│   └─────────┘        │ ┌───────┐ │         └─────────┘         │
//	│        │             │ │Manager│ │              │               │
//	│        │             │ │       │ │              │               │
//	│        ▼             │ └───────┘ │              │               │
//	│   ┌─────────┐        │     │     │              │               │
//	│   │ Update  │────────│─────┼─────│──────────────▼               │
//	│   │Progress │        │     │     │         ┌─────────┐         │
//	│   └─────────┘        │     ▼     │         │Subscribe│         │
//	│        │             │ ┌───────┐ │         │ Channel │         │
//	│        ▼             │ │ Store │ │         └─────────┘         │
//	│   ┌─────────┐        │ │(memory)│ │              ▲               │
//	│   │Complete │────────│ └───────┘ │──────────────┘               │
//	│   │/Fail    │        └───────────┘   notify                     │
//	│   └─────────┘                                                   │
//	│                                                                 │
//	└─────────────────────────────────────────────────────────────────┘
//
// # State Machine
//
// Tasks follow a state machine with the following states and transitions:
//
//	                    ┌─────────────────────────────────────┐
//	                    │          State Machine              │
//	                    └─────────────────────────────────────┘
//
//	     ┌─────────┐      Update       ┌─────────┐
//	     │ pending │─────────────────▶│ running │
//	     └────┬────┘                   └────┬────┘
//	          │                              │
//	          │ Cancel                       │ Complete/Fail/Cancel
//	          │                              │
//	          ▼                              ▼
//	    ┌───────────┐                  ┌───────────┐
//	    │ cancelled │                  │ complete  │
//	    └───────────┘                  │  failed   │
//	                                   │ cancelled │
//	                                   └───────────┘
//	                                    (terminal)
//
// Terminal states (complete, failed, cancelled) allow no further transitions.
//
// # Core Components
//
//   - [Task]: Task with ID, state, progress, message, result, and timestamps
//   - [State]: Task state constants (pending, running, complete, failed, cancelled)
//   - [Manager]: Interface for task lifecycle (Create/Get/Update/Complete/Fail/Cancel)
//   - [DefaultManager]: Thread-safe Manager implementation with subscription support
//   - [Store]: Interface for task persistence
//   - [MemoryStore]: Thread-safe in-memory Store implementation
//
// # Quick Start
//
//	// Create a manager
//	mgr := task.NewManager()
//
//	// Create a task
//	t, err := mgr.Create(ctx, "task-1")
//	if err != nil {
//	    return err
//	}
//
//	// Update progress (transitions pending → running)
//	err = mgr.Update(ctx, "task-1", 0.5, "Processing...")
//
//	// Subscribe to updates
//	ch, _ := mgr.Subscribe(ctx, "task-1")
//	go func() {
//	    for update := range ch {
//	        fmt.Printf("Progress: %.0f%%\n", update.Progress*100)
//	    }
//	}()
//
//	// Complete the task (closes subscription channel)
//	err = mgr.Complete(ctx, "task-1", result)
//
// # Subscriptions
//
// Subscribe returns a channel that receives task updates:
//
//	ch, err := mgr.Subscribe(ctx, "task-1")
//	if err != nil {
//	    return err
//	}
//
//	for update := range ch {
//	    // Process update
//	    if update.State.IsTerminal() {
//	        break
//	    }
//	}
//
// The channel is closed when:
//   - Task reaches terminal state (complete, failed, cancelled)
//   - Context is cancelled
//
// # Thread Safety
//
// All exported types are safe for concurrent use:
//
//   - [DefaultManager]: sync.RWMutex protects all operations
//     - Get/List: Uses RLock for concurrent reads
//     - Create/Update/Complete/Fail/Cancel: Uses Lock for exclusive access
//   - [MemoryStore]: sync.RWMutex protects all operations
//   - [Task]: Not thread-safe; use Manager methods for safe mutations
//   - Subscription channels: Buffered (10) for non-blocking sends
//
// # Error Handling
//
// Sentinel errors (use errors.Is for checking):
//
//   - [ErrTaskNotFound]: Task does not exist
//   - [ErrTaskExists]: Task with ID already exists
//   - [ErrInvalidState]: Invalid task state
//   - [ErrInvalidTransition]: Cannot transition from current state
//   - [ErrEmptyID]: Task ID is empty
//
// The [TaskError] type wraps errors with task context:
//
//	err := &TaskError{
//	    TaskID: "task-123",
//	    Op:     "update",
//	    Err:    ErrTaskNotFound,
//	}
//	// errors.Is(err, ErrTaskNotFound) = true
//
// # Configuration Options
//
// DefaultManager supports functional options:
//
//   - [WithStore]: Use custom Store implementation
//
// # Integration with ApertureStack
//
// task integrates with other ApertureStack packages:
//
//   - session: Tasks may be associated with client sessions
//   - stream: Progress updates may be streamed to clients
//   - wire: Task status maps to protocol-specific formats (MCP progress, A2A status)
package task
