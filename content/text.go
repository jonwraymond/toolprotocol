package content

// TextContent represents plain text content.
type TextContent struct {
	// Text is the text content.
	Text string

	// Mime is the MIME type (default: text/plain).
	Mime string
}

// NewText creates a new text content.
func NewText(text string) *TextContent {
	return &TextContent{
		Text: text,
		Mime: "text/plain",
	}
}

// Type returns TypeText.
func (c *TextContent) Type() Type {
	return TypeText
}

// MIMEType returns the MIME type (default: text/plain).
func (c *TextContent) MIMEType() string {
	if c.Mime == "" {
		return "text/plain"
	}
	return c.Mime
}

// Bytes returns the text as bytes.
func (c *TextContent) Bytes() ([]byte, error) {
	return []byte(c.Text), nil
}

// String returns the text content.
func (c *TextContent) String() string {
	return c.Text
}
