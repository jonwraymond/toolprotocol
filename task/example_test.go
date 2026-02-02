package task_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonwraymond/toolprotocol/task"
)

func ExampleNewManager() {
	mgr := task.NewManager()

	ctx := context.Background()
	t, err := mgr.Create(ctx, "task-1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Task created:", t.ID)
	fmt.Println("Initial state:", t.State)
	// Output:
	// Task created: task-1
	// Initial state: pending
}

func ExampleNewManager_withStore() {
	// Use custom store
	store := task.NewMemoryStore()
	mgr := task.NewManager(task.WithStore(store))

	ctx := context.Background()
	t, _ := mgr.Create(ctx, "task-1")

	fmt.Println("Task ID:", t.ID)
	// Output:
	// Task ID: task-1
}

func ExampleDefaultManager_Create() {
	mgr := task.NewManager()
	ctx := context.Background()

	t, err := mgr.Create(ctx, "my-task")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("ID:", t.ID)
	fmt.Println("State:", t.State)
	fmt.Println("Progress:", t.Progress)
	// Output:
	// ID: my-task
	// State: pending
	// Progress: 0
}

func ExampleDefaultManager_Create_duplicate() {
	mgr := task.NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")
	_, err := mgr.Create(ctx, "task-1")

	fmt.Println("Error is ErrTaskExists:", errors.Is(err, task.ErrTaskExists))
	// Output:
	// Error is ErrTaskExists: true
}

func ExampleDefaultManager_Create_emptyID() {
	mgr := task.NewManager()
	ctx := context.Background()

	_, err := mgr.Create(ctx, "")
	fmt.Println("Error is ErrEmptyID:", errors.Is(err, task.ErrEmptyID))
	// Output:
	// Error is ErrEmptyID: true
}

