package elicit

import "errors"

// Sentinel errors for elicitation operations.
var (
	// ErrInvalidRequest is returned when a request is invalid.
	ErrInvalidRequest = errors.New("invalid elicitation request")

	// ErrTimeout is returned when a request times out.
	ErrTimeout = errors.New("elicitation request timed out")

	// ErrCancelled is returned when a request is cancelled.
	ErrCancelled = errors.New("elicitation request cancelled")

	// ErrNoHandler is returned when no handler is available.
	ErrNoHandler = errors.New("no elicitation handler")
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
