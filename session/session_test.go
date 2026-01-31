package session

import (
	"testing"
	"time"
)

func TestSession_Fields(t *testing.T) {
	now := time.Now()
	expires := now.Add(time.Hour)

	s := &Session{
		ID:        "sess-123",
		ClientID:  "client-456",
		State:     map[string]any{"key": "value"},
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: expires,
	}

	if s.ID != "sess-123" {
		t.Errorf("ID = %q, want %q", s.ID, "sess-123")
	}
	if s.ClientID != "client-456" {
		t.Errorf("ClientID = %q, want %q", s.ClientID, "client-456")
	}
	if s.State["key"] != "value" {
		t.Errorf("State[key] = %v, want %q", s.State["key"], "value")
	}
	if !s.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", s.CreatedAt, now)
	}
	if !s.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", s.UpdatedAt, now)
	}
	if !s.ExpiresAt.Equal(expires) {
		t.Errorf("ExpiresAt = %v, want %v", s.ExpiresAt, expires)
	}
}

func TestSession_Clone(t *testing.T) {
	now := time.Now()
	expires := now.Add(time.Hour)

	original := &Session{
		ID:        "sess-123",
		ClientID:  "client-456",
		State:     map[string]any{"key": "value"},
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: expires,
	}

	clone := original.Clone()

	// Verify clone has same values
	if clone.ID != original.ID {
		t.Errorf("clone.ID = %q, want %q", clone.ID, original.ID)
	}
	if clone.ClientID != original.ClientID {
		t.Errorf("clone.ClientID = %q, want %q", clone.ClientID, original.ClientID)
	}
	if clone.State["key"] != original.State["key"] {
		t.Errorf("clone.State[key] = %v, want %v", clone.State["key"], original.State["key"])
	}

	// Verify clone is independent (modifying clone doesn't affect original)
	clone.State["key"] = "modified"
	if original.State["key"] == "modified" {
		t.Error("modifying clone affected original")
	}
}

func TestSession_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "not expired",
			expiresAt: time.Now().Add(time.Hour),
			want:      false,
		},
		{
			name:      "expired",
			expiresAt: time.Now().Add(-time.Hour),
			want:      true,
		},
		{
			name:      "just expired",
			expiresAt: time.Now().Add(-time.Millisecond),
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{ExpiresAt: tt.expiresAt}
			if got := s.IsExpired(); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_SetState(t *testing.T) {
	s := &Session{State: make(map[string]any)}

	s.SetState("key1", "value1")
	s.SetState("key2", 42)

	if s.State["key1"] != "value1" {
		t.Errorf("State[key1] = %v, want %q", s.State["key1"], "value1")
	}
	if s.State["key2"] != 42 {
		t.Errorf("State[key2] = %v, want %d", s.State["key2"], 42)
	}
}

func TestSession_SetState_NilState(t *testing.T) {
	s := &Session{}

	s.SetState("key", "value")

	if s.State["key"] != "value" {
		t.Errorf("State[key] = %v, want %q", s.State["key"], "value")
	}
}

func TestSession_GetState(t *testing.T) {
	s := &Session{
		State: map[string]any{
			"string": "value",
			"int":    42,
		},
	}

	if v, ok := s.GetState("string"); !ok || v != "value" {
		t.Errorf("GetState(string) = %v, %v; want %q, true", v, ok, "value")
	}

	if v, ok := s.GetState("int"); !ok || v != 42 {
		t.Errorf("GetState(int) = %v, %v; want %d, true", v, ok, 42)
	}

	if v, ok := s.GetState("missing"); ok {
		t.Errorf("GetState(missing) = %v, %v; want nil, false", v, ok)
	}
}

func TestSession_GetState_NilState(t *testing.T) {
	s := &Session{}

	if v, ok := s.GetState("key"); ok {
		t.Errorf("GetState(key) = %v, %v; want nil, false", v, ok)
	}
}

func TestStoreInterface_Defined(t *testing.T) {
	// Compile-time check that Store interface is properly defined
	var _ Store = (*MemoryStore)(nil)
}
