package wire

import "fmt"

// ContentType identifies the type of content in a response.
type ContentType string

const (
	// ContentTypeText is plain text content.
	ContentTypeText ContentType = "text"

	// ContentTypeImage is image data.
	ContentTypeImage ContentType = "image"

	// ContentTypeResource is a resource reference.
	ContentTypeResource ContentType = "resource"
)

// Content represents a piece of response content.
type Content struct {
	// Type identifies the content type.
	Type ContentType

	// Text is the text content (for ContentTypeText).
	Text string

	// MIMEType is the MIME type (for images and resources).
	MIMEType string

	// Data is binary data (for images).
	Data []byte

	// URI is the resource URI (for ContentTypeResource).
	URI string
}

// Tool describes a tool's interface.
type Tool struct {
	// Name is the tool identifier.
	Name string

	// Description explains what the tool does.
	Description string

	// InputSchema is the JSON Schema for tool arguments.
	InputSchema map[string]any
}

// Error represents a wire protocol error.
type Error struct {
	// Code is the error code (JSON-RPC style).
	Code int

	// Message is a human-readable error message.
	Message string

	// Data contains additional error details.
	Data any
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Data != nil {
		return fmt.Sprintf("%s (code: %d, data: %v)", e.Message, e.Code, e.Data)
	}
	return fmt.Sprintf("%s (code: %d)", e.Message, e.Code)
}
