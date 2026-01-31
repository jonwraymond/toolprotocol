package prompt

// NewUserMessage creates a new user message with the given content.
func NewUserMessage(content ...Content) Message {
	return Message{
		Role:    RoleUser,
		Content: content,
	}
}

// NewAssistantMessage creates a new assistant message with the given content.
func NewAssistantMessage(content ...Content) Message {
	return Message{
		Role:    RoleAssistant,
		Content: content,
	}
}

// TextContent creates text content.
func TextContent(text string) Content {
	return Content{
		Type: ContentText,
		Text: text,
	}
}

// ImageContent creates image content.
func ImageContent(mimeType string, data []byte) Content {
	return Content{
		Type:     ContentImage,
		MIMEType: mimeType,
		Data:     data,
	}
}

// ResourceContent creates resource reference content.
func ResourceContent(uri string) Content {
	return Content{
		Type:     ContentResource,
		Resource: &ResourceRef{URI: uri},
	}
}

// PromptBuilder provides a fluent interface for building prompts.
type PromptBuilder struct {
	prompt Prompt
}

// NewPromptBuilder creates a new prompt builder.
func NewPromptBuilder(name string) *PromptBuilder {
	return &PromptBuilder{
		prompt: Prompt{Name: name},
	}
}

// WithDescription sets the prompt description.
func (b *PromptBuilder) WithDescription(desc string) *PromptBuilder {
	b.prompt.Description = desc
	return b
}

// WithArgument adds an argument to the prompt.
func (b *PromptBuilder) WithArgument(name, description string, required bool) *PromptBuilder {
	b.prompt.Arguments = append(b.prompt.Arguments, Argument{
		Name:        name,
		Description: description,
		Required:    required,
	})
	return b
}

// WithRequiredArgument adds a required argument to the prompt.
func (b *PromptBuilder) WithRequiredArgument(name, description string) *PromptBuilder {
	return b.WithArgument(name, description, true)
}

// Build returns the built prompt.
func (b *PromptBuilder) Build() Prompt {
	return b.prompt
}
