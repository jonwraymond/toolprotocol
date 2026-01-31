package elicit

import (
	"testing"
	"time"
)

func TestWithDefaultTimeout(t *testing.T) {
	timeout := 5 * time.Minute
	handler := NewHandler(WithDefaultTimeout(timeout))

	if handler.timeout != timeout {
		t.Errorf("timeout = %v, want %v", handler.timeout, timeout)
	}
}

func TestNewHandler_Defaults(t *testing.T) {
	handler := NewHandler()

	// Default timeout should be 30 seconds
	if handler.timeout != 30*time.Second {
		t.Errorf("default timeout = %v, want %v", handler.timeout, 30*time.Second)
	}

	// Handler callback should be nil by default
	if handler.handler != nil {
		t.Error("default handler should be nil")
	}
}

func TestDefaultRequestTimeout(t *testing.T) {
	if DefaultRequestTimeout != 30*time.Second {
		t.Errorf("DefaultRequestTimeout = %v, want %v", DefaultRequestTimeout, 30*time.Second)
	}
}
