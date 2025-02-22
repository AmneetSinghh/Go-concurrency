package main

import (
	"fmt"
	"sync"
	"time"
)

/*
problem


*/

var rwMu sync.RWMutex
var value int
var count int

func read(id int) {
	rwMu.RLock()
	defer rwMu.RUnlock()

	fmt.Printf("Reader %d: Reading value %d\n", id, value)
	count++
	time.Sleep(100 * time.Millisecond) // Simulate a read
}

func write(id int, newValue int) {
	rwMu.Lock()
	defer rwMu.Unlock()

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
	time.Sleep(10 * time.Second)
}
