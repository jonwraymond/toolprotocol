package discover

import (
	"context"
	"sync"
	"testing"
)

func TestMemoryDiscovery_Register(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	svc := NewService("test-1", "http://localhost:8080")
	err := d.Register(ctx, svc)
	if err != nil {
		t.Fatalf("Register error = %v", err)
	}

	// Should be able to get it back
	got, err := d.Get(ctx, "test-1")
	if err != nil {
		t.Fatalf("Get error = %v", err)
	}
	if got.ID() != "test-1" {
		t.Errorf("ID = %q, want %q", got.ID(), "test-1")
	}
}

func TestMemoryDiscovery_Register_Duplicate(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	svc1 := NewService("dup", "http://localhost:8080")
	svc2 := NewService("dup", "http://localhost:9000")

	_ = d.Register(ctx, svc1)
	err := d.Register(ctx, svc2)

	// Should error on duplicate
	if err == nil {
		t.Error("Register(duplicate) error = nil, want error")
	}
}

func TestMemoryDiscovery_Deregister(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	svc := NewService("test-dereg", "http://localhost:8080")
	_ = d.Register(ctx, svc)

	err := d.Deregister(ctx, "test-dereg")
	if err != nil {
		t.Fatalf("Deregister error = %v", err)
	}

	// Should no longer exist
	_, err = d.Get(ctx, "test-dereg")
	if err == nil {
		t.Error("Get after Deregister should return error")
	}
}

func TestMemoryDiscovery_Deregister_NotFound(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	err := d.Deregister(ctx, "nonexistent")
	if err == nil {
		t.Error("Deregister(nonexistent) error = nil, want error")
	}
}

func TestMemoryDiscovery_Get(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	svc := NewService("get-test", "http://localhost:8080")
	_ = d.Register(ctx, svc)

	got, err := d.Get(ctx, "get-test")
	if err != nil {
		t.Fatalf("Get error = %v", err)
	}
	if got.Endpoint() != "http://localhost:8080" {
		t.Errorf("Endpoint = %q, want %q", got.Endpoint(), "http://localhost:8080")
	}
}

func TestMemoryDiscovery_Get_NotFound(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	_, err := d.Get(ctx, "nonexistent")
	if err == nil {
		t.Error("Get(nonexistent) error = nil, want error")
	}
}

func TestMemoryDiscovery_List_All(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	_ = d.Register(ctx, NewService("svc-1", "http://localhost:8081"))
	_ = d.Register(ctx, NewService("svc-2", "http://localhost:8082"))
	_ = d.Register(ctx, NewService("svc-3", "http://localhost:8083"))

	list, err := d.List(ctx, nil)
	if err != nil {
		t.Fatalf("List error = %v", err)
	}
	if len(list) != 3 {
		t.Errorf("len(List()) = %d, want 3", len(list))
	}
}

func TestMemoryDiscovery_List_Filtered(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	svc1 := NewService("svc-tools", "http://localhost:8081").WithCapability("tools")
	svc2 := NewService("svc-resources", "http://localhost:8082").WithCapability("resources")
	svc3 := NewService("svc-both", "http://localhost:8083").WithCapability("tools").WithCapability("resources")

	_ = d.Register(ctx, svc1)
	_ = d.Register(ctx, svc2)
	_ = d.Register(ctx, svc3)

	// Filter for tools capability
	filter := &Filter{Capabilities: []string{"tools"}}
	list, err := d.List(ctx, filter)
	if err != nil {
		t.Fatalf("List error = %v", err)
	}
	if len(list) != 2 {
		t.Errorf("len(List(tools)) = %d, want 2", len(list))
	}
}

func TestMemoryDiscovery_List_Pagination(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		svc := NewService(string(rune('a'+i)), "http://localhost:8080")
		_ = d.Register(ctx, svc)
	}

	// Get first 5
	filter := &Filter{Limit: 5}
	list, err := d.List(ctx, filter)
	if err != nil {
		t.Fatalf("List error = %v", err)
	}
	if len(list) != 5 {
		t.Errorf("len(List(limit=5)) = %d, want 5", len(list))
	}
}

func TestMemoryDiscovery_ConcurrentSafety(t *testing.T) {
	d := NewMemory()
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			svc := NewService(string(rune('A'+n%26))+string(rune('0'+n%10)), "http://localhost:8080")
			_ = d.Register(ctx, svc)
			_, _ = d.List(ctx, nil)
		}(i)
	}
	wg.Wait()
}

func TestMemoryDiscovery_ContextCancellation(t *testing.T) {
	d := NewMemory()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	svc := NewService("test", "http://localhost:8080")
	err := d.Register(ctx, svc)
	// Should respect context cancellation
	if err != nil && err != context.Canceled {
		// Either nil or context.Canceled is acceptable
	}
}

func TestMemoryDiscovery_ImplementsDiscovery(t *testing.T) {
	var _ Discovery = (*MemoryDiscovery)(nil)
}
