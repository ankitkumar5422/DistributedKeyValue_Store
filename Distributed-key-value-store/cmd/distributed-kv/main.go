package main

import (
	"distributed/internal/network"
	"distributed/internal/store/distributed"
	"distributed/internal/store/inmemory"
	"log"
)

func main() {
	store := inmemory.NewInMemoryStore()

	// Create a new distributed node with the in-memory store
	node := distributed.NewNode("node1", store, "localhost:4006", 3)

	// Initialize the server with the node
	server := network.NewServer(node)

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
