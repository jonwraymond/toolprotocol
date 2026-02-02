// Package discover provides service discovery and capability negotiation
// for tool providers across protocols.
//
// This package enables registration and lookup of discoverable services,
// as well as capability negotiation between clients and servers.
//
// # Ecosystem Position
//
// discover provides service registration and capability negotiation:
//
//	┌─────────────────────────────────────────────────────────────────┐
//	│                      Discovery Flow                             │
//	├─────────────────────────────────────────────────────────────────┤
//	│                                                                 │
//	│   Service              discover                  Client         │
//	│   ┌─────────┐        ┌───────────┐         ┌─────────┐         │
//	│   │Provider │────────│ Register  │─────────│ Lookup  │         │
//	│   │         │        │           │         │         │         │
//	│   └─────────┘        │ ┌───────┐ │         └─────────┘         │
//	│        │             │ │Memory │ │              │               │
//	│        │             │ │ Store │ │              │               │
//	│        ▼             │ └───────┘ │              │               │
//	│   ┌─────────┐        │     │     │              │               │
//	│   │  Caps   │────────│─────┼─────│──────────────▼               │
//	│   │announce │        │     │     │         ┌─────────┐         │
//	│   └─────────┘        │     ▼     │         │ Filter  │         │
//	│        │             │ ┌───────┐ │         │  List   │         │
//	│        ▼             │ │Negoti-│ │         └─────────┘         │
//	│   ┌─────────┐        │ │ ate   │ │              ▲               │
//	│   │ Ready   │        │ └───────┘ │──────────────┘               │
//	│   └─────────┘        └───────────┘   capabilities               │
//	│                                                                 │
//	└─────────────────────────────────────────────────────────────────┘
//
// # Core Components
//
//   - [Discoverable]: Interface for services that can be discovered
//   - [Discovery]: Interface for service registry operations
//   - [Service]: Default Discoverable implementation with fluent builder
//   - [MemoryDiscovery]: In-memory Discovery implementation
//   - [Capabilities]: Protocol feature flags (tools, resources, streaming, etc.)
//   - [Filter]: Criteria for listing services
//   - [Negotiator]: Configurable capability negotiation
//
// # Quick Start
//
//	// Create discovery registry
//	discovery := discover.NewMemory()
//
//	// Register a service
//	svc := discover.NewService("my-tool", "http://localhost:8080").
//	    SetName("My Tool Server").
//	    SetVersion("1.0.0").
//	    WithCapability("tools").
//	    WithCapability("streaming")
//
//	err := discovery.Register(ctx, svc)
//
//	// Find services
//	filter := &discover.Filter{Capabilities: []string{"tools"}}
//	services, _ := discovery.List(ctx, filter)
//
//	// Negotiate capabilities
//	result, _ := discovery.Negotiate(ctx, clientCaps, serverCaps)
//
// # Service Builder
//
// Service provides a fluent API for building discoverable services:
//
//	svc := discover.NewService("svc-id", "http://endpoint").
//	    SetName("Human Readable Name").
//	    SetDescription("What this service does").
//	    SetVersion("1.0.0").
//	    WithCapability("tools").
//	    WithCapability("streaming").
//	    WithCapability("custom-extension")
//
// # Capabilities
//
// Capabilities represent protocol features:
//
//	caps := &discover.Capabilities{
//	    Tools:      true,  // Tool invocation support
//	    Resources:  true,  // Resource access support
//	    Prompts:    true,  // Prompt template support
//	    Streaming:  true,  // Streaming response support
//	    Sampling:   true,  // LLM sampling support
//	    Progress:   true,  // Progress notification support
//	    Extensions: []string{"custom-feature"},
//	}
//
// # Capability Negotiation
//
// Two negotiation strategies are available:
//
//   - NegotiateIntersect (default): AND logic, only common capabilities
//   - NegotiateMerge: OR logic, combined capabilities
//
// Example:
//
//	negotiator := &discover.Negotiator{Strategy: discover.NegotiateIntersect}
//	result := negotiator.Negotiate(clientCaps, serverCaps)
//
// # Filtering
//
// Filter services by capabilities:
//
//	filter := &discover.Filter{
//	    Capabilities: []string{"tools", "streaming"},
//	    Limit:        10,
//	}
//	services, _ := discovery.List(ctx, filter)
//
// # Thread Safety
//
// All exported types are safe for concurrent use:
//
//   - [MemoryDiscovery]: sync.RWMutex protects all operations
//   - [Service]: Immutable after registration (builder pattern)
//   - [Capabilities]: Value type, safe to share
//   - [Filter]: Value type, safe to share
//
// # Error Handling
//
// Sentinel errors (use errors.Is for checking):
//
//   - [ErrNotFound]: Service not found in registry
//   - [ErrDuplicate]: Service with ID already registered
//   - [ErrInvalidService]: Service fails validation
//
// # Integration with ApertureStack
//
// discover integrates with other ApertureStack packages:
//
//   - transport: Services registered with transport endpoints
//   - wire: Capabilities map to protocol-specific features
//   - session: Client sessions track negotiated capabilities
package discover
