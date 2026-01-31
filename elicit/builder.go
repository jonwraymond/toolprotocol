package elicit

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// generateID generates a random request ID.
func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// NewTextRequest creates a new text input request.
func NewTextRequest(message string) *Request {
	return &Request{
		ID:      generateID(),
		Type:    TypeText,
		Message: message,
	}
}

// NewConfirmationRequest creates a new confirmation request.
func NewConfirmationRequest(message string) *Request {
	return &Request{
		ID:      generateID(),
		Type:    TypeConfirmation,
		Message: message,
	}
}

// NewChoiceRequest creates a new choice selection request.
func NewChoiceRequest(message string, choices []Choice) *Request {
	return &Request{
		ID:      generateID(),
		Type:    TypeChoice,
		Message: message,
		Choices: choices,
	}
}

// NewFormRequest creates a new form input request.
func NewFormRequest(message string, schema any) *Request {
	return &Request{
		ID:      generateID(),
		Type:    TypeForm,
		Message: message,
		Schema:  schema,
	}
}

// Builder provides a fluent interface for building requests.
type Builder struct {
	req *Request
}

// NewBuilder creates a new request builder.
func NewBuilder(reqType RequestType, message string) *Builder {
	return &Builder{
		req: &Request{
			ID:      generateID(),
			Type:    reqType,
			Message: message,
		},
	}
}

// WithTimeout sets the request timeout.
func (b *Builder) WithTimeout(timeout time.Duration) *Builder {
	b.req.Timeout = timeout
	return b
}

// WithDefault sets the default value.
func (b *Builder) WithDefault(defaultValue any) *Builder {
	b.req.Default = defaultValue
	return b
}

// WithSchema sets the JSON Schema for form requests.
func (b *Builder) WithSchema(schema any) *Builder {
	b.req.Schema = schema
	return b
}

// WithChoices sets the choices for choice requests.
func (b *Builder) WithChoices(choices []Choice) *Builder {
	b.req.Choices = choices
	return b
}

// Build returns the built request.
func (b *Builder) Build() *Request {
	return b.req
}
