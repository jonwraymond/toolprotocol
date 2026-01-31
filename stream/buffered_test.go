package stream

import (
	"context"
	"testing"
	"time"
)

func TestBufferedStream_Send_Buffered(t *testing.T) {
	s := newBufferedStream(10, BackpressureBlock)
	ctx := context.Background()

	// Send events without consumer (buffered)
	for i := 0; i < 10; i++ {
		err := s.Send(ctx, Event{Type: EventProgress, Data: i})
		if err != nil {
			t.Fatalf("Send() #%d error = %v", i, err)
		}
	}

	// Verify events are buffered
	count := 0
	for range s.Events() {
		count++
		if count == 10 {
			if err := s.Close(); err != nil {
				t.Fatalf("Close() error = %v", err)
			}
		}
	}

	if count != 10 {
		t.Errorf("received %d events, want 10", count)
	}
}

func TestBufferedStream_Send_BackpressureBlock(t *testing.T) {
	s := newBufferedStream(1, BackpressureBlock)
	ctx := context.Background()

	// Fill buffer
	err := s.Send(ctx, Event{Type: EventProgress})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	// Next send should block
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err = s.Send(ctx, Event{Type: EventProgress})
	if err != context.DeadlineExceeded {
		t.Errorf("Send() error = %v, want context.DeadlineExceeded", err)
	}

	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestBufferedStream_Send_BackpressureDrop(t *testing.T) {
	s := newBufferedStream(1, BackpressureDrop)
	ctx := context.Background()

	// Fill buffer
	err := s.Send(ctx, Event{Type: EventProgress})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	// Next send should return ErrBufferFull
	err = s.Send(ctx, Event{Type: EventProgress})
	if err != ErrBufferFull {
		t.Errorf("Send() error = %v, want ErrBufferFull", err)
	}

	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestBufferedStream_Drain(t *testing.T) {
	s := newBufferedStream(10, BackpressureBlock)
	ctx := context.Background()

	// Fill buffer
	for i := 0; i < 5; i++ {
		_ = s.Send(ctx, Event{Type: EventProgress, Data: i})
	}

	// Drain
	s.Drain()

	// Buffer should be empty
	select {
	case <-s.Events():
		t.Error("buffer should be empty after Drain()")
	default:
		// expected
	}

	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestBufferedStream_Close_DrainsBuffer(t *testing.T) {
	s := newBufferedStream(10, BackpressureBlock)
	ctx := context.Background()

	// Fill buffer
	for i := 0; i < 5; i++ {
		_ = s.Send(ctx, Event{Type: EventProgress, Data: i})
	}

	// Close
	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// Should still be able to read remaining events
	count := 0
	for range s.Events() {
		count++
	}

	if count != 5 {
		t.Errorf("received %d events, want 5", count)
	}
}

func TestBufferedStream_DefaultSize(t *testing.T) {
	// Size less than 1 should use 1
	s := newBufferedStream(0, BackpressureBlock)
	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})

	ctx := context.Background()
	err := s.Send(ctx, Event{Type: EventProgress})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
}
