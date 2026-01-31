package prompt

// Option configures prompt package types.
type Option func(any)

// ValidatorFunc is a function that validates prompt arguments.
type ValidatorFunc func(args map[string]string) error

// WithValidator returns an option that sets a validator for the registry.
// Note: This is a placeholder for future enhancement.
func WithValidator(fn ValidatorFunc) Option {
	return func(v any) {
		// Future: Apply validator to registry
	}
}
