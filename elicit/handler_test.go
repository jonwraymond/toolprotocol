package elicit

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestDefaultHandler_Handle_Text(t *testing.T) {
	handler := NewHandler(WithCallback(func(ctx context.Context, req *Request) (*Response, error) {
		return &Response{
			RequestID: req.ID,
			Value:     "test response",
		}, nil
	}))

	req := NewTextRequest("Enter name:")
	resp, err := handler.Handle(context.Background(), req)

	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if resp.Value != "test response" {
		t.Errorf("Value = %v, want %q", resp.Value, "test response")
	}
}

func TestDefaultHandler_Handle_Confirmation(t *testing.T) {
	handler := NewHandler(WithCallback(func(ctx context.Context, req *Request) (*Response, error) {
		return &Response{
			RequestID: req.ID,
			Value:     true,
		}, nil
	}))

	req := NewConfirmationRequest("Are you sure?")
	resp, err := handler.Handle(context.Background(), req)

	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if resp.Value != true {
		t.Errorf("Value = %v, want true", resp.Value)
	}
}

func TestDefaultHandler_Handle_Choice(t *testing.T) {
	handler := NewHandler(WithCallback(func(ctx context.Context, req *Request) (*Response, error) {
		return &Response{
			RequestID: req.ID,
			Value:     "option-a",
		}, nil
	}))

	choices := []Choice{{ID: "option-a", Label: "A"}}
	req := NewChoiceRequest("Select:", choices)
	resp, err := handler.Handle(context.Background(), req)

	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if resp.Value != "option-a" {
		t.Errorf("Value = %v, want %q", resp.Value, "option-a")
	}
}

func TestDefaultHandler_Handle_Form(t *testing.T) {
	formData := map[string]any{"name": "Alice", "age": 30}
	handler := NewHandler(WithCallback(func(ctx context.Context, req *Request) (*Response, error) {
		return &Response{
			RequestID: req.ID,
			Value:     formData,
		}, nil
	}))

	schema := map[string]any{"type": "object"}
	req := NewFormRequest("Fill form:", schema)
	resp, err := handler.Handle(context.Background(), req)

	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if resp.Value == nil {
		t.Error("Value should not be nil")
	}
}

func TestDefaultHandler_Handle_Timeout(t *testing.T) {
	handler := NewHandler(
		WithDefaultTimeout(10*time.Millisecond),
		WithCallback(func(ctx context.Context, req *Request) (*Response, error) {
			time.Sleep(100 * time.Millisecond)
			return &Response{RequestID: req.ID, Value: "late"}, nil
		}),
	)

	req := NewTextRequest("Enter:")
	resp, err := handler.Handle(context.Background(), req)

	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if !resp.TimedOut {
		t.Error("TimedOut = false, want true")
	}
}

func TestDefaultHandler_Handle_Cancelled(t *testing.T) {
	handler := NewHandler(WithCallback(func(ctx context.Context, req *Request) (*Response, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	req := NewTextRequest("Enter:")
	resp, err := handler.Handle(ctx, req)

	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if !resp.Cancelled {
		t.Error("Cancelled = false, want true")
	}
}

func TestDefaultHandler_Handle_InvalidRequest(t *testing.T) {
	handler := NewHandler(WithCallback(func(ctx context.Context, req *Request) (*Response, error) {
		return &Response{RequestID: req.ID}, nil
	}))

	// Request with empty ID is invalid
	req := &Request{Type: TypeText, Message: "test"}
	_, err := handler.Handle(context.Background(), req)

	if err == nil {
		t.Error("Handle() should fail for invalid request")
	}
	var elicitErr *ElicitError
	if !errors.As(err, &elicitErr) {
		t.Errorf("error should be *ElicitError, got %T", err)
	}
}

func TestDefaultHandler_Handle_NoHandler(t *testing.T) {
	handler := NewHandler() // No callback set

	req := NewTextRequest("Enter:")
	_, err := handler.Handle(context.Background(), req)

	if err == nil {
		t.Error("Handle() should fail when no handler set")
	}
	if !errors.Is(err, ErrNoHandler) {
		t.Errorf("error should wrap ErrNoHandler, got %v", err)
	}
}

func TestDefaultHandler_Handle_HandlerError(t *testing.T) {
	testErr := errors.New("handler failed")
	handler := NewHandler(WithCallback(func(ctx context.Context, req *Request) (*Response, error) {
		return nil, testErr
	}))

	req := NewTextRequest("Enter:")
	_, err := handler.Handle(context.Background(), req)

	if err == nil {
		t.Error("Handle() should propagate handler error")
	}
	if !errors.Is(err, testErr) {
		t.Errorf("error should wrap testErr, got %v", err)
	}
}

func TestDefaultHandler_ConcurrentRequests(t *testing.T) {
	var count int
	var mu sync.Mutex

	handler := NewHandler(WithCallback(func(ctx context.Context, req *Request) (*Response, error) {
		mu.Lock()
		count++
		mu.Unlock()
		return &Response{RequestID: req.ID, Value: "ok"}, nil
	}))

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := NewTextRequest("Enter:")
			_, _ = handler.Handle(context.Background(), req)
		}()
	}
	wg.Wait()

	mu.Lock()
	if count != 10 {
		t.Errorf("count = %d, want 10", count)
	}
	mu.Unlock()
}

func TestDefaultHandler_RequestTimeout(t *testing.T) {
	handler := NewHandler(
		WithDefaultTimeout(time.Hour), // Long default
		WithCallback(func(ctx context.Context, req *Request) (*Response, error) {
			time.Sleep(100 * time.Millisecond)
			return &Response{RequestID: req.ID, Value: "late"}, nil
		}),
	)

	// Request with short timeout should override default
	req := NewBuilder(TypeText, "Enter:").
		WithTimeout(10 * time.Millisecond).
		Build()

	resp, err := handler.Handle(context.Background(), req)

	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if !resp.TimedOut {
		t.Error("TimedOut = false, want true (request timeout should override)")
	}
}
