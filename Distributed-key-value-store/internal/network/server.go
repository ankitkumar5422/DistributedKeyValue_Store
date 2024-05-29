package network

import (
	"distributed/internal/store/distributed"
	"fmt"
	"net/http"
)

// Server represents the key-value store server.
type Server struct {
	Node *distributed.Node
}

// NewServer creates a new instance of Server.
func NewServer(node *distributed.Node) *Server {
	return &Server{
		Node: node,
	}
}

// Start starts the server.
func (s *Server) Start() error {
	http.HandleFunc("/get", s.Node.HandleGetRequest)
	http.HandleFunc("/set", s.Node.HandleSetRequest)
	http.HandleFunc("/delete", s.Node.HandleDeleteRequest)
	fmt.Println("Server started on port 9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		return err
	}
	return nil
}
