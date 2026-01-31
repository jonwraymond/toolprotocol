package discover

// Filter specifies criteria for listing services.
type Filter struct {
	// Namespace filters services by namespace.
	Namespace string

	// Tags filters services by tags (must have all).
	Tags []string

	// Capabilities filters services by required capabilities.
	Capabilities []string

	// Limit is the maximum number of results to return.
	Limit int

	// Cursor is the pagination cursor from a previous request.
	Cursor string
}

// Matches returns true if the discoverable matches the filter criteria.
func (f *Filter) Matches(d Discoverable) bool {
	if f == nil {
		return true
	}

	// Check namespace (if specified)
	// Note: Discoverable doesn't have Namespace(), so we skip this check
	// unless we add it to the interface later

	// Check capabilities (if specified)
	if len(f.Capabilities) > 0 {
		caps := d.Capabilities()
		if caps == nil {
			return false
		}
		for _, required := range f.Capabilities {
			if !hasCapability(caps, required) {
				return false
			}
		}
	}

	return true
}

// hasCapability checks if capabilities include the named capability.
func hasCapability(caps *Capabilities, name string) bool {
	switch name {
	case "tools":
		return caps.Tools
	case "resources":
		return caps.Resources
	case "prompts":
		return caps.Prompts
	case "streaming":
		return caps.Streaming
	case "sampling":
		return caps.Sampling
	case "progress":
		return caps.Progress
	default:
		// Check extensions
		for _, ext := range caps.Extensions {
			if ext == name {
				return true
			}
		}
		return false
	}
}
