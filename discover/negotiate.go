package discover

// NegotiateStrategy defines how capabilities are negotiated.
type NegotiateStrategy int

const (
	// NegotiateIntersect uses AND logic (only common capabilities).
	NegotiateIntersect NegotiateStrategy = iota

	// NegotiateMerge uses OR logic (combined capabilities).
	NegotiateMerge
)

// Negotiator handles capability negotiation with configurable strategy.
type Negotiator struct {
	Strategy NegotiateStrategy
}

// Negotiate determines compatible capabilities between client and server.
func (n *Negotiator) Negotiate(client, server *Capabilities) *Capabilities {
	if client == nil && server == nil {
		return &Capabilities{}
	}
	if client == nil {
		return server.clone()
	}
	if server == nil {
		return client.clone()
	}

	switch n.Strategy {
	case NegotiateMerge:
		return client.Merge(server)
	default:
		return client.Intersect(server)
	}
}
