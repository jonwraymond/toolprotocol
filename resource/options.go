package resource

import "time"

// Option configures resource package types.
type Option func(any)

// CacheOption configures caching behavior.
type CacheOption struct {
	Enabled bool
	TTL     time.Duration
}

// WithCache returns a cache option with the given TTL.
func WithCache(enabled bool) CacheOption {
	return CacheOption{
		Enabled: enabled,
		TTL:     5 * time.Minute,
	}
}

// WithCacheTTL returns a cache option with the given TTL.
func WithCacheTTL(ttl time.Duration) CacheOption {
	return CacheOption{
		Enabled: true,
		TTL:     ttl,
	}
}

// DefaultCacheTTL is the default cache TTL.
const DefaultCacheTTL = 5 * time.Minute
