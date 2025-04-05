package distributed

import (
	"crypto/sha256"
	"encoding/binary"
	"sort"
	"strconv"
	"sync"
)

// PartitionSize defines the number of partitions in the system.
const PartitionSize = 10

// VirtualNodes defines the number of virtual nodes per partition for consistent hashing.
const VirtualNodes = 100

// Partition represents a single partition with a thread-safe map to store key-value pairs.
type Partition struct {
	sync.RWMutex                   // Read-Write mutex for thread-safe access to the partition.
	Data         map[string]string // Map to store key-value pairs.
}

// PartitionMap represents the distributed key-value store's partitioning mechanism.
type PartitionMap struct {
	partitions []*Partition   // Array of partitions.
	ring       []uint32       // Sorted list of hash values for consistent hashing.
	vnodeMap   map[uint32]int // Maps hash values to partition indices.
	once       sync.Once      // Ensures the ring is initialized only once.
}

// NewPartitionMap initializes and returns a new PartitionMap.
func NewPartitionMap() *PartitionMap {
	// Create the partitions.
	partitions := make([]*Partition, PartitionSize)
	for i := range partitions {
		partitions[i] = &Partition{
			Data: make(map[string]string), // Initialize the map for each partition.
		}
	}
	// Initialize the PartitionMap structure.
	pm := &PartitionMap{
		partitions: partitions,
		ring:       []uint32{},
		vnodeMap:   map[uint32]int{},
	}
	// Initialize the consistent hashing ring.
	pm.initRing()
	return pm
}

// hashKey computes a consistent hash for a given key using SHA256.
func hashKey(key string) uint32 {
	sum := sha256.Sum256([]byte(key))       // Compute SHA256 hash of the key.
	return binary.BigEndian.Uint32(sum[:4]) // Use the first 4 bytes as the hash value.
}

// initRing initializes the consistent hashing ring with virtual nodes.
func (p *PartitionMap) initRing() {
	for i := 0; i < PartitionSize; i++ { // Iterate over all partitions.
		for v := 0; v < VirtualNodes; v++ { // Create virtual nodes for each partition.
			vkey := "P" + strconv.Itoa(i) + "-V" + strconv.Itoa(v) // Generate a unique virtual node key.
			h := hashKey(vkey)                                     // Compute the hash for the virtual node.
			p.ring = append(p.ring, h)                             // Add the hash to the ring.
			p.vnodeMap[h] = i                                      // Map the hash to the partition index.
		}
	}
	// Sort the ring to enable binary search for consistent hashing.
	sort.Slice(p.ring, func(i, j int) bool { return p.ring[i] < p.ring[j] })
}

// GetPartitionIndex finds the partition index for a given key using consistent hashing.
func (p *PartitionMap) GetPartitionIndex(key string) int {
	p.once.Do(p.initRing) // Ensure the ring is initialized only once.
	h := hashKey(key)     // Compute the hash of the key.
	// Find the smallest hash in the ring that is greater than or equal to the key's hash.
	idx := sort.Search(len(p.ring), func(i int) bool {
		return p.ring[i] >= h
	})
	if idx == len(p.ring) { // If no such hash exists, wrap around to the first hash.
		idx = 0
	}
	return p.vnodeMap[p.ring[idx]] // Return the partition index corresponding to the hash.
}

// Get retrieves the value associated with a key from the appropriate partition.
func (p *PartitionMap) Get(key string) (string, bool) {
	index := p.GetPartitionIndex(key)          // Find the partition index for the key.
	p.partitions[index].RLock()                // Acquire a read lock on the partition.
	defer p.partitions[index].RUnlock()        // Release the lock after the operation.
	value, ok := p.partitions[index].Data[key] // Retrieve the value from the partition.
	return value, ok
}

// Put stores a key-value pair in the appropriate partition.
func (p *PartitionMap) Put(key, value string) {
	index := p.GetPartitionIndex(key)     // Find the partition index for the key.
	p.partitions[index].Lock()            // Acquire a write lock on the partition.
	defer p.partitions[index].Unlock()    // Release the lock after the operation.
	p.partitions[index].Data[key] = value // Store the key-value pair in the partition.
}

// Delete removes a key-value pair from the appropriate partition.
func (p *PartitionMap) Delete(key string) {
	index := p.GetPartitionIndex(key)     // Find the partition index for the key.
	p.partitions[index].Lock()            // Acquire a write lock on the partition.
	defer p.partitions[index].Unlock()    // Release the lock after the operation.
	delete(p.partitions[index].Data, key) // Remove the key-value pair from the partition.
}
