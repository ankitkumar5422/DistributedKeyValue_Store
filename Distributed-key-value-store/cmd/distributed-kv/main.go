package main

import (
	"distributed/internal/network"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello, Distributed Key-Value Store!")

	// Create a new server instance
	server := network.NewServer()

	// Start the key-value store server
	err := server.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)

	}
}
