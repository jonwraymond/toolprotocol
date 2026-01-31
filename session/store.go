package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// MemoryStore is an in-memory implementation of Store.
type MemoryStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	ttl      time.Duration
	idGen    func() string
}

// NewMemoryStore creates a new in-memory session store.
func NewMemoryStore(opts ...Option) *MemoryStore {
	m := &MemoryStore{
		sessions: make(map[string]*Session),
		ttl:      time.Hour, // default TTL
		idGen:    defaultIDGenerator,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// defaultIDGenerator generates a random session ID.
func defaultIDGenerator() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Create creates a new session for the given client ID.
func (m *MemoryStore) Create(ctx context.Context, clientID string) (*Session, error) {
	if clientID == "" {
		return nil, ErrInvalidClientID
	}

	now := time.Now()
	s := &Session{
		ID:        m.idGen(),
		ClientID:  clientID,
		State:     make(map[string]any),
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(m.ttl),
	}

	m.mu.Lock()
	m.sessions[s.ID] = s
	m.mu.Unlock()

	return s.Clone(), nil
}

// Get retrieves a session by ID.
func (m *MemoryStore) Get(ctx context.Context, id string) (*Session, error) {
	m.mu.RLock()
	s, ok := m.sessions[id]
	m.mu.RUnlock()

	if !ok {
		return nil, ErrSessionNotFound
	}

	if s.IsExpired() {
		return nil, ErrSessionExpired
	}

	return s.Clone(), nil
}

// Update updates an existing session.
func (m *MemoryStore) Update(ctx context.Context, session *Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[session.ID]; !ok {
		return ErrSessionNotFound
	}

	session.UpdatedAt = time.Now()
	m.sessions[session.ID] = session.Clone()
	return nil
}

// Delete removes a session by ID.
func (m *MemoryStore) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[id]; !ok {
		return ErrSessionNotFound
	}

	delete(m.sessions, id)
	return nil
}

// Cleanup removes all expired sessions.
func (m *MemoryStore) Cleanup(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, s := range m.sessions {
		if s.IsExpired() {
			delete(m.sessions, id)
		}
	}
	return nil
}
