// Package stream provides streaming response support for protocol servers.
//
// This package enables streaming of events from servers to clients,
// supporting progress updates, partial results, and completion signals.
//
// # Event Types
//
// The package supports five event types:
//
//   - progress: Progress updates during execution
//   - partial: Partial results as they become available
//   - complete: Final completion with result
//   - error: Error occurred during execution
//   - heartbeat: Keep-alive signal for long-running operations
//
// # Stream Interface
//
// The Stream interface represents a streaming response:
//
//	type Stream interface {
//	    Send(ctx context.Context, event Event) error
//	    Close() error
//	    Done() <-chan struct{}
//	}
//
// # Source and Sink
//
// Source creates streams, Sink consumes them:
//
//	// Create a stream
//	source := stream.NewSource()
//	s := source.NewStream(ctx)
//
//	// Send events
//	s.Send(ctx, stream.Event{Type: stream.EventProgress, Data: 0.5})
//	s.Send(ctx, stream.Event{Type: stream.EventComplete, Data: result})
//	s.Close()
//
//	// Consume a stream
//	sink := stream.NewSink()
//	sink.Consume(ctx, s, func(event stream.Event) error {
//	    fmt.Printf("Event: %s\n", event.Type)
//	    return nil
//	})
//
// # Buffered Streams
//
// For high-throughput scenarios, use buffered streams:
//
//	s := source.NewBufferedStream(ctx, 100) // buffer size 100
package stream
