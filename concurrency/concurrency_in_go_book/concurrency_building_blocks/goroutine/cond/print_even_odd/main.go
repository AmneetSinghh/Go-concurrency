package main

import (
	"fmt"
	"sync"
	"time"
)

var count int

func conditionTrue(c *sync.Cond, number int) {
	// Check if 1 second has passed since startTime
	fmt.Println("entering conditionTrue ", number)
	time.Sleep(2 * time.Second)
	c.Signal()
	//fmt.Println("singal done ", number)

}

func run(c *sync.Cond, number int) {
	fmt.Println("gouritng enters", number, " ", count)
	c.L.Lock()
	count++
	fmt.Println("wait is starting ", number, " ", count)
	c.Wait() // goroutine suspended, niche wala part not runs. allowing other gorouitne to run and take lock mutex released.
	/*
	 * PART- AFTER  wait will not execute until we don't get the signal.
	 */
	fmt.Println("wait is ended because c.signal() completed ", number)

	fmt.Println("do processing ", number)
	c.L.Unlock()
}

func run1() {

	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock() // can only aquire, if someone unlocks or suspended or waits.
		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}
}

func run2() {
	c := sync.NewCond(&sync.Mutex{})
	for i := 1; i <= 3; i++ {
		go run(c, i)
		// Wait until the condition is satisfied (signaled)
		go conditionTrue(c, i)
	}
	time.Sleep(3 * time.Second)
}

func run3() {
	fmt.Println("-------------------------------------------even and odd--------------------------------------------------")
	c := sync.NewCond(&sync.Mutex{})
	var waitingState = 0
	go even(c, &waitingState)
	go odd(c, &waitingState)
	time.Sleep(4 * time.Second)
}

func main() {

	operation := 1

	switch operation {
	case 1:
		run1()
	case 2:
		/*
		* how cond working.
		 */
		run2()
	case 3:
		/*
		*print even odd
		 */
		run3()
	}

}

func even(c *sync.Cond, waitingState *int) {
	for i := 0; i <= 10; i += 2 {
		c.L.Lock()
		if *waitingState == 1 { // is odd
			c.Wait() // it means release the lock so that second can go into executing.
		}

		fmt.Println("even -> ", i)
		*waitingState = 1
		c.Signal()
		c.L.Unlock()
	}
}

// if odd gorouitne runs little late, then c.signal() which we get from even will be wasted.
func odd(c *sync.Cond, waitingState *int) {
	for i := 1; i <= 10; i += 2 {
		//time.Sleep(300 * time.Millisecond)
		c.L.Lock()
		if *waitingState == 0 { // is odd
			c.Wait() // it means release the lock so that second can go into executing.
		}
		fmt.Println("odd -> ", i)
		*waitingState = 0
		c.Signal()
		c.L.Unlock()

	}
}
