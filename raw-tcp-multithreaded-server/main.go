package main

import (
	"log"
	"net"
	"sync"
	"time"
)

const (
	port              = ":1729"          // Port the server will listen on
	maxWorkers        = 10               // Maximum number of worker goroutines
	connectionTimeout = 10 * time.Second // Connection timeout duration
)

// Worker pool structure
type WorkerPool struct {
	jobQueue chan net.Conn // Channel to handle incoming connections
	wg       sync.WaitGroup
}

// Worker function to process connections
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for conn := range wp.jobQueue {
		log.Printf("[Worker %d] Processing new connection\n", id)
		handleConnection(conn)
		log.Printf("[Worker %d] Finished processing\n", id)
	}
}

// Add a connection to the worker pool
func (wp *WorkerPool) addJob(conn net.Conn) {
	wp.jobQueue <- conn
}

// Initialize a worker pool
func NewWorkerPool(workerCount int) *WorkerPool {
	wp := &WorkerPool{
		jobQueue: make(chan net.Conn, 100), // Buffered channel for queued connections
	}
	wp.wg.Add(workerCount)

	// Start workers
	for i := 1; i <= workerCount; i++ {
		go wp.worker(i)
	}

	return wp
}

// Shutdown the worker pool
func (wp *WorkerPool) shutdown() {
	close(wp.jobQueue)
	wp.wg.Wait()
	log.Println("[Server] Worker pool has shut down")
}

// Handle individual connections
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Set a deadline for the connection to prevent indefinite blocking
	conn.SetDeadline(time.Now().Add(connectionTimeout))

	// Read data from the client
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Printf("[Error] Failed to read from connection: %v\n", err)
		return
	}

	log.Println("[Server] Received data from client. Processing request...")

	// Simulate long-running processing
	time.Sleep(8 * time.Second)

	// Write response to the client
	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nhello world!\r\n"))
	if err != nil {
		log.Printf("[Error] Failed to write to connection: %v\n", err)
		return
	}

	log.Println("[Server] Successfully responded to client")
}

func main() {
	// Start listening for connections
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("[Error] Failed to start server: %v\n", err)
	}
	defer listener.Close()

	log.Printf("[Server] Listening on port %s\n", port)

	// Initialize a worker pool
	workerPool := NewWorkerPool(maxWorkers)
	defer workerPool.shutdown()

	for {
		log.Println("[Server] Waiting for a client to connect...")

		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[Error] Failed to accept connection: %v\n", err)
			continue
		}

		log.Printf("[Server] Client connected: %s\n", conn.RemoteAddr())

		// Add connection to worker pool for processing
		workerPool.addJob(conn)
	}
}
