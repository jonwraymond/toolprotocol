package stream_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonwraymond/toolprotocol/stream"
)

func ExampleNewSource() {
	source := stream.NewSource()

	fmt.Printf("Type: %T\n", source)
	// Output:
	// Type: *stream.DefaultSource
}

func ExampleNewSource_withBackpressure() {
	source := stream.NewSource(stream.WithBackpressure(stream.BackpressureDrop))

	fmt.Printf("Type: %T\n", source)
	// Output:
	// Type: *stream.DefaultSource
}

func ExampleDefaultSource_NewStream() {
	source := stream.NewSource()
	ctx := context.Background()

	s := source.NewStream(ctx)
	defer s.Close()

	fmt.Printf("Stream type: %T\n", s)
	// Output:
	// Stream type: *stream.DefaultStream
}

func ExampleDefaultSource_NewBufferedStream() {
	source := stream.NewSource()
	ctx := context.Background()

	s := source.NewBufferedStream(ctx, 100)
	defer s.Close()

	fmt.Printf("Stream type: %T\n", s)
	// Output:
	// Stream type: *stream.BufferedStream
}

func ExampleDefaultStream_Send() {
	source := stream.NewSource()
	ctx := context.Background()
	s := source.NewStream(ctx)

	// Use goroutine since unbuffered stream blocks until consumed
	done := make(chan struct{})
	go func() {
		defer close(done)
		// Consume the event
		if es, ok := s.(interface{ Events() <-chan stream.Event }); ok {
			<-es.Events()
		}
	}()

	err := s.Send(ctx, stream.Event{
		Type: stream.EventProgress,
		Data: 0.5,
	})
	fmt.Println("Send error:", err)

	s.Close()
	<-done
	// Output:
	// Send error: <nil>
}

func ExampleDefaultStream_Send_closed() {
	source := stream.NewSource()
	ctx := context.Background()
	s := source.NewStream(ctx)

	// Close the stream first
	s.Close()

	// Try to send to closed stream
	err := s.Send(ctx, stream.Event{Type: stream.EventProgress})
	fmt.Println("Is ErrStreamClosed:", errors.Is(err, stream.ErrStreamClosed))
	// Output:
	// Is ErrStreamClosed: true
}

func ExampleDefaultStream_Close() {
	source := stream.NewSource()
	ctx := context.Background()
	s := source.NewStream(ctx)

	// Close is idempotent
	err1 := s.Close()
	err2 := s.Close()

	fmt.Println("First close error:", err1)
	fmt.Println("Second close error:", err2)
	// Output:
	// First close error: <nil>
	// Second close error: <nil>
}

func ExampleDefaultStream_Done() {
	source := stream.NewSource()
	ctx := context.Background()
	s := source.NewStream(ctx)

	// Done channel is open before close
	select {
	case <-s.Done():
		fmt.Println("Done before close: closed")
	default:
		fmt.Println("Done before close: open")
	}

	s.Close()

	// Done channel is closed after close
	select {
	case <-s.Done():
		fmt.Println("Done after close: closed")
	default:
		fmt.Println("Done after close: open")
	}
	// Output:
	// Done before close: open
	// Done after close: closed
}

func ExampleBufferedStream_Send() {
	source := stream.NewSource()
	ctx := context.Background()
	s := source.NewBufferedStream(ctx, 10)

	// Buffered stream doesn't block
	err := s.Send(ctx, stream.Event{
		Type: stream.EventProgress,
		Data: 0.5,
	})
	fmt.Println("Send error:", err)

	s.Close()
	// Output:
	// Send error: <nil>
}

func ExampleBufferedStream_Send_backpressureDrop() {
	source := stream.NewSource(stream.WithBackpressure(stream.BackpressureDrop))
	ctx := context.Background()
	s := source.NewBufferedStream(ctx, 1) // tiny buffer

	// First send succeeds
	_ = s.Send(ctx, stream.Event{Type: stream.EventProgress})

	// Second send fails with buffer full (drop mode)
	err := s.Send(ctx, stream.Event{Type: stream.EventProgress})
	fmt.Println("Is ErrBufferFull:", errors.Is(err, stream.ErrBufferFull))

	s.Close()
	// Output:
	// Is ErrBufferFull: true
}

func ExampleNewSink() {
	sink := stream.NewSink()

	fmt.Printf("Type: %T\n", sink)
	// Output:
	// Type: *stream.DefaultSink
}

func ExampleDefaultSink_Consume() {
	source := stream.NewSource()
	sink := stream.NewSink()
	ctx := context.Background()

	s := source.NewBufferedStream(ctx, 10)

	// Send some events
	_ = s.Send(ctx, stream.Event{Type: stream.EventProgress, Data: 0.5})
	_ = s.Send(ctx, stream.Event{Type: stream.EventComplete, Data: "done"})
	s.Close()

	// Consume events
	var count int
	err := sink.Consume(ctx, s, func(event stream.Event) error {
		count++
		return nil
	})

	fmt.Println("Consume error:", err)
	fmt.Println("Events consumed:", count)
	// Output:
	// Consume error: <nil>
	// Events consumed: 2
}

