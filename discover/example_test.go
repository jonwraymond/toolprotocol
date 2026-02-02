package discover_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonwraymond/toolprotocol/discover"
)

func ExampleNewService() {
	svc := discover.NewService("my-service", "http://localhost:8080")

	fmt.Println("ID:", svc.ID())
	fmt.Println("Endpoint:", svc.Endpoint())
	// Output:
	// ID: my-service
	// Endpoint: http://localhost:8080
}

func ExampleService_SetName() {
	svc := discover.NewService("svc-1", "http://localhost:8080").
		SetName("My Service")

	fmt.Println("Name:", svc.Name())
	// Output:
	// Name: My Service
}

func ExampleService_SetDescription() {
	svc := discover.NewService("svc-1", "http://localhost:8080").
		SetDescription("A helpful service")

	fmt.Println("Description:", svc.Description())
	// Output:
	// Description: A helpful service
}

func ExampleService_SetVersion() {
	svc := discover.NewService("svc-1", "http://localhost:8080").
		SetVersion("1.0.0")

	fmt.Println("Version:", svc.Version())
	// Output:
	// Version: 1.0.0
}

func ExampleService_WithCapability() {
	svc := discover.NewService("svc-1", "http://localhost:8080").
		WithCapability("tools").
		WithCapability("streaming")

	caps := svc.Capabilities()
	fmt.Println("Tools:", caps.Tools)
	fmt.Println("Streaming:", caps.Streaming)
	// Output:
	// Tools: true
	// Streaming: true
}

func ExampleService_Validate() {
	svc := discover.NewService("svc-1", "http://localhost:8080")
	err := svc.Validate()
	fmt.Println("Valid:", err == nil)

	// Invalid service (empty ID)
	invalid := discover.NewService("", "")
	err = invalid.Validate()
	fmt.Println("Invalid error:", err != nil)
	// Output:
	// Valid: true
	// Invalid error: true
}

func ExampleNewMemory() {
	discovery := discover.NewMemory()

	fmt.Printf("Type: %T\n", discovery)
	// Output:
	// Type: *discover.MemoryDiscovery
}

func ExampleMemoryDiscovery_Register() {
	discovery := discover.NewMemory()
	ctx := context.Background()

	svc := discover.NewService("svc-1", "http://localhost:8080")
	err := discovery.Register(ctx, svc)

	fmt.Println("Register error:", err)
	// Output:
	// Register error: <nil>
}

func ExampleMemoryDiscovery_Register_duplicate() {
	discovery := discover.NewMemory()
	ctx := context.Background()

	svc := discover.NewService("svc-1", "http://localhost:8080")
	_ = discovery.Register(ctx, svc)

	// Try to register again
	err := discovery.Register(ctx, svc)
	fmt.Println("Is duplicate:", errors.Is(err, discover.ErrDuplicate))
	// Output:
	// Is duplicate: true
}

func ExampleMemoryDiscovery_Deregister() {
	discovery := discover.NewMemory()
	ctx := context.Background()

	svc := discover.NewService("svc-1", "http://localhost:8080")
	_ = discovery.Register(ctx, svc)

	err := discovery.Deregister(ctx, "svc-1")
	fmt.Println("Deregister error:", err)

	// Verify it's gone
	_, err = discovery.Get(ctx, "svc-1")
	fmt.Println("Is not found:", errors.Is(err, discover.ErrNotFound))
	// Output:
	// Deregister error: <nil>
	// Is not found: true
}

func ExampleMemoryDiscovery_Get() {
	discovery := discover.NewMemory()
	ctx := context.Background()

	svc := discover.NewService("svc-1", "http://localhost:8080").
		SetName("My Service")
	_ = discovery.Register(ctx, svc)

	found, err := discovery.Get(ctx, "svc-1")
	fmt.Println("Error:", err)
	fmt.Println("Name:", found.Name())
	// Output:
	// Error: <nil>
	// Name: My Service
}

func ExampleMemoryDiscovery_Get_notFound() {
	discovery := discover.NewMemory()
	ctx := context.Background()

	_, err := discovery.Get(ctx, "nonexistent")
	fmt.Println("Is not found:", errors.Is(err, discover.ErrNotFound))
	// Output:
	// Is not found: true
}

func ExampleMemoryDiscovery_List() {
	discovery := discover.NewMemory()
	ctx := context.Background()

	_ = discovery.Register(ctx, discover.NewService("svc-a", "http://a:8080"))
	_ = discovery.Register(ctx, discover.NewService("svc-b", "http://b:8080"))
	_ = discovery.Register(ctx, discover.NewService("svc-c", "http://c:8080"))

	services, err := discovery.List(ctx, nil)
	fmt.Println("Error:", err)
	fmt.Println("Count:", len(services))
	// Output:
	// Error: <nil>
	// Count: 3
}

func ExampleMemoryDiscovery_List_withFilter() {
	discovery := discover.NewMemory()
	ctx := context.Background()

	_ = discovery.Register(ctx, discover.NewService("svc-1", "http://a:8080").WithCapability("tools"))
	_ = discovery.Register(ctx, discover.NewService("svc-2", "http://b:8080").WithCapability("resources"))
	_ = discovery.Register(ctx, discover.NewService("svc-3", "http://c:8080").WithCapability("tools"))

	filter := &discover.Filter{
		Capabilities: []string{"tools"},
	}
	services, _ := discovery.List(ctx, filter)
	fmt.Println("Services with tools:", len(services))
	// Output:
	// Services with tools: 2
}

