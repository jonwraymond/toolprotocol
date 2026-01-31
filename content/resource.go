package content

// ResourceContent represents an external resource reference.
type ResourceContent struct {
	// URI is the resource URI.
	URI string

	// Mime is the MIME type.
	Mime string

	// Text is a text representation of the resource.
	Text string

	// Blob is binary data for the resource.
	Blob []byte
}

// NewResource creates a new resource content.
func NewResource(uri string) *ResourceContent {
	return &ResourceContent{
		URI: uri,
	}
}

// Type returns TypeResource.
func (c *ResourceContent) Type() Type {
	return TypeResource
}

// MIMEType returns the resource MIME type.
func (c *ResourceContent) MIMEType() string {
	if c.Mime == "" {
		return "application/octet-stream"
	}
	return c.Mime
}

// Bytes returns the resource content as bytes.
// Returns Blob if set, otherwise Text as bytes.
func (c *ResourceContent) Bytes() ([]byte, error) {
	if len(c.Blob) > 0 {
		return c.Blob, nil
	}
	return []byte(c.Text), nil
}

// String returns a string representation of the resource.
func (c *ResourceContent) String() string {
	if c.Text != "" {
		return c.Text
	}
	return c.URI
}
