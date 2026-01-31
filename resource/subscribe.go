package resource

import (
	"context"
	"sync"
)

// SubscriptionManager manages resource subscriptions.
type SubscriptionManager struct {
	mu      sync.RWMutex
	subs    map[string][]chan *Contents
	bufSize int
}

// NewSubscriptionManager creates a new SubscriptionManager.
func NewSubscriptionManager(opts ...SubscriptionOption) *SubscriptionManager {
	m := &SubscriptionManager{
		subs:    make(map[string][]chan *Contents),
		bufSize: 10,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// SubscriptionOption configures a SubscriptionManager.
type SubscriptionOption func(*SubscriptionManager)

// WithBufferSize sets the channel buffer size for subscriptions.
func WithBufferSize(size int) SubscriptionOption {
	return func(m *SubscriptionManager) {
		if size > 0 {
			m.bufSize = size
		}
	}
}

// Subscribe subscribes to updates for a resource.
func (m *SubscriptionManager) Subscribe(ctx context.Context, uri string) (<-chan *Contents, error) {
	if uri == "" {
		return nil, &ResourceError{
			URI: uri,
			Op:  "subscribe",
			Err: ErrInvalidURI,
		}
	}

	ch := make(chan *Contents, m.bufSize)

	m.mu.Lock()
	m.subs[uri] = append(m.subs[uri], ch)
	m.mu.Unlock()

	// Handle context cancellation
	go func() {
		<-ctx.Done()
		m.removeSubscription(uri, ch)
		close(ch)
	}()

	return ch, nil
}

// Unsubscribe unsubscribes from a resource.
func (m *SubscriptionManager) Unsubscribe(ctx context.Context, uri string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	subs, ok := m.subs[uri]
	if !ok || len(subs) == 0 {
		return &ResourceError{
			URI: uri,
			Op:  "unsubscribe",
			Err: ErrNotSubscribed,
		}
	}

	// Close all channels for this URI
	for _, ch := range subs {
		close(ch)
	}
	delete(m.subs, uri)

	return nil
}

// Notify sends contents to all subscribers of a resource.
func (m *SubscriptionManager) Notify(uri string, contents *Contents) {
	m.mu.RLock()
	subs := m.subs[uri]
	m.mu.RUnlock()

	for _, ch := range subs {
		select {
		case ch <- contents:
		default:
			// Skip if channel is full (non-blocking)
		}
	}
}

// removeSubscription removes a specific subscription channel.
func (m *SubscriptionManager) removeSubscription(uri string, target chan *Contents) {
	m.mu.Lock()
	defer m.mu.Unlock()

	subs := m.subs[uri]
	for i, ch := range subs {
		if ch == target {
			m.subs[uri] = append(subs[:i], subs[i+1:]...)
			break
		}
	}

	if len(m.subs[uri]) == 0 {
		delete(m.subs, uri)
	}
}
