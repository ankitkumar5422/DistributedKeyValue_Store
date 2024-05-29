package distributed

import (
	"distributed/internal/store/inmemory"
	"sync"
)

// ReplicatedStore represents a replicated key-value store.
type ReplicatedStore struct {
	Nodes []*Node // Array of nodes
	mu    sync.RWMutex
}

// NewReplicatedStore creates a new instance of ReplicatedStore with the specified number of nodes.
func NewReplicatedStore(numNodes int, replicas int) *ReplicatedStore {
	nodes := make([]*Node, numNodes)
	for i := range nodes {
		nodes[i] = NewNode("node"+string(rune(i)), inmemory.NewInMemoryStore(), "", replicas) // Initialize each node with an in-memory store and replication factor
	}
	return &ReplicatedStore{
		Nodes: nodes,
	}
}

// Get retrieves the value associated with the given key.
func (s *ReplicatedStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, node := range s.Nodes {
		value, _ := node.Store.Get(key)
		if value != "" {
			return value, true
		}
	}
	return "", false
}

// Put sets the value associated with the given key.
func (s *ReplicatedStore) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, node := range s.Nodes {
		node.Store.Set(key, value)
	}
}

// Delete removes the value associated with the given key.
func (s *ReplicatedStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, node := range s.Nodes {
		node.Store.Delete(key)
	}
}
