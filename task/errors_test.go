package task

import (
	"errors"
	"testing"
)

func TestErrTaskNotFound(t *testing.T) {
	if ErrTaskNotFound.Error() != "task: not found" {
		t.Errorf("ErrTaskNotFound.Error() = %v, want 'task: not found'", ErrTaskNotFound.Error())
	}
}

func TestErrTaskExists(t *testing.T) {
	if ErrTaskExists.Error() != "task: already exists" {
		t.Errorf("ErrTaskExists.Error() = %v, want 'task: already exists'", ErrTaskExists.Error())
	}
}

func TestErrInvalidState(t *testing.T) {
	if ErrInvalidState.Error() != "task: invalid state" {
		t.Errorf("ErrInvalidState.Error() = %v, want 'task: invalid state'", ErrInvalidState.Error())
	}
}

func TestErrInvalidTransition(t *testing.T) {
	if ErrInvalidTransition.Error() != "task: invalid transition" {
		t.Errorf("ErrInvalidTransition.Error() = %v, want 'task: invalid transition'", ErrInvalidTransition.Error())
	}
}

func TestErrEmptyID(t *testing.T) {
	if ErrEmptyID.Error() != "task: empty ID" {
		t.Errorf("ErrEmptyID.Error() = %v, want 'task: empty ID'", ErrEmptyID.Error())
	}
}

func TestTaskError_Error(t *testing.T) {
	err := &TaskError{
		TaskID: "task-1",
		Op:     "update",
		Err:    ErrTaskNotFound,
	}

	got := err.Error()
	want := "task task-1: update: task: not found"
	if got != want {
		t.Errorf("TaskError.Error() = %v, want %v", got, want)
	}
}

func TestTaskError_Unwrap(t *testing.T) {
	err := &TaskError{
		TaskID: "task-1",
		Op:     "update",
		Err:    ErrTaskNotFound,
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped != ErrTaskNotFound {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, ErrTaskNotFound)
	}
}

func TestTaskError_Is(t *testing.T) {
	err := &TaskError{
		TaskID: "task-1",
		Op:     "update",
		Err:    ErrTaskNotFound,
	}

	if !errors.Is(err, ErrTaskNotFound) {
		t.Error("errors.Is(err, ErrTaskNotFound) should be true")
	}

	if errors.Is(err, ErrTaskExists) {
		t.Error("errors.Is(err, ErrTaskExists) should be false")
	}
}

func TestTaskError_NilErr(t *testing.T) {
	err := &TaskError{
		TaskID: "task-1",
		Op:     "update",
		Err:    nil,
	}

	got := err.Error()
	want := "task task-1: update"
	if got != want {
		t.Errorf("TaskError.Error() with nil Err = %v, want %v", got, want)
	}
}
