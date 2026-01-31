package task

// Option configures a DefaultManager.
type Option func(*DefaultManager)

// WithStore configures the manager to use the given store.
func WithStore(store Store) Option {
	return func(m *DefaultManager) {
		m.store = store
	}
}
