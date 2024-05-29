# Distributed Key-Value Store

This project implements a simple distributed key-value store with replication and partitioning. The system consists of multiple nodes that store key-value pairs, replicate data across nodes, and partition data for efficient access and scalability.

## Features

- **Replication**: Ensures that data is replicated across multiple nodes for fault tolerance.
- **Partitioning**: Distributes data across multiple partitions for efficient access and scalability.
- **Concurrency**: Handles concurrent access to the store with thread-safe operations.

## Project Structure

The project is organized into the following packages:

- `distributed`: Contains the core logic for the distributed key-value store, including node management, replication, and partitioning.
- `inmemory`: Implements an in-memory key-value store used by each node.
- `network`: Implements the HTTP server that handles client requests to interact with the distributed store.
- `pkg/models`: Defines the data models used in the project.

## Setup and Running the Project

### Prerequisites

- Go 1.16 or later
- A web browser to access the frontend

### Getting Started

1. **Clone the Repository**:

   ```sh
   git clone https://github.com/yourusername/distributed-kv-store.git
   cd distributed-kv-store
2.**Build the Project**:
    go build -o kvstore ./cmd/server

**How to Use**
1.Set Key-Value Pair: Enter a key and value, then click "Set". This will store the key-value pair in the distributed store and replicate it to other nodes.
2.Get Value by Key: Enter a key and click "Get". This will retrieve the value for the given key from the distributed store.
3.Node Status: The node status section shows the current status of all nodes in the system.
Example Usage
1.Set a Key-Value Pair
   ```sh
  curl -X POST http://localhost:9090/set -H "Content-Type: application/json" -d '{"key":"name","value":"John Doe"}'
2.Get Value by Key:
      curl http://localhost:9090/get?key=name

3.Delete Key-Value Pair:
    curl -X DELETE http://localhost:9090/delete?key=name    
