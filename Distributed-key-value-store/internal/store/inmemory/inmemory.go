package inmemory

import (
	"sync"
)

// InMemoryStore implements an in-memory key-value store.
type InMemoryStore struct {
	data  map[string]string
	mutex sync.RWMutex
}

// NewInMemoryStore creates a new instance of InMemoryStore.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]string),
	}
}

// Get retrieves the value associated with the given key.
func (s *InMemoryStore) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

// Set sets the value associated with the given key.
func (s *InMemoryStore) Put(key, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[key] = value
}

// Delete removes the value associated with the given key.
func (s *InMemoryStore) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, key)
}
