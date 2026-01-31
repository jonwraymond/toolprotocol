package prompt

import "context"

// Role represents the role of a message sender.
type Role string

const (
	// RoleUser indicates a user message.
	RoleUser Role = "user"

	// RoleAssistant indicates an assistant message.
	RoleAssistant Role = "assistant"
)

// String returns the string representation of the role.
func (r Role) String() string {
	return string(r)
}

// Valid returns true if the role is a known valid role.
func (r Role) Valid() bool {
	switch r {
	case RoleUser, RoleAssistant:
		return true
	default:
		return false
	}
}

// ContentType represents the type of message content.
type ContentType string

const (
	// ContentText indicates text content.
	ContentText ContentType = "text"

	// ContentImage indicates image content.
	ContentImage ContentType = "image"

	// ContentResource indicates an embedded resource reference.
	ContentResource ContentType = "resource"
)

// String returns the string representation of the content type.
func (ct ContentType) String() string {
	return string(ct)
}

// Valid returns true if the content type is a known valid type.
func (ct ContentType) Valid() bool {
	switch ct {
	case ContentText, ContentImage, ContentResource:
		return true
	default:
		return false
	}
}

// Argument describes a prompt argument.
type Argument struct {
	// Name is the argument name.
	Name string

	// Description describes the argument.
	Description string

	// Required indicates if the argument is required.
	Required bool
}

// Prompt represents an MCP prompt template.
type Prompt struct {
	// Name is the unique prompt name.
	Name string

	// Description describes the prompt.
	Description string

	// Arguments defines the prompt arguments.
	Arguments []Argument
}

// Validate validates the prompt.
func (p *Prompt) Validate() error {
	if p.Name == "" {
		return ErrInvalidPrompt
	}
	return nil
}

// RequiredArgs returns the names of required arguments.
func (p *Prompt) RequiredArgs() []string {
	var required []string
	for _, arg := range p.Arguments {
		if arg.Required {
			required = append(required, arg.Name)
		}
	}
	return required
}

// Content represents message content.
type Content struct {
	// Type is the content type.
	Type ContentType

	// Text is the text content (for ContentText).
	Text string

	// MIMEType is the MIME type (for ContentImage).
	MIMEType string

	// Data is binary data (for ContentImage).
	Data []byte

	// Resource is a resource reference (for ContentResource).
	Resource *ResourceRef
}

// IsText returns true if this is text content.
func (c *Content) IsText() bool {
	return c.Type == ContentText
}

// ResourceRef references an MCP resource.
type ResourceRef struct {
	// URI is the resource URI.
	URI string
}

// Message represents a generated message.
type Message struct {
	// Role is the message role.
	Role Role

	// Content is the message content.
	Content []Content
}

// Provider serves prompts.
type Provider interface {
	// List returns all available prompts.
	List(ctx context.Context) ([]Prompt, error)

	// Get returns messages for a prompt with the given arguments.
	Get(ctx context.Context, name string, args map[string]string) ([]Message, error)
}

// PromptHandler generates messages from arguments.
type PromptHandler func(ctx context.Context, args map[string]string) ([]Message, error)
