package discover

import (
	"context"
	"testing"
)

func TestCapabilities_Defaults(t *testing.T) {
	var caps Capabilities

	if caps.Tools {
		t.Error("default Tools = true, want false")
	}
	if caps.Resources {
		t.Error("default Resources = true, want false")
	}
	if caps.Prompts {
		t.Error("default Prompts = true, want false")
	}
	if caps.Streaming {
		t.Error("default Streaming = true, want false")
	}
	if caps.Sampling {
		t.Error("default Sampling = true, want false")
	}
	if caps.Progress {
		t.Error("default Progress = true, want false")
	}
}

func TestCapabilities_AllEnabled(t *testing.T) {
	caps := Capabilities{
		Tools:      true,
		Resources:  true,
		Prompts:    true,
		Streaming:  true,
		Sampling:   true,
		Progress:   true,
		Extensions: []string{"custom.ext"},
	}

	if !caps.Tools {
		t.Error("Tools = false, want true")
	}
	if !caps.Resources {
		t.Error("Resources = false, want true")
	}
	if len(caps.Extensions) != 1 {
		t.Errorf("len(Extensions) = %d, want 1", len(caps.Extensions))
	}
}

func TestCapabilities_Merge(t *testing.T) {
	a := &Capabilities{
		Tools:      true,
		Resources:  false,
		Extensions: []string{"ext.a"},
	}
	b := &Capabilities{
		Tools:      false,
		Resources:  true,
		Extensions: []string{"ext.b"},
	}

	merged := a.Merge(b)

	// Merge should combine: a OR b for bool, union for extensions
	if !merged.Tools {
		t.Error("merged.Tools = false, want true")
	}
	if !merged.Resources {
		t.Error("merged.Resources = false, want true")
	}
	if len(merged.Extensions) != 2 {
		t.Errorf("len(merged.Extensions) = %d, want 2", len(merged.Extensions))
	}
}

func TestCapabilities_Intersect(t *testing.T) {
	a := &Capabilities{
		Tools:      true,
		Resources:  true,
		Prompts:    false,
		Extensions: []string{"ext.a", "ext.common"},
	}
	b := &Capabilities{
		Tools:      true,
		Resources:  false,
		Prompts:    true,
		Extensions: []string{"ext.b", "ext.common"},
	}

	intersected := a.Intersect(b)

	// Intersect: a AND b for bool, common for extensions
	if !intersected.Tools {
		t.Error("intersected.Tools = false, want true")
	}
	if intersected.Resources {
		t.Error("intersected.Resources = true, want false")
	}
	if intersected.Prompts {
		t.Error("intersected.Prompts = true, want false")
	}
	// Only ext.common should be in both
	if len(intersected.Extensions) != 1 || intersected.Extensions[0] != "ext.common" {
		t.Errorf("intersected.Extensions = %v, want [ext.common]", intersected.Extensions)
	}
}

func TestCapabilities_Merge_NilOther(t *testing.T) {
	a := &Capabilities{
		Tools:      true,
		Extensions: []string{"ext.a"},
	}

	merged := a.Merge(nil)
	if !merged.Tools {
		t.Error("merged.Tools = false, want true")
	}
	if len(merged.Extensions) != 1 {
		t.Errorf("len(merged.Extensions) = %d, want 1", len(merged.Extensions))
	}
}

func TestCapabilities_Intersect_NilOther(t *testing.T) {
	a := &Capabilities{
		Tools: true,
	}

	intersected := a.Intersect(nil)
	// Should return empty capabilities when other is nil
	if intersected.Tools {
		t.Error("intersected.Tools = true, want false")
	}
}

func TestFilter_Matches_NilFilter(t *testing.T) {
	svc := NewService("test", "http://localhost:8080")
	var filter *Filter = nil

	if !filter.Matches(svc) {
		t.Error("nil filter should match all services")
	}
}

func TestFilter_Matches_WithCapabilities(t *testing.T) {
	svc := NewService("test", "http://localhost:8080").
		WithCapability("tools").
		WithCapability("streaming")

	filter := &Filter{Capabilities: []string{"tools"}}
	if !filter.Matches(svc) {
		t.Error("service with tools should match filter for tools")
	}

	filter = &Filter{Capabilities: []string{"resources"}}
	if filter.Matches(svc) {
		t.Error("service without resources should not match filter for resources")
	}
}

