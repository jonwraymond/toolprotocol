package resource

import (
	"context"
	"sync"
)

// StaticProvider is a provider backed by static resources.
type StaticProvider struct {
	mu        sync.RWMutex
	resources map[string]*Resource
	contents  map[string]*Contents
	templates []Template
}

// NewStaticProvider creates a new StaticProvider.
func NewStaticProvider() *StaticProvider {
	return &StaticProvider{
		resources: make(map[string]*Resource),
		contents:  make(map[string]*Contents),
	}
}

// Add adds a resource with its contents.
func (p *StaticProvider) Add(res Resource, contents Contents) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.resources[res.URI] = &res
	p.contents[res.URI] = &contents
}

// Remove removes a resource.
func (p *StaticProvider) Remove(uri string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.resources, uri)
	delete(p.contents, uri)
}

// AddTemplate adds a resource template.
func (p *StaticProvider) AddTemplate(tmpl Template) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.templates = append(p.templates, tmpl)
}

// List returns all resources.
func (p *StaticProvider) List(ctx context.Context) ([]Resource, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	resources := make([]Resource, 0, len(p.resources))
	for _, r := range p.resources {
		resources = append(resources, *r)
	}
	return resources, nil
}

// Read returns the contents of a resource.
func (p *StaticProvider) Read(ctx context.Context, uri string) (*Contents, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	contents, ok := p.contents[uri]
	if !ok {
		return nil, ErrResourceNotFound
	}

	// Return a copy
	return &Contents{
		URI:      contents.URI,
		MIMEType: contents.MIMEType,
		Text:     contents.Text,
		Blob:     contents.Blob,
	}, nil
}

// Templates returns all templates.
func (p *StaticProvider) Templates(ctx context.Context) ([]Template, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	templates := make([]Template, len(p.templates))
	copy(templates, p.templates)
	return templates, nil
}
