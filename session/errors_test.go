package session

import (
	"errors"
	"testing"
)

func TestErrSessionNotFound(t *testing.T) {
	err := ErrSessionNotFound
	if err.Error() != "session not found" {
		t.Errorf("Error() = %q, want %q", err.Error(), "session not found")
	}
}

func TestErrSessionExpired(t *testing.T) {
	err := ErrSessionExpired
	if err.Error() != "session expired" {
		t.Errorf("Error() = %q, want %q", err.Error(), "session expired")
	}
}

func TestErrInvalidClientID(t *testing.T) {
	err := ErrInvalidClientID
	if err.Error() != "invalid client ID" {
		t.Errorf("Error() = %q, want %q", err.Error(), "invalid client ID")
	}
}

func TestSessionError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     *SessionError
		want    string
	}{
		{
			name: "with underlying error",
			err: &SessionError{
				SessionID: "sess-123",
				Op:        "get",
				Err:       ErrSessionNotFound,
			},
			want: "session sess-123: get: session not found",
		},
		{
			name: "without underlying error",
			err: &SessionError{
				SessionID: "sess-123",
				Op:        "create",
				Err:       nil,
			},
			want: "session sess-123: create",
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

func TestSessionError_Unwrap(t *testing.T) {
	underlying := ErrSessionNotFound
	err := &SessionError{
		SessionID: "sess-123",
		Op:        "get",
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

func TestSessionError_Unwrap_Nil(t *testing.T) {
	err := &SessionError{
		SessionID: "sess-123",
		Op:        "create",
		Err:       nil,
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped != nil {
		t.Errorf("Unwrap() = %v, want nil", unwrapped)
	}
}
