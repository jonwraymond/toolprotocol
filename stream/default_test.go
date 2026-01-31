package stream

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestDefaultStream_Send(t *testing.T) {
	s := newDefaultStream()
	ctx := context.Background()

	// Start a consumer
	var received Event
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		received = <-s.Events()
	}()

	event := Event{Type: EventProgress, Data: 0.5}
	err := s.Send(ctx, event)
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	wg.Wait()
	if received.Type != EventProgress {
		t.Errorf("received.Type = %v, want %v", received.Type, EventProgress)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestDefaultStream_Send_Closed(t *testing.T) {
	s := newDefaultStream()
	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	err := s.Send(context.Background(), Event{Type: EventProgress})
	if err != ErrStreamClosed {
		t.Errorf("Send() error = %v, want ErrStreamClosed", err)
	}
}

func TestDefaultStream_Send_ContextCancelled(t *testing.T) {
	s := newDefaultStream()
	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := s.Send(ctx, Event{Type: EventProgress})
	if err != context.Canceled {
		t.Errorf("Send() error = %v, want context.Canceled", err)
	}
}

func TestDefaultStream_Close(t *testing.T) {
	s := newDefaultStream()

	err := s.Close()
	if err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// Verify stream is closed
	select {
	case <-s.Done():
		// expected
	default:
		t.Error("Done() should be closed")
	}
}

func TestDefaultStream_Close_Idempotent(t *testing.T) {
	s := newDefaultStream()

	// Close multiple times
	for i := 0; i < 3; i++ {
		err := s.Close()
		if err != nil {
			t.Fatalf("Close() #%d error = %v", i+1, err)
		}
	}
}

func TestDefaultStream_Done(t *testing.T) {
	s := newDefaultStream()

	select {
	case <-s.Done():
		t.Error("Done() should not be closed initially")
	default:
		// expected
	}
}

func TestDefaultStream_Done_ClosedOnClose(t *testing.T) {
	s := newDefaultStream()
	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	select {
	case <-s.Done():
		// expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Done() should be closed after Close()")
	}
}

func TestDefaultStream_ConcurrentSafety(t *testing.T) {
	s := newDefaultStream()
	ctx := context.Background()

	var wg sync.WaitGroup

	// Start consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range s.Events() {
			// consume events
		}
	}()

	// Concurrent sends
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = s.Send(ctx, Event{Type: EventProgress, Data: i})
		}(i)
	}

	// Give time for some sends
	time.Sleep(10 * time.Millisecond)
	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	wg.Wait()
}
