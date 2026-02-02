package elicit

import (
	"context"
	"testing"
	"time"
)

// BenchmarkNewTextRequest measures text request creation.
func BenchmarkNewTextRequest(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewTextRequest("Enter your name:")
	}
}

// BenchmarkNewConfirmationRequest measures confirmation request creation.
func BenchmarkNewConfirmationRequest(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewConfirmationRequest("Do you want to continue?")
	}
}

// BenchmarkNewChoiceRequest measures choice request creation.
func BenchmarkNewChoiceRequest(b *testing.B) {
	choices := []Choice{
		{ID: "a", Label: "Option A"},
		{ID: "b", Label: "Option B"},
		{ID: "c", Label: "Option C"},
	}

	b.ResetTimer()
	for b.Loop() {
		_ = NewChoiceRequest("Select an option:", choices)
	}
}

// BenchmarkNewFormRequest measures form request creation.
func BenchmarkNewFormRequest(b *testing.B) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string"},
		},
	}

	b.ResetTimer()
	for b.Loop() {
		_ = NewFormRequest("Fill out form:", schema)
	}
}

// BenchmarkNewBuilder measures builder creation.
func BenchmarkNewBuilder(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewBuilder(TypeText, "Enter value:")
	}
}

// BenchmarkBuilder_Build measures building a request with fluent API.
func BenchmarkBuilder_Build(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = NewBuilder(TypeText, "Enter value:").
			WithTimeout(10 * time.Second).
			WithDefault("default").
			Build()
	}
}

// BenchmarkBuilder_BuildChoice measures building a choice request.
func BenchmarkBuilder_BuildChoice(b *testing.B) {
	choices := []Choice{
		{ID: "a", Label: "Option A"},
		{ID: "b", Label: "Option B"},
	}

	b.ResetTimer()
	for b.Loop() {
		_ = NewBuilder(TypeChoice, "Select:").
			WithChoices(choices).
			WithTimeout(30 * time.Second).
			Build()
	}
}

// BenchmarkRequest_Validate measures request validation.
func BenchmarkRequest_Validate(b *testing.B) {
	req := NewTextRequest("Enter name:")

	b.ResetTimer()
	for b.Loop() {
		_ = req.Validate()
	}
}

// BenchmarkRequest_Validate_Choice measures choice request validation.
func BenchmarkRequest_Validate_Choice(b *testing.B) {
	req := NewChoiceRequest("Select:", []Choice{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
	})

	b.ResetTimer()
	for b.Loop() {
		_ = req.Validate()
	}
}

// BenchmarkResponse_IsSuccess measures success checking.
func BenchmarkResponse_IsSuccess(b *testing.B) {
	resp := &Response{
		RequestID: "req-1",
		Value:     "test",
	}

	b.ResetTimer()
	for b.Loop() {
		_ = resp.IsSuccess()
	}
}

// BenchmarkRequestType_Valid measures request type validation.
func BenchmarkRequestType_Valid(b *testing.B) {
	types := []RequestType{TypeText, TypeConfirmation, TypeChoice, TypeForm}

	b.ResetTimer()
	for b.Loop() {
		for _, rt := range types {
			_ = rt.Valid()
		}
	}
}

// BenchmarkNewHandler measures handler creation.
func BenchmarkNewHandler(b *testing.B) {
	callback := func(ctx context.Context, req *Request) (*Response, error) {
		return &Response{RequestID: req.ID, Value: "test"}, nil
	}

	b.ResetTimer()
	for b.Loop() {
		_ = NewHandler(
			WithDefaultTimeout(30*time.Second),
			WithCallback(callback),
		)
	}
}

// BenchmarkDefaultHandler_Handle measures request handling.
func BenchmarkDefaultHandler_Handle(b *testing.B) {
	handler := NewHandler(
		WithDefaultTimeout(30*time.Second),
		WithCallback(func(ctx context.Context, req *Request) (*Response, error) {
			return &Response{RequestID: req.ID, Value: "response"}, nil
		}),
	)

	req := NewTextRequest("Enter name:")
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = handler.Handle(ctx, req)
	}
}

// BenchmarkDefaultHandler_Handle_Concurrent measures concurrent handling.
func BenchmarkDefaultHandler_Handle_Concurrent(b *testing.B) {
	handler := NewHandler(
		WithDefaultTimeout(30*time.Second),
		WithCallback(func(ctx context.Context, req *Request) (*Response, error) {
			return &Response{RequestID: req.ID, Value: "response"}, nil
		}),
	)

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		req := NewTextRequest("Enter name:")
		for pb.Next() {
			_, _ = handler.Handle(ctx, req)
		}
	})
}

// BenchmarkGenerateID measures ID generation overhead.
func BenchmarkGenerateID(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = generateID()
	}
}
