package task

import (
	"context"
	"sync"
)

// MemoryStore is an in-memory implementation of Store.
type MemoryStore struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

// NewMemoryStore creates a new in-memory task store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks: make(map[string]*Task),
	}
}

// Save saves or updates a task.
func (s *MemoryStore) Save(ctx context.Context, task *Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store a clone to prevent external modifications
	s.tasks[task.ID] = task.Clone()
	return nil
}

// Load retrieves a task by ID.
func (s *MemoryStore) Load(ctx context.Context, id string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, ErrTaskNotFound
	}

	// Return a clone to prevent external modifications
	return task.Clone(), nil
}

// LoadAll retrieves all tasks.
func (s *MemoryStore) LoadAll(ctx context.Context) ([]*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task.Clone())
	}
	return result, nil
}

// Delete removes a task by ID.
func (s *MemoryStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}
