package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	op := 5
	switch op {
	case 1:
		batchProcessing()
	case 2:
		streamProcessing()
	case 3:
		pipelineChannels()
	case 4:
		fanoutFanIn()

	case 5:
		queuePipeline()
	default:
	}
}

/*
* Batch processing
  - Original data to remain unaltered, each stage has to make a new slice of equal length to store the
    results of its calcula‐ tions. That means that the memory footprint of our program at any one time 2*size of the slice
    we send into the start of our pipeline.
  - lets figure out later.
  -

* Stream processing
  - Each stage is receiving and emitting a discrete value, and the memory footprint of our program is back down to only the size of the pipeline’s input
  - Limit to scale, when each stage will run concurrenctly. lets figure out later.
  -
*/

func batchProcessing() {

	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}
		return multipliedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}

	ints := []int{1, 2, 3, 4}
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}

}

func streamProcessing() {

	multiply := func(value, multiplier int) int {
		return value * multiplier
	}
	add := func(value, additive int) int {
		return value + additive
	}
	ints := []int{1, 2, 3, 4}
	for _, v := range ints {
		fmt.Println(add(multiply(v, 2), 1))
	}

}

// generator -> add -> multiplly, channels.
// val -> generator -> add -> multiply -> printed in main goroutine ( stream of intergers flowing into channels from one to another same as stream processing)
// separation of concerns

func pipelineChannels() {
	generator := func(done <-chan int) <-chan int {
		output := make(chan int)
		go func() {
			defer close(output)
			for val := range 100 {
				select {
				case <-done:
					fmt.Println("generator terminated")
					return
				case output <- val:
				}
			}
		}()
		return output
	}

	// intput:: generator channel, output : transformation
	add := func(done <-chan int, generatorOutput <-chan int, add int) <-chan int {
		addOutput := make(chan int)
		go func() {
			defer close(addOutput)
			for val := range generatorOutput {
				select {
				case <-done:
					fmt.Println("add terminated")
					return
				case addOutput <- val + add:
				}
			}
		}()
		return addOutput
	}

	// intput:: add channel, output : transformation
	multiply := func(done <-chan int, addOutput <-chan int, multiply int) <-chan int {
		multiplyOutput := make(chan int)
		go func() {
			defer close(multiplyOutput)
			for val := range addOutput {
				select {
				case <-done:
					fmt.Println("multiply terminated")
					return
				case multiplyOutput <- val + multiply:
				}
			}
		}()
		return multiplyOutput
	}

	done := make(chan int)
	// done used for stopping each channel, so done is common to all channels.
	// each pipeline is independent from each other

	generatorChan := generator(done)

	pipeline := add(done, multiply(done, generatorChan, 2), 1)

	go func() {
		time.Sleep(100 * time.Microsecond)
		close(done)
		time.Sleep(100 * time.Microsecond)
	}()

	for v := range pipeline {
		fmt.Println(v)
	}

	fmt.Println("After done")
	///	time.Sleep(100 * time.Millisecond)
	fmt.Println("waiting completed")
}

/*
Time complexity of this : completion of all pipeline stages

issue :

If one of your stages is computationally expensive, this will certainly eclipse this per‐ formance overhead.
Speaking of one stage being computationally expensive, how can we help mitigate this?
Won’t it rate-limit the entire pipeline?

*/

