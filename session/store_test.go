package session

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestMemoryStore_Create(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	sess, err := store.Create(ctx, "client-123")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if sess.ClientID != "client-123" {
		t.Errorf("ClientID = %q, want %q", sess.ClientID, "client-123")
	}
	if sess.State == nil {
		t.Error("State is nil, want initialized map")
	}
}

func TestMemoryStore_Create_GeneratesID(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	sess1, _ := store.Create(ctx, "client-1")
	sess2, _ := store.Create(ctx, "client-2")

	if sess1.ID == "" {
		t.Error("ID is empty")
	}
	if sess1.ID == sess2.ID {
		t.Error("generated IDs should be unique")
	}
}

func TestMemoryStore_Create_SetsTimestamps(t *testing.T) {
	store := NewMemoryStore(WithTTL(time.Hour))
	ctx := context.Background()

	before := time.Now()
	sess, _ := store.Create(ctx, "client-123")
	after := time.Now()

	if sess.CreatedAt.Before(before) || sess.CreatedAt.After(after) {
		t.Errorf("CreatedAt = %v, want between %v and %v", sess.CreatedAt, before, after)
	}
	if sess.UpdatedAt.Before(before) || sess.UpdatedAt.After(after) {
		t.Errorf("UpdatedAt = %v, want between %v and %v", sess.UpdatedAt, before, after)
	}
	if sess.ExpiresAt.Before(after.Add(time.Hour - time.Second)) {
		t.Errorf("ExpiresAt = %v, want ~1 hour from now", sess.ExpiresAt)
	}
}

func TestMemoryStore_Create_InvalidClientID(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	_, err := store.Create(ctx, "")
	if err != ErrInvalidClientID {
		t.Errorf("Create() error = %v, want ErrInvalidClientID", err)
	}
}

func TestMemoryStore_Get(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	created, _ := store.Create(ctx, "client-123")
	retrieved, err := store.Get(ctx, created.ID)

	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if retrieved.ID != created.ID {
		t.Errorf("ID = %q, want %q", retrieved.ID, created.ID)
	}
	if retrieved.ClientID != created.ClientID {
		t.Errorf("ClientID = %q, want %q", retrieved.ClientID, created.ClientID)
	}
}

func TestMemoryStore_Get_NotFound(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	_, err := store.Get(ctx, "nonexistent")
	if err != ErrSessionNotFound {
		t.Errorf("Get() error = %v, want ErrSessionNotFound", err)
	}
}

func TestMemoryStore_Get_Expired(t *testing.T) {
	store := NewMemoryStore(WithTTL(time.Millisecond))
	ctx := context.Background()

	sess, _ := store.Create(ctx, "client-123")
	time.Sleep(5 * time.Millisecond)

	_, err := store.Get(ctx, sess.ID)
	if err != ErrSessionExpired {
		t.Errorf("Get() error = %v, want ErrSessionExpired", err)
	}
}

func TestMemoryStore_Update(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	sess, _ := store.Create(ctx, "client-123")
	sess.State["key"] = "value"

	err := store.Update(ctx, sess)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	retrieved, _ := store.Get(ctx, sess.ID)
	if retrieved.State["key"] != "value" {
		t.Errorf("State[key] = %v, want %q", retrieved.State["key"], "value")
	}
}

func TestMemoryStore_Update_NotFound(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	sess := &Session{ID: "nonexistent"}
	err := store.Update(ctx, sess)
	if err != ErrSessionNotFound {
		t.Errorf("Update() error = %v, want ErrSessionNotFound", err)
	}
}

func TestMemoryStore_Update_RefreshesUpdatedAt(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	sess, _ := store.Create(ctx, "client-123")
	originalUpdatedAt := sess.UpdatedAt

	time.Sleep(time.Millisecond)
	sess.State["key"] = "value"
	_ = store.Update(ctx, sess)

	retrieved, _ := store.Get(ctx, sess.ID)
	if !retrieved.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("UpdatedAt should be refreshed on update")
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	sess, _ := store.Create(ctx, "client-123")
	err := store.Delete(ctx, sess.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = store.Get(ctx, sess.ID)
	if err != ErrSessionNotFound {
		t.Errorf("Get() after Delete() error = %v, want ErrSessionNotFound", err)
	}
}

func TestMemoryStore_Delete_NotFound(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	err := store.Delete(ctx, "nonexistent")
	if err != ErrSessionNotFound {
		t.Errorf("Delete() error = %v, want ErrSessionNotFound", err)
	}
}

func TestMemoryStore_Cleanup(t *testing.T) {
	store := NewMemoryStore(WithTTL(time.Millisecond))
	ctx := context.Background()

	sess1, _ := store.Create(ctx, "client-1")
	time.Sleep(5 * time.Millisecond)

	// Create a fresh session that won't be expired
	store2 := NewMemoryStore(WithTTL(time.Hour))
	sess2, _ := store2.Create(ctx, "client-2")

	// Copy sess2 to original store for the test
	store.mu.Lock()
	store.sessions[sess2.ID] = sess2
	store.mu.Unlock()

	err := store.Cleanup(ctx)
	if err != nil {
		t.Fatalf("Cleanup() error = %v", err)
	}

	// Expired session should be removed
	_, err = store.Get(ctx, sess1.ID)
	if err != ErrSessionNotFound {
		t.Errorf("expired session should be removed, got error = %v", err)
	}

	// Non-expired session should remain
	_, err = store.Get(ctx, sess2.ID)
	if err != nil {
		t.Errorf("non-expired session should remain, got error = %v", err)
	}
}

func TestMemoryStore_Cleanup_RemovesExpired(t *testing.T) {
	store := NewMemoryStore(WithTTL(time.Millisecond))
	ctx := context.Background()

	// Create multiple sessions
	for i := 0; i < 5; i++ {
		store.Create(ctx, "client")
	}

	time.Sleep(5 * time.Millisecond)

	err := store.Cleanup(ctx)
	if err != nil {
		t.Fatalf("Cleanup() error = %v", err)
	}

	store.mu.RLock()
	count := len(store.sessions)
	store.mu.RUnlock()

	if count != 0 {
		t.Errorf("Cleanup should remove all expired sessions, got %d remaining", count)
	}
}

func TestMemoryStore_ConcurrentSafety(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	var wg sync.WaitGroup
	sessions := make(chan *Session, 100)

	// Concurrent creates
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sess, err := store.Create(ctx, "client")
			if err != nil {
				t.Errorf("Create() error = %v", err)
				return
			}
			sessions <- sess
		}(i)
	}

	wg.Wait()
	close(sessions)

	// Concurrent reads and updates
	var wg2 sync.WaitGroup
	for sess := range sessions {
		wg2.Add(2)
		go func(s *Session) {
			defer wg2.Done()
			_, _ = store.Get(ctx, s.ID)
		}(sess)
		go func(s *Session) {
			defer wg2.Done()
			s.State["test"] = true
			_ = store.Update(ctx, s)
		}(sess)
	}
	wg2.Wait()
}
