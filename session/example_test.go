package session_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonwraymond/toolprotocol/session"
)

func ExampleNewMemoryStore() {
	store := session.NewMemoryStore()

	ctx := context.Background()
	sess, err := store.Create(ctx, "client-123")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Session created:", sess.ID != "")
	fmt.Println("Client ID:", sess.ClientID)
	// Output:
	// Session created: true
	// Client ID: client-123
}

func ExampleNewMemoryStore_withTTL() {
	// Create store with custom TTL
	store := session.NewMemoryStore(
		session.WithTTL(30 * time.Minute),
	)

	ctx := context.Background()
	sess, _ := store.Create(ctx, "client-123")

	// Session expires in 30 minutes
	expiresIn := time.Until(sess.ExpiresAt)
	fmt.Println("Expires in ~30min:", expiresIn > 29*time.Minute && expiresIn <= 30*time.Minute)
	// Output:
	// Expires in ~30min: true
}

func ExampleNewMemoryStore_withIDGenerator() {
	counter := 0
	// Create store with custom ID generator
	store := session.NewMemoryStore(
		session.WithIDGenerator(func() string {
			counter++
			return fmt.Sprintf("sess-%d", counter)
		}),
	)

	ctx := context.Background()
	sess, _ := store.Create(ctx, "client-123")

	fmt.Println("Custom ID:", sess.ID)
	// Output:
	// Custom ID: sess-1
}

func ExampleMemoryStore_Create() {
	store := session.NewMemoryStore()
	ctx := context.Background()

	sess, err := store.Create(ctx, "client-123")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("ID set:", sess.ID != "")
	fmt.Println("ClientID:", sess.ClientID)
	fmt.Println("State initialized:", sess.State != nil)
	fmt.Println("Not expired:", !sess.IsExpired())
	// Output:
	// ID set: true
	// ClientID: client-123
	// State initialized: true
	// Not expired: true
}

func ExampleMemoryStore_Create_emptyClientID() {
	store := session.NewMemoryStore()
	ctx := context.Background()

	_, err := store.Create(ctx, "")
	fmt.Println("Error is ErrInvalidClientID:", errors.Is(err, session.ErrInvalidClientID))
	// Output:
	// Error is ErrInvalidClientID: true
}

