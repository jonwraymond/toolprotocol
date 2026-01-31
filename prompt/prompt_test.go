package prompt

import (
	"testing"
)

func TestRole_String(t *testing.T) {
	tests := []struct {
		role Role
		want string
	}{
		{RoleUser, "user"},
		{RoleAssistant, "assistant"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.role.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRole_Valid(t *testing.T) {
	tests := []struct {
		role Role
		want bool
	}{
		{RoleUser, true},
		{RoleAssistant, true},
		{Role("system"), false},
		{Role(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			if got := tt.role.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentType_String(t *testing.T) {
	tests := []struct {
		ct   ContentType
		want string
	}{
		{ContentText, "text"},
		{ContentImage, "image"},
		{ContentResource, "resource"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.ct.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestContentType_Valid(t *testing.T) {
	tests := []struct {
		ct   ContentType
		want bool
	}{
		{ContentText, true},
		{ContentImage, true},
		{ContentResource, true},
		{ContentType("audio"), false},
		{ContentType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.ct), func(t *testing.T) {
			if got := tt.ct.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgument_Fields(t *testing.T) {
	arg := Argument{
		Name:        "username",
		Description: "The user's name",
		Required:    true,
	}

	if arg.Name != "username" {
		t.Errorf("Name = %q, want %q", arg.Name, "username")
	}
	if arg.Description != "The user's name" {
		t.Errorf("Description = %q, want %q", arg.Description, "The user's name")
	}
	if !arg.Required {
		t.Error("Required = false, want true")
	}
}

func TestPrompt_Fields(t *testing.T) {
	p := Prompt{
		Name:        "greeting",
		Description: "A friendly greeting",
		Arguments: []Argument{
			{Name: "name", Required: true},
		},
	}

	if p.Name != "greeting" {
		t.Errorf("Name = %q, want %q", p.Name, "greeting")
	}
	if p.Description != "A friendly greeting" {
		t.Errorf("Description = %q, want %q", p.Description, "A friendly greeting")
	}
	if len(p.Arguments) != 1 {
		t.Errorf("Arguments length = %d, want 1", len(p.Arguments))
	}
}

func TestPrompt_Validate(t *testing.T) {
	tests := []struct {
		name    string
		prompt  Prompt
		wantErr bool
	}{
		{
			name:    "valid",
			prompt:  Prompt{Name: "test"},
			wantErr: false,
		},
		{
			name:    "empty name",
			prompt:  Prompt{Name: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.prompt.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPrompt_RequiredArgs(t *testing.T) {
	p := Prompt{
		Name: "test",
		Arguments: []Argument{
			{Name: "required1", Required: true},
			{Name: "optional", Required: false},
			{Name: "required2", Required: true},
		},
	}

	required := p.RequiredArgs()
	if len(required) != 2 {
		t.Errorf("RequiredArgs() length = %d, want 2", len(required))
	}
	if required[0] != "required1" || required[1] != "required2" {
		t.Errorf("RequiredArgs() = %v, want [required1, required2]", required)
	}
}

func TestMessage_Fields(t *testing.T) {
	msg := Message{
		Role: RoleUser,
		Content: []Content{
			{Type: ContentText, Text: "Hello"},
		},
	}

	if msg.Role != RoleUser {
		t.Errorf("Role = %v, want %v", msg.Role, RoleUser)
	}
	if len(msg.Content) != 1 {
		t.Errorf("Content length = %d, want 1", len(msg.Content))
	}
}

func TestContent_Fields(t *testing.T) {
	c := Content{
		Type: ContentText,
		Text: "Hello, World!",
	}

	if c.Type != ContentText {
		t.Errorf("Type = %v, want %v", c.Type, ContentText)
	}
	if c.Text != "Hello, World!" {
		t.Errorf("Text = %q, want %q", c.Text, "Hello, World!")
	}
}

func TestContent_IsText(t *testing.T) {
	tests := []struct {
		name string
		c    Content
		want bool
	}{
		{
			name: "text content",
			c:    Content{Type: ContentText},
			want: true,
		},
		{
			name: "image content",
			c:    Content{Type: ContentImage},
			want: false,
		},
		{
			name: "resource content",
			c:    Content{Type: ContentResource},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsText(); got != tt.want {
				t.Errorf("IsText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceRef_Fields(t *testing.T) {
	ref := ResourceRef{
		URI: "file:///path/to/file.txt",
	}

	if ref.URI != "file:///path/to/file.txt" {
		t.Errorf("URI = %q, want %q", ref.URI, "file:///path/to/file.txt")
	}
}

func TestProviderInterface_Defined(t *testing.T) {
	// Compile-time check that interfaces are properly defined
	var _ Provider = (*Registry)(nil)
}
