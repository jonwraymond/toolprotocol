package stream

import "errors"

// Sentinel errors for stream operations.
var (
	// ErrStreamClosed is returned when sending to a closed stream.
	ErrStreamClosed = errors.New("stream closed")

	// ErrBufferFull is returned when the buffer is full and backpressure is drop.
	ErrBufferFull = errors.New("stream buffer full")
)

// StreamError wraps an error with stream context.
type StreamError struct {
	// StreamID is the ID of the stream that caused the error.
	StreamID string

	// Op is the operation that failed.
	Op string

	// Err is the underlying error.
	Err error
}

// Error returns the error message.
func (e *StreamError) Error() string {
	if e.Err == nil {
		return "stream " + e.StreamID + ": " + e.Op
	}
	return "stream " + e.StreamID + ": " + e.Op + ": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *StreamError) Unwrap() error {
	return e.Err
}
