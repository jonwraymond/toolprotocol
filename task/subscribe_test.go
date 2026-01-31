package task

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestManager_Subscribe(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	ch, err := mgr.Subscribe(ctx, "task-1")
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	if ch == nil {
		t.Fatal("Subscribe() returned nil channel")
	}
}

func TestManager_Subscribe_NotFound(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, err := mgr.Subscribe(ctx, "nonexistent")
	if err == nil {
		t.Fatal("Subscribe() expected error for nonexistent task")
	}
	if err != ErrTaskNotFound {
		t.Errorf("Subscribe() error = %v, want %v", err, ErrTaskNotFound)
	}
}

func TestManager_Subscribe_ReceivesUpdates(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	ch, err := mgr.Subscribe(ctx, "task-1")
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	// Update the task
	_ = mgr.Update(ctx, "task-1", 0.5, "halfway")

	// Should receive update
	select {
	case update := <-ch:
		if update.Progress != 0.5 {
			t.Errorf("update.Progress = %v, want 0.5", update.Progress)
		}
		if update.Message != "halfway" {
			t.Errorf("update.Message = %v, want 'halfway'", update.Message)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for update")
	}
}

func TestManager_Subscribe_ReceivesComplete(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	ch, err := mgr.Subscribe(ctx, "task-1")
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	// Complete the task
	_ = mgr.Complete(ctx, "task-1", "result")

	// Should receive completion
	select {
	case update := <-ch:
		if update.State != StateComplete {
			t.Errorf("update.State = %v, want %v", update.State, StateComplete)
		}
		if update.Result != "result" {
			t.Errorf("update.Result = %v, want 'result'", update.Result)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for complete notification")
	}

	// Channel should be closed
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("channel should be closed after terminal state")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("channel not closed after terminal state")
	}
}

func TestManager_Subscribe_ReceivesFail(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	ch, err := mgr.Subscribe(ctx, "task-1")
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	// Fail the task
	testErr := errors.New("test error")
	_ = mgr.Fail(ctx, "task-1", testErr)

	// Should receive failure
	select {
	case update := <-ch:
		if update.State != StateFailed {
			t.Errorf("update.State = %v, want %v", update.State, StateFailed)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for fail notification")
	}

	// Channel should be closed
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("channel should be closed after terminal state")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("channel not closed after terminal state")
	}
}

func TestManager_Subscribe_ReceivesCancel(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	ch, err := mgr.Subscribe(ctx, "task-1")
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	// Cancel the task
	_ = mgr.Cancel(ctx, "task-1")

	// Should receive cancellation
	select {
	case update := <-ch:
		if update.State != StateCancelled {
			t.Errorf("update.State = %v, want %v", update.State, StateCancelled)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for cancel notification")
	}

	// Channel should be closed
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("channel should be closed after terminal state")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("channel not closed after terminal state")
	}
}

func TestManager_Subscribe_MultipleSubscribers(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	ch1, _ := mgr.Subscribe(ctx, "task-1")
	ch2, _ := mgr.Subscribe(ctx, "task-1")
	ch3, _ := mgr.Subscribe(ctx, "task-1")

	// Update the task
	_ = mgr.Update(ctx, "task-1", 0.5, "halfway")

	// All subscribers should receive the update
	channels := []<-chan *Task{ch1, ch2, ch3}
	for i, ch := range channels {
		select {
		case update := <-ch:
			if update.Progress != 0.5 {
				t.Errorf("subscriber %d: progress = %v, want 0.5", i, update.Progress)
			}
		case <-time.After(time.Second):
			t.Fatalf("subscriber %d: timeout waiting for update", i)
		}
	}
}

func TestManager_Subscribe_ContextCancellation(t *testing.T) {
	mgr := NewManager()
	bgCtx := context.Background()

	_, _ = mgr.Create(bgCtx, "task-1")

	// Create a cancellable context for subscription
	subCtx, cancel := context.WithCancel(bgCtx)

	ch, err := mgr.Subscribe(subCtx, "task-1")
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	// Cancel the subscription context
	cancel()

	// Give time for cleanup goroutine
	time.Sleep(50 * time.Millisecond)

	// Channel should be closed
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("channel should be closed after context cancellation")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("channel not closed after context cancellation")
	}
}

func TestManager_Subscribe_ClosesOnTerminal(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	ch, _ := mgr.Subscribe(ctx, "task-1")

	// Complete to terminal state
	_ = mgr.Complete(ctx, "task-1", nil)

	// Drain the complete notification
	<-ch

	// Channel should now be closed
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("channel should be closed after terminal state")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("channel not closed after terminal state")
	}
}

func TestManager_Subscribe_AlreadyTerminal(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")
	_ = mgr.Complete(ctx, "task-1", nil)

	// Subscribe to already completed task
	ch, err := mgr.Subscribe(ctx, "task-1")
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}

	// Channel should be immediately closed
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("channel should be closed for terminal task")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("channel not closed for terminal task")
	}
}

func TestManager_Subscribe_ConcurrentUpdates(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	ch, _ := mgr.Subscribe(ctx, "task-1")

	// Send many concurrent updates
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = mgr.Update(ctx, "task-1", float64(n)/20, "update")
		}(i)
	}

	// Complete after all updates
	go func() {
		wg.Wait()
		_ = mgr.Complete(ctx, "task-1", nil)
	}()

	// Drain all notifications until closed
	count := 0
	for range ch {
		count++
	}

	// Should have received some updates (may not be all due to buffer)
	if count == 0 {
		t.Error("should have received at least some updates")
	}
}

func TestManager_Subscribe_ReturnsClonedTasks(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	ch, _ := mgr.Subscribe(ctx, "task-1")

	_ = mgr.Update(ctx, "task-1", 0.5, "original")

	update := <-ch
	update.Message = "modified"

	// Verify stored task is unchanged
	task, _ := mgr.Get(ctx, "task-1")
	if task.Message != "original" {
		t.Errorf("modifying update affected stored task: got %v, want 'original'", task.Message)
	}

	_ = mgr.Complete(ctx, "task-1", nil)
}
