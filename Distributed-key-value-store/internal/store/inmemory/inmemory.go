package inmemory

import (
	"errors"
	"log"
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
func (s *InMemoryStore) Get(key string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, ok := s.data[key]
	if !ok {
		log.Printf("Get: key %s not found", key)
		return "", errors.New("key not found")
	}

	log.Printf("Get: key %s, value %s", key, value)
	return value, nil
}

// Set sets the value associated with the given key.
func (s *InMemoryStore) Set(key, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = value
	log.Printf("Set: key %s, value %s", key, value)
	return nil
}

// Delete removes the value associated with the given key.
func (s *InMemoryStore) Delete(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.data[key]; !ok {
		log.Printf("Delete: key %s not found", key)
		return errors.New("key not found")
	}

	delete(s.data, key)
	log.Printf("Delete: key %s", key)
	return nil
}
