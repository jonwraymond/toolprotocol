package session

import "time"

// Option configures a MemoryStore.
type Option func(*MemoryStore)

// WithTTL configures the session time-to-live.
func WithTTL(ttl time.Duration) Option {
	return func(m *MemoryStore) {
		m.ttl = ttl
	}
}

// WithIDGenerator configures a custom session ID generator.
func WithIDGenerator(gen func() string) Option {
	return func(m *MemoryStore) {
		m.idGen = gen
	}
}
