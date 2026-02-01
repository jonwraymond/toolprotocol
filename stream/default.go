package stream

import (
	"context"
	"sync"
)

// DefaultStream is a basic unbuffered stream implementation.
type DefaultStream struct {
	mu     sync.RWMutex
	events chan Event
	done   chan struct{}
	closed bool
	wg     sync.WaitGroup
}

// newDefaultStream creates a new DefaultStream.
func newDefaultStream() *DefaultStream {
	return &DefaultStream{
		events: make(chan Event),
		done:   make(chan struct{}),
	}
}

// Send sends an event to the stream.
func (s *DefaultStream) Send(ctx context.Context, event Event) error {
	s.mu.RLock()
	if s.closed {
		s.mu.RUnlock()
		return ErrStreamClosed
	}
	s.wg.Add(1)
	events := s.events
	done := s.done
	s.mu.RUnlock()
	defer s.wg.Done()

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
func (s *DefaultStream) Close() error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil // idempotent
	}

	s.closed = true
	close(s.done)
	s.mu.Unlock()

	s.wg.Wait()
	close(s.events)
	return nil
}

// Done returns a channel that is closed when the stream is closed.
func (s *DefaultStream) Done() <-chan struct{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.done
}

// Events returns the events channel for consuming.
// The channel is closed after Close() returns (once in-flight sends complete).
func (s *DefaultStream) Events() <-chan Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.events
}
