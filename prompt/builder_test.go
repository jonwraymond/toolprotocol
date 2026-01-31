package prompt

import (
	"testing"
)

func TestNewMessage_User(t *testing.T) {
	msg := NewUserMessage(TextContent("Hello"))

	if msg.Role != RoleUser {
		t.Errorf("Role = %v, want %v", msg.Role, RoleUser)
	}
	if len(msg.Content) != 1 {
		t.Errorf("Content length = %d, want 1", len(msg.Content))
	}
}

func TestNewMessage_Assistant(t *testing.T) {
	msg := NewAssistantMessage(TextContent("Hi there!"))

	if msg.Role != RoleAssistant {
		t.Errorf("Role = %v, want %v", msg.Role, RoleAssistant)
	}
	if len(msg.Content) != 1 {
		t.Errorf("Content length = %d, want 1", len(msg.Content))
	}
}

func TestMessage_AddTextContent(t *testing.T) {
	msg := NewUserMessage(
		TextContent("Hello"),
		TextContent("World"),
	)

	if len(msg.Content) != 2 {
		t.Errorf("Content length = %d, want 2", len(msg.Content))
	}
	if msg.Content[0].Text != "Hello" {
		t.Errorf("Content[0].Text = %q, want %q", msg.Content[0].Text, "Hello")
	}
	if msg.Content[1].Text != "World" {
		t.Errorf("Content[1].Text = %q, want %q", msg.Content[1].Text, "World")
	}
}

func TestMessage_AddImageContent(t *testing.T) {
	data := []byte{0x89, 0x50, 0x4E, 0x47} // PNG magic bytes
	content := ImageContent("image/png", data)

	if content.Type != ContentImage {
		t.Errorf("Type = %v, want %v", content.Type, ContentImage)
	}
	if content.MIMEType != "image/png" {
		t.Errorf("MIMEType = %q, want %q", content.MIMEType, "image/png")
	}
	if len(content.Data) != 4 {
		t.Errorf("Data length = %d, want 4", len(content.Data))
	}
}

func TestMessage_AddResourceContent(t *testing.T) {
	content := ResourceContent("file:///path/to/doc.md")

	if content.Type != ContentResource {
		t.Errorf("Type = %v, want %v", content.Type, ContentResource)
	}
	if content.Resource == nil {
		t.Fatal("Resource is nil")
	}
	if content.Resource.URI != "file:///path/to/doc.md" {
		t.Errorf("Resource.URI = %q, want %q", content.Resource.URI, "file:///path/to/doc.md")
	}
}

func TestBuilder_BuildPrompt(t *testing.T) {
	p := NewPromptBuilder("test-prompt").
		WithDescription("A test prompt").
		Build()

	if p.Name != "test-prompt" {
		t.Errorf("Name = %q, want %q", p.Name, "test-prompt")
	}
	if p.Description != "A test prompt" {
		t.Errorf("Description = %q, want %q", p.Description, "A test prompt")
	}
}

func TestBuilder_WithArgument(t *testing.T) {
	p := NewPromptBuilder("test").
		WithArgument("optional", "An optional arg", false).
		Build()

	if len(p.Arguments) != 1 {
		t.Fatalf("Arguments length = %d, want 1", len(p.Arguments))
	}
	if p.Arguments[0].Name != "optional" {
		t.Errorf("Arguments[0].Name = %q, want %q", p.Arguments[0].Name, "optional")
	}
	if p.Arguments[0].Required {
		t.Error("Arguments[0].Required = true, want false")
	}
}

func TestBuilder_WithRequiredArgument(t *testing.T) {
	p := NewPromptBuilder("test").
		WithRequiredArgument("name", "The user name").
		Build()

	if len(p.Arguments) != 1 {
		t.Fatalf("Arguments length = %d, want 1", len(p.Arguments))
	}
	if p.Arguments[0].Name != "name" {
		t.Errorf("Arguments[0].Name = %q, want %q", p.Arguments[0].Name, "name")
	}
	if !p.Arguments[0].Required {
		t.Error("Arguments[0].Required = false, want true")
	}
	if p.Arguments[0].Description != "The user name" {
		t.Errorf("Arguments[0].Description = %q, want %q", p.Arguments[0].Description, "The user name")
	}
}
