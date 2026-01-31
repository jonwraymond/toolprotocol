package discover

import (
	"context"
	"testing"
)

func TestNegotiate_BothFull(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	client := &Capabilities{
		Tools:      true,
		Resources:  true,
		Prompts:    true,
		Streaming:  true,
		Sampling:   true,
		Progress:   true,
		Extensions: []string{"ext.a", "ext.b"},
	}
	server := &Capabilities{
		Tools:      true,
		Resources:  true,
		Prompts:    true,
		Streaming:  true,
		Sampling:   true,
		Progress:   true,
		Extensions: []string{"ext.a", "ext.b"},
	}

	result, err := d.Negotiate(ctx, client, server)
	if err != nil {
		t.Fatalf("Negotiate error = %v", err)
	}

	if !result.Tools {
		t.Error("Tools = false, want true")
	}
	if !result.Resources {
		t.Error("Resources = false, want true")
	}
	if len(result.Extensions) != 2 {
		t.Errorf("len(Extensions) = %d, want 2", len(result.Extensions))
	}
}

func TestNegotiate_ClientSubset(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	client := &Capabilities{
		Tools:     true,
		Resources: false,
		Streaming: true,
	}
	server := &Capabilities{
		Tools:     true,
		Resources: true,
		Streaming: true,
		Prompts:   true,
	}

	result, err := d.Negotiate(ctx, client, server)
	if err != nil {
		t.Fatalf("Negotiate error = %v", err)
	}

	// Result should be intersection
	if !result.Tools {
		t.Error("Tools = false, want true")
	}
	if result.Resources {
		t.Error("Resources = true, want false (client doesn't support)")
	}
	if !result.Streaming {
		t.Error("Streaming = false, want true")
	}
	if result.Prompts {
		t.Error("Prompts = true, want false (client doesn't support)")
	}
}

func TestNegotiate_ServerSubset(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	client := &Capabilities{
		Tools:     true,
		Resources: true,
		Prompts:   true,
	}
	server := &Capabilities{
		Tools:     true,
		Resources: false,
	}

	result, err := d.Negotiate(ctx, client, server)
	if err != nil {
		t.Fatalf("Negotiate error = %v", err)
	}

	if !result.Tools {
		t.Error("Tools = false, want true")
	}
	if result.Resources {
		t.Error("Resources = true, want false (server doesn't support)")
	}
}

func TestNegotiate_NoOverlap(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	client := &Capabilities{
		Tools:      true,
		Extensions: []string{"ext.a"},
	}
	server := &Capabilities{
		Resources:  true,
		Extensions: []string{"ext.b"},
	}

	result, err := d.Negotiate(ctx, client, server)
	if err != nil {
		t.Fatalf("Negotiate error = %v", err)
	}

	if result.Tools {
		t.Error("Tools = true, want false")
	}
	if result.Resources {
		t.Error("Resources = true, want false")
	}
	if len(result.Extensions) != 0 {
		t.Errorf("len(Extensions) = %d, want 0", len(result.Extensions))
	}
}

func TestNegotiate_Extensions(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	client := &Capabilities{
		Extensions: []string{"ext.a", "ext.common", "ext.b"},
	}
	server := &Capabilities{
		Extensions: []string{"ext.c", "ext.common", "ext.d"},
	}

	result, err := d.Negotiate(ctx, client, server)
	if err != nil {
		t.Fatalf("Negotiate error = %v", err)
	}

	if len(result.Extensions) != 1 {
		t.Fatalf("len(Extensions) = %d, want 1", len(result.Extensions))
	}
	if result.Extensions[0] != "ext.common" {
		t.Errorf("Extensions[0] = %q, want %q", result.Extensions[0], "ext.common")
	}
}

func TestNegotiate_NilCapabilities(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	// Client nil
	result, err := d.Negotiate(ctx, nil, &Capabilities{Tools: true})
	if err != nil {
		t.Fatalf("Negotiate(nil, caps) error = %v", err)
	}
	if result == nil {
		t.Fatal("Negotiate(nil, caps) returned nil")
	}

	// Server nil
	result, err = d.Negotiate(ctx, &Capabilities{Tools: true}, nil)
	if err != nil {
		t.Fatalf("Negotiate(caps, nil) error = %v", err)
	}
	if result == nil {
		t.Fatal("Negotiate(caps, nil) returned nil")
	}

	// Both nil
	result, err = d.Negotiate(ctx, nil, nil)
	if err != nil {
		t.Fatalf("Negotiate(nil, nil) error = %v", err)
	}
	if result == nil {
		t.Fatal("Negotiate(nil, nil) returned nil")
	}
}

func TestNegotiator_Strategy(t *testing.T) {
	// Test that Negotiator can use different strategies
	negotiator := &Negotiator{
		Strategy: NegotiateIntersect,
	}

	client := &Capabilities{Tools: true, Resources: true}
	server := &Capabilities{Tools: true, Resources: false}

	result := negotiator.Negotiate(client, server)
	if !result.Tools {
		t.Error("Tools = false, want true")
	}
	if result.Resources {
		t.Error("Resources = true, want false")
	}
}

func TestNegotiator_Strategy_Merge(t *testing.T) {
	negotiator := &Negotiator{
		Strategy: NegotiateMerge,
	}

	client := &Capabilities{Tools: true, Resources: false}
	server := &Capabilities{Tools: false, Resources: true}

	result := negotiator.Negotiate(client, server)
	if !result.Tools {
		t.Error("Tools = false, want true")
	}
	if !result.Resources {
		t.Error("Resources = false, want true")
	}
}
