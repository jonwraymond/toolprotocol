package stream

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestSink_Consume(t *testing.T) {
	source := NewSource()
	sink := NewSink()
	ctx := context.Background()

	s := source.NewBufferedStream(ctx, 10).(*BufferedStream)

	// Send some events
	events := []Event{
		{Type: EventProgress, Data: 0.25},
		{Type: EventProgress, Data: 0.5},
		{Type: EventComplete, Data: "done"},
	}
	for _, e := range events {
		_ = s.Send(ctx, e)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// Consume
	var received []Event
	err := sink.Consume(ctx, s, func(event Event) error {
		received = append(received, event)
		return nil
	})

	if err != nil {
		t.Fatalf("Consume() error = %v", err)
	}
	if len(received) != len(events) {
		t.Errorf("received %d events, want %d", len(received), len(events))
	}
}

func TestSink_Consume_AllEvents(t *testing.T) {
	source := NewSource()
	sink := NewSink()
	ctx := context.Background()

	s := source.NewBufferedStream(ctx, 100).(*BufferedStream)

	// Send events and close
	for i := 0; i < 50; i++ {
		_ = s.Send(ctx, Event{Type: EventProgress, Data: i})
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	count := 0
	err := sink.Consume(ctx, s, func(event Event) error {
		count++
		return nil
	})

	if err != nil {
		t.Fatalf("Consume() error = %v", err)
	}
	if count != 50 {
		t.Errorf("received %d events, want 50", count)
	}
}

func TestSink_Consume_HandlerError(t *testing.T) {
	source := NewSource()
	sink := NewSink()
	ctx := context.Background()

	s := source.NewBufferedStream(ctx, 10).(*BufferedStream)

	_ = s.Send(ctx, Event{Type: EventProgress})
	_ = s.Send(ctx, Event{Type: EventError})

	testErr := errors.New("handler error")
	err := sink.Consume(ctx, s, func(event Event) error {
		if event.Type == EventError {
			return testErr
		}
		return nil
	})

	if err != testErr {
		t.Errorf("Consume() error = %v, want testErr", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestSink_Consume_ContextCancellation(t *testing.T) {
	source := NewSource()
	sink := NewSink()
	ctx, cancel := context.WithCancel(context.Background())

	s := source.NewStream(ctx).(*DefaultStream)
	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})

	var wg sync.WaitGroup
	wg.Add(1)

	var consumeErr error
	go func() {
		defer wg.Done()
		consumeErr = sink.Consume(ctx, s, func(event Event) error {
			return nil
		})
	}()

	// Cancel context
	time.Sleep(10 * time.Millisecond)
	cancel()
	wg.Wait()

	if consumeErr != context.Canceled {
		t.Errorf("Consume() error = %v, want context.Canceled", consumeErr)
	}
}

func TestSink_Consume_StreamClosed(t *testing.T) {
	source := NewSource()
	sink := NewSink()
	ctx := context.Background()

	s := source.NewStream(ctx).(*DefaultStream)

	var wg sync.WaitGroup
	wg.Add(1)

	var consumeErr error
	go func() {
		defer wg.Done()
		consumeErr = sink.Consume(ctx, s, func(event Event) error {
			return nil
		})
	}()

	// Close stream
	time.Sleep(10 * time.Millisecond)
	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	wg.Wait()

	if consumeErr != nil {
		t.Errorf("Consume() error = %v, want nil", consumeErr)
	}
}
