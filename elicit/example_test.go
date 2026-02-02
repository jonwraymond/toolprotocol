package elicit_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonwraymond/toolprotocol/elicit"
)

func ExampleNewTextRequest() {
	req := elicit.NewTextRequest("Enter your name:")

	fmt.Println("Type:", req.Type)
	fmt.Println("Message:", req.Message)
	fmt.Println("Has ID:", req.ID != "")
	// Output:
	// Type: text
	// Message: Enter your name:
	// Has ID: true
}

func ExampleNewConfirmationRequest() {
	req := elicit.NewConfirmationRequest("Do you want to continue?")

	fmt.Println("Type:", req.Type)
	fmt.Println("Message:", req.Message)
	// Output:
	// Type: confirmation
	// Message: Do you want to continue?
}

func ExampleNewChoiceRequest() {
	choices := []elicit.Choice{
		{ID: "a", Label: "Option A", Description: "First option"},
		{ID: "b", Label: "Option B", Description: "Second option"},
	}
	req := elicit.NewChoiceRequest("Select an option:", choices)

	fmt.Println("Type:", req.Type)
	fmt.Println("Message:", req.Message)
	fmt.Println("Choices count:", len(req.Choices))
	fmt.Println("First choice:", req.Choices[0].Label)
	// Output:
	// Type: choice
	// Message: Select an option:
	// Choices count: 2
	// First choice: Option A
}

func ExampleNewFormRequest() {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string"},
			"age":  map[string]any{"type": "integer"},
		},
	}
	req := elicit.NewFormRequest("Fill out the form:", schema)

	fmt.Println("Type:", req.Type)
	fmt.Println("Message:", req.Message)
	fmt.Println("Has schema:", req.Schema != nil)
	// Output:
	// Type: form
	// Message: Fill out the form:
	// Has schema: true
}

func ExampleNewBuilder() {
	builder := elicit.NewBuilder(elicit.TypeText, "Enter value:")

	fmt.Printf("Type: %T\n", builder)
	// Output:
	// Type: *elicit.Builder
}

func ExampleBuilder_Build() {
	req := elicit.NewBuilder(elicit.TypeText, "Enter value:").
		WithTimeout(10 * time.Second).
		WithDefault("default value").
		Build()

	fmt.Println("Type:", req.Type)
	fmt.Println("Message:", req.Message)
	fmt.Println("Timeout:", req.Timeout)
	fmt.Println("Default:", req.Default)
	// Output:
	// Type: text
	// Message: Enter value:
	// Timeout: 10s
	// Default: default value
}

func ExampleBuilder_WithChoices() {
	choices := []elicit.Choice{
		{ID: "yes", Label: "Yes"},
		{ID: "no", Label: "No"},
	}
	req := elicit.NewBuilder(elicit.TypeChoice, "Choose:").
		WithChoices(choices).
		Build()

	fmt.Println("Type:", req.Type)
	fmt.Println("Choices:", len(req.Choices))
	// Output:
	// Type: choice
	// Choices: 2
}

func ExampleBuilder_WithSchema() {
	schema := map[string]any{"type": "string"}
	req := elicit.NewBuilder(elicit.TypeForm, "Enter data:").
		WithSchema(schema).
		Build()

	fmt.Println("Type:", req.Type)
	fmt.Println("Has schema:", req.Schema != nil)
	// Output:
	// Type: form
	// Has schema: true
}

func ExampleRequest_Validate() {
	// Valid text request
	validReq := elicit.NewTextRequest("Enter name:")
	fmt.Println("Valid text error:", validReq.Validate())

	// Invalid request (no ID)
	invalidReq := &elicit.Request{Type: elicit.TypeText}
	fmt.Println("Invalid has error:", invalidReq.Validate() != nil)
	// Output:
	// Valid text error: <nil>
	// Invalid has error: true
}

func ExampleRequest_Validate_choice() {
	// Choice request without choices is invalid
	req := &elicit.Request{
		ID:   "test-id",
		Type: elicit.TypeChoice,
	}
	err := req.Validate()
	fmt.Println("Choice without options error:", errors.Is(err, elicit.ErrInvalidRequest))
	// Output:
	// Choice without options error: true
}

func ExampleResponse_IsSuccess() {
	// Successful response
	success := &elicit.Response{
		RequestID: "req-1",
		Value:     "user input",
	}
	fmt.Println("Success:", success.IsSuccess())

	// Cancelled response
	cancelled := &elicit.Response{
		RequestID: "req-2",
		Cancelled: true,
	}
	fmt.Println("Cancelled is success:", cancelled.IsSuccess())

	// Timed out response
	timedOut := &elicit.Response{
		RequestID: "req-3",
		TimedOut:  true,
	}
	fmt.Println("TimedOut is success:", timedOut.IsSuccess())
	// Output:
	// Success: true
	// Cancelled is success: false
	// TimedOut is success: false
}

