package main

import (
	"fmt"
	"sync"
)

func main() {
	var i int
	var memoryAcceess sync.Mutex
	go func() {
		memoryAcceess.Lock()
		i++
		memoryAcceess.Unlock()
	}()

	memoryAcceess.Lock() // if we not write lock here,  '1 first' is possible.
	if i == 0 {
		fmt.Printf("%d first\n", i)
	} else {
		fmt.Printf("%d second \n", i)
	}
	memoryAcceess.Unlock()
}

/*
critical section : section which has exclusive access to shared resource.

- i++,
- if condition, checking the variable
- fmt.printf whihc is printing the variable


- data race problem is solved, but still we don't know whihc statement will run first, go func() or if condition.
- The order of operations still non-deterministic.


Lock :

Mutex Lock Behavior:

The mutex guarantees mutual exclusion. Only one goroutine can hold the lock at a time. This means that:
When the main goroutine holds the lock, the i++ operation in the other goroutine cannot proceed.
Similarly, when the other goroutine holds the lock, the main goroutine will wait.
*/
