package main

import (
	"fmt"
	"net/http"
	"time"
)

/*
Confinement is the simple yet powerful idea of ensuring information is only ever available
from one concurrent process. When this is achieved, a concurrent program is implicitly safe and
 no synchronization is needed. There are two kinds of confine‐ ment possible: ad hoc and lexical.
*/

func main() {
	op := 3
	switch op {
	case 1:
		preventGoroutineLeak() // goroutine response for creating it and stoping it
	case 2:
		preventGoroutineLeak2() // goroutine response for creating it and stoping it
	case 3:
		errorhandlingPattern() // how parent main function, iteract with goroutine, to check errors.
	}
}

func preventGoroutineLeak() {
	doWork := func(strings <-chan string, done <-chan string) <-chan interface{} {
		completed := make(chan interface{})

		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for {
				select {
				case <-strings: // blocking
					fmt.Println("read successful")
				case <-done:
					fmt.Println("Channel done runs")
					return
				default:

				}
			}
		}()
		return completed
	}

	terminate := make(chan string)
	signal := make(chan string)
	postTermination := doWork(signal, terminate) // final channel for waiting at the end.

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("cancelling goroutine because of timeout")
		terminate <- ""
	}()

	<-postTermination // used final here for waiting purpose.
	///time.Sleep(1 * time.Second)
	fmt.Println("Done.")

}

func preventGoroutineLeak2() {

	newRandStream := func(termination <-chan int) (<-chan int, <-chan int) {
		randStream := make(chan int)
		done := make(chan int)
		go func() {
			defer fmt.Println("newRandStream exited.")
			defer close(randStream)
			defer close(done)
			var count int
			for {
				select {
				case randStream <- count: // writing into channel
				case <-termination:
					fmt.Println("termimation runs")
					return
				}
				count++
			}
		}()
		return randStream, done
	}

	termination := make(chan int)
	randStream, done := newRandStream(termination) // read only channel.
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(termination) // termination will run in above goroutine.
	fmt.Println("finished")
	<-done
}

func errorhandlingPattern() {

	type Result struct {
		Error    error
		Response *http.Response
	}

	dowork := func(urls []string) <-chan Result {
		result := make(chan Result)
		go func() {
			defer close(result)
			for _, url := range urls {
				resp, err := http.Get(url)
				res := Result{
					Error:    err,
					Response: resp,
				}
				result <- res
			}
		}()
		return result
	}

	urls := []string{"https://www.google.com", "https://badhost"}
	//termination := make(chan string)
	//defer close(termination)
	resultChannel := dowork(urls)

	for result := range resultChannel { // it will block reading channnel, no need for DONE channel
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}

	//<-resultChannel
	//time.Sleep(1 * time.Second)
}
