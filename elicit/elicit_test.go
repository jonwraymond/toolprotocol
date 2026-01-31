package elicit

import (
	"testing"
	"time"
)

func TestRequestType_String(t *testing.T) {
	tests := []struct {
		rt   RequestType
		want string
	}{
		{TypeText, "text"},
		{TypeConfirmation, "confirmation"},
		{TypeChoice, "choice"},
		{TypeForm, "form"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.rt.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRequestType_Valid(t *testing.T) {
	tests := []struct {
		rt   RequestType
		want bool
	}{
		{TypeText, true},
		{TypeConfirmation, true},
		{TypeChoice, true},
		{TypeForm, true},
		{RequestType("unknown"), false},
		{RequestType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.rt), func(t *testing.T) {
			if got := tt.rt.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequest_Fields(t *testing.T) {
	schema := map[string]any{"type": "object"}
	choices := []Choice{{ID: "a", Label: "A"}}

	req := &Request{
		ID:      "req-123",
		Type:    TypeChoice,
		Message: "Select an option:",
		Schema:  schema,
		Choices: choices,
		Default: "a",
		Timeout: 30 * time.Second,
	}

	if req.ID != "req-123" {
		t.Errorf("ID = %q, want %q", req.ID, "req-123")
	}
	if req.Type != TypeChoice {
		t.Errorf("Type = %v, want %v", req.Type, TypeChoice)
	}
	if req.Message != "Select an option:" {
		t.Errorf("Message = %q, want %q", req.Message, "Select an option:")
	}
	if req.Schema == nil {
		t.Error("Schema is nil")
	}
	if len(req.Choices) != 1 {
		t.Errorf("Choices length = %d, want 1", len(req.Choices))
	}
	if req.Default != "a" {
		t.Errorf("Default = %v, want %q", req.Default, "a")
	}
	if req.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want %v", req.Timeout, 30*time.Second)
	}
}

func TestRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *Request
		wantErr bool
	}{
		{
			name: "valid text request",
			req: &Request{
				ID:      "req-1",
				Type:    TypeText,
				Message: "Enter name:",
			},
			wantErr: false,
		},
		{
			name: "valid confirmation request",
			req: &Request{
				ID:      "req-2",
				Type:    TypeConfirmation,
				Message: "Are you sure?",
			},
			wantErr: false,
		},
		{
			name: "valid choice request",
			req: &Request{
				ID:      "req-3",
				Type:    TypeChoice,
				Message: "Select:",
				Choices: []Choice{{ID: "a", Label: "A"}},
			},
			wantErr: false,
		},
		{
			name: "valid form request",
			req: &Request{
				ID:      "req-4",
				Type:    TypeForm,
				Message: "Fill form:",
				Schema:  map[string]any{"type": "object"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequest_Validate_EmptyID(t *testing.T) {
	req := &Request{
		ID:      "",
		Type:    TypeText,
		Message: "Enter name:",
	}

	err := req.Validate()
	if err == nil {
		t.Error("Validate() should fail for empty ID")
	}
}

func TestRequest_Validate_InvalidType(t *testing.T) {
	req := &Request{
		ID:      "req-1",
		Type:    RequestType("invalid"),
		Message: "Enter name:",
	}

	err := req.Validate()
	if err == nil {
		t.Error("Validate() should fail for invalid type")
	}
}

func TestRequest_Validate_ChoiceWithoutChoices(t *testing.T) {
	req := &Request{
		ID:      "req-1",
		Type:    TypeChoice,
		Message: "Select:",
		Choices: nil,
	}

	err := req.Validate()
	if err == nil {
		t.Error("Validate() should fail for choice request without choices")
	}
}

func TestRequest_Validate_FormWithoutSchema(t *testing.T) {
	req := &Request{
		ID:      "req-1",
		Type:    TypeForm,
		Message: "Fill form:",
		Schema:  nil,
	}

	err := req.Validate()
	if err == nil {
		t.Error("Validate() should fail for form request without schema")
	}
}

func TestChoice_Fields(t *testing.T) {
	c := Choice{
		ID:          "opt-1",
		Label:       "Option One",
		Description: "The first option",
	}

	if c.ID != "opt-1" {
		t.Errorf("ID = %q, want %q", c.ID, "opt-1")
	}
	if c.Label != "Option One" {
		t.Errorf("Label = %q, want %q", c.Label, "Option One")
	}
	if c.Description != "The first option" {
		t.Errorf("Description = %q, want %q", c.Description, "The first option")
	}
}

func TestResponse_Fields(t *testing.T) {
	resp := &Response{
		RequestID: "req-123",
		Value:     "test value",
		Cancelled: false,
		TimedOut:  false,
	}

	if resp.RequestID != "req-123" {
		t.Errorf("RequestID = %q, want %q", resp.RequestID, "req-123")
	}
	if resp.Value != "test value" {
		t.Errorf("Value = %v, want %q", resp.Value, "test value")
	}
	if resp.Cancelled {
		t.Error("Cancelled = true, want false")
	}
	if resp.TimedOut {
		t.Error("TimedOut = true, want false")
	}
}

func TestResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		name string
		resp *Response
		want bool
	}{
		{
			name: "success",
			resp: &Response{Value: "value"},
			want: true,
		},
		{
			name: "cancelled",
			resp: &Response{Cancelled: true},
			want: false,
		},
		{
			name: "timed out",
			resp: &Response{TimedOut: true},
			want: false,
		},
		{
			name: "both cancelled and timed out",
			resp: &Response{Cancelled: true, TimedOut: true},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.resp.IsSuccess(); got != tt.want {
				t.Errorf("IsSuccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElicitorInterface_Defined(t *testing.T) {
	// Compile-time check that interfaces are properly defined
	var _ = Elicitor(nil)
	var _ = Handler(nil)
}
