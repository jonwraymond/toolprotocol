package task

import (
	"context"
	"fmt"
	"testing"
)

// BenchmarkManager_Create measures task creation performance.
func BenchmarkManager_Create(b *testing.B) {
	mgr := NewManager()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mgr.Create(ctx, fmt.Sprintf("task-%d", i))
	}
}

// BenchmarkManager_Get measures task retrieval performance.
func BenchmarkManager_Get(b *testing.B) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mgr.Get(ctx, "task-1")
	}
}

// BenchmarkManager_Get_Miss measures miss performance.
func BenchmarkManager_Get_Miss(b *testing.B) {
	mgr := NewManager()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mgr.Get(ctx, "nonexistent")
	}
}

// BenchmarkManager_List measures list performance.
func BenchmarkManager_List(b *testing.B) {
	mgr := NewManager()
	ctx := context.Background()

	// Create some tasks
	for i := 0; i < 100; i++ {
		_, _ = mgr.Create(ctx, fmt.Sprintf("task-%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mgr.List(ctx)
	}
}

// BenchmarkManager_Update measures update performance.
func BenchmarkManager_Update(b *testing.B) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mgr.Update(ctx, "task-1", float64(i%100)/100.0, "Processing...")
	}
}

// BenchmarkManager_Complete measures completion performance.
func BenchmarkManager_Complete(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		mgr := NewManager()
		_, _ = mgr.Create(ctx, "task-1")
		b.StartTimer()

		_ = mgr.Complete(ctx, "task-1", "result")
	}
}

// BenchmarkManager_Fail measures fail performance.
func BenchmarkManager_Fail(b *testing.B) {
	ctx := context.Background()
	testErr := fmt.Errorf("test error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		mgr := NewManager()
		_, _ = mgr.Create(ctx, "task-1")
		b.StartTimer()

		_ = mgr.Fail(ctx, "task-1", testErr)
	}
}

// BenchmarkManager_Cancel measures cancellation performance.
func BenchmarkManager_Cancel(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		mgr := NewManager()
		_, _ = mgr.Create(ctx, "task-1")
		b.StartTimer()

		_ = mgr.Cancel(ctx, "task-1")
	}
}

// BenchmarkManager_Subscribe measures subscription setup performance.
func BenchmarkManager_Subscribe(b *testing.B) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch, _ := mgr.Subscribe(ctx, "task-1")
		// Don't leave hanging goroutines
		go func() {
			for range ch {
			}
		}()
	}
}

// BenchmarkManager_Concurrent measures concurrent operations.
func BenchmarkManager_Concurrent(b *testing.B) {
	mgr := NewManager()
	ctx := context.Background()

	// Pre-create some tasks
	for i := 0; i < 100; i++ {
		_, _ = mgr.Create(ctx, fmt.Sprintf("task-%d", i))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			switch i % 5 {
			case 0:
				_, _ = mgr.Create(ctx, fmt.Sprintf("new-task-%d", i))
			case 1:
				_, _ = mgr.Get(ctx, fmt.Sprintf("task-%d", i%100))
			case 2:
				_ = mgr.Update(ctx, fmt.Sprintf("task-%d", i%100), 0.5, "progress")
			case 3:
				_, _ = mgr.List(ctx)
			case 4:
				_ = mgr.Cancel(ctx, fmt.Sprintf("task-%d", i%100))
			}
			i++
		}
	})
}

// BenchmarkManager_ConcurrentReadHeavy measures read-heavy workload.
func BenchmarkManager_ConcurrentReadHeavy(b *testing.B) {
	mgr := NewManager()
	ctx := context.Background()

	// Pre-create tasks
	for i := 0; i < 100; i++ {
		_, _ = mgr.Create(ctx, fmt.Sprintf("task-%d", i))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = mgr.Get(ctx, fmt.Sprintf("task-%d", i%100))
			i++
		}
	})
}

// BenchmarkTask_Clone measures clone performance.
func BenchmarkTask_Clone(b *testing.B) {
	task := &Task{
		ID:       "task-123",
		State:    StateRunning,
		Progress: 0.5,
		Message:  "Processing...",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = task.Clone()
	}
}

// BenchmarkState_IsTerminal measures terminal check performance.
func BenchmarkState_IsTerminal(b *testing.B) {
	states := []State{StatePending, StateRunning, StateComplete, StateFailed, StateCancelled}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = states[i%5].IsTerminal()
	}
}

// BenchmarkState_Valid measures validation performance.
func BenchmarkState_Valid(b *testing.B) {
	state := StateRunning

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = state.Valid()
	}
}

// BenchmarkState_String measures string conversion performance.
func BenchmarkState_String(b *testing.B) {
	state := StateRunning

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = state.String()
	}
}

// BenchmarkMemoryStore_Save measures store save performance.
func BenchmarkMemoryStore_Save(b *testing.B) {
	store := NewMemoryStore()
	ctx := context.Background()
	task := &Task{ID: "task-1", State: StatePending}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.Save(ctx, task)
	}
}

// BenchmarkMemoryStore_Load measures store load performance.
func BenchmarkMemoryStore_Load(b *testing.B) {
	store := NewMemoryStore()
	ctx := context.Background()
	_ = store.Save(ctx, &Task{ID: "task-1", State: StatePending})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Load(ctx, "task-1")
	}
}

// BenchmarkMemoryStore_LoadAll measures store load all performance.
func BenchmarkMemoryStore_LoadAll(b *testing.B) {
	store := NewMemoryStore()
	ctx := context.Background()

	for i := 0; i < 100; i++ {
		_ = store.Save(ctx, &Task{ID: fmt.Sprintf("task-%d", i), State: StatePending})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.LoadAll(ctx)
	}
}

// BenchmarkMemoryStore_Delete measures store delete performance.
func BenchmarkMemoryStore_Delete(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		store := NewMemoryStore()
		_ = store.Save(ctx, &Task{ID: "task-1", State: StatePending})
		b.StartTimer()

		_ = store.Delete(ctx, "task-1")
	}
}

// BenchmarkMemoryStore_Concurrent measures concurrent store access.
func BenchmarkMemoryStore_Concurrent(b *testing.B) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Pre-populate
	for i := 0; i < 100; i++ {
		_ = store.Save(ctx, &Task{ID: fmt.Sprintf("task-%d", i), State: StatePending})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%4 == 0 {
				_ = store.Save(ctx, &Task{ID: fmt.Sprintf("task-%d", i%100), State: StateRunning})
			} else {
				_, _ = store.Load(ctx, fmt.Sprintf("task-%d", i%100))
			}
			i++
		}
	})
}
