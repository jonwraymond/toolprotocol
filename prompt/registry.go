package prompt

import (
	"context"
	"sync"
)

// Registry manages prompts.
type Registry struct {
	mu       sync.RWMutex
	prompts  map[string]Prompt
	handlers map[string]PromptHandler
}

// NewRegistry creates a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		prompts:  make(map[string]Prompt),
		handlers: make(map[string]PromptHandler),
	}
}

// Register registers a prompt with its handler.
func (r *Registry) Register(prompt Prompt, handler PromptHandler) error {
	if err := prompt.Validate(); err != nil {
		return &PromptError{
			PromptName: prompt.Name,
			Op:         "register",
			Err:        err,
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prompts[prompt.Name]; exists {
		return &PromptError{
			PromptName: prompt.Name,
			Op:         "register",
			Err:        ErrDuplicatePrompt,
		}
	}

	r.prompts[prompt.Name] = prompt
	r.handlers[prompt.Name] = handler
	return nil
}

// List returns all registered prompts.
func (r *Registry) List(ctx context.Context) ([]Prompt, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prompts := make([]Prompt, 0, len(r.prompts))
	for _, p := range r.prompts {
		prompts = append(prompts, p)
	}
	return prompts, nil
}

// Get returns messages for a prompt with the given arguments.
func (r *Registry) Get(ctx context.Context, name string, args map[string]string) ([]Message, error) {
	r.mu.RLock()
	prompt, exists := r.prompts[name]
	handler := r.handlers[name]
	r.mu.RUnlock()

	if !exists {
		return nil, &PromptError{
			PromptName: name,
			Op:         "get",
			Err:        ErrPromptNotFound,
		}
	}

	// Check required arguments
	for _, arg := range prompt.Arguments {
		if arg.Required {
			if _, ok := args[arg.Name]; !ok {
				return nil, &PromptError{
					PromptName: name,
					Op:         "get",
					Err:        ErrMissingArgument,
				}
			}
		}
	}

	// Call handler
	messages, err := handler(ctx, args)
	if err != nil {
		return nil, &PromptError{
			PromptName: name,
			Op:         "get",
			Err:        err,
		}
	}

	return messages, nil
}
