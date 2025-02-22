package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var dataStream chan interface{}

func main() {

	var op = 5

	switch op {
	case 1:
		practiceChannel()
	case 2:
		nilChannels()
	case 3:
		closeChannel()
	case 18:
		rangeWillUnblockIfChannelClose()
	case 4:
		bestPracticeForWriteChannels()
	case 5:
		badDesign()
	case 6:
		goodDesign()
	case 7:
		printEvenOdd()
	default:
		fmt.Println("operation not defined")
	}
}

/*
Instantiate the channel.
2. Perform writes, or pass ownership to another goroutine.
3. Close the channel.
4. Ecapsulate the previous three things in this list and expose them via a reader channel.

notes:
- making channel outside caller function, can be dangerous, as he can also write into channel.

best thing:
* lifecycle of the resultStream channel is encapsulated within the chan Owner function.
* It’s very clear that the writes will not happen on a nil or closed chan‐ nel, and that the close will always happen once
*/

func bestPracticeForWriteChannels() {
	channelOwner := func() <-chan int {
		ch := make(chan int, 5)
		go func() {
			defer close(ch)
			for i := 1; i <= 5; i++ {
				ch <- i
			}
		}()
		return ch
	}

	resultStream := channelOwner()
	for value := range resultStream {
		fmt.Println("channel output :-> ", value)
	}
}

/*
 * Reading from close channel is allowed, but writing is not allowed, it will PANIC
 */
func closeChannel() {

	begin := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(5)

	for i := 1; i <= 5; i++ {
		go func(j int) {
			defer wg.Done()
			<-begin // reading
			fmt.Printf("%v has begun\n", j)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(begin) // close will un-block all channels which are waiting to receive a event that is zero.
	// without close they will wait forever deadock will occur.
	wg.Wait()
}

func rangeWillUnblockIfChannelClose() {
	stream := make(chan int)

	go func() {
		defer close(stream)
		stream <- 1
		stream <- 2
	}()

	for val := range stream {
		fmt.Println(val)
	}
}

func nilChannels() {
	var dataStream chan interface{}
	//<-dataStream
	dataStream <- "write in nil channel"
	//close(dataStream)

	/*
	 all 3 lines will make panic
	*/
}

func printEvenOdd() {
	var wg sync.WaitGroup
	wg.Add(2)
	var state = 0
	signal := make(chan interface{})
	even := func() {
		defer wg.Done()
		for i := 0; i <= 10; i += 2 {
			if state == 1 {
				<-signal // blocking state.
			}
			fmt.Println(i)
			state = 1
			if i == 10 {
				close(signal)
			} else {
				signal <- struct{}{}
			}

		}
	}

	odd := func() {
		defer wg.Done()
		for i := 1; i <= 10; i += 2 {
			if state == 0 {
				<-signal
			}
			fmt.Println(i)
			state = 0
			signal <- struct{}{}
		}
	}

	go even()
	go odd()
	wg.Wait()
}

func practiceChannel() {
	dataStream := make(chan interface{}, 10)

	///dataStream <- "amneet singh powerful "
	run := func() {
		defer close(dataStream)
		for i := 1; i <= 100; i++ {
			fmt.Println("here")
			dataStream <- "amneet singh powerful " + strconv.Itoa(i)
			fmt.Println("here1")
		}
	}

	go run()

	for value := range dataStream {
		fmt.Println(value, " ->------")
	}
	close(dataStream) // Close the channel after sending all data
	dataStream <- "amneet"
	fmt.Println(<-dataStream)
	fmt.Println(<-dataStream)
	fmt.Println(<-dataStream)
	time.Sleep(3 * time.Second)
}

func badDesign() {
	data := make([]int, 4)
	handleData := func(hData chan<- int) {
		defer close(hData)
		for value := range data {
			hData <- value
		}
	}
	hData := make(chan int)
	go handleData(hData)

	for value := range hData {
		fmt.Println(value)
	}

	/*
	 * channel is created outside, so caller can write into the channel, goroutine don't have proper ownership of the channel
	 * data is available for both caller and goroutine function.
	 */
}

func goodDesign() {
	// so owner has the ownership of creating channel, and can pass readonly write only or normal channel to us.
	channelOwner := func() <-chan int {
		results := make(chan int, 4)
		go func() {
			defer close(results)
			for i := 0; i < 5; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Println(result)
		}
	}

	owner := channelOwner()
	consumer(owner)
}
