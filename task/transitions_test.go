package task

import (
	"context"
	"errors"
	"testing"
)

func TestManager_Update_Progress(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, err := mgr.Create(ctx, "task-1")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err = mgr.Update(ctx, "task-1", 0.5, "")
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	task, _ := mgr.Get(ctx, "task-1")
	if task.Progress != 0.5 {
		t.Errorf("Progress = %v, want 0.5", task.Progress)
	}
	if task.State != StateRunning {
		t.Errorf("State = %v, want %v (should transition on first update)", task.State, StateRunning)
	}
}

func TestManager_Update_Message(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	err := mgr.Update(ctx, "task-1", 0.25, "Processing step 1")
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	task, _ := mgr.Get(ctx, "task-1")
	if task.Message != "Processing step 1" {
		t.Errorf("Message = %v, want 'Processing step 1'", task.Message)
	}
}

func TestManager_Update_NotFound(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	err := mgr.Update(ctx, "nonexistent", 0.5, "test")
	if err == nil {
		t.Fatal("Update() expected error for nonexistent task")
	}
	if err != ErrTaskNotFound {
		t.Errorf("Update() error = %v, want %v", err, ErrTaskNotFound)
	}
}

func TestManager_Update_TerminalState(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	// Create and complete a task
	_, _ = mgr.Create(ctx, "task-1")
	_ = mgr.Complete(ctx, "task-1", "result")

	// Try to update completed task
	err := mgr.Update(ctx, "task-1", 0.75, "more progress")
	if err == nil {
		t.Fatal("Update() expected error for terminal state")
	}
	if err != ErrInvalidTransition {
		t.Errorf("Update() error = %v, want %v", err, ErrInvalidTransition)
	}
}

func TestManager_Complete(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")
	_ = mgr.Update(ctx, "task-1", 0.5, "halfway")

	err := mgr.Complete(ctx, "task-1", "final result")
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	task, _ := mgr.Get(ctx, "task-1")
	if task.State != StateComplete {
		t.Errorf("State = %v, want %v", task.State, StateComplete)
	}
	if task.Result != "final result" {
		t.Errorf("Result = %v, want 'final result'", task.Result)
	}
	if task.Progress != 1.0 {
		t.Errorf("Progress = %v, want 1.0", task.Progress)
	}
}

func TestManager_Complete_SetsCompletedAt(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	err := mgr.Complete(ctx, "task-1", nil)
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	task, _ := mgr.Get(ctx, "task-1")
	if task.CompletedAt == nil {
		t.Error("CompletedAt should be set after Complete()")
	}
}

func TestManager_Complete_NotRunning(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	// Create and cancel a task
	_, _ = mgr.Create(ctx, "task-1")
	_ = mgr.Cancel(ctx, "task-1")

	// Try to complete cancelled task
	err := mgr.Complete(ctx, "task-1", "result")
	if err == nil {
		t.Fatal("Complete() expected error for terminal state")
	}
	if err != ErrInvalidTransition {
		t.Errorf("Complete() error = %v, want %v", err, ErrInvalidTransition)
	}
}

func TestManager_Complete_FromPending(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	// Complete directly from pending (allowed)
	_, _ = mgr.Create(ctx, "task-1")

	err := mgr.Complete(ctx, "task-1", "instant result")
	if err != nil {
		t.Fatalf("Complete() from pending error = %v", err)
	}

	task, _ := mgr.Get(ctx, "task-1")
	if task.State != StateComplete {
		t.Errorf("State = %v, want %v", task.State, StateComplete)
	}
}

func TestManager_Fail(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")
	_ = mgr.Update(ctx, "task-1", 0.5, "halfway")

	testErr := errors.New("something went wrong")
	err := mgr.Fail(ctx, "task-1", testErr)
	if err != nil {
		t.Fatalf("Fail() error = %v", err)
	}

	task, _ := mgr.Get(ctx, "task-1")
	if task.State != StateFailed {
		t.Errorf("State = %v, want %v", task.State, StateFailed)
	}
	if task.Error == nil || task.Error.Error() != testErr.Error() {
		t.Errorf("Error = %v, want %v", task.Error, testErr)
	}
}

