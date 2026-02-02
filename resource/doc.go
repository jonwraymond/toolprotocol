// Package resource provides MCP resource management.
//
// This package enables servers to expose resources that clients can
// read and subscribe to for updates. Resources represent data with
// URIs, supporting text and binary content.
//
// # Ecosystem Position
//
// resource provides URI-based resource access and subscriptions:
//
//	┌─────────────────────────────────────────────────────────────────┐
//	│                      Resource Flow                              │
//	├─────────────────────────────────────────────────────────────────┤
//	│                                                                 │
//	│   Provider             resource                  Client         │
//	│   ┌─────────┐        ┌───────────┐         ┌─────────┐         │
//	│   │ Static  │────────│ Registry  │─────────│  Read   │         │
//	│   │  File   │ Register           │  Read   │         │         │
//	│   └─────────┘        │ ┌───────┐ │         └─────────┘         │
//	│        │             │ │Provider│ │              │               │
//	│        │             │ │  Map  │ │              │               │
//	│        ▼             │ └───────┘ │              │               │
//	│   ┌─────────┐        │     │     │              │               │
//	│   │  List   │────────│─────┼─────│──────────────▼               │
//	│   │Templates│        │     │     │         ┌─────────┐         │
//	│   └─────────┘        │     ▼     │         │Subscribe│         │
//	│        │             │ ┌───────┐ │         │ Manager │         │
//	│        ▼             │ │ Subs  │ │         └─────────┘         │
//	│   ┌─────────┐        │ └───────┘ │              ▲               │
//	│   │ Update  │────────│───────────│──────────────┘               │
//	│   │ Notify  │        └───────────┘   Notify                     │
//	│   └─────────┘                                                   │
//	│                                                                 │
//	└─────────────────────────────────────────────────────────────────┘
//
// # Core Components
//
//   - [Resource]: Resource metadata (URI, name, description, MIME type)
//   - [Contents]: Resource contents (text or binary)
//   - [Template]: URI template for parameterized resources
//   - [Provider]: Interface for serving resources
//   - [Subscriber]: Interface for receiving resource updates
//   - [Registry]: Multi-provider resource manager
//   - [StaticProvider]: In-memory static resource provider
//   - [SubscriptionManager]: Resource update subscription handler
//
// # Quick Start
//
//	// Create registry and provider
//	registry := resource.NewRegistry()
//	provider := resource.NewStaticProvider()
//
//	// Add resources
//	provider.Add(
//	    resource.Resource{URI: "file:///config.json", Name: "Config"},
//	    resource.Contents{URI: "file:///config.json", Text: `{"debug":true}`},
//	)
//
//	// Register provider for scheme
//	registry.Register("file", provider)
//
//	// Read resource
//	contents, err := registry.Read(ctx, "file:///config.json")
//
// # Resource Structure
//
// A resource has a URI, name, description, and MIME type:
//
//	res := resource.Resource{
//	    URI:         "file:///path/to/doc.md",
//	    Name:        "Documentation",
//	    Description: "Project documentation",
//	    MIMEType:    "text/markdown",
//	    Annotations: map[string]any{"priority": "high"},
//	}
//
// # Contents
//
// Resource contents can be text or binary:
//
//	// Text content
//	textContents := resource.Contents{
//	    URI:      "file:///doc.md",
//	    MIMEType: "text/markdown",
//	    Text:     "# Hello World",
//	}
//
//	// Binary content
//	binaryContents := resource.Contents{
//	    URI:      "file:///image.png",
//	    MIMEType: "image/png",
//	    Blob:     imageBytes,
//	}
//
// # URI Templates
//
// Templates support parameterized resource URIs:
//
//	tmpl := resource.Template{
//	    URITemplate: "file:///users/{userId}/profile",
//	    Name:        "User Profile",
//	}
//	uri := tmpl.Expand(map[string]string{"userId": "123"})
//	// Result: "file:///users/123/profile"
//
// # Subscriptions
//
// Clients can subscribe to resource updates:
//
//	mgr := resource.NewSubscriptionManager()
//	ch, _ := mgr.Subscribe(ctx, "file:///config.json")
//
//	go func() {
//	    for contents := range ch {
//	        fmt.Printf("Updated: %s\n", contents.Text)
//	    }
//	}()
//
//	// Notify subscribers of update
//	mgr.Notify("file:///config.json", &resource.Contents{Text: "new config"})
//
// # Thread Safety
//
// All exported types are safe for concurrent use:
//
//   - [Registry]: sync.RWMutex protects provider map
//   - [StaticProvider]: sync.RWMutex protects resources and contents
//   - [SubscriptionManager]: sync.RWMutex protects subscription map
//   - [Resource], [Contents], [Template]: Value types, safe to copy
//
// # Error Handling
//
// Sentinel errors (use errors.Is for checking):
//
//   - [ErrResourceNotFound]: Resource not found
//   - [ErrProviderNotFound]: No provider for URI scheme
//   - [ErrInvalidURI]: URI format is invalid
//   - [ErrDuplicateProvider]: Provider already registered for scheme
//   - [ErrNotSubscribed]: Not subscribed to resource
//
// The [ResourceError] type wraps errors with context:
//
//	err := &ResourceError{
//	    URI: "file:///missing.txt",
//	    Op:  "read",
//	    Err: ErrResourceNotFound,
//	}
//	// errors.Is(err, ErrResourceNotFound) = true
//
// # Integration with ApertureStack
//
// resource integrates with other ApertureStack packages:
//
//   - content: Resource contents map to Content types
//   - wire: Resources encoded to protocol-specific formats
//   - stream: Resource updates delivered via streaming
package resource
