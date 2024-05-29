package distributed

import (
	"hash/fnv"
	"sync"
)

// PartitionSize defines the number of partitions.
const PartitionSize = 10

// Partition represents a partition in the distributed key-value store.
type Partition struct {
	sync.RWMutex
	Data map[string]string // Data stored in the partition
}

// PartitionMap represents the mapping of keys to partitions.
type PartitionMap struct {
	partitions []*Partition // Array of partitions
}

// NewPartitionMap creates a new instance of PartitionMap with the specified number of partitions.
func NewPartitionMap() *PartitionMap {
	partitions := make([]*Partition, PartitionSize)
	for i := range partitions {
		partitions[i] = &Partition{
			Data: make(map[string]string),
		}
	}
	return &PartitionMap{
		partitions: partitions,
	}
}

// GetPartitionIndex returns the index of the partition for the given key.
func (p *PartitionMap) GetPartitionIndex(key string) int {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return int(hash.Sum32()) % PartitionSize
}

// Get retrieves the value associated with the given key.
func (p *PartitionMap) Get(key string) (string, bool) {
	index := p.GetPartitionIndex(key)
	p.partitions[index].RLock()
	defer p.partitions[index].RUnlock()
	value, ok := p.partitions[index].Data[key]
	return value, ok
}

// Put sets the value associated with the given key.
func (p *PartitionMap) Put(key, value string) {
	index := p.GetPartitionIndex(key)
	p.partitions[index].Lock()
	defer p.partitions[index].Unlock()
	p.partitions[index].Data[key] = value
}

// Delete removes the value associated with the given key.
func (p *PartitionMap) Delete(key string) {
	index := p.GetPartitionIndex(key)
	p.partitions[index].Lock()
	defer p.partitions[index].Unlock()
	delete(p.partitions[index].Data, key)
}
