package main

import (
	"fmt"
	"sync"
	"time"
)

/*
solution
*/

var rwMu sync.RWMutex
var mu sync.Mutex // Additional mutex to block new readers if writer is waiting
var value int

func read(id int) {
	mu.Lock() // will block if write is ongoing
	rwMu.RLock()
	mu.Unlock() //     // Allow other readers once the lock is acquired
	defer rwMu.RUnlock()

	fmt.Printf("Reader %d: Reading value %d\n", id, value)
	time.Sleep(1 * time.Second) // Simulate a read
}

func write(id int, newValue int) {
	mu.Lock() // prevents other readers from starting, write t krke hi unlock kruga,
	rwMu.Lock()
	defer rwMu.Unlock()
	defer mu.Unlock()

	fmt.Printf("--------- Writer %d: Writing value %d ------------- \n", id, newValue)
	value = newValue
	time.Sleep(1 * time.Second) // Simulate a write
}

func main() {
	// Start multiple readers
	for i := 1; i <= 100; i++ {
		go func(id int) {
			for {
				read(id) // Readers keep acquiring the lock
			}
		}(i)
	}

	// Start a single writer
	go func() {
		for {
			write(1, 42) // Writer waits indefinitely
		}
	}()

	// Run for a while to observe starvation
	time.Sleep(20 * time.Second)
}

/*
* Writers are given priority by the mu lock, but the Go scheduler enforces fairness.
* After two writers run consecutively, the scheduler gives readers a chance to proceed, ensuring that they don’t get starved.
* The scheduler's fairness mechanism ensures that a batch of readers executes before switching back to the writers.

 */
