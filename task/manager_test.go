package task

import (
	"context"
	"testing"
	"time"
)

func TestManager_Create(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	task, err := mgr.Create(ctx, "task-1")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if task.ID != "task-1" {
		t.Errorf("task.ID = %v, want task-1", task.ID)
	}
	if task.State != StatePending {
		t.Errorf("task.State = %v, want %v", task.State, StatePending)
	}
	if task.Progress != 0 {
		t.Errorf("task.Progress = %v, want 0", task.Progress)
	}
	if task.CreatedAt.IsZero() {
		t.Error("task.CreatedAt should not be zero")
	}
	if task.UpdatedAt.IsZero() {
		t.Error("task.UpdatedAt should not be zero")
	}
}

func TestManager_Create_DuplicateID(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, err := mgr.Create(ctx, "task-1")
	if err != nil {
		t.Fatalf("Create() first call error = %v", err)
	}

	_, err = mgr.Create(ctx, "task-1")
	if err == nil {
		t.Fatal("Create() expected error for duplicate ID")
	}
	if err != ErrTaskExists {
		t.Errorf("Create() error = %v, want %v", err, ErrTaskExists)
	}
}

func TestManager_Create_EmptyID(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, err := mgr.Create(ctx, "")
	if err == nil {
		t.Fatal("Create() expected error for empty ID")
	}
	if err != ErrEmptyID {
		t.Errorf("Create() error = %v, want %v", err, ErrEmptyID)
	}
}

func TestManager_Get(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	created, err := mgr.Create(ctx, "task-1")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := mgr.Get(ctx, "task-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.ID != created.ID {
		t.Errorf("got.ID = %v, want %v", got.ID, created.ID)
	}
	if got.State != created.State {
		t.Errorf("got.State = %v, want %v", got.State, created.State)
	}
}

func TestManager_Get_NotFound(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, err := mgr.Get(ctx, "nonexistent")
	if err == nil {
		t.Fatal("Get() expected error for nonexistent task")
	}
	if err != ErrTaskNotFound {
		t.Errorf("Get() error = %v, want %v", err, ErrTaskNotFound)
	}
}

func TestManager_Get_ReturnsClone(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, err := mgr.Create(ctx, "task-1")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Get and modify
	got1, _ := mgr.Get(ctx, "task-1")
	got1.Message = "modified"

	// Get again and verify original unchanged
	got2, _ := mgr.Get(ctx, "task-1")
	if got2.Message == "modified" {
		t.Error("modifying returned task affected stored task")
	}
}

func TestManager_List(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")
	_, _ = mgr.Create(ctx, "task-2")
	_, _ = mgr.Create(ctx, "task-3")

	tasks, err := mgr.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("len(tasks) = %v, want 3", len(tasks))
	}

	// Verify all tasks present
	ids := make(map[string]bool)
	for _, task := range tasks {
		ids[task.ID] = true
	}
	for _, id := range []string{"task-1", "task-2", "task-3"} {
		if !ids[id] {
			t.Errorf("task %s not found", id)
		}
	}
}

func TestManager_List_Empty(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	tasks, err := mgr.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("len(tasks) = %v, want 0", len(tasks))
	}
}

func TestManager_ContextCancellation(t *testing.T) {
	mgr := NewManager()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Operations should still work with cancelled context
	// (context cancellation is for long-running operations)
	_, err := mgr.Create(ctx, "task-1")
	if err != nil {
		t.Logf("Create() with cancelled context: %v", err)
	}
}

func TestManager_ImplementsInterface(t *testing.T) {
	var _ Manager = (*DefaultManager)(nil)
}

func TestManager_Create_SetsTimestamps(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	before := time.Now()
	task, err := mgr.Create(ctx, "task-1")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	after := time.Now()

	if task.CreatedAt.Before(before) || task.CreatedAt.After(after) {
		t.Errorf("CreatedAt = %v, want between %v and %v", task.CreatedAt, before, after)
	}
	if task.UpdatedAt.Before(before) || task.UpdatedAt.After(after) {
		t.Errorf("UpdatedAt = %v, want between %v and %v", task.UpdatedAt, before, after)
	}
	if task.CompletedAt != nil {
		t.Error("CompletedAt should be nil for new task")
	}
}
