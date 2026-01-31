package task

import (
	"context"
	"testing"
)

func TestWithStore(t *testing.T) {
	customStore := NewMemoryStore()
	mgr := NewManager(WithStore(customStore))

	ctx := context.Background()

	// Create a task through the manager
	_, err := mgr.Create(ctx, "task-1")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify it's stored in the custom store
	task, err := customStore.Load(ctx, "task-1")
	if err != nil {
		t.Fatalf("customStore.Load() error = %v", err)
	}
	if task.ID != "task-1" {
		t.Errorf("task.ID = %v, want task-1", task.ID)
	}
}

func TestNewManager_Defaults(t *testing.T) {
	mgr := NewManager()

	// Should have a default store
	if mgr.store == nil {
		t.Error("manager should have default store")
	}

	// Should have initialized subscribers map
	if mgr.subscribers == nil {
		t.Error("manager should have initialized subscribers map")
	}

	// Should be usable
	ctx := context.Background()
	_, err := mgr.Create(ctx, "task-1")
	if err != nil {
		t.Fatalf("Create() with defaults error = %v", err)
	}
}

func TestNewManager_MultipleOptions(t *testing.T) {
	store1 := NewMemoryStore()
	store2 := NewMemoryStore()

	// Last option should win
	mgr := NewManager(
		WithStore(store1),
		WithStore(store2),
	)

	ctx := context.Background()
	_, _ = mgr.Create(ctx, "task-1")

	// Should be in store2, not store1
	_, err := store2.Load(ctx, "task-1")
	if err != nil {
		t.Error("task should be in store2 (last option)")
	}

	_, err = store1.Load(ctx, "task-1")
	if err == nil {
		t.Error("task should not be in store1")
	}
}

func TestOption_FunctionalPattern(t *testing.T) {
	// Verify Option is a function type
	var opt Option = func(m *DefaultManager) {
		// Custom configuration
	}

	// Should be applicable
	mgr := NewManager(opt)
	if mgr == nil {
		t.Error("manager should not be nil")
	}
}
