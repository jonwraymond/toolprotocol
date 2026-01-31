package resource

import (
	"context"
	"strings"
)

// Resource represents an MCP resource.
type Resource struct {
	// URI is the unique resource identifier.
	URI string

	// Name is the display name.
	Name string

	// Description describes the resource.
	Description string

	// MIMEType is the content MIME type.
	MIMEType string

	// Annotations holds additional metadata.
	Annotations map[string]any
}

// Clone returns a deep copy of the resource.
func (r *Resource) Clone() *Resource {
	clone := &Resource{
		URI:         r.URI,
		Name:        r.Name,
		Description: r.Description,
		MIMEType:    r.MIMEType,
	}
	if r.Annotations != nil {
		clone.Annotations = make(map[string]any, len(r.Annotations))
		for k, v := range r.Annotations {
			clone.Annotations[k] = v
		}
	}
	return clone
}

// Contents represents resource contents.
type Contents struct {
	// URI is the resource URI.
	URI string

	// MIMEType is the content MIME type.
	MIMEType string

	// Text is the text content (for text resources).
	Text string

	// Blob is the binary content (for binary resources).
	Blob []byte
}

// IsText returns true if this is text content.
func (c *Contents) IsText() bool {
	return c.Text != "" || len(c.Blob) == 0
}

// IsBinary returns true if this is binary content.
func (c *Contents) IsBinary() bool {
	return len(c.Blob) > 0
}

// Template represents a resource template.
type Template struct {
	// URITemplate is the URI template pattern.
	URITemplate string

	// Name is the display name.
	Name string

	// Description describes the template.
	Description string

	// MIMEType is the expected content MIME type.
	MIMEType string
}

// Expand expands the template with the given values.
// Template variables use {name} format.
func (t *Template) Expand(values map[string]string) string {
	result := t.URITemplate
	for k, v := range values {
		result = strings.ReplaceAll(result, "{"+k+"}", v)
	}
	return result
}

// Provider serves resources.
type Provider interface {
	// List returns all available resources.
	List(ctx context.Context) ([]Resource, error)

	// Read returns the contents of a resource.
	Read(ctx context.Context, uri string) (*Contents, error)

	// Templates returns available resource templates.
	Templates(ctx context.Context) ([]Template, error)
}

// Subscriber receives resource updates.
type Subscriber interface {
	// Subscribe subscribes to updates for a resource.
	Subscribe(ctx context.Context, uri string) (<-chan *Contents, error)

	// Unsubscribe unsubscribes from a resource.
	Unsubscribe(ctx context.Context, uri string) error
}
