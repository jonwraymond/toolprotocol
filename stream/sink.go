package stream

import "context"

// eventStream is an interface for streams that expose their events channel.
type eventStream interface {
	Stream
	Events() <-chan Event
}

// DefaultSink is the default Sink implementation.
type DefaultSink struct{}

// NewSink creates a new DefaultSink.
func NewSink() *DefaultSink {
	return &DefaultSink{}
}

// Consume reads events from the stream and calls the handler for each.
func (s *DefaultSink) Consume(ctx context.Context, stream Stream, handler EventHandler) error {
	// Get events channel if available
	var events <-chan Event
	if es, ok := stream.(eventStream); ok {
		events = es.Events()
	} else {
		// For streams without Events(), we can only wait for Done()
		<-stream.Done()
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-stream.Done():
			// Drain remaining events
			for event := range events {
				if err := handler(event); err != nil {
					return err
				}
			}
			return nil
		case event, ok := <-events:
			if !ok {
				return nil
			}
			if err := handler(event); err != nil {
				return err
			}
		}
	}
}
