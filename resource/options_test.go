package resource

import (
	"testing"
	"time"
)

func TestWithCache(t *testing.T) {
	opt := WithCache(true)
	if !opt.Enabled {
		t.Error("Enabled = false, want true")
	}
	if opt.TTL != 5*time.Minute {
		t.Errorf("TTL = %v, want %v", opt.TTL, 5*time.Minute)
	}
}

func TestWithCacheTTL(t *testing.T) {
	opt := WithCacheTTL(10 * time.Minute)
	if !opt.Enabled {
		t.Error("Enabled = false, want true")
	}
	if opt.TTL != 10*time.Minute {
		t.Errorf("TTL = %v, want %v", opt.TTL, 10*time.Minute)
	}
}

func TestNewRegistry_Defaults(t *testing.T) {
	r := NewRegistry()

	if r.providers == nil {
		t.Error("providers map not initialized")
	}
}

func TestDefaultCacheTTL(t *testing.T) {
	if DefaultCacheTTL != 5*time.Minute {
		t.Errorf("DefaultCacheTTL = %v, want %v", DefaultCacheTTL, 5*time.Minute)
	}
}