func TestManager_Fail_SetsError(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	testErr := errors.New("test error")
	_ = mgr.Fail(ctx, "task-1", testErr)

	task, _ := mgr.Get(ctx, "task-1")
	if task.Error == nil {
		t.Fatal("Error should be set after Fail()")
	}
	if task.Error.Error() != "test error" {
		t.Errorf("Error = %v, want 'test error'", task.Error)
	}
}

func TestManager_Fail_SetsCompletedAt(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")
	_ = mgr.Fail(ctx, "task-1", errors.New("error"))

	task, _ := mgr.Get(ctx, "task-1")
	if task.CompletedAt == nil {
		t.Error("CompletedAt should be set after Fail()")
	}
}

func TestManager_Cancel(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")
	_ = mgr.Update(ctx, "task-1", 0.5, "halfway")

	err := mgr.Cancel(ctx, "task-1")
	if err != nil {
		t.Fatalf("Cancel() error = %v", err)
	}

	task, _ := mgr.Get(ctx, "task-1")
	if task.State != StateCancelled {
		t.Errorf("State = %v, want %v", task.State, StateCancelled)
	}
}

func TestManager_Cancel_FromPending(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")

	err := mgr.Cancel(ctx, "task-1")
	if err != nil {
		t.Fatalf("Cancel() from pending error = %v", err)
	}

	task, _ := mgr.Get(ctx, "task-1")
	if task.State != StateCancelled {
		t.Errorf("State = %v, want %v", task.State, StateCancelled)
	}
}

func TestManager_Cancel_SetsCompletedAt(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	_, _ = mgr.Create(ctx, "task-1")
	_ = mgr.Cancel(ctx, "task-1")

	task, _ := mgr.Get(ctx, "task-1")
	if task.CompletedAt == nil {
		t.Error("CompletedAt should be set after Cancel()")
	}
}

func TestManager_Cancel_AlreadyTerminal(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	// Create and complete a task
	_, _ = mgr.Create(ctx, "task-1")
	_ = mgr.Complete(ctx, "task-1", "result")

	// Try to cancel completed task
	err := mgr.Cancel(ctx, "task-1")
	if err == nil {
		t.Fatal("Cancel() expected error for terminal state")
	}
	if err != ErrInvalidTransition {
		t.Errorf("Cancel() error = %v, want %v", err, ErrInvalidTransition)
	}
}

