// Package stream provides streaming response support for protocol servers.
//
// This package enables streaming of events from servers to clients,
// supporting progress updates, partial results, and completion signals.
//
// # Ecosystem Position
//
// stream provides event streaming for long-running operations:
//
//	┌─────────────────────────────────────────────────────────────────┐
//	│                      Streaming Flow                             │
//	├─────────────────────────────────────────────────────────────────┤
//	│                                                                 │
//	│   Producer              Stream                  Consumer        │
//	│   ┌─────────┐        ┌───────────┐         ┌─────────┐         │
//	│   │  task   │────────│  Source   │─────────│  Sink   │         │
//	│   │ handler │  Send  │           │ Events  │         │         │
//	│   └─────────┘        │ ┌───────┐ │         └─────────┘         │
//	│        │             │ │Stream │ │              │               │
//	│        │             │ │       │ │              │               │
//	│        ▼             │ └───────┘ │              │               │
//	│   ┌─────────┐        │     │     │              │               │
//	│   │progress │────────│─────┼─────│──────────────▼               │
//	│   │ updates │        │     │     │         ┌─────────┐         │
//	│   └─────────┘        │     ▼     │         │ Handler │         │
//	│        │             │ ┌───────┐ │         │(callback)│         │
//	│        ▼             │ │Buffer │ │         └─────────┘         │
//	│   ┌─────────┐        │ │(opt.) │ │              ▲               │
//	│   │complete │────────│ └───────┘ │──────────────┘               │
//	│   │  /error │  Close └───────────┘   Consume                    │
//	│   └─────────┘                                                   │
//	│                                                                 │
//	└─────────────────────────────────────────────────────────────────┘
//
// # Event Types
//
// The package supports five event types:
//
//   - progress: Progress updates during execution (0.0 to 1.0)
//   - partial: Partial results as they become available
//   - complete: Final completion with result
//   - error: Error occurred during execution
//   - heartbeat: Keep-alive signal for long-running operations
//
// # Core Components
//
//   - [Stream]: Interface for sending events and managing lifecycle
//   - [Source]: Factory interface for creating streams (unbuffered/buffered)
//   - [Sink]: Interface for consuming streams with event handlers
//   - [Event]: Event structure with type, ID, data, and retry hint
//   - [EventType]: Type constants (progress, partial, complete, error, heartbeat)
//   - [DefaultStream]: Unbuffered stream implementation
//   - [BufferedStream]: Buffered stream with configurable backpressure
//
// # Quick Start
//
//	// Create a source and stream
//	source := stream.NewSource()
//	s := source.NewBufferedStream(ctx, 100)
//
//	// Producer sends events
//	go func() {
//	    s.Send(ctx, stream.Event{Type: stream.EventProgress, Data: 0.5})
//	    s.Send(ctx, stream.Event{Type: stream.EventComplete, Data: result})
//	    s.Close()
//	}()
//
//	// Consumer handles events
//	sink := stream.NewSink()
//	sink.Consume(ctx, s, func(event stream.Event) error {
//	    fmt.Printf("Event: %s\n", event.Type)
//	    return nil
//	})
//
// # Buffered vs Unbuffered Streams
//
// Unbuffered streams (NewStream) block on Send until consumed:
//
//	s := source.NewStream(ctx)
//	// Send blocks until event is received
//
// Buffered streams (NewBufferedStream) have a queue:
//
//	s := source.NewBufferedStream(ctx, 100)
//	// Send returns immediately if buffer has space
//
// # Backpressure Modes
//
// Buffered streams support two backpressure modes:
//
//   - BackpressureBlock (default): Send blocks when buffer is full
//   - BackpressureDrop: Send returns ErrBufferFull when buffer is full
//
// Example:
//
//	source := stream.NewSource(stream.WithBackpressure(stream.BackpressureDrop))
//	s := source.NewBufferedStream(ctx, 10)
//
// # Thread Safety
//
// All exported types are safe for concurrent use:
//
//   - [DefaultStream]: sync.RWMutex protects all operations; wg ensures
//     in-flight sends complete before events channel closes
//   - [BufferedStream]: sync.RWMutex protects all operations
//   - [DefaultSource]: Stateless after construction (backpressure config only)
//   - [DefaultSink]: Stateless, safe for concurrent Consume calls
//   - Event channels: Properly synchronized via mutexes
//
// # Error Handling
//
// Sentinel errors (use errors.Is for checking):
//
//   - [ErrStreamClosed]: Stream is closed, cannot send
//   - [ErrBufferFull]: Buffer full in drop mode
//
// The [StreamError] type wraps errors with stream context:
//
//	err := &StreamError{
//	    StreamID: "stream-123",
//	    Op:       "send",
//	    Err:      ErrStreamClosed,
//	}
//	// errors.Is(err, ErrStreamClosed) = true
//
// # Integration with ApertureStack
//
// stream integrates with other ApertureStack packages:
//
//   - task: Progress updates streamed to clients
//   - wire: Events encoded to protocol-specific formats
//   - transport: Streams delivered over HTTP/SSE/stdio
package stream
