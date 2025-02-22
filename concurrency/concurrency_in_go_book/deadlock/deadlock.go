package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	type value struct {
		val int
		mu  sync.Mutex
	}

	var wg *sync.WaitGroup
	wg = &sync.WaitGroup{}

	printSum := func(v1, v2 *value) {
		//defer wg.Done()
		v1.mu.Lock()
		defer v1.mu.Unlock()

		time.Sleep(1 * time.Second)
		v2.mu.Lock()
		defer v2.mu.Unlock()
		fmt.Printf(" value1 %+v, value2 %+v\n", v1.val, v2.val)
	}

	var a, b value
	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait() // need to wait for all goroutines to run, so added waitgroup
}

/*
- All 4 conditions met for a deadlock

- mutual exclusion
- wait for conditio
- no preemption
- circular wait

removing one of these, wil remove deadlock.
*/
