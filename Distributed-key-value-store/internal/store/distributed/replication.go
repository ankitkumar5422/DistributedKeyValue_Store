package distributed

import (
	"sync"
)

// ReplicationFactor defines the number of replicas for each key.
const ReplicationFactor = 3

// ReplicatedStore represents a replicated key-value store.
type ReplicatedStore struct {
	Nodes []*DistributedStore
	mu    sync.RWMutex
}

// NewReplicatedStore creates a new instance of ReplicatedStore with the specified number of nodes.
func NewReplicatedStore(numNodes int) *ReplicatedStore {
	nodes := make([]*DistributedStore, numNodes)
	for i := range nodes {
		nodes[i] = NewDistributedStore()
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
		if value, ok := node.Get(key); ok {
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
		node.Put(key, value)
	}
}

// Delete removes the value associated with the given key.
func (s *ReplicatedStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, node := range s.Nodes {
		node.Delete(key)
	}
}
