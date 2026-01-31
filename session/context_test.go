package session

import (
	"context"
	"testing"
	"time"
)

func TestWithSession(t *testing.T) {
	ctx := context.Background()
	sess := &Session{ID: "test-123"}

	newCtx := WithSession(ctx, sess)

	if newCtx == ctx {
		t.Error("WithSession should return a new context")
	}
}

func TestFromContext(t *testing.T) {
	sess := &Session{ID: "test-123", ClientID: "client-456"}
	ctx := WithSession(context.Background(), sess)

	retrieved, ok := FromContext(ctx)
	if !ok {
		t.Fatal("FromContext returned ok = false")
	}
	if retrieved.ID != sess.ID {
		t.Errorf("ID = %q, want %q", retrieved.ID, sess.ID)
	}
	if retrieved.ClientID != sess.ClientID {
		t.Errorf("ClientID = %q, want %q", retrieved.ClientID, sess.ClientID)
	}
}

func TestFromContext_NotFound(t *testing.T) {
	ctx := context.Background()

	sess, ok := FromContext(ctx)
	if ok {
		t.Error("FromContext should return ok = false for empty context")
	}
	if sess != nil {
		t.Error("FromContext should return nil session for empty context")
	}
}

func TestMustFromContext(t *testing.T) {
	sess := &Session{ID: "test-123"}
	ctx := WithSession(context.Background(), sess)

	retrieved := MustFromContext(ctx)
	if retrieved.ID != sess.ID {
		t.Errorf("ID = %q, want %q", retrieved.ID, sess.ID)
	}
}

func TestMustFromContext_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustFromContext should panic when no session in context")
		}
	}()

	ctx := context.Background()
	_ = MustFromContext(ctx)
}

func TestContext_RoundTrip(t *testing.T) {
	now := time.Now()
	original := &Session{
		ID:        "sess-123",
		ClientID:  "client-456",
		State:     map[string]any{"user": "alice"},
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(time.Hour),
	}

	ctx := WithSession(context.Background(), original)
	retrieved, ok := FromContext(ctx)

	if !ok {
		t.Fatal("FromContext returned ok = false")
	}
	if retrieved.ID != original.ID {
		t.Errorf("ID = %q, want %q", retrieved.ID, original.ID)
	}
	if retrieved.ClientID != original.ClientID {
		t.Errorf("ClientID = %q, want %q", retrieved.ClientID, original.ClientID)
	}
	if retrieved.State["user"] != original.State["user"] {
		t.Errorf("State[user] = %v, want %v", retrieved.State["user"], original.State["user"])
	}
}
