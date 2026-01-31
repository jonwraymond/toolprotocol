package prompt

import (
	"errors"
	"testing"
)

func TestExpandTemplate(t *testing.T) {
	template := "Hello, {{name}}! Welcome to {{place}}."
	args := map[string]string{
		"name":  "Alice",
		"place": "Wonderland",
	}

	result, err := ExpandTemplate(template, args)
	if err != nil {
		t.Fatalf("ExpandTemplate() error = %v", err)
	}

	expected := "Hello, Alice! Welcome to Wonderland."
	if result != expected {
		t.Errorf("result = %q, want %q", result, expected)
	}
}

func TestExpandTemplate_MissingArg(t *testing.T) {
	template := "Hello, {{name}}!"
	args := map[string]string{} // Missing "name"

	_, err := ExpandTemplate(template, args)
	if !errors.Is(err, ErrMissingArgument) {
		t.Errorf("ExpandTemplate() error = %v, want ErrMissingArgument", err)
	}
}

func TestExpandTemplate_ExtraArgs(t *testing.T) {
	template := "Hello, {{name}}!"
	args := map[string]string{
		"name":  "Alice",
		"extra": "ignored",
	}

	result, err := ExpandTemplate(template, args)
	if err != nil {
		t.Fatalf("ExpandTemplate() error = %v", err)
	}

	expected := "Hello, Alice!"
	if result != expected {
		t.Errorf("result = %q, want %q", result, expected)
	}
}

func TestExpandTemplate_NoArgs(t *testing.T) {
	template := "Hello, World!"
	args := map[string]string{}

	result, err := ExpandTemplate(template, args)
	if err != nil {
		t.Fatalf("ExpandTemplate() error = %v", err)
	}

	if result != template {
		t.Errorf("result = %q, want %q", result, template)
	}
}

func TestExpandTemplate_NestedPlaceholders(t *testing.T) {
	// Single placeholder with spaces
	template := "Hello, {{ name }}!"
	args := map[string]string{
		"name": "Alice",
	}

	result, err := ExpandTemplate(template, args)
	if err != nil {
		t.Fatalf("ExpandTemplate() error = %v", err)
	}

	expected := "Hello, Alice!"
	if result != expected {
		t.Errorf("result = %q, want %q", result, expected)
	}
}

func TestExpandTemplate_MultipleSamePlaceholder(t *testing.T) {
	template := "{{name}} says: Hello, {{name}}!"
	args := map[string]string{
		"name": "Alice",
	}

	result, err := ExpandTemplate(template, args)
	if err != nil {
		t.Fatalf("ExpandTemplate() error = %v", err)
	}

	expected := "Alice says: Hello, Alice!"
	if result != expected {
		t.Errorf("result = %q, want %q", result, expected)
	}
}

func TestExpandTemplateWithDefaults(t *testing.T) {
	template := "Hello, {{name}}! You have {{count}} messages."
	args := map[string]string{
		"name": "Alice",
	}
	defaults := map[string]string{
		"name":  "User",
		"count": "0",
	}

	result, err := ExpandTemplateWithDefaults(template, args, defaults)
	if err != nil {
		t.Fatalf("ExpandTemplateWithDefaults() error = %v", err)
	}

	// "name" from args should override default, "count" from defaults
	expected := "Hello, Alice! You have 0 messages."
	if result != expected {
		t.Errorf("result = %q, want %q", result, expected)
	}
}
