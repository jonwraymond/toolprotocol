package transport

import (
	"fmt"
	"sync"
)

// Factory creates a Transport from configuration.
type Factory func(cfg any) (Transport, error)

// Registry manages transport factories.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

// NewRegistry creates a new empty registry.
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]Factory),
	}
}

// Register adds a transport factory to the registry.
// If a factory with the same name exists, it is replaced.
func (r *Registry) Register(name string, factory Factory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[name] = factory
}

// Get returns the factory for the given transport name, or nil if not found.
func (r *Registry) Get(name string) Factory {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.factories[name]
}

// List returns the names of all registered transports.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	return names
}

// New creates a transport using the registered factory.
func (r *Registry) New(name string, cfg any) (Transport, error) {
	factory := r.Get(name)
	if factory == nil {
		return nil, fmt.Errorf("unknown transport: %s", name)
	}
	return factory(cfg)
}

// defaultRegistry is the global registry with standard transports.
var defaultRegistry = func() *Registry {
	r := NewRegistry()

	r.Register("stdio", func(cfg any) (Transport, error) {
		return &StdioTransport{}, nil
	})

	r.Register("sse", func(cfg any) (Transport, error) {
		t := &SSETransport{}
		if c, ok := cfg.(*SSEConfig); ok && c != nil {
			t.Config = *c
		}
		return t, nil
	})

	r.Register("streamable", func(cfg any) (Transport, error) {
		t := &StreamableHTTPTransport{}
		if c, ok := cfg.(*StreamableConfig); ok && c != nil {
			t.Config = *c
		}
		return t, nil
	})

	return r
}()

// DefaultRegistry returns the default transport registry.
func DefaultRegistry() *Registry {
	return defaultRegistry
}

// New creates a transport using the default registry.
func New(name string, cfg any) (Transport, error) {
	return defaultRegistry.New(name, cfg)
}