func ExampleMemoryStore_Get() {
	store := session.NewMemoryStore()
	ctx := context.Background()

	// Create a session
	created, _ := store.Create(ctx, "client-123")

	// Retrieve it
	retrieved, err := store.Get(ctx, created.ID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Same ID:", retrieved.ID == created.ID)
	fmt.Println("Same ClientID:", retrieved.ClientID == created.ClientID)
	// Output:
	// Same ID: true
	// Same ClientID: true
}

func ExampleMemoryStore_Get_notFound() {
	store := session.NewMemoryStore()
	ctx := context.Background()

	_, err := store.Get(ctx, "nonexistent-id")
	fmt.Println("Error is ErrSessionNotFound:", errors.Is(err, session.ErrSessionNotFound))
	// Output:
	// Error is ErrSessionNotFound: true
}

func ExampleMemoryStore_Update() {
	store := session.NewMemoryStore()
	ctx := context.Background()

	// Create a session
	sess, _ := store.Create(ctx, "client-123")

	// Modify and update
	sess.SetState("user", "alice")
	err := store.Update(ctx, sess)
	fmt.Println("Update error:", err)

	// Retrieve and verify
	updated, _ := store.Get(ctx, sess.ID)
	user, _ := updated.GetState("user")
	fmt.Println("User state:", user)
	// Output:
	// Update error: <nil>
	// User state: alice
}

func ExampleMemoryStore_Update_notFound() {
	store := session.NewMemoryStore()
	ctx := context.Background()

	// Try to update non-existent session
	sess := &session.Session{ID: "nonexistent"}
	err := store.Update(ctx, sess)
	fmt.Println("Error is ErrSessionNotFound:", errors.Is(err, session.ErrSessionNotFound))
	// Output:
	// Error is ErrSessionNotFound: true
}

func ExampleMemoryStore_Delete() {
	store := session.NewMemoryStore()
	ctx := context.Background()

	// Create and then delete
	sess, _ := store.Create(ctx, "client-123")

	err := store.Delete(ctx, sess.ID)
	fmt.Println("Delete error:", err)

	// Verify deleted
	_, err = store.Get(ctx, sess.ID)
	fmt.Println("After delete, error is ErrSessionNotFound:", errors.Is(err, session.ErrSessionNotFound))
	// Output:
	// Delete error: <nil>
	// After delete, error is ErrSessionNotFound: true
}

func ExampleMemoryStore_Delete_notFound() {
	store := session.NewMemoryStore()
	ctx := context.Background()

	err := store.Delete(ctx, "nonexistent-id")
	fmt.Println("Error is ErrSessionNotFound:", errors.Is(err, session.ErrSessionNotFound))
	// Output:
	// Error is ErrSessionNotFound: true
}

func ExampleMemoryStore_Cleanup() {
	// Create store with very short TTL
	store := session.NewMemoryStore(
		session.WithTTL(1 * time.Millisecond),
	)
	ctx := context.Background()

	// Create a session
	sess, _ := store.Create(ctx, "client-123")

	// Wait for expiration
	time.Sleep(5 * time.Millisecond)

	// Cleanup expired sessions
	err := store.Cleanup(ctx)
	fmt.Println("Cleanup error:", err)

	// Session should be cleaned up
	_, err = store.Get(ctx, sess.ID)
	fmt.Println("After cleanup, error is ErrSessionNotFound:", errors.Is(err, session.ErrSessionNotFound))
	// Output:
	// Cleanup error: <nil>
	// After cleanup, error is ErrSessionNotFound: true
}

func ExampleSession_Clone() {
	sess := &session.Session{
		ID:       "sess-123",
		ClientID: "client-123",
		State:    map[string]any{"key": "value"},
	}

	clone := sess.Clone()

	fmt.Println("Same ID:", clone.ID == sess.ID)
	fmt.Println("Same State key:", clone.State["key"] == sess.State["key"])

	// Modifying clone doesn't affect original
	clone.State["key"] = "modified"
	fmt.Println("Original unchanged:", sess.State["key"] == "value")
	// Output:
	// Same ID: true
	// Same State key: true
	// Original unchanged: true
}

func ExampleSession_IsExpired() {
	// Not expired
	sess1 := &session.Session{
		ExpiresAt: time.Now().Add(time.Hour),
	}
	fmt.Println("Future expiry, IsExpired:", sess1.IsExpired())

	// Already expired
	sess2 := &session.Session{
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	fmt.Println("Past expiry, IsExpired:", sess2.IsExpired())
	// Output:
	// Future expiry, IsExpired: false
	// Past expiry, IsExpired: true
}

func ExampleSession_SetState() {
	sess := &session.Session{}

	// SetState initializes State map if nil
	sess.SetState("user", "alice")
	sess.SetState("role", "admin")

	fmt.Println("User:", sess.State["user"])
	fmt.Println("Role:", sess.State["role"])
	// Output:
	// User: alice
	// Role: admin
}

func ExampleSession_GetState() {
	sess := &session.Session{
		State: map[string]any{
			"user": "alice",
		},
	}

	// Get existing key
	user, ok := sess.GetState("user")
	fmt.Println("User:", user, "found:", ok)

	// Get non-existent key
	role, ok := sess.GetState("role")
	fmt.Println("Role:", role, "found:", ok)
	// Output:
	// User: alice found: true
	// Role: <nil> found: false
}

func ExampleSession_GetState_nilState() {
	sess := &session.Session{} // State is nil

	value, ok := sess.GetState("key")
	fmt.Println("Value:", value, "found:", ok)
	// Output:
	// Value: <nil> found: false
}

func ExampleWithSession() {
	sess := &session.Session{
		ID:       "sess-123",
		ClientID: "client-123",
	}

	ctx := context.Background()
	ctx = session.WithSession(ctx, sess)

	// Session is now available in context
	retrieved, ok := session.FromContext(ctx)
	fmt.Println("Found in context:", ok)
	fmt.Println("Session ID:", retrieved.ID)
	// Output:
	// Found in context: true
	// Session ID: sess-123
}

func ExampleFromContext() {
	sess := &session.Session{ID: "sess-123"}
	ctx := session.WithSession(context.Background(), sess)

	// Retrieve session
	retrieved, ok := session.FromContext(ctx)
	fmt.Println("Found:", ok)
	fmt.Println("ID:", retrieved.ID)
	// Output:
	// Found: true
	// ID: sess-123
}

func ExampleFromContext_noSession() {
	ctx := context.Background() // No session attached

	retrieved, ok := session.FromContext(ctx)
	fmt.Println("Found:", ok)
	fmt.Println("Session is nil:", retrieved == nil)
	// Output:
	// Found: false
	// Session is nil: true
}

func ExampleMustFromContext() {
	sess := &session.Session{ID: "sess-123"}
	ctx := session.WithSession(context.Background(), sess)

	// MustFromContext returns session directly
	retrieved := session.MustFromContext(ctx)
	fmt.Println("ID:", retrieved.ID)
	// Output:
	// ID: sess-123
}

func ExampleWithTTL() {
	store := session.NewMemoryStore(
		session.WithTTL(2 * time.Hour),
	)

	ctx := context.Background()
	sess, _ := store.Create(ctx, "client-123")

	// Check TTL was applied
	expiresIn := time.Until(sess.ExpiresAt)
	fmt.Println("Expires in ~2hr:", expiresIn > time.Hour && expiresIn <= 2*time.Hour)
	// Output:
	// Expires in ~2hr: true
}

func ExampleWithIDGenerator() {
	counter := 0
	store := session.NewMemoryStore(
		session.WithIDGenerator(func() string {
			counter++
			return fmt.Sprintf("custom-id-%d", counter)
		}),
	)

	ctx := context.Background()
	sess1, _ := store.Create(ctx, "client-1")
	sess2, _ := store.Create(ctx, "client-2")

	fmt.Println("Session 1 ID:", sess1.ID)
	fmt.Println("Session 2 ID:", sess2.ID)
	// Output:
	// Session 1 ID: custom-id-1
	// Session 2 ID: custom-id-2
}

func ExampleSessionError() {
	err := &session.SessionError{
		SessionID: "sess-123",
		Op:        "update",
		Err:       session.ErrSessionNotFound,
	}

	fmt.Println(err.Error())
	fmt.Println("Unwraps to ErrSessionNotFound:", errors.Is(err, session.ErrSessionNotFound))
	// Output:
	// session sess-123: update: session: not found
	// Unwraps to ErrSessionNotFound: true
}

func ExampleSessionError_noUnderlying() {
	err := &session.SessionError{
		SessionID: "sess-123",
		Op:        "validate",
	}

	fmt.Println(err.Error())
	// Output:
	// session sess-123: validate
}
