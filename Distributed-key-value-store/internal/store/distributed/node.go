package distributed

import (
	"bytes"
	"distributed/internal/store/inmemory"
	model "distributed/pkg/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Node represents a single node in the distributed key-value store.
type Node struct {
	ID       string
	Store    *inmemory.InMemoryStore
	Peers    map[string]string // Map of peer node IDs to their addresses
	mutex    sync.RWMutex      // Mutex for synchronizing access to Peers
	address  string            // Address of this node
	replicas int               // Number of replicas
}

// NewNode creates a new Node instance.
func NewNode(id string, store *inmemory.InMemoryStore, address string, replicas int) *Node {
	return &Node{
		ID:       id,
		Store:    store,
		Peers:    make(map[string]string),
		address:  address,
		replicas: replicas,
		mutex:    sync.RWMutex{}, // Initialize the mutex
	}
}

// AddPeer adds a new peer node to the node's peer list.
func (n *Node) AddPeer(id, address string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.Peers[id] = address
}

// RemovePeer removes a peer node from the node's peer list.
func (n *Node) RemovePeer(id string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	delete(n.Peers, id)
}

// HandleGetRequest handles GET requests to retrieve values.
func (n *Node) HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is missing", http.StatusBadRequest)
		return
	}

	value, err := n.Store.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if value == "" {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	response := model.KeyValue{Key: key, Value: value}
	json.NewEncoder(w).Encode(response)
}

// HandleSetRequest handles POST requests to set values.
func (n *Node) HandleSetRequest(w http.ResponseWriter, r *http.Request) {
	kv := model.KeyValue{}
	if err := json.NewDecoder(r.Body).Decode(&kv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := n.Store.Set(kv.Key, kv.Value); err != nil {
		log.Println("Error setting key-value pair in store:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	n.replicateSet(kv.Key, kv.Value)

	w.WriteHeader(http.StatusOK)
}

// HandleDeleteRequest handles DELETE requests to delete values.
func (n *Node) HandleDeleteRequest(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is missing", http.StatusBadRequest)
		return
	}

	if err := n.Store.Delete(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	n.replicateDelete(key)

	w.WriteHeader(http.StatusOK)
}

// replicateSet replicates a set operation to peer nodes.
func (n *Node) replicateSet(key, value string) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	for _, address := range n.Peers {
		go func(addr string) {
			client := &http.Client{}
			kv := model.KeyValue{Key: key, Value: value}
			data, _ := json.Marshal(kv)
			req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s/set", addr), bytes.NewBuffer(data))
			client.Do(req)
		}(address)
	}
}

// replicateDelete replicates a delete operation to peer nodes.
func (n *Node) replicateDelete(key string) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	for _, address := range n.Peers {
		go func(addr string) {
			client := &http.Client{}
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("http://%s/delete?key=%s", addr, key), nil)
			client.Do(req)
		}(address)
	}
}
