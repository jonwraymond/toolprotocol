package stream

import (
	"context"
	"sync"
)

// BackpressureMode defines how to handle buffer overflow.
type BackpressureMode int

const (
	// BackpressureBlock blocks until space is available.
	BackpressureBlock BackpressureMode = iota

	// BackpressureDrop drops the event if buffer is full.
	BackpressureDrop
)

// BufferedStream is a stream with a buffer for events.
type BufferedStream struct {
	mu           sync.RWMutex
	events       chan Event
	done         chan struct{}
	closed       bool
	backpressure BackpressureMode
}

// newBufferedStream creates a new BufferedStream with the given buffer size.
func newBufferedStream(size int, backpressure BackpressureMode) *BufferedStream {
	if size < 1 {
		size = 1
	}
	return &BufferedStream{
		events:       make(chan Event, size),
		done:         make(chan struct{}),
		backpressure: backpressure,
	}
}

// Send sends an event to the stream.
func (s *BufferedStream) Send(ctx context.Context, event Event) error {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return ErrStreamClosed
	}
	events := s.events
	done := s.done
	backpressure := s.backpressure
	s.mu.RUnlock()

	if backpressure == BackpressureDrop {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-done:
			return ErrStreamClosed
		case events <- event:
			return nil
		default:
			return ErrBufferFull
		}
	}

	// BackpressureBlock
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return ErrStreamClosed
	case events <- event:
		return nil
	}
}

// Close closes the stream.
func (s *BufferedStream) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	close(s.done)
	close(s.events)
	return nil
}

// Done returns a channel that is closed when the stream is closed.
func (s *BufferedStream) Done() <-chan struct{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.done
}

// Events returns the events channel for consuming.
func (s *BufferedStream) Events() <-chan Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.events
}

// Drain reads and discards all pending events.
func (s *BufferedStream) Drain() {
	for {
		select {
		case _, ok := <-s.events:
			if !ok {
				return
			}
		default:
			return
		}
	}
}
