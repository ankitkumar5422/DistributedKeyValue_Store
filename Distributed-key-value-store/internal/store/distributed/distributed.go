package distributed

import "sync"

// DistributedStore implements a distributed key-value store.
type DistributedStore struct {
	data map[string]string
	mu   sync.RWMutex // Mutex for concurrent access to data
}

// NewDistributedStore creates a new instance of DistributedStore.
func NewDistributedStore() *DistributedStore {
	return &DistributedStore{
		data: make(map[string]string),
	}
}

// Get retrieves the value associated with the given key.
func (s *DistributedStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

// Put sets the value associated with the given key.
func (s *DistributedStore) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Delete removes the value associated with the given key.
func (s *DistributedStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}