func TestFilter_Matches_NilCapabilities(t *testing.T) {
	svc := &mockDiscoverable{
		id:   "test",
		caps: nil,
	}

	filter := &Filter{Capabilities: []string{"tools"}}
	if filter.Matches(svc) {
		t.Error("service with nil capabilities should not match capability filter")
	}
}

func TestFilter_Matches_ExtensionCapability(t *testing.T) {
	svc := NewService("test", "http://localhost:8080").
		WithCapability("custom.ext")

	filter := &Filter{Capabilities: []string{"custom.ext"}}
	if !filter.Matches(svc) {
		t.Error("service with custom extension should match filter for that extension")
	}
}

func TestFilter_Matches_AllCapabilityTypes(t *testing.T) {
	svc := NewService("test", "http://localhost:8080").
		WithCapability("tools").
		WithCapability("resources").
		WithCapability("prompts").
		WithCapability("streaming").
		WithCapability("sampling").
		WithCapability("progress")

	caps := []string{"tools", "resources", "prompts", "streaming", "sampling", "progress"}
	for _, cap := range caps {
		filter := &Filter{Capabilities: []string{cap}}
		if !filter.Matches(svc) {
			t.Errorf("service should match filter for %q", cap)
		}
	}
}

func TestFilter_Defaults(t *testing.T) {
	var filter Filter

	if filter.Namespace != "" {
		t.Errorf("default Namespace = %q, want empty", filter.Namespace)
	}
	if filter.Tags != nil {
		t.Errorf("default Tags = %v, want nil", filter.Tags)
	}
	if filter.Limit != 0 {
		t.Errorf("default Limit = %d, want 0", filter.Limit)
	}
	if filter.Cursor != "" {
		t.Errorf("default Cursor = %q, want empty", filter.Cursor)
	}
}

func TestFilter_CustomValues(t *testing.T) {
	filter := Filter{
		Namespace:    "myns",
		Tags:         []string{"tag1", "tag2"},
		Capabilities: []string{"tools"},
		Limit:        10,
		Cursor:       "abc123",
	}

	if filter.Namespace != "myns" {
		t.Errorf("Namespace = %q, want %q", filter.Namespace, "myns")
	}
	if len(filter.Tags) != 2 {
		t.Errorf("len(Tags) = %d, want 2", len(filter.Tags))
	}
	if filter.Limit != 10 {
		t.Errorf("Limit = %d, want 10", filter.Limit)
	}
}

func TestDiscoverableInterface_Contract(t *testing.T) {
	var _ Discoverable = (*mockDiscoverable)(nil)
}

func TestDiscoveryInterface_Contract(t *testing.T) {
	var _ Discovery = (*mockDiscovery)(nil)
}

// Mock implementations for interface contracts
type mockDiscoverable struct {
	id          string
	name        string
	description string
	version     string
	endpoint    string
	caps        *Capabilities
}

func (m *mockDiscoverable) ID() string                  { return m.id }
func (m *mockDiscoverable) Name() string                { return m.name }
func (m *mockDiscoverable) Description() string         { return m.description }
func (m *mockDiscoverable) Version() string             { return m.version }
func (m *mockDiscoverable) Endpoint() string            { return m.endpoint }
func (m *mockDiscoverable) Capabilities() *Capabilities { return m.caps }

// mockDiscovery implements Discovery for testing interface contract
type mockDiscovery struct{}

func (m *mockDiscovery) Register(ctx context.Context, svc Discoverable) error { return nil }
func (m *mockDiscovery) Deregister(ctx context.Context, id string) error      { return nil }
func (m *mockDiscovery) Get(ctx context.Context, id string) (Discoverable, error) {
	return nil, nil
}
func (m *mockDiscovery) List(ctx context.Context, filter *Filter) ([]Discoverable, error) {
	return nil, nil
}
func (m *mockDiscovery) Negotiate(ctx context.Context, client, server *Capabilities) (*Capabilities, error) {
	return nil, nil
}
