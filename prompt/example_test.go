package prompt_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonwraymond/toolprotocol/prompt"
)

func ExampleNewRegistry() {
	registry := prompt.NewRegistry()

	fmt.Printf("Type: %T\n", registry)
	// Output:
	// Type: *prompt.Registry
}

func ExampleRegistry_Register() {
	registry := prompt.NewRegistry()

	p := prompt.Prompt{
		Name:        "greeting",
		Description: "A friendly greeting",
		Arguments: []prompt.Argument{
			{Name: "name", Description: "Name to greet", Required: true},
		},
	}

	handler := func(ctx context.Context, args map[string]string) ([]prompt.Message, error) {
		return []prompt.Message{
			prompt.NewUserMessage(prompt.TextContent("Hello, " + args["name"] + "!")),
		}, nil
	}

	err := registry.Register(p, handler)
	fmt.Println("Register error:", err)
	// Output:
	// Register error: <nil>
}

func ExampleRegistry_Register_duplicate() {
	registry := prompt.NewRegistry()
	p := prompt.Prompt{Name: "greeting"}
	handler := func(ctx context.Context, args map[string]string) ([]prompt.Message, error) {
		return nil, nil
	}

	_ = registry.Register(p, handler)
	err := registry.Register(p, handler)
	fmt.Println("Is duplicate:", errors.Is(err, prompt.ErrDuplicatePrompt))
	// Output:
	// Is duplicate: true
}

func ExampleRegistry_List() {
	registry := prompt.NewRegistry()
	handler := func(ctx context.Context, args map[string]string) ([]prompt.Message, error) {
		return nil, nil
	}

	_ = registry.Register(prompt.Prompt{Name: "greeting"}, handler)
	_ = registry.Register(prompt.Prompt{Name: "farewell"}, handler)

	ctx := context.Background()
	prompts, _ := registry.List(ctx)
	fmt.Println("Count:", len(prompts))
	// Output:
	// Count: 2
}

func ExampleRegistry_Get() {
	registry := prompt.NewRegistry()

	p := prompt.Prompt{
		Name: "greeting",
		Arguments: []prompt.Argument{
			{Name: "name", Required: true},
		},
	}

	handler := func(ctx context.Context, args map[string]string) ([]prompt.Message, error) {
		return []prompt.Message{
			prompt.NewUserMessage(prompt.TextContent("Hello, " + args["name"] + "!")),
		}, nil
	}

	_ = registry.Register(p, handler)

	ctx := context.Background()
	messages, err := registry.Get(ctx, "greeting", map[string]string{"name": "Alice"})
	fmt.Println("Error:", err)
	fmt.Println("Message count:", len(messages))
	fmt.Println("Content:", messages[0].Content[0].Text)
	// Output:
	// Error: <nil>
	// Message count: 1
	// Content: Hello, Alice!
}

func ExampleRegistry_Get_notFound() {
	registry := prompt.NewRegistry()
	ctx := context.Background()

	_, err := registry.Get(ctx, "nonexistent", nil)
	fmt.Println("Is not found:", errors.Is(err, prompt.ErrPromptNotFound))
	// Output:
	// Is not found: true
}

func ExampleRegistry_Get_missingArgument() {
	registry := prompt.NewRegistry()
	p := prompt.Prompt{
		Name: "greeting",
		Arguments: []prompt.Argument{
			{Name: "name", Required: true},
		},
	}
	handler := func(ctx context.Context, args map[string]string) ([]prompt.Message, error) {
		return nil, nil
	}
	_ = registry.Register(p, handler)

	ctx := context.Background()
	_, err := registry.Get(ctx, "greeting", map[string]string{})
	fmt.Println("Is missing arg:", errors.Is(err, prompt.ErrMissingArgument))
	// Output:
	// Is missing arg: true
}

func ExampleNewPromptBuilder() {
	builder := prompt.NewPromptBuilder("greeting")

	fmt.Printf("Type: %T\n", builder)
	// Output:
	// Type: *prompt.PromptBuilder
}

func ExamplePromptBuilder_Build() {
	p := prompt.NewPromptBuilder("greeting").
		WithDescription("A friendly greeting").
		WithRequiredArgument("name", "Name to greet").
		WithArgument("title", "Optional title", false).
		Build()

	fmt.Println("Name:", p.Name)
	fmt.Println("Description:", p.Description)
	fmt.Println("Args count:", len(p.Arguments))
	// Output:
	// Name: greeting
	// Description: A friendly greeting
	// Args count: 2
}

