// Package prompt provides MCP prompt template management.
//
// This package enables servers to expose prompt templates that clients
// can use to generate messages with dynamic argument substitution.
//
// # Ecosystem Position
//
// prompt provides template-based message generation:
//
//	┌─────────────────────────────────────────────────────────────────┐
//	│                      Prompt Flow                                │
//	├─────────────────────────────────────────────────────────────────┤
//	│                                                                 │
//	│   Server               prompt                   Client          │
//	│   ┌─────────┐        ┌───────────┐         ┌─────────┐         │
//	│   │Register │────────│ Registry  │─────────│  List   │         │
//	│   │ Prompt  │        │           │         │         │         │
//	│   └─────────┘        │ ┌───────┐ │         └─────────┘         │
//	│        │             │ │Prompts│ │              │               │
//	│        │             │ │  Map  │ │              │               │
//	│        ▼             │ └───────┘ │              │               │
//	│   ┌─────────┐        │     │     │              │               │
//	│   │ Handler │────────│─────┼─────│──────────────▼               │
//	│   │   fn    │        │     │     │         ┌─────────┐         │
//	│   └─────────┘        │     ▼     │         │   Get   │         │
//	│        │             │ ┌───────┐ │         │  (args) │         │
//	│        ▼             │ │Handler│ │         └─────────┘         │
//	│   ┌─────────┐        │ │  Map  │ │              ▲               │
//	│   │Messages │◄───────│ └───────┘ │──────────────┘               │
//	│   └─────────┘        └───────────┘   []Message                  │
//	│                                                                 │
//	└─────────────────────────────────────────────────────────────────┘
//
// # Core Components
//
//   - [Prompt]: Prompt template with name, description, and arguments
//   - [Argument]: Prompt argument definition (name, description, required)
//   - [Message]: Generated message with role and content
//   - [Content]: Message content (text, image, or resource reference)
//   - [Provider]: Interface for serving prompts
//   - [Registry]: Prompt registry with handlers
//   - [PromptBuilder]: Fluent API for building prompts
//
// # Quick Start
//
//	// Create registry
//	registry := prompt.NewRegistry()
//
//	// Build and register prompt
//	p := prompt.NewPromptBuilder("greeting").
//	    WithDescription("A friendly greeting").
//	    WithRequiredArgument("name", "Name to greet").
//	    Build()
//
//	handler := func(ctx context.Context, args map[string]string) ([]prompt.Message, error) {
//	    return []prompt.Message{
//	        prompt.NewUserMessage(prompt.TextContent("Hello, " + args["name"] + "!")),
//	    }, nil
//	}
//
//	registry.Register(p, handler)
//
//	// Get messages
//	messages, err := registry.Get(ctx, "greeting", map[string]string{"name": "Alice"})
//
// # Message Building
//
// Helper functions for creating messages:
//
//	// User message
//	msg := prompt.NewUserMessage(prompt.TextContent("Hello!"))
//
//	// Assistant message
//	msg := prompt.NewAssistantMessage(prompt.TextContent("Hi there!"))
//
//	// Multiple content items
//	msg := prompt.NewUserMessage(
//	    prompt.TextContent("Check this image:"),
//	    prompt.ImageContent("image/png", imageData),
//	)
//
// # Content Types
//
// Three content types are supported:
//
//   - ContentText: Plain text content
//   - ContentImage: Image data with MIME type
//   - ContentResource: Reference to an MCP resource
//
// Example:
//
//	text := prompt.TextContent("Hello, world!")
//	image := prompt.ImageContent("image/png", pngData)
//	resource := prompt.ResourceContent("file:///doc.txt")
//
// # Prompt Builder
//
// Fluent API for building prompts:
//
//	p := prompt.NewPromptBuilder("code-review").
//	    WithDescription("Generate a code review").
//	    WithRequiredArgument("code", "Code to review").
//	    WithArgument("style", "Review style", false).
//	    Build()
//
// # Thread Safety
//
// All exported types are safe for concurrent use:
//
//   - [Registry]: sync.RWMutex protects prompt and handler maps
//   - [PromptBuilder]: Not thread-safe; use one builder per goroutine
//   - [Prompt], [Message], [Content]: Value types, safe to copy
//
// # Error Handling
//
// Sentinel errors (use errors.Is for checking):
//
//   - [ErrPromptNotFound]: Prompt not found in registry
//   - [ErrMissingArgument]: Required argument not provided
//   - [ErrInvalidPrompt]: Prompt validation failed
//   - [ErrHandlerFailed]: Handler returned an error
//   - [ErrDuplicatePrompt]: Prompt already registered
//
// The [PromptError] type wraps errors with context:
//
//	err := &PromptError{
//	    PromptName: "greeting",
//	    Op:         "get",
//	    Err:        ErrPromptNotFound,
//	}
//	// errors.Is(err, ErrPromptNotFound) = true
//
// # Integration with ApertureStack
//
// prompt integrates with other ApertureStack packages:
//
//   - content: Message content maps to Content types
//   - resource: ResourceContent references MCP resources
//   - wire: Prompts encoded to protocol-specific formats
package prompt
