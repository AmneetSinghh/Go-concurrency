package main

import (
	"fmt"
	"sync"
	"time"
)

var sharedLock sync.Mutex
var wg sync.WaitGroup

const runtime = 1 * time.Second

func main() {
	greedyWorker := func() {
		var count int
		defer wg.Done()
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(3 * time.Nanosecond)
			sharedLock.Unlock()
			count++
		}

		fmt.Printf("greedy worker work done is : %+v\n", count)
	}

	politeWorker := func() {
		var count int
		defer wg.Done()
		for begin := time.Now(); time.Since(begin) <= runtime; {
			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			sharedLock.Lock()
			time.Sleep(1 * time.Nanosecond)
			sharedLock.Unlock()

			count++
		}
		fmt.Printf("polite worker work done is : %+v\n", count)
	}

	wg.Add(2)
	go greedyWorker()
	go politeWorker()
	wg.Wait()

}

/*

The greedy worker greedily holds onto the shared lock for the entirety of its work loop,
whereas the polite worker attempts to only lock when it needs to. Both workers do the same amount of
simulated work (sleeping for three nanoseconds), but as you can see in the same amount of time, the greedy worker got almost twice the amount of work done!


If we assume both workers have the same-sized critical section,
rather than conclud‐ ing that the greedy worker’s algorithm is more efficient (or that the calls
to Lock and Unlock are slow—they aren’t), we instead conclude that the greedy worker has unnec‐ essarily expanded its hold on the shared lock beyond
its critical section and is pre‐ venting (via starvation) the polite worker’s goroutine from performing work efficiently.



make notes after getting ipad.

- if critical section is big, we can run the problem of starvation as you can see above
- start the problem, with smaller critical sections. If sync ( mutex cost) time is increasing, then do balance broad critical sections.

*/
