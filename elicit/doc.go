// Package elicit provides user input elicitation for MCP protocol servers.
//
// This package enables servers to request input from users during tool
// execution, supporting text input, confirmations, choices, and forms.
//
// # Ecosystem Position
//
// elicit provides interactive user input collection:
//
//	┌─────────────────────────────────────────────────────────────────┐
//	│                     Elicitation Flow                            │
//	├─────────────────────────────────────────────────────────────────┤
//	│                                                                 │
//	│   Server               elicit                    Client         │
//	│   ┌─────────┐        ┌───────────┐         ┌─────────┐         │
//	│   │  Tool   │────────│  Request  │─────────│ Handler │         │
//	│   │Execution│        │           │         │         │         │
//	│   └─────────┘        │ ┌───────┐ │         └─────────┘         │
//	│        │             │ │ Type  │ │              │               │
//	│        │             │ │Message│ │              │               │
//	│        │             │ └───────┘ │              ▼               │
//	│        │             │     │     │         ┌─────────┐         │
//	│        ▼             │     │     │         │  User   │         │
//	│   ┌─────────┐        │     ▼     │         │  Input  │         │
//	│   │Elicitor │────────│ Validate │         └─────────┘         │
//	│   │Interface│        │     │     │              │               │
//	│   └─────────┘        │     │     │              ▼               │
//	│        │             │     ▼     │         ┌─────────┐         │
//	│        │             │ ┌───────┐ │         │Response │         │
//	│        ▼             │ │Handle │◄│─────────│  Value  │         │
//	│   ┌─────────┐        │ └───────┘ │         └─────────┘         │
//	│   │Response │◄───────│           │                              │
//	│   └─────────┘        └───────────┘                              │
//	│                                                                 │
//	└─────────────────────────────────────────────────────────────────┘
//
// # Request Types
//
// The package supports four types of elicitation requests:
//
//   - text: Free-form text input
//   - confirmation: Yes/no confirmation
//   - choice: Selection from predefined options
//   - form: Structured input with JSON Schema validation
//
// # Core Components
//
//   - [Request]: Elicitation request with type, message, and options
//   - [Response]: User response with value or cancellation status
//   - [Elicitor]: Interface for sending requests (server-side)
//   - [Handler]: Interface for processing requests (client-side)
//   - [DefaultHandler]: Basic handler implementation with callbacks
//   - [Builder]: Fluent API for building requests
//
// # Quick Start
//
// Create and send elicitation requests:
//
//	// Create a text request
//	req := elicit.NewTextRequest("Enter your name:")
//
//	// Create a choice request
//	req := elicit.NewChoiceRequest("Select an option:", []elicit.Choice{
//	    {ID: "a", Label: "Option A"},
//	    {ID: "b", Label: "Option B"},
//	})
//
//	// Send the request
//	resp, err := elicitor.Elicit(ctx, req)
//	if err != nil {
//	    return err
//	}
//
//	if resp.Cancelled {
//	    // User cancelled
//	}
//	if resp.TimedOut {
//	    // Request timed out
//	}
//
//	// Use the response value
//	name := resp.Value.(string)
//
// # Request Building
//
// Use the fluent builder for complex requests:
//
//	req := elicit.NewBuilder(elicit.TypeChoice, "Select language:").
//	    WithChoices([]elicit.Choice{
//	        {ID: "go", Label: "Go", Description: "Fast compiled language"},
//	        {ID: "py", Label: "Python", Description: "Dynamic scripting"},
//	    }).
//	    WithTimeout(30 * time.Second).
//	    WithDefault("go").
//	    Build()
//
// # Handler Implementation
//
// Create a handler to process requests:
//
//	handler := elicit.NewHandler(
//	    elicit.WithDefaultTimeout(30 * time.Second),
//	    elicit.WithCallback(func(ctx context.Context, req *elicit.Request) (*elicit.Response, error) {
//	        // Process request and get user input
//	        value := getUserInput(req.Message)
//	        return &elicit.Response{
//	            RequestID: req.ID,
//	            Value:     value,
//	        }, nil
//	    }),
//	)
//
// # Thread Safety
//
// All exported types are safe for concurrent use:
//
//   - [DefaultHandler]: sync.RWMutex protects handler and timeout fields
//   - [Builder]: Not thread-safe; use one builder per goroutine
//   - [Request], [Response]: Value types, safe to copy
//
// # Error Handling
//
// Sentinel errors (use errors.Is for checking):
//
//   - [ErrInvalidRequest]: Request validation failed
//   - [ErrTimeout]: Request timed out waiting for response
//   - [ErrCancelled]: User cancelled the request
//   - [ErrNoHandler]: No handler configured to process requests
//
// The [ElicitError] type wraps errors with context:
//
//	err := &ElicitError{
//	    RequestID: "req-123",
//	    Op:        "validate",
//	    Err:       ErrInvalidRequest,
//	}
//	// errors.Is(err, ErrInvalidRequest) = true
//
// # Response States
//
// Check response states before using values:
//
//	resp, err := handler.Handle(ctx, req)
//	if err != nil {
//	    return err
//	}
//
//	if !resp.IsSuccess() {
//	    if resp.Cancelled {
//	        return errors.New("user cancelled")
//	    }
//	    if resp.TimedOut {
//	        return errors.New("request timed out")
//	    }
//	}
//
//	// Safe to use resp.Value
//
// # Integration with ApertureStack
//
// elicit integrates with other ApertureStack packages:
//
//   - transport: Requests sent over wire protocols
//   - wire: Elicitation messages encoded to protocol formats
//   - session: Request IDs can be tracked per session
package elicit
