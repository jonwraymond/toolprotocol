package task

import (
	"testing"
	"time"
)

func TestState_String(t *testing.T) {
	tests := []struct {
		state State
		want  string
	}{
		{StatePending, "pending"},
		{StateRunning, "running"},
		{StateComplete, "complete"},
		{StateFailed, "failed"},
		{StateCancelled, "cancelled"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("State.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_Valid(t *testing.T) {
	tests := []struct {
		state State
		valid bool
	}{
		{StatePending, true},
		{StateRunning, true},
		{StateComplete, true},
		{StateFailed, true},
		{StateCancelled, true},
		{State("invalid"), false},
		{State(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if got := tt.state.Valid(); got != tt.valid {
				t.Errorf("State.Valid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestState_IsTerminal(t *testing.T) {
	tests := []struct {
		state    State
		terminal bool
	}{
		{StatePending, false},
		{StateRunning, false},
		{StateComplete, true},
		{StateFailed, true},
		{StateCancelled, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if got := tt.state.IsTerminal(); got != tt.terminal {
				t.Errorf("State.IsTerminal() = %v, want %v", got, tt.terminal)
			}
		})
	}
}

func TestTask_Fields(t *testing.T) {
	now := time.Now()
	completedAt := now.Add(time.Minute)

	task := &Task{
		ID:          "task-1",
		State:       StateRunning,
		Progress:    0.5,
		Message:     "Processing...",
		Result:      "some result",
		Error:       nil,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: &completedAt,
	}

	if task.ID != "task-1" {
		t.Errorf("ID = %v, want task-1", task.ID)
	}
	if task.State != StateRunning {
		t.Errorf("State = %v, want %v", task.State, StateRunning)
	}
	if task.Progress != 0.5 {
		t.Errorf("Progress = %v, want 0.5", task.Progress)
	}
	if task.Message != "Processing..." {
		t.Errorf("Message = %v, want Processing...", task.Message)
	}
	if task.Result != "some result" {
		t.Errorf("Result = %v, want some result", task.Result)
	}
	if task.Error != nil {
		t.Errorf("Error = %v, want nil", task.Error)
	}
	if !task.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", task.CreatedAt, now)
	}
	if !task.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", task.UpdatedAt, now)
	}
	if task.CompletedAt == nil || !task.CompletedAt.Equal(completedAt) {
		t.Errorf("CompletedAt = %v, want %v", task.CompletedAt, completedAt)
	}
}

func TestTask_Clone(t *testing.T) {
	now := time.Now()
	completedAt := now.Add(time.Minute)

	original := &Task{
		ID:          "task-1",
		State:       StateComplete,
		Progress:    1.0,
		Message:     "Done",
		Result:      "result",
		Error:       nil,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: &completedAt,
	}

	clone := original.Clone()

	// Verify clone has same values
	if clone.ID != original.ID {
		t.Errorf("clone.ID = %v, want %v", clone.ID, original.ID)
	}
	if clone.State != original.State {
		t.Errorf("clone.State = %v, want %v", clone.State, original.State)
	}
	if clone.Progress != original.Progress {
		t.Errorf("clone.Progress = %v, want %v", clone.Progress, original.Progress)
	}
	if clone.Message != original.Message {
		t.Errorf("clone.Message = %v, want %v", clone.Message, original.Message)
	}

	// Verify clone is independent (modifying clone doesn't affect original)
	clone.Message = "Modified"
	if original.Message == clone.Message {
		t.Error("modifying clone affected original")
	}

	// Verify CompletedAt is a deep copy
	if clone.CompletedAt == original.CompletedAt {
		t.Error("CompletedAt should be a new pointer, not shared")
	}
	if clone.CompletedAt == nil || !clone.CompletedAt.Equal(*original.CompletedAt) {
		t.Errorf("clone.CompletedAt = %v, want %v", clone.CompletedAt, original.CompletedAt)
	}
}

func TestTask_Clone_NilCompletedAt(t *testing.T) {
	original := &Task{
		ID:          "task-1",
		State:       StateRunning,
		CompletedAt: nil,
	}

	clone := original.Clone()

	if clone.CompletedAt != nil {
		t.Errorf("clone.CompletedAt = %v, want nil", clone.CompletedAt)
	}
}

func TestManagerInterface_Defined(t *testing.T) {
	// This test ensures the Manager interface is correctly defined
	// The actual implementation contract is tested in manager_test.go
	var m Manager
	_ = m // verify interface is defined
}