func ExampleDefaultSink_Consume_handlerError() {
	source := stream.NewSource()
	sink := stream.NewSink()
	ctx := context.Background()

	s := source.NewBufferedStream(ctx, 10)

	// Send events
	_ = s.Send(ctx, stream.Event{Type: stream.EventProgress})
	_ = s.Send(ctx, stream.Event{Type: stream.EventProgress})
	s.Close()

	// Handler returns error on first event
	handlerErr := errors.New("handler failed")
	err := sink.Consume(ctx, s, func(event stream.Event) error {
		return handlerErr
	})

	fmt.Println("Is handler error:", errors.Is(err, handlerErr))
	// Output:
	// Is handler error: true
}

func ExampleEventType_String() {
	fmt.Println(stream.EventProgress.String())
	fmt.Println(stream.EventPartial.String())
	fmt.Println(stream.EventComplete.String())
	fmt.Println(stream.EventError.String())
	fmt.Println(stream.EventHeartbeat.String())
	// Output:
	// progress
	// partial
	// complete
	// error
	// heartbeat
}

func ExampleEventType_Valid() {
	fmt.Println("progress valid:", stream.EventProgress.Valid())
	fmt.Println("partial valid:", stream.EventPartial.Valid())
	fmt.Println("complete valid:", stream.EventComplete.Valid())
	fmt.Println("error valid:", stream.EventError.Valid())
	fmt.Println("heartbeat valid:", stream.EventHeartbeat.Valid())
	fmt.Println("unknown valid:", stream.EventType("unknown").Valid())
	// Output:
	// progress valid: true
	// partial valid: true
	// complete valid: true
	// error valid: true
	// heartbeat valid: true
	// unknown valid: false
}

func ExampleEvent_Clone() {
	original := stream.Event{
		Type:  stream.EventProgress,
		ID:    "evt-1",
		Data:  0.75,
		Retry: 1000,
	}

	clone := original.Clone()

	fmt.Println("Type matches:", clone.Type == original.Type)
	fmt.Println("ID matches:", clone.ID == original.ID)
	fmt.Println("Data matches:", clone.Data == original.Data)
	fmt.Println("Retry matches:", clone.Retry == original.Retry)
	// Output:
	// Type matches: true
	// ID matches: true
	// Data matches: true
	// Retry matches: true
}

func ExampleWithBackpressure() {
	// Block mode (default) - waits for buffer space
	sourceBlock := stream.NewSource(stream.WithBackpressure(stream.BackpressureBlock))
	_ = sourceBlock

	// Drop mode - drops events when buffer full
	sourceDrop := stream.NewSource(stream.WithBackpressure(stream.BackpressureDrop))
	_ = sourceDrop

	fmt.Println("Configured both modes")
	// Output:
	// Configured both modes
}

func ExampleWithHeartbeat() {
	opt := stream.WithHeartbeat(30 * time.Second)

	fmt.Println("Enabled:", opt.Enabled)
	fmt.Println("Interval:", opt.Interval)
	// Output:
	// Enabled: true
	// Interval: 30s
}

func ExampleStreamError() {
	err := &stream.StreamError{
		StreamID: "stream-123",
		Op:       "send",
		Err:      stream.ErrStreamClosed,
	}

	fmt.Println(err.Error())
	fmt.Println("Unwraps to ErrStreamClosed:", errors.Is(err, stream.ErrStreamClosed))
	// Output:
	// stream stream-123: send: stream: closed
	// Unwraps to ErrStreamClosed: true
}

func ExampleStreamError_noUnderlying() {
	err := &stream.StreamError{
		StreamID: "stream-123",
		Op:       "validate",
	}

	fmt.Println(err.Error())
	// Output:
	// stream stream-123: validate
}

func Example_streamingWorkflow() {
	source := stream.NewSource()
	sink := stream.NewSink()
	ctx := context.Background()

	// Create buffered stream for the workflow
	s := source.NewBufferedStream(ctx, 10)

	// Simulate producer in goroutine
	go func() {
		// Send progress updates
		_ = s.Send(ctx, stream.Event{Type: stream.EventProgress, Data: 0.25})
		_ = s.Send(ctx, stream.Event{Type: stream.EventProgress, Data: 0.50})
		_ = s.Send(ctx, stream.Event{Type: stream.EventProgress, Data: 0.75})

		// Send partial result
		_ = s.Send(ctx, stream.Event{Type: stream.EventPartial, Data: "partial data"})

		// Complete
		_ = s.Send(ctx, stream.Event{Type: stream.EventComplete, Data: "final result"})

		s.Close()
	}()

	// Consume all events
	var lastProgress float64
	var result any
	_ = sink.Consume(ctx, s, func(event stream.Event) error {
		switch event.Type {
		case stream.EventProgress:
			if p, ok := event.Data.(float64); ok {
				lastProgress = p
			}
		case stream.EventComplete:
			result = event.Data
		}
		return nil
	})

	fmt.Println("Last progress:", lastProgress)
	fmt.Println("Result:", result)
	// Output:
	// Last progress: 0.75
	// Result: final result
}
