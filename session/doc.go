// Package session provides client session management for protocol servers.
//
// This package enables tracking of client sessions with state persistence,
// expiration handling, and context integration.
//
// # Session Lifecycle
//
// Sessions follow a simple lifecycle:
//
//   - Create: New session created with unique ID and TTL
//   - Get: Retrieve existing session by ID
//   - Update: Modify session state and refresh timestamp
//   - Delete: Explicitly remove session
//   - Cleanup: Remove expired sessions
//
// # Store Interface
//
// The Store interface provides session persistence:
//
//	type Store interface {
//	    Create(ctx context.Context, clientID string) (*Session, error)
//	    Get(ctx context.Context, id string) (*Session, error)
//	    Update(ctx context.Context, session *Session) error
//	    Delete(ctx context.Context, id string) error
//	    Cleanup(ctx context.Context) error
//	}
//
// # Context Integration
//
// Sessions can be attached to context for request-scoped access:
//
//	// Attach session to context
//	ctx = session.WithSession(ctx, sess)
//
//	// Retrieve session from context
//	sess, ok := session.FromContext(ctx)
//
//	// Retrieve session or panic
//	sess := session.MustFromContext(ctx)
//
// # Usage
//
//	// Create a store
//	store := session.NewMemoryStore()
//
//	// Create a session
//	sess, err := store.Create(ctx, "client-123")
//
//	// Set state
//	sess.State["user"] = "alice"
//	store.Update(ctx, sess)
//
//	// Retrieve session
//	sess, err := store.Get(ctx, sess.ID)
//
//	// Cleanup expired sessions
//	store.Cleanup(ctx)
package session
