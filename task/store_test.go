package task

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestMemoryStore_Save(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	task := &Task{
		ID:        "task-1",
		State:     StatePending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := store.Save(ctx, task)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify task was saved
	loaded, err := store.Load(ctx, "task-1")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.ID != task.ID {
		t.Errorf("loaded.ID = %v, want %v", loaded.ID, task.ID)
	}
}

func TestMemoryStore_Save_Update(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	task := &Task{
		ID:        "task-1",
		State:     StatePending,
		Message:   "initial",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := store.Save(ctx, task)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Update the task
	task.State = StateRunning
	task.Message = "updated"
	err = store.Save(ctx, task)
	if err != nil {
		t.Fatalf("Save() update error = %v", err)
	}

	// Verify update
	loaded, err := store.Load(ctx, "task-1")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.State != StateRunning {
		t.Errorf("loaded.State = %v, want %v", loaded.State, StateRunning)
	}
	if loaded.Message != "updated" {
		t.Errorf("loaded.Message = %v, want updated", loaded.Message)
	}
}

func TestMemoryStore_Load(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	now := time.Now()
	task := &Task{
		ID:        "task-1",
		State:     StateRunning,
		Progress:  0.5,
		Message:   "Processing",
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := store.Save(ctx, task)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := store.Load(ctx, "task-1")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.ID != "task-1" {
		t.Errorf("ID = %v, want task-1", loaded.ID)
	}
	if loaded.State != StateRunning {
		t.Errorf("State = %v, want %v", loaded.State, StateRunning)
	}
	if loaded.Progress != 0.5 {
		t.Errorf("Progress = %v, want 0.5", loaded.Progress)
	}
	if loaded.Message != "Processing" {
		t.Errorf("Message = %v, want Processing", loaded.Message)
	}
}

func TestMemoryStore_Load_NotFound(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	_, err := store.Load(ctx, "nonexistent")
	if err == nil {
		t.Fatal("Load() expected error for nonexistent task")
	}
}

func TestMemoryStore_LoadAll(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	now := time.Now()
	tasks := []*Task{
		{ID: "task-1", State: StatePending, CreatedAt: now, UpdatedAt: now},
		{ID: "task-2", State: StateRunning, CreatedAt: now, UpdatedAt: now},
		{ID: "task-3", State: StateComplete, CreatedAt: now, UpdatedAt: now},
	}

	for _, task := range tasks {
		if err := store.Save(ctx, task); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	loaded, err := store.LoadAll(ctx)
	if err != nil {
		t.Fatalf("LoadAll() error = %v", err)
	}

	if len(loaded) != 3 {
		t.Errorf("len(loaded) = %v, want 3", len(loaded))
	}

	// Verify all tasks are present
	ids := make(map[string]bool)
	for _, task := range loaded {
		ids[task.ID] = true
	}
	for _, task := range tasks {
		if !ids[task.ID] {
			t.Errorf("task %s not found in loaded tasks", task.ID)
		}
	}
}

func TestMemoryStore_LoadAll_Empty(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	loaded, err := store.LoadAll(ctx)
	if err != nil {
		t.Fatalf("LoadAll() error = %v", err)
	}

	if len(loaded) != 0 {
		t.Errorf("len(loaded) = %v, want 0", len(loaded))
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	task := &Task{
		ID:        "task-1",
		State:     StatePending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := store.Save(ctx, task); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify task exists
	_, err := store.Load(ctx, "task-1")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Delete task
	err = store.Delete(ctx, "task-1")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify task is gone
	_, err = store.Load(ctx, "task-1")
	if err == nil {
		t.Fatal("Load() expected error after delete")
	}
}

func TestMemoryStore_Delete_NotFound(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	err := store.Delete(ctx, "nonexistent")
	if err == nil {
		t.Fatal("Delete() expected error for nonexistent task")
	}
}

func TestMemoryStore_ConcurrentSafety(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Create initial task
	task := &Task{
		ID:        "task-1",
		State:     StatePending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := store.Save(ctx, task); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 100)

	// Concurrent saves
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			t := &Task{
				ID:        "task-1",
				State:     StateRunning,
				Progress:  float64(n) / 100,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := store.Save(ctx, t); err != nil {
				errChan <- err
			}
		}(i)
	}

	// Concurrent loads
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := store.Load(ctx, "task-1")
			if err != nil {
				errChan <- err
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Errorf("concurrent operation error: %v", err)
	}
}

func TestMemoryStore_Load_ReturnsClone(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	task := &Task{
		ID:        "task-1",
		State:     StatePending,
		Message:   "original",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := store.Save(ctx, task); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load and modify
	loaded1, _ := store.Load(ctx, "task-1")
	loaded1.Message = "modified"

	// Load again and verify original is unchanged
	loaded2, _ := store.Load(ctx, "task-1")
	if loaded2.Message != "original" {
		t.Errorf("modifying loaded task affected stored task: got %v, want original", loaded2.Message)
	}
}

func TestMemoryStore_ImplementsStore(t *testing.T) {
	var _ Store = (*MemoryStore)(nil)
}
