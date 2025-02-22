package main

import (
	"fmt"
	"sync"
	"time"
)

/*
* using mutex on each goroutine.
* Execution time: 965 ms -1 second


Contention occurs because all goroutines try to access the same lock, which results in high context switching overhead.
The mutex lock and unlock operations themselves are relatively expensive when there is heavy contention. So at each time
only 1 goroutine can access mutex lock, so 499999 will be waiting for the lock to be released.

*/

func main() {
	start := time.Now()
	var mu sync.Mutex
	counter := 0

	// Increment counter concurrently using mutex
	var wg sync.WaitGroup
	for i := 0; i < 5000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("Final Counter:", counter)

	elapsed := time.Since(start)
	fmt.Printf("Program completed in %v\n", elapsed)
}

/*
* using channels
* Execution time: 630-640 ms

- Channels introduce some overhead due to message passing but generally avoid the contention issues seen with mutexes.
*/

func ConcurrencyUsingChannels() {
	start := time.Now()
	counter := make(chan int)
	done := make(chan struct{})

	// Goroutine to manage the counter
	go func() {
		count := 0
		for {
			select {
			case c, ok := <-counter:
				if !ok {
					// Exit if the counter channel is closed
					done <- struct{}{}
					fmt.Println("Final Counter:", count)
					return
				}
				count += c
			}
		}
	}()

	// Send 5 increments to the counter channel
	for i := 0; i < 5000000; i++ {
		counter <- 1
	}

	// Close the counter channel to signal no more data will be sent
	close(counter)

	// Wait for the goroutine to complete
	<-done

	elapsed := time.Since(start)
	fmt.Printf("Program completed in %v\n", elapsed)
}

/*
* using mutex + batch processing
* Execution time: <=1 ms

This approach uses batch processing: instead of locking and unlocking the mutex for each increment, each goroutine processes a batch of 50,000 increments at once.
By batching the work, we significantly reduce the number of times the mutex is locked and unlocked, which greatly reduces the contention and context-switching overhead.
Worker pool: A fixed number of 100 goroutines are used, which makes the program more efficient by avoiding the overhead of creating too many goroutines.

- mutex lock/unlock minimized so less contention, execusion becomes faster.
*/

func GoBatchProcessing() {
	start := time.Now()

	var mu sync.Mutex
	counter := 0
	numGoroutines := 100 // Reduce the number of goroutines
	batchSize := 50000   // Number of increments each goroutine processes
	var wg sync.WaitGroup

	// Creating a fixed number of workers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Each worker will perform batchSize increments
			localCount := 0
			for j := 0; j < batchSize; j++ {
				localCount++
			}

			// Lock only once per batch of work
			mu.Lock()
			counter += localCount
			mu.Unlock()
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("Final Counter:", counter)

	elapsed := time.Since(start)
	fmt.Printf("Program completed in %v\n", elapsed)
}