func ExampleMemoryDiscovery_List_withLimit() {
	discovery := discover.NewMemory()
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		_ = discovery.Register(ctx, discover.NewService(fmt.Sprintf("svc-%d", i), "http://localhost"))
	}

	filter := &discover.Filter{Limit: 3}
	services, _ := discovery.List(ctx, filter)
	fmt.Println("Limited count:", len(services))
	// Output:
	// Limited count: 3
}

func ExampleMemoryDiscovery_Negotiate() {
	discovery := discover.NewMemory()
	ctx := context.Background()

	client := &discover.Capabilities{
		Tools:     true,
		Streaming: true,
		Resources: true,
	}
	server := &discover.Capabilities{
		Tools:     true,
		Streaming: false,
		Prompts:   true,
	}

	result, err := discovery.Negotiate(ctx, client, server)
	fmt.Println("Error:", err)
	fmt.Println("Tools:", result.Tools)
	fmt.Println("Streaming:", result.Streaming)
	fmt.Println("Resources:", result.Resources)
	fmt.Println("Prompts:", result.Prompts)
	// Output:
	// Error: <nil>
	// Tools: true
	// Streaming: false
	// Resources: false
	// Prompts: false
}

func ExampleCapabilities_Merge() {
	caps1 := &discover.Capabilities{
		Tools:      true,
		Extensions: []string{"ext-a"},
	}
	caps2 := &discover.Capabilities{
		Resources:  true,
		Extensions: []string{"ext-b"},
	}

	merged := caps1.Merge(caps2)
	fmt.Println("Tools:", merged.Tools)
	fmt.Println("Resources:", merged.Resources)
	fmt.Println("Extensions:", len(merged.Extensions))
	// Output:
	// Tools: true
	// Resources: true
	// Extensions: 2
}

func ExampleCapabilities_Intersect() {
	caps1 := &discover.Capabilities{
		Tools:      true,
		Resources:  true,
		Extensions: []string{"ext-a", "ext-b"},
	}
	caps2 := &discover.Capabilities{
		Tools:      true,
		Streaming:  true,
		Extensions: []string{"ext-b", "ext-c"},
	}

	intersected := caps1.Intersect(caps2)
	fmt.Println("Tools:", intersected.Tools)
	fmt.Println("Resources:", intersected.Resources)
	fmt.Println("Streaming:", intersected.Streaming)
	fmt.Println("Extensions:", intersected.Extensions)
	// Output:
	// Tools: true
	// Resources: false
	// Streaming: false
	// Extensions: [ext-b]
}

func ExampleFilter_Matches() {
	filter := &discover.Filter{
		Capabilities: []string{"tools", "streaming"},
	}

	svc1 := discover.NewService("svc-1", "http://localhost").
		WithCapability("tools").
		WithCapability("streaming")
	svc2 := discover.NewService("svc-2", "http://localhost").
		WithCapability("tools")

	fmt.Println("svc1 matches:", filter.Matches(svc1))
	fmt.Println("svc2 matches:", filter.Matches(svc2))
	// Output:
	// svc1 matches: true
	// svc2 matches: false
}

func ExampleNegotiator() {
	// Default strategy: Intersect
	negotiator := &discover.Negotiator{
		Strategy: discover.NegotiateIntersect,
	}

	client := &discover.Capabilities{Tools: true, Resources: true}
	server := &discover.Capabilities{Tools: true, Streaming: true}

	result := negotiator.Negotiate(client, server)
	fmt.Println("Tools:", result.Tools)
	fmt.Println("Resources:", result.Resources)
	fmt.Println("Streaming:", result.Streaming)
	// Output:
	// Tools: true
	// Resources: false
	// Streaming: false
}

func ExampleNegotiator_merge() {
	negotiator := &discover.Negotiator{
		Strategy: discover.NegotiateMerge,
	}

	client := &discover.Capabilities{Tools: true}
	server := &discover.Capabilities{Resources: true}

	result := negotiator.Negotiate(client, server)
	fmt.Println("Tools:", result.Tools)
	fmt.Println("Resources:", result.Resources)
	// Output:
	// Tools: true
	// Resources: true
}

func Example_discoveryWorkflow() {
	discovery := discover.NewMemory()
	ctx := context.Background()

	// Register services
	toolService := discover.NewService("tool-server", "http://tools:8080").
		SetName("Tool Server").
		SetVersion("1.0.0").
		WithCapability("tools").
		WithCapability("streaming")

	resourceService := discover.NewService("resource-server", "http://resources:8080").
		SetName("Resource Server").
		SetVersion("2.0.0").
		WithCapability("resources")

	_ = discovery.Register(ctx, toolService)
	_ = discovery.Register(ctx, resourceService)

	// Find services with tools capability
	filter := &discover.Filter{Capabilities: []string{"tools"}}
	services, _ := discovery.List(ctx, filter)

	fmt.Println("Tool services found:", len(services))
	fmt.Println("First service:", services[0].Name())
	// Output:
	// Tool services found: 1
	// First service: Tool Server
}
