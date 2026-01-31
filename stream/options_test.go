package stream

import (
	"testing"
	"time"
)

func TestWithBackpressure_Block(t *testing.T) {
	source := NewSource(WithBackpressure(BackpressureBlock))
	if source.backpressure != BackpressureBlock {
		t.Errorf("backpressure = %v, want BackpressureBlock", source.backpressure)
	}
}

func TestWithBackpressure_Drop(t *testing.T) {
	source := NewSource(WithBackpressure(BackpressureDrop))
	if source.backpressure != BackpressureDrop {
		t.Errorf("backpressure = %v, want BackpressureDrop", source.backpressure)
	}
}

func TestWithHeartbeat(t *testing.T) {
	opt := WithHeartbeat(15 * time.Second)
	if opt.Interval != 15*time.Second {
		t.Errorf("Interval = %v, want %v", opt.Interval, 15*time.Second)
	}
	if !opt.Enabled {
		t.Error("Enabled = false, want true")
	}
}

func TestDefaultHeartbeatInterval(t *testing.T) {
	if DefaultHeartbeatInterval != 30*time.Second {
		t.Errorf("DefaultHeartbeatInterval = %v, want %v", DefaultHeartbeatInterval, 30*time.Second)
	}
}

func TestDefaultBufferSize(t *testing.T) {
	if DefaultBufferSize != 100 {
		t.Errorf("DefaultBufferSize = %d, want 100", DefaultBufferSize)
	}
}
