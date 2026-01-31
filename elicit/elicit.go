package elicit

import (
	"context"
	"time"
)

// RequestType represents the type of elicitation request.
type RequestType string

const (
	// TypeText requests free-form text input.
	TypeText RequestType = "text"

	// TypeConfirmation requests a yes/no confirmation.
	TypeConfirmation RequestType = "confirmation"

	// TypeChoice requests selection from predefined options.
	TypeChoice RequestType = "choice"

	// TypeForm requests structured input with JSON Schema validation.
	TypeForm RequestType = "form"
)

// String returns the string representation of the request type.
func (rt RequestType) String() string {
	return string(rt)
}

// Valid returns true if the request type is a known valid type.
func (rt RequestType) Valid() bool {
	switch rt {
	case TypeText, TypeConfirmation, TypeChoice, TypeForm:
		return true
	default:
		return false
	}
}

// Request represents an elicitation request.
type Request struct {
	// ID is the unique identifier for the request.
	ID string

	// Type is the type of elicitation request.
	Type RequestType

	// Message is the prompt message to display.
	Message string

	// Schema is the JSON Schema for form type requests.
	Schema any

	// Choices are the available options for choice type requests.
	Choices []Choice

	// Default is the default value for the request.
	Default any

	// Timeout is the maximum time to wait for a response.
	Timeout time.Duration
}

// Validate validates the request.
func (r *Request) Validate() error {
	if r.ID == "" {
		return ErrInvalidRequest
	}
	if !r.Type.Valid() {
		return ErrInvalidRequest
	}
	if r.Type == TypeChoice && len(r.Choices) == 0 {
		return ErrInvalidRequest
	}
	if r.Type == TypeForm && r.Schema == nil {
		return ErrInvalidRequest
	}
	return nil
}

// Choice represents a selection option for choice type requests.
type Choice struct {
	// ID is the unique identifier for the choice.
	ID string

	// Label is the display text for the choice.
	Label string

	// Description is additional context for the choice.
	Description string
}

// Response represents a user response to an elicitation request.
type Response struct {
	// RequestID is the ID of the request this responds to.
	RequestID string

	// Value is the user's response value.
	Value any

	// Cancelled indicates the user cancelled the request.
	Cancelled bool

	// TimedOut indicates the request timed out.
	TimedOut bool
}

// IsSuccess returns true if the response is successful (not cancelled or timed out).
func (r *Response) IsSuccess() bool {
	return !r.Cancelled && !r.TimedOut
}

// Elicitor sends elicitation requests to clients.
type Elicitor interface {
	// Elicit sends an elicitation request and waits for a response.
	Elicit(ctx context.Context, req *Request) (*Response, error)
}

// Handler processes elicitation requests (client-side).
type Handler interface {
	// Handle processes an elicitation request and returns a response.
	Handle(ctx context.Context, req *Request) (*Response, error)
}
