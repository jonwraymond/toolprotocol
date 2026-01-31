package discover

// Capabilities describes protocol features supported by a service.
type Capabilities struct {
	// Tools indicates support for tool invocation.
	Tools bool

	// Resources indicates support for resource access.
	Resources bool

	// Prompts indicates support for prompt templates.
	Prompts bool

	// Streaming indicates support for streaming responses.
	Streaming bool

	// Sampling indicates support for LLM sampling.
	Sampling bool

	// Progress indicates support for progress notifications.
	Progress bool

	// Extensions lists additional protocol extensions supported.
	Extensions []string
}

// Merge combines capabilities from two sources (OR logic for bools, union for extensions).
func (c *Capabilities) Merge(other *Capabilities) *Capabilities {
	if other == nil {
		return c.clone()
	}

	merged := &Capabilities{
		Tools:     c.Tools || other.Tools,
		Resources: c.Resources || other.Resources,
		Prompts:   c.Prompts || other.Prompts,
		Streaming: c.Streaming || other.Streaming,
		Sampling:  c.Sampling || other.Sampling,
		Progress:  c.Progress || other.Progress,
	}

	// Union of extensions (deduplicated)
	seen := make(map[string]bool)
	for _, ext := range c.Extensions {
		if !seen[ext] {
			seen[ext] = true
			merged.Extensions = append(merged.Extensions, ext)
		}
	}
	for _, ext := range other.Extensions {
		if !seen[ext] {
			seen[ext] = true
			merged.Extensions = append(merged.Extensions, ext)
		}
	}

	return merged
}

// Intersect finds common capabilities (AND logic for bools, intersection for extensions).
func (c *Capabilities) Intersect(other *Capabilities) *Capabilities {
	if other == nil {
		return &Capabilities{}
	}

	intersected := &Capabilities{
		Tools:     c.Tools && other.Tools,
		Resources: c.Resources && other.Resources,
		Prompts:   c.Prompts && other.Prompts,
		Streaming: c.Streaming && other.Streaming,
		Sampling:  c.Sampling && other.Sampling,
		Progress:  c.Progress && other.Progress,
	}

	// Intersection of extensions
	otherSet := make(map[string]bool)
	for _, ext := range other.Extensions {
		otherSet[ext] = true
	}
	for _, ext := range c.Extensions {
		if otherSet[ext] {
			intersected.Extensions = append(intersected.Extensions, ext)
		}
	}

	return intersected
}

// clone creates a copy of capabilities.
func (c *Capabilities) clone() *Capabilities {
	if c == nil {
		return nil
	}
	cloned := &Capabilities{
		Tools:     c.Tools,
		Resources: c.Resources,
		Prompts:   c.Prompts,
		Streaming: c.Streaming,
		Sampling:  c.Sampling,
		Progress:  c.Progress,
	}
	if len(c.Extensions) > 0 {
		cloned.Extensions = make([]string, len(c.Extensions))
		copy(cloned.Extensions, c.Extensions)
	}
	return cloned
}
