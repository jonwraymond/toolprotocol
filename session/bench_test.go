package session

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// BenchmarkMemoryStore_Create measures session creation performance.
func BenchmarkMemoryStore_Create(b *testing.B) {
	store := NewMemoryStore()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Create(ctx, fmt.Sprintf("client-%d", i))
	}
}

// BenchmarkMemoryStore_Create_SameClient measures creation for same client.
func BenchmarkMemoryStore_Create_SameClient(b *testing.B) {
	store := NewMemoryStore()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Create(ctx, "client-123")
	}
}

// BenchmarkMemoryStore_Get measures session retrieval performance.
func BenchmarkMemoryStore_Get(b *testing.B) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Create a session
	sess, _ := store.Create(ctx, "client-123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Get(ctx, sess.ID)
	}
}

// BenchmarkMemoryStore_Get_Miss measures miss performance.
func BenchmarkMemoryStore_Get_Miss(b *testing.B) {
	store := NewMemoryStore()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Get(ctx, "nonexistent")
	}
}

// BenchmarkMemoryStore_Update measures session update performance.
func BenchmarkMemoryStore_Update(b *testing.B) {
	store := NewMemoryStore()
	ctx := context.Background()

	sess, _ := store.Create(ctx, "client-123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sess.State["counter"] = i
		_ = store.Update(ctx, sess)
	}
}

// BenchmarkMemoryStore_Delete measures session deletion performance.
func BenchmarkMemoryStore_Delete(b *testing.B) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Pre-create sessions
	ids := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		sess, _ := store.Create(ctx, fmt.Sprintf("client-%d", i))
		ids[i] = sess.ID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.Delete(ctx, ids[i])
	}
}

// BenchmarkMemoryStore_Cleanup measures cleanup performance.
func BenchmarkMemoryStore_Cleanup(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		store := NewMemoryStore(WithTTL(time.Nanosecond))
		// Create some expired sessions
		for j := 0; j < 100; j++ {
			_, _ = store.Create(ctx, fmt.Sprintf("client-%d", j))
		}
		time.Sleep(time.Microsecond) // Ensure expiration
		b.StartTimer()

		_ = store.Cleanup(ctx)
	}
}

// BenchmarkMemoryStore_Concurrent measures concurrent operations.
func BenchmarkMemoryStore_Concurrent(b *testing.B) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Pre-create some sessions
	for i := 0; i < 100; i++ {
		_, _ = store.Create(ctx, fmt.Sprintf("client-%d", i))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			switch i % 4 {
			case 0:
				_, _ = store.Create(ctx, fmt.Sprintf("new-client-%d", i))
			case 1:
				sess, err := store.Get(ctx, fmt.Sprintf("client-%d", i%100))
				if err == nil {
					_ = store.Update(ctx, sess)
				}
			case 2:
				_, _ = store.Get(ctx, fmt.Sprintf("client-%d", i%100))
			case 3:
				_ = store.Cleanup(ctx)
			}
			i++
		}
	})
}

// BenchmarkMemoryStore_ConcurrentReadHeavy measures read-heavy workload.
func BenchmarkMemoryStore_ConcurrentReadHeavy(b *testing.B) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Pre-create sessions
	sessions := make([]*Session, 100)
	for i := 0; i < 100; i++ {
		sess, _ := store.Create(ctx, fmt.Sprintf("client-%d", i))
		sessions[i] = sess
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = store.Get(ctx, sessions[i%100].ID)
			i++
		}
	})
}

// BenchmarkSession_Clone measures clone performance.
func BenchmarkSession_Clone(b *testing.B) {
	sess := &Session{
		ID:       "sess-123",
		ClientID: "client-123",
		State: map[string]any{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sess.Clone()
	}
}

// BenchmarkSession_Clone_LargeState measures clone with large state.
func BenchmarkSession_Clone_LargeState(b *testing.B) {
	state := make(map[string]any, 100)
	for i := 0; i < 100; i++ {
		state[fmt.Sprintf("key-%d", i)] = fmt.Sprintf("value-%d", i)
	}

	sess := &Session{
		ID:       "sess-123",
		ClientID: "client-123",
		State:    state,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sess.Clone()
	}
}

// BenchmarkSession_SetState measures state mutation performance.
func BenchmarkSession_SetState(b *testing.B) {
	sess := &Session{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sess.SetState(fmt.Sprintf("key-%d", i%10), i)
	}
}

// BenchmarkSession_GetState measures state retrieval performance.
func BenchmarkSession_GetState(b *testing.B) {
	sess := &Session{
		State: map[string]any{
			"key": "value",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = sess.GetState("key")
	}
}

// BenchmarkSession_IsExpired measures expiration check performance.
func BenchmarkSession_IsExpired(b *testing.B) {
	sess := &Session{
		ExpiresAt: time.Now().Add(time.Hour),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sess.IsExpired()
	}
}

// BenchmarkWithSession measures context attachment performance.
func BenchmarkWithSession(b *testing.B) {
	ctx := context.Background()
	sess := &Session{ID: "sess-123"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WithSession(ctx, sess)
	}
}

// BenchmarkFromContext measures context retrieval performance.
func BenchmarkFromContext(b *testing.B) {
	sess := &Session{ID: "sess-123"}
	ctx := WithSession(context.Background(), sess)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromContext(ctx)
	}
}

// BenchmarkDefaultIDGenerator measures ID generation performance.
func BenchmarkDefaultIDGenerator(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = defaultIDGenerator()
	}
}
