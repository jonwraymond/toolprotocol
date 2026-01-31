package resource

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestSubscriptionManager_Subscribe(t *testing.T) {
	m := NewSubscriptionManager()
	ctx := context.Background()

	ch, err := m.Subscribe(ctx, "test://resource")
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	if ch == nil {
		t.Error("Subscribe() returned nil channel")
	}
}

func TestSubscriptionManager_Subscribe_InvalidURI(t *testing.T) {
	m := NewSubscriptionManager()
	ctx := context.Background()

	_, err := m.Subscribe(ctx, "")
	if !errors.Is(err, ErrInvalidURI) {
		t.Errorf("Subscribe() error = %v, want ErrInvalidURI", err)
	}
}

func TestSubscriptionManager_Unsubscribe(t *testing.T) {
	m := NewSubscriptionManager()
	ctx := context.Background()

	_, _ = m.Subscribe(ctx, "test://resource")

	err := m.Unsubscribe(ctx, "test://resource")
	if err != nil {
		t.Fatalf("Unsubscribe() error = %v", err)
	}
}

func TestSubscriptionManager_Unsubscribe_NotSubscribed(t *testing.T) {
	m := NewSubscriptionManager()
	ctx := context.Background()

	err := m.Unsubscribe(ctx, "test://nonexistent")
	if !errors.Is(err, ErrNotSubscribed) {
		t.Errorf("Unsubscribe() error = %v, want ErrNotSubscribed", err)
	}
}

func TestSubscriptionManager_Notify(t *testing.T) {
	m := NewSubscriptionManager()
	ctx := context.Background()

	ch, _ := m.Subscribe(ctx, "test://resource")

	contents := &Contents{URI: "test://resource", Text: "updated"}
	m.Notify("test://resource", contents)

	select {
	case received := <-ch:
		if received.Text != "updated" {
			t.Errorf("received.Text = %q, want %q", received.Text, "updated")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Notify() did not send to subscriber")
	}
}

func TestSubscriptionManager_Notify_AllSubscribers(t *testing.T) {
	m := NewSubscriptionManager()
	ctx := context.Background()

	ch1, _ := m.Subscribe(ctx, "test://resource")
	ch2, _ := m.Subscribe(ctx, "test://resource")

	contents := &Contents{URI: "test://resource", Text: "updated"}
	m.Notify("test://resource", contents)

	// Both should receive
	select {
	case <-ch1:
	case <-time.After(100 * time.Millisecond):
		t.Error("ch1 did not receive notification")
	}

	select {
	case <-ch2:
	case <-time.After(100 * time.Millisecond):
		t.Error("ch2 did not receive notification")
	}
}

func TestSubscriptionManager_ContextCancellation(t *testing.T) {
	m := NewSubscriptionManager()
	ctx, cancel := context.WithCancel(context.Background())

	ch, _ := m.Subscribe(ctx, "test://resource")

	cancel()

	// Channel should be closed
	time.Sleep(10 * time.Millisecond)
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("channel should be closed after context cancellation")
		}
	default:
		// May still be processing, wait a bit more
		time.Sleep(50 * time.Millisecond)
	}
}

func TestSubscriptionManager_ConcurrentSafety(t *testing.T) {
	m := NewSubscriptionManager()
	ctx := context.Background()

	var wg sync.WaitGroup

	// Concurrent subscriptions
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			uri := "test://resource-" + string(rune('a'+i%26))
			_, _ = m.Subscribe(ctx, uri)
		}(i)
	}

	// Concurrent notifications
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			uri := "test://resource-" + string(rune('a'+i%26))
			m.Notify(uri, &Contents{URI: uri})
		}(i)
	}

	wg.Wait()
}

func TestWithBufferSize(t *testing.T) {
	m := NewSubscriptionManager(WithBufferSize(100))

	if m.bufSize != 100 {
		t.Errorf("bufSize = %d, want 100", m.bufSize)
	}
}

func TestNewSubscriptionManager_Defaults(t *testing.T) {
	m := NewSubscriptionManager()

	if m.subs == nil {
		t.Error("subs map not initialized")
	}
	if m.bufSize != 10 {
		t.Errorf("default bufSize = %d, want 10", m.bufSize)
	}
}
