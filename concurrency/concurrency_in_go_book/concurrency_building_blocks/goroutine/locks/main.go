package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

/*
- Rlocks :
- multiple reads can aquire locks at same time, write will wait for aquiring lock while all readers finishes.
- use only in read heavy scenarios
Write Starvation:
* If readers (RLock) are continuous and frequent, the writer (Lock) may wait indefinitely because sync.RWMutex prioritizes readers.
* If writes are frequent, the benefits of using RLock diminish because every Lock operation will block all readers, resulting in reduced concurrency.
*/

func producer(wg *sync.WaitGroup, l sync.Locker) {
	defer wg.Done()
	for i := 5; i > 0; i-- {
		l.Lock()
		l.Unlock()
		time.Sleep(1)
	}
}

func observer(wg *sync.WaitGroup, l sync.Locker) {
	defer wg.Done() // runs second
	l.Lock()
	defer l.Unlock() // runs first.
}

func test(count int, mutex, rwMutex sync.Locker) time.Duration {
	var wg sync.WaitGroup
	wg.Add(count + 1)
	beginTestTime := time.Now()
	go producer(&wg, mutex) // write lock
	for i := count; i > 0; i-- {
		go observer(&wg, rwMutex)
	}
	wg.Wait()
	return time.Since(beginTestTime)
}

func main() {
	var m sync.RWMutex
	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()
	fmt.Fprintf(tw, "Readers\tRWMutext\tMutex\n")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &m, m.RLocker()),
			test(count, &m, &m),
		)
	}
}