func ExampleDefaultManager_Get() {
	mgr := task.NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	t, err := mgr.Get(ctx, "task-1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("ID:", t.ID)
	fmt.Println("State:", t.State)
	// Output:
	// ID: task-1
	// State: pending
}

func ExampleDefaultManager_Get_notFound() {
	mgr := task.NewManager()
	ctx := context.Background()

	_, err := mgr.Get(ctx, "nonexistent")
	fmt.Println("Error is ErrTaskNotFound:", errors.Is(err, task.ErrTaskNotFound))
	// Output:
	// Error is ErrTaskNotFound: true
}

func ExampleDefaultManager_List() {
	mgr := task.NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")
	_, _ = mgr.Create(ctx, "task-2")
	_, _ = mgr.Create(ctx, "task-3")

	tasks, err := mgr.List(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Task count:", len(tasks))
	// Output:
	// Task count: 3
}

func ExampleDefaultManager_Update() {
	mgr := task.NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	// Update progress
	err := mgr.Update(ctx, "task-1", 0.5, "Processing...")
	fmt.Println("Update error:", err)

	// Check state changed to running
	t, _ := mgr.Get(ctx, "task-1")
	fmt.Println("State after update:", t.State)
	fmt.Println("Progress:", t.Progress)
	fmt.Println("Message:", t.Message)
	// Output:
	// Update error: <nil>
	// State after update: running
	// Progress: 0.5
	// Message: Processing...
}

func ExampleDefaultManager_Update_notFound() {
	mgr := task.NewManager()
	ctx := context.Background()

	err := mgr.Update(ctx, "nonexistent", 0.5, "...")
	fmt.Println("Error is ErrTaskNotFound:", errors.Is(err, task.ErrTaskNotFound))
	// Output:
	// Error is ErrTaskNotFound: true
}

func ExampleDefaultManager_Complete() {
	mgr := task.NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")
	_ = mgr.Update(ctx, "task-1", 0.5, "Processing")

	// Complete the task
	result := map[string]string{"status": "done"}
	err := mgr.Complete(ctx, "task-1", result)
	fmt.Println("Complete error:", err)

	// Check state
	t, _ := mgr.Get(ctx, "task-1")
	fmt.Println("State:", t.State)
	fmt.Println("Progress:", t.Progress)
	fmt.Println("Is terminal:", t.State.IsTerminal())
	// Output:
	// Complete error: <nil>
	// State: complete
	// Progress: 1
	// Is terminal: true
}

func ExampleDefaultManager_Fail() {
	mgr := task.NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")
	_ = mgr.Update(ctx, "task-1", 0.5, "Processing")

	// Fail the task
	err := mgr.Fail(ctx, "task-1", errors.New("something went wrong"))
	fmt.Println("Fail error:", err)

	// Check state
	t, _ := mgr.Get(ctx, "task-1")
	fmt.Println("State:", t.State)
	fmt.Println("Task error:", t.Error)
	fmt.Println("Is terminal:", t.State.IsTerminal())
	// Output:
	// Fail error: <nil>
	// State: failed
	// Task error: something went wrong
	// Is terminal: true
}

func ExampleDefaultManager_Cancel() {
	mgr := task.NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	// Cancel the task
	err := mgr.Cancel(ctx, "task-1")
	fmt.Println("Cancel error:", err)

	// Check state
	t, _ := mgr.Get(ctx, "task-1")
	fmt.Println("State:", t.State)
	fmt.Println("Is terminal:", t.State.IsTerminal())
	// Output:
	// Cancel error: <nil>
	// State: cancelled
	// Is terminal: true
}

func ExampleDefaultManager_Cancel_terminal() {
	mgr := task.NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")
	_ = mgr.Complete(ctx, "task-1", nil) // Already complete

	// Try to cancel terminal task
	err := mgr.Cancel(ctx, "task-1")
	fmt.Println("Error is ErrInvalidTransition:", errors.Is(err, task.ErrInvalidTransition))
	// Output:
	// Error is ErrInvalidTransition: true
}

func ExampleDefaultManager_Subscribe() {
	mgr := task.NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	// Subscribe to updates
	ch, err := mgr.Subscribe(ctx, "task-1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Subscribed:", ch != nil)

	// Complete task (sends update then closes channel)
	_ = mgr.Complete(ctx, "task-1", "done")

	// Drain the channel to get the final update
	count := 0
	for range ch {
		count++
	}
	fmt.Println("Received updates before close:", count > 0)
	// Output:
	// Subscribed: true
	// Received updates before close: true
}

func ExampleTask_Clone() {
	t := &task.Task{
		ID:       "task-1",
		State:    task.StateRunning,
		Progress: 0.5,
		Message:  "Processing",
	}

	clone := t.Clone()

	fmt.Println("Same ID:", clone.ID == t.ID)
	fmt.Println("Same State:", clone.State == t.State)
	fmt.Println("Same Progress:", clone.Progress == t.Progress)
	// Output:
	// Same ID: true
	// Same State: true
	// Same Progress: true
}

func ExampleState_String() {
	states := []task.State{
		task.StatePending,
		task.StateRunning,
		task.StateComplete,
		task.StateFailed,
		task.StateCancelled,
	}

	for _, s := range states {
		fmt.Println(s.String())
	}
	// Output:
	// pending
	// running
	// complete
	// failed
	// cancelled
}

func ExampleState_Valid() {
	fmt.Println("pending valid:", task.StatePending.Valid())
	fmt.Println("running valid:", task.StateRunning.Valid())
	fmt.Println("invalid valid:", task.State("unknown").Valid())
	// Output:
	// pending valid: true
	// running valid: true
	// invalid valid: false
}

func ExampleState_IsTerminal() {
	fmt.Println("pending terminal:", task.StatePending.IsTerminal())
	fmt.Println("running terminal:", task.StateRunning.IsTerminal())
	fmt.Println("complete terminal:", task.StateComplete.IsTerminal())
	fmt.Println("failed terminal:", task.StateFailed.IsTerminal())
	fmt.Println("cancelled terminal:", task.StateCancelled.IsTerminal())
	// Output:
	// pending terminal: false
	// running terminal: false
	// complete terminal: true
	// failed terminal: true
	// cancelled terminal: true
}

func ExampleNewMemoryStore() {
	store := task.NewMemoryStore()
	ctx := context.Background()

	// Save a task
	t := &task.Task{ID: "task-1", State: task.StatePending}
	err := store.Save(ctx, t)
	fmt.Println("Save error:", err)

	// Load it back
	loaded, err := store.Load(ctx, "task-1")
	if err != nil {
		fmt.Println("Load error:", err)
		return
	}

	fmt.Println("Loaded ID:", loaded.ID)
	// Output:
	// Save error: <nil>
	// Loaded ID: task-1
}

func ExampleMemoryStore_Save() {
	store := task.NewMemoryStore()
	ctx := context.Background()

	t := &task.Task{ID: "task-1", State: task.StatePending}
	err := store.Save(ctx, t)
	fmt.Println("Save error:", err)
	// Output:
	// Save error: <nil>
}

func ExampleMemoryStore_Load() {
	store := task.NewMemoryStore()
	ctx := context.Background()

	// Save first
	_ = store.Save(ctx, &task.Task{ID: "task-1", State: task.StatePending})

	// Load
	t, err := store.Load(ctx, "task-1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("ID:", t.ID)
	fmt.Println("State:", t.State)
	// Output:
	// ID: task-1
	// State: pending
}

func ExampleMemoryStore_Load_notFound() {
	store := task.NewMemoryStore()
	ctx := context.Background()

	_, err := store.Load(ctx, "nonexistent")
	fmt.Println("Error is ErrTaskNotFound:", errors.Is(err, task.ErrTaskNotFound))
	// Output:
	// Error is ErrTaskNotFound: true
}

func ExampleMemoryStore_LoadAll() {
	store := task.NewMemoryStore()
	ctx := context.Background()

	_ = store.Save(ctx, &task.Task{ID: "task-1"})
	_ = store.Save(ctx, &task.Task{ID: "task-2"})

	tasks, err := store.LoadAll(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Task count:", len(tasks))
	// Output:
	// Task count: 2
}

func ExampleMemoryStore_Delete() {
	store := task.NewMemoryStore()
	ctx := context.Background()

	_ = store.Save(ctx, &task.Task{ID: "task-1"})

	err := store.Delete(ctx, "task-1")
	fmt.Println("Delete error:", err)

	// Verify deleted
	_, err = store.Load(ctx, "task-1")
	fmt.Println("After delete, error is ErrTaskNotFound:", errors.Is(err, task.ErrTaskNotFound))
	// Output:
	// Delete error: <nil>
	// After delete, error is ErrTaskNotFound: true
}

func ExampleWithStore() {
	// Create custom store
	store := task.NewMemoryStore()

	// Use it with manager
	mgr := task.NewManager(task.WithStore(store))

	ctx := context.Background()
	_, _ = mgr.Create(ctx, "task-1")

	// Verify task is in our custom store
	t, _ := store.Load(ctx, "task-1")
	fmt.Println("Task in custom store:", t != nil)
	// Output:
	// Task in custom store: true
}

func ExampleTaskError() {
	err := &task.TaskError{
		TaskID: "task-123",
		Op:     "update",
		Err:    task.ErrTaskNotFound,
	}

	fmt.Println(err.Error())
	fmt.Println("Unwraps to ErrTaskNotFound:", errors.Is(err, task.ErrTaskNotFound))
	// Output:
	// task task-123: update: task: not found
	// Unwraps to ErrTaskNotFound: true
}

func ExampleTaskError_noUnderlying() {
	err := &task.TaskError{
		TaskID: "task-123",
		Op:     "validate",
	}

	fmt.Println(err.Error())
	// Output:
	// task task-123: validate
}
