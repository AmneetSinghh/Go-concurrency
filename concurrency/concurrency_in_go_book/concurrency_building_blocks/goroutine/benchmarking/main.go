package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {

	var c <-chan interface{}
	var wg sync.WaitGroup
	noob := func() {
		wg.Done()
		<-c // blocking operation
	}

	const numOfGoroutines = 1e4
	wg.Add(numOfGoroutines)
	before := memConsumed()

	for i := numOfGoroutines; i > 0; i-- {
		go noob()
	}
	wg.Wait()
	after := memConsumed()

	fmt.Printf("%.3f kb", float64(after-before)/numOfGoroutines/1000)

}

func memConsumed() uint64 {

	runtime.GC()
	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	return s.Sys
}

/*
- time taken for channel passing very less ( nano-seconds)
- do make proper notes.
- context switching in goroutines means message passing in a channel, as channels are used widely in golang for concurrency.


types of goroutines:

Empty goroutines: Use minimal memory and almost no CPU.
Busy goroutines: Use more memory (stack + heap) and compete for CPU time.

notes :
- context switches
- number of goroutines
- closures

*/
