package elicit

import (
	"testing"
	"time"
)

func TestNewTextRequest(t *testing.T) {
	req := NewTextRequest("Enter your name:")

	if req.ID == "" {
		t.Error("ID should be generated")
	}
	if req.Type != TypeText {
		t.Errorf("Type = %v, want %v", req.Type, TypeText)
	}
	if req.Message != "Enter your name:" {
		t.Errorf("Message = %q, want %q", req.Message, "Enter your name:")
	}
}

func TestNewConfirmationRequest(t *testing.T) {
	req := NewConfirmationRequest("Are you sure?")

	if req.ID == "" {
		t.Error("ID should be generated")
	}
	if req.Type != TypeConfirmation {
		t.Errorf("Type = %v, want %v", req.Type, TypeConfirmation)
	}
	if req.Message != "Are you sure?" {
		t.Errorf("Message = %q, want %q", req.Message, "Are you sure?")
	}
}

func TestNewChoiceRequest(t *testing.T) {
	choices := []Choice{
		{ID: "a", Label: "Option A"},
		{ID: "b", Label: "Option B"},
	}
	req := NewChoiceRequest("Select an option:", choices)

	if req.ID == "" {
		t.Error("ID should be generated")
	}
	if req.Type != TypeChoice {
		t.Errorf("Type = %v, want %v", req.Type, TypeChoice)
	}
	if req.Message != "Select an option:" {
		t.Errorf("Message = %q, want %q", req.Message, "Select an option:")
	}
	if len(req.Choices) != 2 {
		t.Errorf("Choices length = %d, want 2", len(req.Choices))
	}
}

func TestNewChoiceRequest_WithDefault(t *testing.T) {
	choices := []Choice{
		{ID: "a", Label: "Option A"},
		{ID: "b", Label: "Option B"},
	}
	req := NewChoiceRequest("Select an option:", choices)
	req.Default = "b"

	if req.Default != "b" {
		t.Errorf("Default = %v, want %q", req.Default, "b")
	}
}

func TestNewFormRequest(t *testing.T) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string"},
		},
	}
	req := NewFormRequest("Fill out the form:", schema)

	if req.ID == "" {
		t.Error("ID should be generated")
	}
	if req.Type != TypeForm {
		t.Errorf("Type = %v, want %v", req.Type, TypeForm)
	}
	if req.Message != "Fill out the form:" {
		t.Errorf("Message = %q, want %q", req.Message, "Fill out the form:")
	}
	if req.Schema == nil {
		t.Error("Schema should be set")
	}
}

func TestBuilder_WithTimeout(t *testing.T) {
	req := NewBuilder(TypeText, "Enter name:").
		WithTimeout(30 * time.Second).
		Build()

	if req.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want %v", req.Timeout, 30*time.Second)
	}
}

func TestBuilder_WithDefault(t *testing.T) {
	req := NewBuilder(TypeText, "Enter name:").
		WithDefault("John").
		Build()

	if req.Default != "John" {
		t.Errorf("Default = %v, want %q", req.Default, "John")
	}
}

func TestBuilder_WithSchema(t *testing.T) {
	schema := map[string]any{"type": "object"}
	req := NewBuilder(TypeForm, "Fill form:").
		WithSchema(schema).
		Build()

	if req.Schema == nil {
		t.Error("Schema should be set")
	}
}

func TestBuilder_Chaining(t *testing.T) {
	choices := []Choice{{ID: "a", Label: "A"}}
	req := NewBuilder(TypeChoice, "Select:").
		WithChoices(choices).
		WithDefault("a").
		WithTimeout(time.Minute).
		Build()

	if req.Type != TypeChoice {
		t.Errorf("Type = %v, want %v", req.Type, TypeChoice)
	}
	if len(req.Choices) != 1 {
		t.Errorf("Choices length = %d, want 1", len(req.Choices))
	}
	if req.Default != "a" {
		t.Errorf("Default = %v, want %q", req.Default, "a")
	}
	if req.Timeout != time.Minute {
		t.Errorf("Timeout = %v, want %v", req.Timeout, time.Minute)
	}
}
