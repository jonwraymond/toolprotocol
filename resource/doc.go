// Package resource provides MCP resource management.
//
// This package enables servers to expose resources that clients can
// read and subscribe to for updates. Resources represent data with
// URIs, supporting text and binary content.
//
// # Resource Structure
//
// A resource has a URI, name, description, and MIME type:
//
//	res := Resource{
//	    URI:         "file:///path/to/doc.md",
//	    Name:        "Documentation",
//	    Description: "Project documentation",
//	    MIMEType:    "text/markdown",
//	}
//
// # Contents
//
// Resource contents can be text or binary:
//
//	// Text content
//	contents := Contents{
//	    URI:      "file:///doc.md",
//	    MIMEType: "text/markdown",
//	    Text:     "# Hello World",
//	}
//
//	// Binary content
//	contents := Contents{
//	    URI:      "file:///image.png",
//	    MIMEType: "image/png",
//	    Blob:     imageBytes,
//	}
//
// # Provider Interface
//
// The Provider interface serves resources:
//
//	type Provider interface {
//	    List(ctx context.Context) ([]Resource, error)
//	    Read(ctx context.Context, uri string) (*Contents, error)
//	    Templates(ctx context.Context) ([]Template, error)
//	}
//
// # Registry Usage
//
//	// Create registry
//	registry := resource.NewRegistry()
//
//	// Register a provider for a URI scheme
//	registry.Register("file", fileProvider)
//
//	// Read a resource
//	contents, err := registry.Read(ctx, "file:///path/to/file.txt")
//
// # Subscriptions
//
// Clients can subscribe to resource updates:
//
//	sub := resource.NewSubscriptionManager()
//	ch, err := sub.Subscribe(ctx, "file:///config.json")
//	for contents := range ch {
//	    fmt.Printf("Config updated: %s\n", contents.Text)
//	}
package resource
