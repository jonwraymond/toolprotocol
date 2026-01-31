package session

import "errors"

// Sentinel errors for session operations.
var (
	// ErrSessionNotFound is returned when a session cannot be found.
	ErrSessionNotFound = errors.New("session not found")

	// ErrSessionExpired is returned when a session has expired.
	ErrSessionExpired = errors.New("session expired")

	// ErrInvalidClientID is returned when a client ID is empty or invalid.
	ErrInvalidClientID = errors.New("invalid client ID")
)

// SessionError wraps an error with session context.
type SessionError struct {
	// SessionID is the ID of the session that caused the error.
	SessionID string

	// Op is the operation that failed.
	Op string

	// Err is the underlying error.
	Err error
}

// Error returns the error message.
func (e *SessionError) Error() string {
	if e.Err == nil {
		return "session " + e.SessionID + ": " + e.Op
	}
	return "session " + e.SessionID + ": " + e.Op + ": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *SessionError) Unwrap() error {
	return e.Err
}
