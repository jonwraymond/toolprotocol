package stream

import (
	"context"
	"testing"
)

func TestSource_NewStream(t *testing.T) {
	source := NewSource()
	ctx := context.Background()

	s := source.NewStream(ctx)
	if s == nil {
		t.Fatal("NewStream() returned nil")
	}

	// Should be a DefaultStream
	if _, ok := s.(*DefaultStream); !ok {
		t.Errorf("NewStream() returned %T, want *DefaultStream", s)
	}

	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestSource_NewBufferedStream(t *testing.T) {
	source := NewSource()
	ctx := context.Background()

	s := source.NewBufferedStream(ctx, 50)
	if s == nil {
		t.Fatal("NewBufferedStream() returned nil")
	}

	// Should be a BufferedStream
	if _, ok := s.(*BufferedStream); !ok {
		t.Errorf("NewBufferedStream() returned %T, want *BufferedStream", s)
	}

	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestSource_NewBufferedStream_DefaultSize(t *testing.T) {
	source := NewSource()
	ctx := context.Background()

	s := source.NewBufferedStream(ctx, DefaultBufferSize)
	if s == nil {
		t.Fatal("NewBufferedStream() returned nil")
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestSource_WithBackpressure(t *testing.T) {
	source := NewSource(WithBackpressure(BackpressureDrop))
	ctx := context.Background()

	s := source.NewBufferedStream(ctx, 1).(*BufferedStream)
	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})

	// Fill buffer
	_ = s.Send(ctx, Event{Type: EventProgress})

	// Should drop instead of block
	err := s.Send(ctx, Event{Type: EventProgress})
	if err != ErrBufferFull {
		t.Errorf("Send() error = %v, want ErrBufferFull", err)
	}
}