func ExampleNewHandler() {
	handler := elicit.NewHandler(
		elicit.WithDefaultTimeout(5*time.Second),
		elicit.WithCallback(func(ctx context.Context, req *elicit.Request) (*elicit.Response, error) {
			return &elicit.Response{
				RequestID: req.ID,
				Value:     "test response",
			}, nil
		}),
	)

	fmt.Printf("Type: %T\n", handler)
	// Output:
	// Type: *elicit.DefaultHandler
}

func ExampleDefaultHandler_Handle() {
	handler := elicit.NewHandler(
		elicit.WithCallback(func(ctx context.Context, req *elicit.Request) (*elicit.Response, error) {
			return &elicit.Response{
				RequestID: req.ID,
				Value:     "Hello, " + req.Message,
			}, nil
		}),
	)

	req := elicit.NewTextRequest("World")
	ctx := context.Background()
	resp, err := handler.Handle(ctx, req)

	fmt.Println("Error:", err)
	fmt.Println("Value:", resp.Value)
	// Output:
	// Error: <nil>
	// Value: Hello, World
}

func ExampleDefaultHandler_Handle_noHandler() {
	handler := elicit.NewHandler() // No callback set

	req := elicit.NewTextRequest("test")
	ctx := context.Background()
	_, err := handler.Handle(ctx, req)

	fmt.Println("Has error:", err != nil)
	fmt.Println("Is no handler:", errors.Is(err, elicit.ErrNoHandler))
	// Output:
	// Has error: true
	// Is no handler: true
}

func ExampleRequestType_String() {
	fmt.Println(elicit.TypeText.String())
	fmt.Println(elicit.TypeConfirmation.String())
	fmt.Println(elicit.TypeChoice.String())
	fmt.Println(elicit.TypeForm.String())
	// Output:
	// text
	// confirmation
	// choice
	// form
}

func ExampleRequestType_Valid() {
	fmt.Println("text valid:", elicit.TypeText.Valid())
	fmt.Println("confirmation valid:", elicit.TypeConfirmation.Valid())
	fmt.Println("choice valid:", elicit.TypeChoice.Valid())
	fmt.Println("form valid:", elicit.TypeForm.Valid())
	fmt.Println("unknown valid:", elicit.RequestType("unknown").Valid())
	// Output:
	// text valid: true
	// confirmation valid: true
	// choice valid: true
	// form valid: true
	// unknown valid: false
}

func ExampleElicitError() {
	err := &elicit.ElicitError{
		RequestID: "req-123",
		Op:        "validate",
		Err:       elicit.ErrInvalidRequest,
	}

	fmt.Println(err.Error())
	fmt.Println("Unwraps:", errors.Is(err, elicit.ErrInvalidRequest))
	// Output:
	// elicit req-123: validate: elicit: invalid request
	// Unwraps: true
}

func Example_elicitationWorkflow() {
	// Create a handler that simulates user input
	handler := elicit.NewHandler(
		elicit.WithDefaultTimeout(5*time.Second),
		elicit.WithCallback(func(ctx context.Context, req *elicit.Request) (*elicit.Response, error) {
			// Simulate processing based on request type
			switch req.Type {
			case elicit.TypeText:
				return &elicit.Response{
					RequestID: req.ID,
					Value:     "Alice",
				}, nil
			case elicit.TypeConfirmation:
				return &elicit.Response{
					RequestID: req.ID,
					Value:     true,
				}, nil
			case elicit.TypeChoice:
				return &elicit.Response{
					RequestID: req.ID,
					Value:     req.Choices[0].ID,
				}, nil
			default:
				return nil, elicit.ErrInvalidRequest
			}
		}),
	)

	ctx := context.Background()

	// Text input
	textReq := elicit.NewTextRequest("Enter name:")
	textResp, _ := handler.Handle(ctx, textReq)
	fmt.Println("Name:", textResp.Value)

	// Confirmation
	confirmReq := elicit.NewConfirmationRequest("Continue?")
	confirmResp, _ := handler.Handle(ctx, confirmReq)
	fmt.Println("Confirmed:", confirmResp.Value)

	// Choice
	choiceReq := elicit.NewChoiceRequest("Select:", []elicit.Choice{
		{ID: "opt1", Label: "First"},
		{ID: "opt2", Label: "Second"},
	})
	choiceResp, _ := handler.Handle(ctx, choiceReq)
	fmt.Println("Selected:", choiceResp.Value)
	// Output:
	// Name: Alice
	// Confirmed: true
	// Selected: opt1
}
