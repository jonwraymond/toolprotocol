package elicit

import "time"

// Option configures elicit package types.
type Option func(any)

// DefaultRequestTimeout is the default timeout for elicitation requests.
const DefaultRequestTimeout = 30 * time.Second
