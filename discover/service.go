package discover

import "errors"

// Service is the default implementation of Discoverable.
type Service struct {
	id          string
	name        string
	description string
	version     string
	endpoint    string
	caps        *Capabilities
}

// NewService creates a new service with the given ID and endpoint.
func NewService(id, endpoint string) *Service {
	return &Service{
		id:       id,
		endpoint: endpoint,
		caps:     &Capabilities{},
	}
}

// ID returns the service identifier.
func (s *Service) ID() string {
	return s.id
}

// Name returns the service name (defaults to ID if not set).
func (s *Service) Name() string {
	if s.name != "" {
		return s.name
	}
	return s.id
}

// SetName sets the service name.
func (s *Service) SetName(name string) *Service {
	s.name = name
	return s
}

// Description returns the service description.
func (s *Service) Description() string {
	return s.description
}

// SetDescription sets the service description.
func (s *Service) SetDescription(desc string) *Service {
	s.description = desc
	return s
}

// Version returns the service version.
func (s *Service) Version() string {
	return s.version
}

// SetVersion sets the service version.
func (s *Service) SetVersion(version string) *Service {
	s.version = version
	return s
}

// Endpoint returns the service endpoint URL.
func (s *Service) Endpoint() string {
	return s.endpoint
}

// Capabilities returns the service capabilities.
func (s *Service) Capabilities() *Capabilities {
	if s.caps == nil {
		return &Capabilities{}
	}
	return s.caps
}

// SetCapabilities sets the service capabilities.
func (s *Service) SetCapabilities(caps *Capabilities) *Service {
	s.caps = caps
	return s
}

// WithCapability adds a capability to the service.
func (s *Service) WithCapability(name string) *Service {
	if s.caps == nil {
		s.caps = &Capabilities{}
	}
	switch name {
	case "tools":
		s.caps.Tools = true
	case "resources":
		s.caps.Resources = true
	case "prompts":
		s.caps.Prompts = true
	case "streaming":
		s.caps.Streaming = true
	case "sampling":
		s.caps.Sampling = true
	case "progress":
		s.caps.Progress = true
	default:
		s.caps.Extensions = append(s.caps.Extensions, name)
	}
	return s
}

// Validate checks that the service has required fields.
func (s *Service) Validate() error {
	if s.id == "" {
		return errors.New("service ID is required")
	}
	if s.endpoint == "" {
		return errors.New("service endpoint is required")
	}
	return nil
}
