package elicit

import "errors"

// Sentinel errors for elicitation operations.
// All errors use the "elicit: " prefix for consistent error identification.
var (
	// ErrInvalidRequest is returned when a request is invalid.
	ErrInvalidRequest = errors.New("elicit: invalid request")

	// ErrTimeout is returned when a request times out.
	ErrTimeout = errors.New("elicit: timeout")

	// ErrCancelled is returned when a request is cancelled.
	ErrCancelled = errors.New("elicit: cancelled")

	// ErrNoHandler is returned when no handler is available.
	ErrNoHandler = errors.New("elicit: no handler")
)

// ElicitError wraps an error with elicitation context.
type ElicitError struct {
	// RequestID is the ID of the request that caused the error.
	RequestID string

	// Op is the operation that failed.
	Op string

	// Err is the underlying error.
	Err error
}

// Error returns the error message.
func (e *ElicitError) Error() string {
	if e.Err == nil {
		return "elicit " + e.RequestID + ": " + e.Op
	}
	return "elicit " + e.RequestID + ": " + e.Op + ": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *ElicitError) Unwrap() error {
	return e.Err
}
