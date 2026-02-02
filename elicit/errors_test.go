package elicit

import (
	"errors"
	"testing"
)

func TestErrInvalidRequest(t *testing.T) {
	if ErrInvalidRequest.Error() != "elicit: invalid request" {
		t.Errorf("Error() = %q, want %q", ErrInvalidRequest.Error(), "elicit: invalid request")
	}
}

func TestErrTimeout(t *testing.T) {
	if ErrTimeout.Error() != "elicit: timeout" {
		t.Errorf("Error() = %q, want %q", ErrTimeout.Error(), "elicit: timeout")
	}
}

func TestErrCancelled(t *testing.T) {
	if ErrCancelled.Error() != "elicit: cancelled" {
		t.Errorf("Error() = %q, want %q", ErrCancelled.Error(), "elicit: cancelled")
	}
}

func TestErrNoHandler(t *testing.T) {
	if ErrNoHandler.Error() != "elicit: no handler" {
		t.Errorf("Error() = %q, want %q", ErrNoHandler.Error(), "elicit: no handler")
	}
}

func TestElicitError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *ElicitError
		want string
	}{
		{
			name: "with underlying error",
			err: &ElicitError{
				RequestID: "req-123",
				Op:        "validate",
				Err:       ErrInvalidRequest,
			},
			want: "elicit req-123: validate: elicit: invalid request",
		},
		{
			name: "without underlying error",
			err: &ElicitError{
				RequestID: "req-123",
				Op:        "handle",
				Err:       nil,
			},
			want: "elicit req-123: handle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestElicitError_Unwrap(t *testing.T) {
	underlying := ErrInvalidRequest
	err := &ElicitError{
		RequestID: "req-123",
		Op:        "validate",
		Err:       underlying,
	}

	if !errors.Is(err, underlying) {
		t.Error("errors.Is should return true for underlying error")
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped != underlying {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlying)
	}
}
