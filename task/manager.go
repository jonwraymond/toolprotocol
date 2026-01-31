package task

import (
	"context"
	"sync"
	"time"
)

// DefaultManager is the default implementation of Manager.
type DefaultManager struct {
	mu          sync.RWMutex
	store       Store
	subscribers map[string][]chan *Task
}

// NewManager creates a new task manager with default configuration.
func NewManager(opts ...Option) *DefaultManager {
	m := &DefaultManager{
		store:       NewMemoryStore(),
		subscribers: make(map[string][]chan *Task),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Create creates a new task with the given ID.
func (m *DefaultManager) Create(ctx context.Context, id string) (*Task, error) {
	if id == "" {
		return nil, ErrEmptyID
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if task already exists
	if _, err := m.store.Load(ctx, id); err == nil {
		return nil, ErrTaskExists
	}

	now := time.Now()
	task := &Task{
		ID:        id,
		State:     StatePending,
		Progress:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := m.store.Save(ctx, task); err != nil {
		return nil, err
	}

	return task.Clone(), nil
}

// Get retrieves a task by ID.
func (m *DefaultManager) Get(ctx context.Context, id string) (*Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.store.Load(ctx, id)
}

// List returns all tasks.
func (m *DefaultManager) List(ctx context.Context) ([]*Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.store.LoadAll(ctx)
}

// Update updates task progress and message.
// Transitions pending → running on first update.
func (m *DefaultManager) Update(ctx context.Context, id string, progress float64, message string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, err := m.store.Load(ctx, id)
	if err != nil {
		return err
	}

	// Check if task is in terminal state
	if task.State.IsTerminal() {
		return ErrInvalidTransition
	}

	// Transition pending → running on first update
	if task.State == StatePending {
		task.State = StateRunning
	}

	task.Progress = progress
	task.Message = message
	task.UpdatedAt = time.Now()

	if err := m.store.Save(ctx, task); err != nil {
		return err
	}

	m.notifySubscribers(id, task)
	return nil
}

// Complete marks the task as complete with the given result.
func (m *DefaultManager) Complete(ctx context.Context, id string, result any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, err := m.store.Load(ctx, id)
	if err != nil {
		return err
	}

	// Can only complete from pending or running state
	if task.State.IsTerminal() {
		return ErrInvalidTransition
	}

	now := time.Now()
	task.State = StateComplete
	task.Progress = 1.0
	task.Result = result
	task.UpdatedAt = now
	task.CompletedAt = &now

	if err := m.store.Save(ctx, task); err != nil {
		return err
	}

	m.notifySubscribers(id, task)
	m.closeSubscribers(id)
	return nil
}

// Fail marks the task as failed with the given error.
func (m *DefaultManager) Fail(ctx context.Context, id string, err error) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, loadErr := m.store.Load(ctx, id)
	if loadErr != nil {
		return loadErr
	}

	// Can only fail from pending or running state
	if task.State.IsTerminal() {
		return ErrInvalidTransition
	}

	now := time.Now()
	task.State = StateFailed
	task.Error = err
	task.UpdatedAt = now
	task.CompletedAt = &now

	if saveErr := m.store.Save(ctx, task); saveErr != nil {
		return saveErr
	}

	m.notifySubscribers(id, task)
	m.closeSubscribers(id)
	return nil
}

// Cancel cancels the task.
func (m *DefaultManager) Cancel(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, err := m.store.Load(ctx, id)
	if err != nil {
		return err
	}

	// Can only cancel from pending or running state
	if task.State.IsTerminal() {
		return ErrInvalidTransition
	}

	now := time.Now()
	task.State = StateCancelled
	task.UpdatedAt = now
	task.CompletedAt = &now

	if err := m.store.Save(ctx, task); err != nil {
		return err
	}

	m.notifySubscribers(id, task)
	m.closeSubscribers(id)
	return nil
}

// Subscribe returns a channel that receives task updates.
// The channel is closed when the task reaches a terminal state.
func (m *DefaultManager) Subscribe(ctx context.Context, id string) (<-chan *Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Verify task exists
	task, err := m.store.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	// If already terminal, return a closed channel
	if task.State.IsTerminal() {
		ch := make(chan *Task)
		close(ch)
		return ch, nil
	}

	// Create buffered channel for non-blocking sends
	ch := make(chan *Task, 10)
	m.subscribers[id] = append(m.subscribers[id], ch)

	// Handle context cancellation
	go func() {
		<-ctx.Done()
		m.mu.Lock()
		defer m.mu.Unlock()
		m.removeSubscriber(id, ch)
	}()

	return ch, nil
}

// notifySubscribers sends task update to all subscribers.
// Must be called with m.mu held.
func (m *DefaultManager) notifySubscribers(id string, task *Task) {
	subs := m.subscribers[id]
	clone := task.Clone()
	for _, ch := range subs {
		select {
		case ch <- clone:
		default:
			// Channel full, skip (non-blocking)
		}
	}
}

// closeSubscribers closes all subscriber channels for a task.
// Must be called with m.mu held.
func (m *DefaultManager) closeSubscribers(id string) {
	for _, ch := range m.subscribers[id] {
		close(ch)
	}
	delete(m.subscribers, id)
}

// removeSubscriber removes a specific subscriber channel.
// Must be called with m.mu held.
func (m *DefaultManager) removeSubscriber(id string, ch chan *Task) {
	subs := m.subscribers[id]
	for i, sub := range subs {
		if sub == ch {
			m.subscribers[id] = append(subs[:i], subs[i+1:]...)
			close(ch)
			break
		}
	}
	if len(m.subscribers[id]) == 0 {
		delete(m.subscribers, id)
	}
}
