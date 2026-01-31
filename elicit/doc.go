// Package elicit provides user input elicitation for MCP protocol servers.
//
// This package enables servers to request input from users during tool
// execution, supporting text input, confirmations, choices, and forms.
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
// # Elicitor Interface
//
// The Elicitor interface sends requests to clients:
//
//	type Elicitor interface {
//	    Elicit(ctx context.Context, req *Request) (*Response, error)
//	}
//
// # Handler Interface
//
// The Handler interface processes requests (client-side):
//
//	type Handler interface {
//	    Handle(ctx context.Context, req *Request) (*Response, error)
//	}
//
// # Usage
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
package elicit
