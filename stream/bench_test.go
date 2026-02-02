package stream

import (
	"context"
	"testing"
)

// BenchmarkDefaultStream_Send measures unbuffered stream send performance.
func BenchmarkDefaultStream_Send(b *testing.B) {
	ctx := context.Background()
	s := newDefaultStream()
	event := Event{Type: EventProgress, Data: 0.5}

	// Start consumer
	go func() {
		for range s.Events() {
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		_ = s.Send(ctx, event)
	}
	b.StopTimer()
	_ = s.Close()
}

// BenchmarkBufferedStream_Send measures buffered stream send performance.
func BenchmarkBufferedStream_Send(b *testing.B) {
	ctx := context.Background()
	s := newBufferedStream(1000, BackpressureBlock)
	event := Event{Type: EventProgress, Data: 0.5}

	// Start consumer
	go func() {
		for range s.Events() {
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		_ = s.Send(ctx, event)
	}
	b.StopTimer()
	_ = s.Close()
}

// BenchmarkBufferedStream_Send_Drop measures drop-mode send performance.
func BenchmarkBufferedStream_Send_Drop(b *testing.B) {
	ctx := context.Background()
	s := newBufferedStream(1000, BackpressureDrop)
	event := Event{Type: EventProgress, Data: 0.5}

	// Start consumer
	go func() {
		for range s.Events() {
		}
	}()

	b.ResetTimer()
	for b.Loop() {
		_ = s.Send(ctx, event)
	}
	b.StopTimer()
	_ = s.Close()
}

// BenchmarkDefaultStream_Close measures stream close performance.
func BenchmarkDefaultStream_Close(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		b.StopTimer()
		s := newDefaultStream()
		b.StartTimer()

		_ = s.Close()
	}
}

// BenchmarkBufferedStream_Close measures buffered stream close performance.
func BenchmarkBufferedStream_Close(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		b.StopTimer()
		s := newBufferedStream(100, BackpressureBlock)
		b.StartTimer()

		_ = s.Close()
	}
}

// BenchmarkDefaultSink_Consume measures sink consume performance.
func BenchmarkDefaultSink_Consume(b *testing.B) {
	ctx := context.Background()
	sink := NewSink()
	handler := func(event Event) error { return nil }

	b.ResetTimer()
	for b.Loop() {
		b.StopTimer()
		s := newBufferedStream(100, BackpressureBlock)
		// Pre-populate with events
		for range 50 {
			_ = s.Send(ctx, Event{Type: EventProgress})
		}
		_ = s.Close()
		b.StartTimer()

		_ = sink.Consume(ctx, s, handler)
	}
}

// BenchmarkEvent_Clone measures event clone performance.
func BenchmarkEvent_Clone(b *testing.B) {
	event := Event{
		Type:  EventProgress,
		ID:    "evt-123",
		Data:  0.5,
		Retry: 1000,
	}

	b.ResetTimer()
	for b.Loop() {
		_ = event.Clone()
	}
}

// BenchmarkEventType_Valid measures event type validation performance.
func BenchmarkEventType_Valid(b *testing.B) {
	types := []EventType{EventProgress, EventPartial, EventComplete, EventError, EventHeartbeat}

	b.ResetTimer()
	for b.Loop() {
		for _, t := range types {
			_ = t.Valid()
		}
	}
}

// BenchmarkEventType_String measures event type string conversion.
func BenchmarkEventType_String(b *testing.B) {
	et := EventProgress

	b.ResetTimer()
	for b.Loop() {
		_ = et.String()
	}
}

// BenchmarkNewSource measures source creation performance.
func BenchmarkNewSource(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewSource()
	}
}

// BenchmarkNewSource_WithOptions measures source creation with options.
func BenchmarkNewSource_WithOptions(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewSource(WithBackpressure(BackpressureDrop))
	}
}

// BenchmarkNewSink measures sink creation performance.
func BenchmarkNewSink(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewSink()
	}
}

// BenchmarkStream_Throughput measures end-to-end throughput.
func BenchmarkStream_Throughput(b *testing.B) {
	ctx := context.Background()
	source := NewSource()
	sink := NewSink()

	b.ResetTimer()
	for b.Loop() {
		b.StopTimer()
		s := source.NewBufferedStream(ctx, 1000)
		done := make(chan struct{})

		// Start consumer
		go func() {
			_ = sink.Consume(ctx, s, func(event Event) error {
				return nil
			})
			close(done)
		}()

		b.StartTimer()
		// Send events
		for range 100 {
			_ = s.Send(ctx, Event{Type: EventProgress, Data: 0.5})
		}
		_ = s.Close()
		<-done
	}
}

// BenchmarkStream_Concurrent measures concurrent send/receive.
func BenchmarkStream_Concurrent(b *testing.B) {
	ctx := context.Background()
	s := newBufferedStream(1000, BackpressureBlock)
	event := Event{Type: EventProgress, Data: 0.5}

	// Start consumer
	go func() {
		for range s.Events() {
		}
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = s.Send(ctx, event)
		}
	})
	b.StopTimer()
	_ = s.Close()
}
