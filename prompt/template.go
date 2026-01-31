package prompt

import (
	"strings"
)

// ExpandTemplate expands placeholders in a template string.
// Placeholders are in the format {{name}} and are replaced with
// the corresponding value from args.
func ExpandTemplate(template string, args map[string]string) (string, error) {
	result := template

	// Process each placeholder
	for {
		start := strings.Index(result, "{{")
		if start == -1 {
			break
		}

		end := strings.Index(result[start:], "}}")
		if end == -1 {
			break
		}
		end += start + 2

		// Extract placeholder name
		name := strings.TrimSpace(result[start+2 : end-2])

		// Look up value
		value, ok := args[name]
		if !ok {
			return "", &PromptError{
				PromptName: name,
				Op:         "expand",
				Err:        ErrMissingArgument,
			}
		}

		// Replace placeholder
		result = result[:start] + value + result[end:]
	}

	return result, nil
}

// ExpandTemplateWithDefaults expands placeholders, using defaults for missing args.
func ExpandTemplateWithDefaults(template string, args, defaults map[string]string) (string, error) {
	merged := make(map[string]string)
	for k, v := range defaults {
		merged[k] = v
	}
	for k, v := range args {
		merged[k] = v
	}
	return ExpandTemplate(template, merged)
}