func fanoutFanIn() {

	rand := func() interface{} { return rand.Intn(50000000) }
	done := make(chan int)
	defer close(done)
	fmt.Println(rand())

	toInt := func(done <-chan int, valueStream <-chan interface{}) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for v := range valueStream {
				select {
				case <-done:
					fmt.Println("\ntoInt finished")
					return
				case intStream <- v.(int):
				}
			}
		}()
		return intStream
	}

	repeatFn := func(done <-chan int, fn func() interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					fmt.Println("\nrepeat fn finished")
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}

	// take first n elements.
	take := func(done <-chan int, valueStream <-chan int, num int) <-chan int {
		takeStream := make(chan int)
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					fmt.Println("\ntake finished")
					return
				case value := <-valueStream:
					takeStream <- value
				}
			}
		}()
		return takeStream
	}

	/*
	 * primeFinder is taking time, so make single prime calculation in multiple goroutines.
	 */
	primeFinder := func(done <-chan int, randStream <-chan int) <-chan int {
		primeStream := make(chan int)
		go func() {
			defer close(primeStream)
			for num := range randStream {
				select {
				case <-done:
					fmt.Println("\nprime finder finished")
					return
				default:
					if isPrime(num) { // isPrime is a helper function to check primality
						primeStream <- num
					}
				}
			}
		}()
		return primeStream
	}

	start := time.Now()
	randomIntStream := toInt(done, repeatFn(done, rand))

	primeChannels := make([]<-chan int, 10)
	for i := 0; i < 10; i++ {
		primeChannels[i] = primeFinder(done, randomIntStream) // 10 channels -> consumers.
	}

	fmt.Println("Primes:")
	for prime := range take(done, fanIn(done, primeChannels), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	fmt.Printf("Search took: %v", time.Since(start))

}

// Helper function to check if a number is prime
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	time.Sleep(500 * time.Millisecond)
	return true
}

func fanIn(done <-chan int, chanList []<-chan int) <-chan int {
	var wg sync.WaitGroup
	multiplexStream := make(chan int)
	multiplex := func(currentChannel <-chan int) {
		go func() {
			defer wg.Done()
			for i := range currentChannel {
				select {
				case <-done:
					return
				case multiplexStream <- i:
				}
			}
		}()
	}

	wg.Add(len(chanList))
	for _, ch := range chanList {
		go multiplex(ch) // 10 consumers each run on separate goroutine.
	}

	/*
		we are doing this in goroutine, because we are writing to multiplexStream, so if no one reads, then deadlock will occurs.
	*/
	go func() {
		wg.Wait()
		close(multiplexStream)
	}()

	return multiplexStream
}

func queuePipeline() {
	/*
			// with the fanIn and fanOut we already have optimized the time complexity, but what about this below example:

			---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

			// done := make(chan interface{})
			// defer close(done)
			// zeros := take(done, 3, repeat(done, 0))
			// short := sleep(done, 1*time.Second, zeros)
			// long := sleep(done, 4*time.Second, shott1)
			// pipeline := long


			// total time for running these will be 13 seconds, we are wasting cpu cycles because short will wait to send event to long, as reading from long will take 4 seconds and short just completed in 1 seconds.
			// so we wil introduce buffer queue into short, so that short can complete its stage efficiently, resource utilization improved. CPU can do other tasks instead of blocking short

			---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

			// 	done := make(chan interface{}) defer close(done)
			// zeros := take(done, 3, repeat(done, 0))
			// short := sleep(done, 1*time.Second, zeros)
			// buffer := buffer(done, 2, short) // Buffers sends from short by 2 long := sleep(done, 4*time.Second, short)
			// pipeline := long


			// time complexity will be the same, but just the resource utilization improved.

			---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

			// 	A pipeline processes incoming tasks, but if one part of the pipeline slows down:

			// The pipeline cannot process tasks quickly enough.
			// Upstream systems (the systems sending tasks into the pipeline) may send even more tasks to "catch up."
			// This increases the load on the already slow pipeline, making it even slower.
			// The cycle repeats, and the pipeline eventually collapses because it cannot handle the increasing load.


			   * rate limiting,
			   * queueing, kafka,
			   * tcp queue -> it also store requests, if I am unable to process or slow to process.

			---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		    So from our examples we can begin to see a pattern emerge; queuing should be implemented either:


			• At the entrance to your pipeline.
			• In stages where batching will lead to higher efficiency.



	*/
}
