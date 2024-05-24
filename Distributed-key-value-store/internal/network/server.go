package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Server represents the key-value store server.
type Server struct {
	store map[string]string // Key-value store
	mutex sync.RWMutex      // Mutex for concurrent access
}

// NewServer creates a new instance of Server.
func NewServer() *Server {
	// Initialize server fields
	return &Server{
		store: make(map[string]string),
		mutex: sync.RWMutex{},
	}
}

// Start starts the key-value store server.
func (s *Server) Start() error {
	// Implement server start logic
	http.HandleFunc("/get", s.HandleGetRequest)
	http.HandleFunc("/set", s.HandleSetRequest)
	http.HandleFunc("/delete", s.HandleDeleteRequest)
	fmt.Println("Server started on port 4006")
	err := http.ListenAndServe(":4006", nil)
	if err != nil {
		return err
	}
	return nil
}

// HandleGetRequest handles GET requests to retrieve values.
func (s *Server) HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	// Implement logic to handle GET requests
	key := r.URL.Query().Get("key")
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, ok := s.store[key]
	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(value)
}

// HandleSetRequest handles POST requests to set values.
func (s *Server) HandleSetRequest(w http.ResponseWriter, r *http.Request) {
	// Implement logic to handle POST requests
	// Example: parse JSON request body
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validate request body
	if data["key"] == "" || data["value"] == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Store key-value pair in key-value store
	s.store[data["key"]] = data["value"]
	fmt.Fprintf(w, "Key-value pair set successfully")
}

// HandleDeleteRequest handles DELETE requests to delete values.
func (s *Server) HandleDeleteRequest(w http.ResponseWriter, r *http.Request) {
	// Implement logic to handle DELETE requests
	key := r.URL.Query().Get("key")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Delete key-value pair from key-value store
	if _, ok := s.store[key]; !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	delete(s.store, key)
	fmt.Fprintf(w, "Key deleted successfully")
}
