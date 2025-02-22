package main

import (
	"fmt"
	"sync"
	"time"
)

type Button struct {
	Clicked *sync.Cond
}

func main() {
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}
	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)

	subscribe := func(c *sync.Cond, fn func()) {
		c.L.Lock()
		defer c.L.Unlock()
		c.Wait()
		fn()
	}

	// Start the goroutines that subscribe to the button click
	go subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	go subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})
	go subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})

	// Ensure that all goroutines have started waiting
	time.Sleep(500 * time.Millisecond)

	// Now broadcast the signal after all goroutines are waiting
	button.Clicked.Broadcast()

	// Wait until all tasks are done
	clickRegistered.Wait()
}
