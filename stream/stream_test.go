package stream

import (
	"testing"
)

func TestEventType_String(t *testing.T) {
	tests := []struct {
		et   EventType
		want string
	}{
		{EventProgress, "progress"},
		{EventPartial, "partial"},
		{EventComplete, "complete"},
		{EventError, "error"},
		{EventHeartbeat, "heartbeat"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.et.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEventType_Valid(t *testing.T) {
	tests := []struct {
		et   EventType
		want bool
	}{
		{EventProgress, true},
		{EventPartial, true},
		{EventComplete, true},
		{EventError, true},
		{EventHeartbeat, true},
		{EventType("unknown"), false},
		{EventType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.et), func(t *testing.T) {
			if got := tt.et.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvent_Fields(t *testing.T) {
	e := Event{
		Type:  EventProgress,
		ID:    "evt-123",
		Data:  map[string]any{"progress": 0.5},
		Retry: 3000,
	}

	if e.Type != EventProgress {
		t.Errorf("Type = %v, want %v", e.Type, EventProgress)
	}
	if e.ID != "evt-123" {
		t.Errorf("ID = %q, want %q", e.ID, "evt-123")
	}
	if e.Data == nil {
		t.Error("Data is nil")
	}
	if e.Retry != 3000 {
		t.Errorf("Retry = %d, want 3000", e.Retry)
	}
}

func TestEvent_Clone(t *testing.T) {
	original := Event{
		Type:  EventComplete,
		ID:    "evt-456",
		Data:  "result",
		Retry: 5000,
	}

	clone := original.Clone()

	if clone.Type != original.Type {
		t.Errorf("clone.Type = %v, want %v", clone.Type, original.Type)
	}
	if clone.ID != original.ID {
		t.Errorf("clone.ID = %q, want %q", clone.ID, original.ID)
	}
	if clone.Data != original.Data {
		t.Errorf("clone.Data = %v, want %v", clone.Data, original.Data)
	}
	if clone.Retry != original.Retry {
		t.Errorf("clone.Retry = %d, want %d", clone.Retry, original.Retry)
	}
}

func TestStreamInterface_Defined(t *testing.T) {
	// Compile-time check that interfaces are properly defined
	var _ Stream = (*DefaultStream)(nil)
	var _ Source = (*DefaultSource)(nil)
	var _ Sink = (*DefaultSink)(nil)
}
