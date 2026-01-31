package discover

import (
	"context"
	"sort"
	"sync"
)

// MemoryDiscovery is an in-memory implementation of Discovery.
type MemoryDiscovery struct {
	mu       sync.RWMutex
	services map[string]Discoverable
}

// NewMemory creates a new in-memory discovery registry.
func NewMemory() *MemoryDiscovery {
	return &MemoryDiscovery{
		services: make(map[string]Discoverable),
	}
}

// Register adds a service to the registry.
func (d *MemoryDiscovery) Register(ctx context.Context, svc Discoverable) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.services[svc.ID()]; exists {
		return ErrDuplicate
	}

	d.services[svc.ID()] = svc
	return nil
}

// Deregister removes a service from the registry.
func (d *MemoryDiscovery) Deregister(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.services[id]; !exists {
		return ErrNotFound
	}

	delete(d.services, id)
	return nil
}

// Get retrieves a service by ID.
func (d *MemoryDiscovery) Get(ctx context.Context, id string) (Discoverable, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	svc, exists := d.services[id]
	if !exists {
		return nil, ErrNotFound
	}

	return svc, nil
}

// List returns services matching the filter.
func (d *MemoryDiscovery) List(ctx context.Context, filter *Filter) ([]Discoverable, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	// Collect matching services
	var result []Discoverable
	for _, svc := range d.services {
		if filter == nil || filter.Matches(svc) {
			result = append(result, svc)
		}
	}

	// Sort by ID for consistent ordering
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID() < result[j].ID()
	})

	// Apply limit
	if filter != nil && filter.Limit > 0 && len(result) > filter.Limit {
		result = result[:filter.Limit]
	}

	return result, nil
}

// Negotiate determines compatible capabilities between client and server.
func (d *MemoryDiscovery) Negotiate(ctx context.Context, client, server *Capabilities) (*Capabilities, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if client == nil || server == nil {
		return &Capabilities{}, nil
	}

	return client.Intersect(server), nil
}
