// Package session provides client session management for protocol servers.
//
// This package enables tracking of client sessions with state persistence,
// expiration handling, and context integration.
//
// # Ecosystem Position
//
// session manages client state across protocol interactions:
//
//	┌─────────────────────────────────────────────────────────────────┐
//	│                    Session Management Flow                      │
//	├─────────────────────────────────────────────────────────────────┤
//	│                                                                 │
//	│   Transport              session               Application     │
//	│   ┌─────────┐         ┌───────────┐         ┌─────────┐       │
//	│   │ Client  │─────────│ Get/Create│─────────│ Handler │       │
//	│   │ Request │         │  Session  │         │         │       │
//	│   └─────────┘         │ ┌───────┐ │         └─────────┘       │
//	│        │              │ │ Store │ │              │              │
//	│        │              │ │(memory)│ │              │              │
//	│        │              │ └───────┘ │              │              │
//	│        │              │     │     │              │              │
//	│        │              │ ┌───────┐ │              │              │
//	│        │              │ │Session│◀───────────────┘              │
//	│        │              │ │ State │ │  update                     │
//	│        │              │ └───────┘ │                             │
//	│        │              └───────────┘                             │
//	│        │                    │                                   │
//	│        │   WithSession ┌────┴────┐                             │
//	│        └───────────────│ Context │─────────────────────────────│
//	│                        └─────────┘                             │
//	│                                                                 │
//	└─────────────────────────────────────────────────────────────────┘
//
// # Core Components
//
//   - [Session]: Client session with ID, state, and expiration
//   - [Store]: Interface for session persistence (Create/Get/Update/Delete)
//   - [MemoryStore]: Thread-safe in-memory Store implementation
//   - [WithSession]: Attach session to context
//   - [FromContext]: Retrieve session from context
//   - [Option]: Functional options for MemoryStore configuration
//
// # Session Lifecycle
//
// Sessions follow a simple lifecycle:
//
//   - Create: New session created with unique ID and TTL
//   - Get: Retrieve existing session by ID (returns copy)
//   - Update: Modify session state and persist changes
//   - Delete: Explicitly remove session
//   - Cleanup: Periodic removal of expired sessions
//
// # Quick Start
//
//	// Create a store with custom TTL
//	store := session.NewMemoryStore(
//	    session.WithTTL(30 * time.Minute),
//	)
//
//	// Create a session for a client
//	sess, err := store.Create(ctx, "client-123")
//	if err != nil {
//	    return err
//	}
//
//	// Store state in the session
//	sess.SetState("user", "alice")
//	sess.SetState("role", "admin")
//	err = store.Update(ctx, sess)
//
//	// Later, retrieve the session
//	sess, err = store.Get(ctx, sessID)
//	user, ok := sess.GetState("user")
//
//	// Attach to context for request handling
//	ctx = session.WithSession(ctx, sess)
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
//	if !ok {
//	    // No session in context
//	}
//
//	// Retrieve session or panic (use in middleware after validation)
//	sess = session.MustFromContext(ctx)
//
// # Thread Safety
//
// All exported types are safe for concurrent use:
//
//   - [MemoryStore]: sync.RWMutex protects all operations
//   - Get: Uses RLock for concurrent reads
//   - Create/Update/Delete/Cleanup: Uses Lock for exclusive access
//   - [Session]: Not thread-safe; use Store.Update() for safe mutations
//   - Context functions: Thread-safe (context.Context is immutable)
//
// # Error Handling
//
// Sentinel errors (use errors.Is for checking):
//
//   - [ErrSessionNotFound]: Session does not exist
//   - [ErrSessionExpired]: Session has expired
//   - [ErrInvalidClientID]: Client ID is empty or invalid
//
// The [SessionError] type wraps errors with session context:
//
//	err := &SessionError{
//	    SessionID: "sess-123",
//	    Op:        "update",
//	    Err:       ErrSessionNotFound,
//	}
//	// err.Error() = "session sess-123: update: session: not found"
//	// errors.Is(err, ErrSessionNotFound) = true
//
// # Configuration Options
//
// MemoryStore supports functional options:
//
//   - [WithTTL]: Configure session time-to-live (default: 1 hour)
//   - [WithIDGenerator]: Custom session ID generation
//
// # Integration with ApertureStack
//
// session integrates with other ApertureStack packages:
//
//   - transport: HTTP transports extract session IDs from headers
//   - task: Long-running tasks may store progress in session state
//   - stream: Streaming connections may be associated with sessions
package session
