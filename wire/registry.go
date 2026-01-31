package wire

import "sync"

// Registry manages wire format handlers.
type Registry struct {
	mu    sync.RWMutex
	wires map[string]Wire
}

// NewRegistry creates a new empty registry.
func NewRegistry() *Registry {
	return &Registry{
		wires: make(map[string]Wire),
	}
}

// Register adds a wire format handler to the registry.
// If a handler with the same name exists, it is replaced.
func (r *Registry) Register(name string, wire Wire) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.wires[name] = wire
}

// Get returns the wire handler for the given format name, or nil if not found.
func (r *Registry) Get(name string) Wire {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.wires[name]
}

// List returns the names of all registered wire formats.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.wires))
	for name := range r.wires {
		names = append(names, name)
	}
	return names
}

// defaultRegistry is the global registry with standard wire formats.
var defaultRegistry = func() *Registry {
	r := NewRegistry()
	r.Register("mcp", NewMCP())
	r.Register("a2a", NewA2A())
	r.Register("acp", NewACP())
	return r
}()

// DefaultRegistry returns the default wire format registry.
func DefaultRegistry() *Registry {
	return defaultRegistry
}
