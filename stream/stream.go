package stream

import "context"

// EventType represents the type of stream event.
type EventType string

const (
	// EventProgress indicates a progress update.
	EventProgress EventType = "progress"

	// EventPartial indicates a partial result.
	EventPartial EventType = "partial"

	// EventComplete indicates completion with final result.
	EventComplete EventType = "complete"

	// EventError indicates an error occurred.
	EventError EventType = "error"

	// EventHeartbeat indicates a keep-alive signal.
	EventHeartbeat EventType = "heartbeat"
)

// String returns the string representation of the event type.
func (et EventType) String() string {
	return string(et)
}

// Valid returns true if the event type is a known valid type.
func (et EventType) Valid() bool {
	switch et {
	case EventProgress, EventPartial, EventComplete, EventError, EventHeartbeat:
		return true
	default:
		return false
	}
}

// Event represents a stream event.
type Event struct {
	// Type is the event type.
	Type EventType

	// ID is an optional event identifier.
	ID string

	// Data is the event payload.
	Data any

	// Retry is the reconnect interval in milliseconds for SSE clients.
	Retry int
}

// Clone returns a copy of the event.
func (e Event) Clone() Event {
	return Event{
		Type:  e.Type,
		ID:    e.ID,
		Data:  e.Data,
		Retry: e.Retry,
	}
}

// Stream represents a streaming response.
//
// Contract:
//   - Concurrency: Implementations must be safe for concurrent use.
//   - Context: Send honors context cancellation.
//   - Errors: Returns ErrStreamClosed for closed streams; use errors.Is for checking.
//   - Ownership: Events are not cloned; caller should not modify after Send.
//   - Idempotent: Close may be called multiple times safely.
type Stream interface {
	// Send sends an event to the stream.
	// Returns ErrStreamClosed if the stream is closed.
	Send(ctx context.Context, event Event) error

	// Close closes the stream.
	// Multiple calls are idempotent.
	Close() error

	// Done returns a channel that is closed when the stream is closed.
	Done() <-chan struct{}
}

// Source creates streams.
//
// Contract:
//   - Concurrency: Implementations must be safe for concurrent use.
//   - Context: Context is associated with the created stream's lifecycle.
//   - Ownership: Caller owns the returned Stream and must Close it.
type Source interface {
	// NewStream creates a new unbuffered stream.
	NewStream(ctx context.Context) Stream

	// NewBufferedStream creates a new buffered stream with the given size.
	NewBufferedStream(ctx context.Context, size int) Stream
}

// Sink consumes streams.
//
// Contract:
//   - Concurrency: Implementations must be safe for concurrent use.
//   - Context: Consume returns ctx.Err() on context cancellation.
//   - Errors: Returns handler errors immediately, stopping consumption.
//   - Completion: Returns nil when stream closes and all events are handled.
type Sink interface {
	// Consume reads events from the stream and calls the handler for each.
	// Returns when the stream is closed or an error occurs.
	Consume(ctx context.Context, stream Stream, handler EventHandler) error
}

// EventHandler is called for each event in a stream.
type EventHandler func(event Event) error
