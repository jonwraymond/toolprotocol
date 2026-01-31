package stream

import "time"

// Option configures stream package types.
type Option func(any)

// HeartbeatOption configures heartbeat behavior.
type HeartbeatOption struct {
	Interval time.Duration
	Enabled  bool
}

// WithHeartbeat returns a HeartbeatOption with the given interval.
func WithHeartbeat(interval time.Duration) HeartbeatOption {
	return HeartbeatOption{
		Interval: interval,
		Enabled:  true,
	}
}

// DefaultHeartbeatInterval is the default interval for heartbeat events.
const DefaultHeartbeatInterval = 30 * time.Second
