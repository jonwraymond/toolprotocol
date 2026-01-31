// Package prompt provides MCP prompt template management.
//
// This package enables servers to expose prompt templates that clients
// can use to generate messages with dynamic argument substitution.
//
// # Prompt Structure
//
// A prompt defines a template with arguments:
//
//	prompt := Prompt{
//	    Name:        "greeting",
//	    Description: "A friendly greeting",
//	    Arguments: []Argument{
//	        {Name: "name", Description: "Name to greet", Required: true},
//	    },
//	}
//
// # Message Generation
//
// Prompts generate messages with content:
//
//	messages := []Message{
//	    {
//	        Role: RoleUser,
//	        Content: []Content{
//	            {Type: ContentText, Text: "Hello, Alice!"},
//	        },
//	    },
//	}
//
// # Provider Interface
//
// The Provider interface serves prompts:
//
//	type Provider interface {
//	    List(ctx context.Context) ([]Prompt, error)
//	    Get(ctx context.Context, name string, args map[string]string) ([]Message, error)
//	}
//
// # Registry Usage
//
//	// Create registry
//	registry := prompt.NewRegistry()
//
//	// Register a prompt
//	registry.Register(prompt.Prompt{
//	    Name: "greeting",
//	    Arguments: []prompt.Argument{{Name: "name", Required: true}},
//	}, func(ctx context.Context, args map[string]string) ([]prompt.Message, error) {
//	    return []prompt.Message{
//	        prompt.NewUserMessage(prompt.TextContent("Hello, " + args["name"] + "!")),
//	    }, nil
//	})
//
//	// Get prompt messages
//	msgs, err := registry.Get(ctx, "greeting", map[string]string{"name": "Alice"})
package prompt
