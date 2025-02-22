package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	op := 4
	switch op {
	case 1:
		heartImplementation()
	case 2:
		heartbeatInBeginOfUnitWork()
	case 3:
		heartBeatToCheckGoroutineHasStartedDoingWork()
	case 4:
		replicatedRequests()
	}
}

/*
------------------------------ HEARTBEAT USED FOR LONG RUNNING GOROUTINES ------------------------------------

		parent goroutine
	/ 					\

500 ms	heartbeat 		result ( 2 seconds)

	\					/
		doWork
*/
func heartImplementation() {
	doWork := func(
		done <-chan interface{},
		pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
		heartBeat := make(chan interface{})
		results := make(chan time.Time)

		go func() {
			defer close(heartBeat)
			defer close(results)

			pulse := time.Tick(pulseInterval)       // if 500 miliseconds
			workgen := time.Tick(2 * pulseInterval) // 1 seconds.
			sendPulse := func() {
				select {
				case heartBeat <- struct{}{}:
				default: // because if heartBeat is not ready to receive, then we don't want to block
				}
			}

			sendResults := func(r time.Time) {
				for {
					// done and pulse can be checked in next iteration as well.
					select {
					case <-pulse:
						sendPulse()
					case <-done:
						return
					case results <- r:
						return
					default: // because if heartBeat is not ready to receive, then we don't want to block
					}
				}

			}

			// main code.
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workgen:
					sendResults(r)
				}
			}
		}()
		return heartBeat, results
	}

	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartBeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartBeat:
			if !ok {
				fmt.Println("heartBeat closed")
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				fmt.Println("result closed")
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout): // after 2 seconds. we get nothing from heartBeat and results. this determines in single select statement how much time we should wait.
			fmt.Println("timeout bcz heart or results is not responded")
			return
		}
	}
}

/*
heartBeats :

- Beautiful! Within two seconds our system realizes something is amiss with our goroutine and breaks the for-select loop. By using a heartbeat,
  we have successfully avoided a deadlock, and we remain deterministic by not having to rely on a longer timeout ( 10 second )
- Also note that heartbeats help with the opposite case: they let us know that long- running goroutines remain up,
  but are just taking a while to produce a value to send on the values channel.

  - case <-time.After(timeout):  ADDING this is necessary
*/

func heartbeatInBeginOfUnitWork() {
	doWork := func(
		done <-chan interface{}) (<-chan interface{}, <-chan int) {
		heartBeat := make(chan interface{}, 1) // even if no one listening, at least we will send 1 heartbeat.
		results := make(chan int)

		go func() {
			defer close(heartBeat)
			defer close(results)
			// main code.
			for i := 1; i <= 10; i++ {

				// heartbeat at the begin of unit of work
				select {
				case heartBeat <- struct{}{}:
				default: // because if heartBeat is not ready to receive, then we don't want to block
				}

				select {
				case <-done:
					return
				case results <- i:
				}
			}
		}()
		return heartBeat, results
	}

	done := make(chan interface{})
	defer close(done)
	heartbeat, results := doWork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}

func heartBeatToCheckGoroutineHasStartedDoingWork() {

	doWork := func(
		done <-chan interface{}) (<-chan interface{}, <-chan int) {
		heartBeat := make(chan interface{}, 1) // even if no one listening, at least we will send 1 heartbeat.
		results := make(chan int)

		go func() {
			defer close(heartBeat)
			defer close(results)
			// main code.
			time.Sleep(3 * time.Second)
			for i := 1; i <= 10; i++ {

				// heartbeat at the begin of unit of work
				select {
				case heartBeat <- struct{}{}:
				default: // because if heartBeat is not ready to receive, then we don't want to block
				}

				select {
				case <-done:
					return
				case results <- i:
				}
			}
		}()
		return heartBeat, results
	}

	done := make(chan interface{})
	defer close(done)
	heartbeat, results := doWork(done)
	fmt.Println("heart beat waiting")
	<-heartbeat //	waiting for first heartbeat.
	fmt.Println("heart beat runs")
	for r := range results {
		fmt.Println("result -=>", r)
	}

	//Because of the heartbeat, we can safely write our test without timeouts. if heartbeat is not implemented in this case, we need to write timeout in test. that timeout will be un-deterministic.
	// We don't know correct value. thats why heartbeat is important in this case.

}

func replicatedRequests() {
	doWork := func(
		done <-chan interface{}, id int,
		wg *sync.WaitGroup, result chan<- int,
	) {
		started := time.Now()
		defer wg.Done()
		// Simulate random load
		simulatedLoadTime := time.Duration(1+rand.Intn(5)) * time.Second
		if id >= 8 {
			simulatedLoadTime = 10 * time.Second
		}
		select {
		case <-done:
		case <-time.After(simulatedLoadTime):
		}
		select {
		case <-done:
		case result <- id:
		}
		took := time.Since(started)
		if took < simulatedLoadTime {
			//	fmt.Printf("%v less than simulated time %v\n", id)
			took = simulatedLoadTime
		}
		fmt.Printf("%v took %v\n", id, took)
	}
	done := make(chan interface{})
	result := make(chan int)
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go doWork(done, i, &wg, result)
	}
	firstReturned := <-result
	fmt.Printf("fast result#%v\n", firstReturned)
	close(done)
	fmt.Printf("done closed#%v\n", firstReturned)

	wg.Wait() // all wil be finished, as done is triggered,just a few nano-seconds away.

	fmt.Printf("Received an answer from #%v\n", firstReturned)
}
