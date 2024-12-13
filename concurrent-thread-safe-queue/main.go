package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// ConcurrentQueue is a thread-safe queue using a mutex for synchronization.
type ConcurrentQueue struct {
	queue []int32  // Stores the queue elements.
	mu    sync.Mutex // Protects access to the queue.
}

// Enqueue adds an item to the queue in a thread-safe manner.
func (q *ConcurrentQueue) Enqueue(item int32) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, item)
}

// Dequeue removes and returns the first item in a thread-safe manner.
// Panics if the queue is empty.
func (q *ConcurrentQueue) Dequeue() int32 {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queue) == 0 {
		panic("cannot dequeue from an empty queue!")
	}

	item := q.queue[0]
	q.queue = q.queue[1:] // Shift remaining elements.
	return item
}

// Size returns the number of elements in the queue.
func (q *ConcurrentQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue)
}

// WaitGroups to synchronize goroutines.
var wgE sync.WaitGroup // Tracks enqueue operations.
var wgD sync.WaitGroup // Tracks dequeue operations.

func main() {
	// Initialize an empty queue.
	q1 := ConcurrentQueue{
		queue: make([]int32, 0),
	}

	// Launch 1,000,000 goroutines to enqueue random integers.
	for i := 0; i < 1000000; i++ {
		wgE.Add(1)
		go func() {
			q1.Enqueue(rand.Int31())
			wgE.Done()
		}()
	}

	// Launch 1,000,000 goroutines to dequeue elements.
	for i := 0; i < 1000000; i++ {
		wgD.Add(1)
		go func() {
			q1.Dequeue()
			wgD.Done()
		}()
	}

	// Wait for all enqueue and dequeue operations to complete.
	wgE.Wait()
	wgD.Wait()

	// Print the final size of the queue (should be 0 if no issues).
	fmt.Println(q1.Size())
}
