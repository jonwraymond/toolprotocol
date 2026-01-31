// Package discover provides service discovery and capability negotiation
// for tool providers across protocols.
//
// This package enables registration and lookup of discoverable services,
// as well as capability negotiation between clients and servers.
//
// # Discoverable Interface
//
// Services implement the Discoverable interface:
//
//	type Discoverable interface {
//	    ID() string
//	    Name() string
//	    Description() string
//	    Version() string
//	    Endpoint() string
//	    Capabilities() *Capabilities
//	}
//
// # Discovery Interface
//
// The Discovery interface manages service registration and lookup:
//
//	type Discovery interface {
//	    Register(ctx context.Context, svc Discoverable) error
//	    Deregister(ctx context.Context, id string) error
//	    Get(ctx context.Context, id string) (Discoverable, error)
//	    List(ctx context.Context, filter *Filter) ([]Discoverable, error)
//	    Negotiate(ctx context.Context, client, server *Capabilities) (*Capabilities, error)
//	}
//
// # Usage
//
//	discovery := discover.NewMemory()
//	svc := discover.NewService("my-service", "http://localhost:8080")
//	err := discovery.Register(ctx, svc)
package discover