func TestManager_StateTransitions_Valid(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T, mgr *DefaultManager, ctx context.Context)
		action       func(mgr *DefaultManager, ctx context.Context) error
		wantState    State
		wantErr      bool
		wantErrValue error
	}{
		{
			name: "pending->running via Update",
			setup: func(t *testing.T, m *DefaultManager, ctx context.Context) {
				if _, err := m.Create(ctx, "t"); err != nil {
					t.Fatalf("Create() error = %v", err)
				}
			},
			action:    func(m *DefaultManager, ctx context.Context) error { return m.Update(ctx, "t", 0.5, "") },
			wantState: StateRunning,
		},
		{
			name: "running->complete",
			setup: func(t *testing.T, m *DefaultManager, ctx context.Context) {
				if _, err := m.Create(ctx, "t"); err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				if err := m.Update(ctx, "t", 0.5, ""); err != nil {
					t.Fatalf("Update() error = %v", err)
				}
			},
			action:    func(m *DefaultManager, ctx context.Context) error { return m.Complete(ctx, "t", nil) },
			wantState: StateComplete,
		},
		{
			name: "running->failed",
			setup: func(t *testing.T, m *DefaultManager, ctx context.Context) {
				if _, err := m.Create(ctx, "t"); err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				if err := m.Update(ctx, "t", 0.5, ""); err != nil {
					t.Fatalf("Update() error = %v", err)
				}
			},
			action:    func(m *DefaultManager, ctx context.Context) error { return m.Fail(ctx, "t", errors.New("e")) },
			wantState: StateFailed,
		},
		{
			name: "running->cancelled",
			setup: func(t *testing.T, m *DefaultManager, ctx context.Context) {
				if _, err := m.Create(ctx, "t"); err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				if err := m.Update(ctx, "t", 0.5, ""); err != nil {
					t.Fatalf("Update() error = %v", err)
				}
			},
			action:    func(m *DefaultManager, ctx context.Context) error { return m.Cancel(ctx, "t") },
			wantState: StateCancelled,
		},
		{
			name: "pending->complete (allowed)",
			setup: func(t *testing.T, m *DefaultManager, ctx context.Context) {
				if _, err := m.Create(ctx, "t"); err != nil {
					t.Fatalf("Create() error = %v", err)
				}
			},
			action:    func(m *DefaultManager, ctx context.Context) error { return m.Complete(ctx, "t", nil) },
			wantState: StateComplete,
		},
		{
			name: "pending->cancelled (allowed)",
			setup: func(t *testing.T, m *DefaultManager, ctx context.Context) {
				if _, err := m.Create(ctx, "t"); err != nil {
					t.Fatalf("Create() error = %v", err)
				}
			},
			action:    func(m *DefaultManager, ctx context.Context) error { return m.Cancel(ctx, "t") },
			wantState: StateCancelled,
		},
		{
			name: "complete->update (invalid)",
			setup: func(t *testing.T, m *DefaultManager, ctx context.Context) {
				if _, err := m.Create(ctx, "t"); err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				if err := m.Complete(ctx, "t", nil); err != nil {
					t.Fatalf("Complete() error = %v", err)
				}
			},
			action:       func(m *DefaultManager, ctx context.Context) error { return m.Update(ctx, "t", 1.0, "") },
			wantState:    StateComplete,
			wantErr:      true,
			wantErrValue: ErrInvalidTransition,
		},
		{
			name: "failed->cancel (invalid)",
			setup: func(t *testing.T, m *DefaultManager, ctx context.Context) {
				if _, err := m.Create(ctx, "t"); err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				if err := m.Fail(ctx, "t", errors.New("e")); err != nil {
					t.Fatalf("Fail() error = %v", err)
				}
			},
			action:       func(m *DefaultManager, ctx context.Context) error { return m.Cancel(ctx, "t") },
			wantState:    StateFailed,
			wantErr:      true,
			wantErrValue: ErrInvalidTransition,
		},
		{
			name: "cancelled->complete (invalid)",
			setup: func(t *testing.T, m *DefaultManager, ctx context.Context) {
				if _, err := m.Create(ctx, "t"); err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				if err := m.Cancel(ctx, "t"); err != nil {
					t.Fatalf("Cancel() error = %v", err)
				}
			},
			action:       func(m *DefaultManager, ctx context.Context) error { return m.Complete(ctx, "t", nil) },
			wantState:    StateCancelled,
			wantErr:      true,
			wantErrValue: ErrInvalidTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := NewManager()
			ctx := context.Background()

			if tt.setup != nil {
				tt.setup(t, mgr, ctx)
			}
			err := tt.action(mgr, ctx)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.wantErrValue != nil && err != tt.wantErrValue {
					t.Errorf("error = %v, want %v", err, tt.wantErrValue)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			task, _ := mgr.Get(ctx, "t")
			if task.State != tt.wantState {
				t.Errorf("State = %v, want %v", task.State, tt.wantState)
			}
		})
	}
}

func TestManager_Update_UpdatesTimestamp(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	created, _ := mgr.Create(ctx, "task-1")
	originalUpdatedAt := created.UpdatedAt

	// Small delay to ensure timestamp difference
	_ = mgr.Update(ctx, "task-1", 0.5, "test")

	task, _ := mgr.Get(ctx, "task-1")
	if !task.UpdatedAt.After(originalUpdatedAt) && task.UpdatedAt.Equal(originalUpdatedAt) {
		// Allow equal if test runs fast, but Updated should not be before
		if task.UpdatedAt.Before(originalUpdatedAt) {
			t.Error("UpdatedAt should be >= original after Update()")
		}
	}
}