func ExampleNewUserMessage() {
	msg := prompt.NewUserMessage(prompt.TextContent("Hello!"))

	fmt.Println("Role:", msg.Role)
	fmt.Println("Content:", msg.Content[0].Text)
	// Output:
	// Role: user
	// Content: Hello!
}

func ExampleNewAssistantMessage() {
	msg := prompt.NewAssistantMessage(prompt.TextContent("Hi there!"))

	fmt.Println("Role:", msg.Role)
	fmt.Println("Content:", msg.Content[0].Text)
	// Output:
	// Role: assistant
	// Content: Hi there!
}

func ExampleTextContent() {
	content := prompt.TextContent("Hello, world!")

	fmt.Println("Type:", content.Type)
	fmt.Println("Text:", content.Text)
	// Output:
	// Type: text
	// Text: Hello, world!
}

func ExampleImageContent() {
	content := prompt.ImageContent("image/png", []byte{0x89, 0x50})

	fmt.Println("Type:", content.Type)
	fmt.Println("MIME:", content.MIMEType)
	// Output:
	// Type: image
	// MIME: image/png
}

func ExampleResourceContent() {
	content := prompt.ResourceContent("file:///doc.txt")

	fmt.Println("Type:", content.Type)
	fmt.Println("URI:", content.Resource.URI)
	// Output:
	// Type: resource
	// URI: file:///doc.txt
}

func ExampleRole_String() {
	fmt.Println(prompt.RoleUser.String())
	fmt.Println(prompt.RoleAssistant.String())
	// Output:
	// user
	// assistant
}

func ExampleRole_Valid() {
	fmt.Println("user valid:", prompt.RoleUser.Valid())
	fmt.Println("assistant valid:", prompt.RoleAssistant.Valid())
	fmt.Println("unknown valid:", prompt.Role("unknown").Valid())
	// Output:
	// user valid: true
	// assistant valid: true
	// unknown valid: false
}

func ExampleContentType_String() {
	fmt.Println(prompt.ContentText.String())
	fmt.Println(prompt.ContentImage.String())
	fmt.Println(prompt.ContentResource.String())
	// Output:
	// text
	// image
	// resource
}

func ExampleContentType_Valid() {
	fmt.Println("text valid:", prompt.ContentText.Valid())
	fmt.Println("image valid:", prompt.ContentImage.Valid())
	fmt.Println("unknown valid:", prompt.ContentType("unknown").Valid())
	// Output:
	// text valid: true
	// image valid: true
	// unknown valid: false
}

func ExamplePrompt_RequiredArgs() {
	p := prompt.Prompt{
		Name: "greeting",
		Arguments: []prompt.Argument{
			{Name: "name", Required: true},
			{Name: "title", Required: false},
			{Name: "greeting", Required: true},
		},
	}

	required := p.RequiredArgs()
	fmt.Println("Required args:", required)
	// Output:
	// Required args: [name greeting]
}

func ExamplePrompt_Validate() {
	valid := prompt.Prompt{Name: "greeting"}
	invalid := prompt.Prompt{Name: ""}

	fmt.Println("Valid error:", valid.Validate())
	fmt.Println("Invalid has error:", invalid.Validate() != nil)
	// Output:
	// Valid error: <nil>
	// Invalid has error: true
}

func ExamplePromptError() {
	err := &prompt.PromptError{
		PromptName: "greeting",
		Op:         "get",
		Err:        prompt.ErrPromptNotFound,
	}

	fmt.Println(err.Error())
	fmt.Println("Unwraps:", errors.Is(err, prompt.ErrPromptNotFound))
	// Output:
	// prompt greeting: get: prompt: not found
	// Unwraps: true
}

func Example_promptWorkflow() {
	registry := prompt.NewRegistry()

	// Build and register prompt
	p := prompt.NewPromptBuilder("code-review").
		WithDescription("Generate a code review").
		WithRequiredArgument("code", "Code to review").
		WithArgument("style", "Review style", false).
		Build()

	handler := func(ctx context.Context, args map[string]string) ([]prompt.Message, error) {
		style := args["style"]
		if style == "" {
			style = "constructive"
		}
		return []prompt.Message{
			prompt.NewUserMessage(
				prompt.TextContent("Please review this code with a "+style+" approach:"),
				prompt.TextContent(args["code"]),
			),
		}, nil
	}

	_ = registry.Register(p, handler)

	// Get prompt messages
	ctx := context.Background()
	messages, _ := registry.Get(ctx, "code-review", map[string]string{
		"code":  "func main() {}",
		"style": "detailed",
	})

	fmt.Println("Role:", messages[0].Role)
	fmt.Println("Content count:", len(messages[0].Content))
	// Output:
	// Role: user
	// Content count: 2
}
