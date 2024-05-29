package distributed

const ReplicationFactor = 3

// DistributedStore represents a distributed key-value store with partitioning and replication.
type DistributedStore struct {
	partitionMap *PartitionMap
	replicas     int
}

// NewDistributedStore creates a new instance of DistributedStore with the specified number of partitions and replicas.
func NewDistributedStore(replicas int) *DistributedStore {
	return &DistributedStore{
		partitionMap: NewPartitionMap(),
		replicas:     replicas,
	}
}

// Get retrieves the value associated with the given key.
func (s *DistributedStore) Get(key string) (string, bool) {
	return s.partitionMap.Get(key)
}

// Put sets the value associated with the given key.
func (s *DistributedStore) Put(key, value string) {
	s.partitionMap.Put(key, value)
	s.replicatePut(key, value)
}

// Delete removes the value associated with the given key.
func (s *DistributedStore) Delete(key string) {
	s.partitionMap.Delete(key)
	s.replicateDelete(key)
}

// replicatePut replicates a Put operation to replica nodes.
func (s *DistributedStore) replicatePut(key, value string) {
	partitionIndex := s.partitionMap.GetPartitionIndex(key)
	for i := 1; i <= s.replicas; i++ {
		replicaIndex := (partitionIndex + i) % len(s.partitionMap.partitions)
		s.partitionMap.partitions[replicaIndex].Data[key] = value
	}
}

// replicateDelete replicates a Delete operation to replica nodes.
func (s *DistributedStore) replicateDelete(key string) {
	partitionIndex := s.partitionMap.GetPartitionIndex(key)
	for i := 1; i <= s.replicas; i++ {
		replicaIndex := (partitionIndex + i) % len(s.partitionMap.partitions)
		delete(s.partitionMap.partitions[replicaIndex].Data, key)
	}
}
