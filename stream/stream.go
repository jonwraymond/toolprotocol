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
type Source interface {
	// NewStream creates a new unbuffered stream.
	NewStream(ctx context.Context) Stream

	// NewBufferedStream creates a new buffered stream with the given size.
	NewBufferedStream(ctx context.Context, size int) Stream
}

// Sink consumes streams.
type Sink interface {
	// Consume reads events from the stream and calls the handler for each.
	// Returns when the stream is closed or an error occurs.
	Consume(ctx context.Context, stream Stream, handler EventHandler) error
}

// EventHandler is called for each event in a stream.
type EventHandler func(event Event) error
