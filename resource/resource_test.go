package resource

import (
	"testing"
)

func TestResource_Fields(t *testing.T) {
	r := Resource{
		URI:         "file:///path/to/doc.md",
		Name:        "Documentation",
		Description: "Project docs",
		MIMEType:    "text/markdown",
		Annotations: map[string]any{"version": "1.0"},
	}

	if r.URI != "file:///path/to/doc.md" {
		t.Errorf("URI = %q", r.URI)
	}
	if r.Name != "Documentation" {
		t.Errorf("Name = %q", r.Name)
	}
	if r.Description != "Project docs" {
		t.Errorf("Description = %q", r.Description)
	}
	if r.MIMEType != "text/markdown" {
		t.Errorf("MIMEType = %q", r.MIMEType)
	}
	if r.Annotations["version"] != "1.0" {
		t.Errorf("Annotations[version] = %v", r.Annotations["version"])
	}
}

func TestResource_Clone(t *testing.T) {
	original := &Resource{
		URI:         "file:///test.txt",
		Name:        "Test",
		Description: "Test resource",
		MIMEType:    "text/plain",
		Annotations: map[string]any{"key": "value"},
	}

	clone := original.Clone()

	if clone.URI != original.URI {
		t.Errorf("clone.URI = %q, want %q", clone.URI, original.URI)
	}
	if clone.Name != original.Name {
		t.Errorf("clone.Name = %q, want %q", clone.Name, original.Name)
	}

	// Verify independence
	clone.Annotations["key"] = "modified"
	if original.Annotations["key"] == "modified" {
		t.Error("modifying clone affected original")
	}
}

func TestResource_Clone_NilAnnotations(t *testing.T) {
	original := &Resource{
		URI: "file:///test.txt",
	}

	clone := original.Clone()
	if clone.Annotations != nil {
		t.Error("clone.Annotations should be nil")
	}
}

func TestContents_Fields(t *testing.T) {
	c := Contents{
		URI:      "file:///doc.md",
		MIMEType: "text/markdown",
		Text:     "# Hello",
	}

	if c.URI != "file:///doc.md" {
		t.Errorf("URI = %q", c.URI)
	}
	if c.MIMEType != "text/markdown" {
		t.Errorf("MIMEType = %q", c.MIMEType)
	}
	if c.Text != "# Hello" {
		t.Errorf("Text = %q", c.Text)
	}
}

func TestContents_IsText(t *testing.T) {
	tests := []struct {
		name string
		c    Contents
		want bool
	}{
		{
			name: "text content",
			c:    Contents{Text: "hello"},
			want: true,
		},
		{
			name: "empty content",
			c:    Contents{},
			want: true,
		},
		{
			name: "binary content",
			c:    Contents{Blob: []byte{1, 2, 3}},
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

func TestContents_IsBinary(t *testing.T) {
	tests := []struct {
		name string
		c    Contents
		want bool
	}{
		{
			name: "binary content",
			c:    Contents{Blob: []byte{1, 2, 3}},
			want: true,
		},
		{
			name: "text content",
			c:    Contents{Text: "hello"},
			want: false,
		},
		{
			name: "empty content",
			c:    Contents{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsBinary(); got != tt.want {
				t.Errorf("IsBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplate_Fields(t *testing.T) {
	tmpl := Template{
		URITemplate: "file:///users/{id}/profile",
		Name:        "User Profile",
		Description: "User profile resource",
		MIMEType:    "application/json",
	}

	if tmpl.URITemplate != "file:///users/{id}/profile" {
		t.Errorf("URITemplate = %q", tmpl.URITemplate)
	}
	if tmpl.Name != "User Profile" {
		t.Errorf("Name = %q", tmpl.Name)
	}
}

func TestTemplate_Expand(t *testing.T) {
	tmpl := Template{
		URITemplate: "file:///users/{id}/docs/{name}",
	}

	result := tmpl.Expand(map[string]string{
		"id":   "123",
		"name": "readme.md",
	})

	expected := "file:///users/123/docs/readme.md"
	if result != expected {
		t.Errorf("Expand() = %q, want %q", result, expected)
	}
}

func TestProviderInterface_Defined(t *testing.T) {
	var _ Provider = (*StaticProvider)(nil)
}
