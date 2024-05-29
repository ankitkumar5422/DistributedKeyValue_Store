package store

import (
	"errors"
	"sync"
)

// KVStore represents a simple in-memory key-value store.
type KVStore struct {
	store map[string]string // Key-value store
	mutex sync.RWMutex      // Mutex for concurrent access
}

// NewKVStore creates a new instance of KVStore.
func NewKVStore() *KVStore {
	// Initialize KVStore fields
	return &KVStore{
		store: make(map[string]string),
	}
}

// Get retrieves the value for a given key.
func (k *KVStore) Get(key string) (string, error) {
	k.mutex.RLock()
	defer k.mutex.RUnlock()

	value, ok := k.store[key]
	if !ok {
		return "", errors.New("key not found: " + key)
	}
	return value, nil
}

// Set sets the value for a given key.
func (k *KVStore) Set(key string, value string) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	k.store[key] = value
	return nil
}

// Delete deletes the value for a given key.
func (k *KVStore) Delete(key string) error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	if _, ok := k.store[key]; !ok {
		return errors.New("key not found: " + key)
	}
	delete(k.store, key)
	return nil
}
