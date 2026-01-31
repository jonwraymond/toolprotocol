package resource

import (
	"context"
	"strings"
	"sync"
)

// Registry manages multiple resource providers.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

// NewRegistry creates a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register registers a provider for a URI scheme.
func (r *Registry) Register(scheme string, provider Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.providers[scheme]; exists {
		return &ResourceError{
			URI: scheme,
			Op:  "register",
			Err: ErrDuplicateProvider,
		}
	}

	r.providers[scheme] = provider
	return nil
}

// Unregister removes a provider for a URI scheme.
func (r *Registry) Unregister(scheme string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.providers[scheme]; !exists {
		return &ResourceError{
			URI: scheme,
			Op:  "unregister",
			Err: ErrProviderNotFound,
		}
	}

	delete(r.providers, scheme)
	return nil
}

// parseScheme extracts the scheme from a URI.
func parseScheme(uri string) string {
	if idx := strings.Index(uri, "://"); idx > 0 {
		return uri[:idx]
	}
	return ""
}

// List returns resources from all providers.
func (r *Registry) List(ctx context.Context) ([]Resource, error) {
	r.mu.RLock()
	providers := make(map[string]Provider, len(r.providers))
	for k, v := range r.providers {
		providers[k] = v
	}
	r.mu.RUnlock()

	var all []Resource
	for _, p := range providers {
		resources, err := p.List(ctx)
		if err != nil {
			return nil, err
		}
		all = append(all, resources...)
	}
	return all, nil
}

// Read reads a resource from the appropriate provider.
func (r *Registry) Read(ctx context.Context, uri string) (*Contents, error) {
	scheme := parseScheme(uri)
	if scheme == "" {
		return nil, &ResourceError{
			URI: uri,
			Op:  "read",
			Err: ErrInvalidURI,
		}
	}

	r.mu.RLock()
	provider, ok := r.providers[scheme]
	r.mu.RUnlock()

	if !ok {
		return nil, &ResourceError{
			URI: uri,
			Op:  "read",
			Err: ErrProviderNotFound,
		}
	}

	return provider.Read(ctx, uri)
}

// Templates returns templates from all providers.
func (r *Registry) Templates(ctx context.Context) ([]Template, error) {
	r.mu.RLock()
	providers := make(map[string]Provider, len(r.providers))
	for k, v := range r.providers {
		providers[k] = v
	}
	r.mu.RUnlock()

	var all []Template
	for _, p := range providers {
		templates, err := p.Templates(ctx)
		if err != nil {
			return nil, err
		}
		all = append(all, templates...)
	}
	return all, nil
}
