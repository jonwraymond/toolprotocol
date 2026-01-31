package discover

import "context"

// Discoverable represents a service that can be discovered.
type Discoverable interface {
	// ID returns the unique service identifier.
	ID() string

	// Name returns the human-readable service name.
	Name() string

	// Description returns what the service does.
	Description() string

	// Version returns the service version.
	Version() string

	// Endpoint returns the service endpoint URL.
	Endpoint() string

	// Capabilities returns the service capabilities.
	Capabilities() *Capabilities
}

// Discovery manages service registration and lookup.
type Discovery interface {
	// Register adds a service to the discovery registry.
	Register(ctx context.Context, svc Discoverable) error

	// Deregister removes a service from the registry.
	Deregister(ctx context.Context, id string) error

	// Get retrieves a service by ID.
	Get(ctx context.Context, id string) (Discoverable, error)

	// List returns services matching the filter.
	List(ctx context.Context, filter *Filter) ([]Discoverable, error)

	// Negotiate determines compatible capabilities between client and server.
	Negotiate(ctx context.Context, client, server *Capabilities) (*Capabilities, error)
}
