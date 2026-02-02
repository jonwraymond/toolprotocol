package session

import (
	"context"
	"time"
)

// Session represents a client session.
type Session struct {
	// ID is the unique session identifier.
	ID string

	// ClientID identifies the client that owns this session.
	ClientID string

	// State holds arbitrary session data.
	State map[string]any

	// CreatedAt is when the session was created.
	CreatedAt time.Time

	// UpdatedAt is when the session was last updated.
	UpdatedAt time.Time

	// ExpiresAt is when the session expires.
	ExpiresAt time.Time
}

// Clone returns a deep copy of the session.
func (s *Session) Clone() *Session {
	clone := &Session{
		ID:        s.ID,
		ClientID:  s.ClientID,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		ExpiresAt: s.ExpiresAt,
	}
	if s.State != nil {
		clone.State = make(map[string]any, len(s.State))
		for k, v := range s.State {
			clone.State[k] = v
		}
	}
	return clone
}

// IsExpired returns true if the session has expired.
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// SetState sets a value in the session state.
// Initializes the state map if nil.
func (s *Session) SetState(key string, value any) {
	if s.State == nil {
		s.State = make(map[string]any)
	}
	s.State[key] = value
}

// GetState retrieves a value from the session state.
// Returns nil, false if the key is not found or state is nil.
func (s *Session) GetState(key string) (any, bool) {
	if s.State == nil {
		return nil, false
	}
	v, ok := s.State[key]
	return v, ok
}

// Store manages session persistence.
//
// Contract:
//   - Concurrency: Implementations must be safe for concurrent use.
//   - Context: All methods should honor context cancellation.
//   - Errors: Returns ErrSessionNotFound, ErrSessionExpired, ErrInvalidClientID
//     as appropriate. Use errors.Is for checking.
//   - Ownership: Returned *Session is a copy; modifications require Update().
//   - Cleanup: Callers should periodically call Cleanup to remove expired sessions.
type Store interface {
	// Create creates a new session for the given client ID.
	Create(ctx context.Context, clientID string) (*Session, error)

	// Get retrieves a session by ID.
	// Returns ErrSessionNotFound if the session does not exist.
	// Returns ErrSessionExpired if the session has expired.
	Get(ctx context.Context, id string) (*Session, error)

	// Update updates an existing session.
	// Returns ErrSessionNotFound if the session does not exist.
	Update(ctx context.Context, session *Session) error

	// Delete removes a session by ID.
	// Returns ErrSessionNotFound if the session does not exist.
	Delete(ctx context.Context, id string) error

	// Cleanup removes all expired sessions.
	Cleanup(ctx context.Context) error
}
