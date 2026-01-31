package session

import (
	"context"
	"testing"
	"time"
)

func TestWithTTL(t *testing.T) {
	ttl := 30 * time.Minute
	store := NewMemoryStore(WithTTL(ttl))

	if store.ttl != ttl {
		t.Errorf("ttl = %v, want %v", store.ttl, ttl)
	}
}

func TestWithIDGenerator(t *testing.T) {
	counter := 0
	gen := func() string {
		counter++
		return "custom-id"
	}

	store := NewMemoryStore(WithIDGenerator(gen))
	ctx := context.Background()

	sess, _ := store.Create(ctx, "client")
	if sess.ID != "custom-id" {
		t.Errorf("ID = %q, want %q", sess.ID, "custom-id")
	}
	if counter != 1 {
		t.Errorf("generator called %d times, want 1", counter)
	}
}

func TestNewMemoryStore_Defaults(t *testing.T) {
	store := NewMemoryStore()

	// Default TTL should be 1 hour
	if store.ttl != time.Hour {
		t.Errorf("default ttl = %v, want %v", store.ttl, time.Hour)
	}

	// Default ID generator should produce non-empty strings
	id := store.idGen()
	if id == "" {
		t.Error("default ID generator produced empty string")
	}

	// Sessions map should be initialized
	if store.sessions == nil {
		t.Error("sessions map not initialized")
	}
}
