package session

import "context"

// contextKey is used for context value storage.
type contextKey struct{}

// WithSession returns a new context with the session attached.
func WithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, contextKey{}, session)
}

// FromContext retrieves the session from the context.
// Returns nil, false if no session is attached.
func FromContext(ctx context.Context) (*Session, bool) {
	s, ok := ctx.Value(contextKey{}).(*Session)
	return s, ok
}

// MustFromContext retrieves the session from the context.
// Panics if no session is attached.
func MustFromContext(ctx context.Context) *Session {
	s, ok := FromContext(ctx)
	if !ok {
		panic("session: no session in context")
	}
	return s
}
