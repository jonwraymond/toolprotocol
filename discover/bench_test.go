package discover

import (
	"context"
	"fmt"
	"testing"
)

// BenchmarkNewService measures service creation performance.
func BenchmarkNewService(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewService("svc-1", "http://localhost:8080")
	}
}

// BenchmarkService_WithCapability measures capability addition.
func BenchmarkService_WithCapability(b *testing.B) {
	svc := NewService("svc-1", "http://localhost:8080")

	b.ResetTimer()
	for b.Loop() {
		_ = svc.WithCapability("tools")
	}
}

// BenchmarkNewMemory measures memory discovery creation.
func BenchmarkNewMemory(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewMemory()
	}
}

// BenchmarkMemoryDiscovery_Register measures registration performance.
func BenchmarkMemoryDiscovery_Register(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		b.StopTimer()
		discovery := NewMemory()
		svc := NewService("svc-1", "http://localhost:8080")
		b.StartTimer()

		_ = discovery.Register(ctx, svc)
	}
}

// BenchmarkMemoryDiscovery_Get measures lookup performance.
func BenchmarkMemoryDiscovery_Get(b *testing.B) {
	discovery := NewMemory()
	ctx := context.Background()
	_ = discovery.Register(ctx, NewService("svc-1", "http://localhost:8080"))

	b.ResetTimer()
	for b.Loop() {
		_, _ = discovery.Get(ctx, "svc-1")
	}
}

// BenchmarkMemoryDiscovery_Get_Miss measures miss performance.
func BenchmarkMemoryDiscovery_Get_Miss(b *testing.B) {
	discovery := NewMemory()
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = discovery.Get(ctx, "nonexistent")
	}
}

// BenchmarkMemoryDiscovery_List measures listing performance.
func BenchmarkMemoryDiscovery_List(b *testing.B) {
	discovery := NewMemory()
	ctx := context.Background()

	// Pre-populate
	for i := range 100 {
		_ = discovery.Register(ctx, NewService(fmt.Sprintf("svc-%d", i), "http://localhost"))
	}

	b.ResetTimer()
	for b.Loop() {
		_, _ = discovery.List(ctx, nil)
	}
}

// BenchmarkMemoryDiscovery_List_WithFilter measures filtered listing.
func BenchmarkMemoryDiscovery_List_WithFilter(b *testing.B) {
	discovery := NewMemory()
	ctx := context.Background()

	// Pre-populate with mixed capabilities
	for i := range 100 {
		svc := NewService(fmt.Sprintf("svc-%d", i), "http://localhost")
		if i%2 == 0 {
			svc.WithCapability("tools")
		}
		_ = discovery.Register(ctx, svc)
	}

	filter := &Filter{Capabilities: []string{"tools"}}

	b.ResetTimer()
	for b.Loop() {
		_, _ = discovery.List(ctx, filter)
	}
}

// BenchmarkMemoryDiscovery_Deregister measures deregistration performance.
func BenchmarkMemoryDiscovery_Deregister(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		b.StopTimer()
		discovery := NewMemory()
		_ = discovery.Register(ctx, NewService("svc-1", "http://localhost"))
		b.StartTimer()

		_ = discovery.Deregister(ctx, "svc-1")
	}
}

// BenchmarkMemoryDiscovery_Negotiate measures negotiation performance.
func BenchmarkMemoryDiscovery_Negotiate(b *testing.B) {
	discovery := NewMemory()
	ctx := context.Background()

	client := &Capabilities{Tools: true, Streaming: true, Resources: true}
	server := &Capabilities{Tools: true, Prompts: true}

	b.ResetTimer()
	for b.Loop() {
		_, _ = discovery.Negotiate(ctx, client, server)
	}
}

// BenchmarkCapabilities_Merge measures merge performance.
func BenchmarkCapabilities_Merge(b *testing.B) {
	caps1 := &Capabilities{
		Tools:      true,
		Extensions: []string{"ext-a", "ext-b"},
	}
	caps2 := &Capabilities{
		Resources:  true,
		Extensions: []string{"ext-c", "ext-d"},
	}

	b.ResetTimer()
	for b.Loop() {
		_ = caps1.Merge(caps2)
	}
}

// BenchmarkCapabilities_Intersect measures intersection performance.
func BenchmarkCapabilities_Intersect(b *testing.B) {
	caps1 := &Capabilities{
		Tools:      true,
		Resources:  true,
		Extensions: []string{"ext-a", "ext-b", "ext-c"},
	}
	caps2 := &Capabilities{
		Tools:      true,
		Streaming:  true,
		Extensions: []string{"ext-b", "ext-c", "ext-d"},
	}

	b.ResetTimer()
	for b.Loop() {
		_ = caps1.Intersect(caps2)
	}
}

// BenchmarkFilter_Matches measures filter matching performance.
func BenchmarkFilter_Matches(b *testing.B) {
	filter := &Filter{
		Capabilities: []string{"tools", "streaming"},
	}
	svc := NewService("svc-1", "http://localhost").
		WithCapability("tools").
		WithCapability("streaming")

	b.ResetTimer()
	for b.Loop() {
		_ = filter.Matches(svc)
	}
}

// BenchmarkNegotiator_Negotiate measures negotiator performance.
func BenchmarkNegotiator_Negotiate(b *testing.B) {
	negotiator := &Negotiator{Strategy: NegotiateIntersect}
	client := &Capabilities{Tools: true, Resources: true}
	server := &Capabilities{Tools: true, Streaming: true}

	b.ResetTimer()
	for b.Loop() {
		_ = negotiator.Negotiate(client, server)
	}
}

// BenchmarkMemoryDiscovery_Concurrent measures concurrent access.
func BenchmarkMemoryDiscovery_Concurrent(b *testing.B) {
	discovery := NewMemory()
	ctx := context.Background()

	// Pre-populate
	for i := range 100 {
		_ = discovery.Register(ctx, NewService(fmt.Sprintf("svc-%d", i), "http://localhost"))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			switch i % 4 {
			case 0:
				_, _ = discovery.Get(ctx, fmt.Sprintf("svc-%d", i%100))
			case 1:
				_, _ = discovery.List(ctx, nil)
			case 2:
				_, _ = discovery.List(ctx, &Filter{Capabilities: []string{"tools"}})
			case 3:
				_, _ = discovery.Negotiate(ctx, &Capabilities{Tools: true}, &Capabilities{Tools: true})
			}
			i++
		}
	})
}
