package main

import (
	"fmt"
	"time"
)

func main() {
	op := 4
	switch op {
	case 1:
		practice()
	case 2:
		timeoutSelect()
	case 3:
		defaultSelect()
	case 4:
		forSelectLoop()
	default:
		fmt.Println("op is not defined")
	}
}

func practice() {
	first := make(chan interface{})
	close(first)
	second := make(chan interface{})
	close(second)

	var count1, count2 int
	for i := 1; i <= 500; i++ {
		select {
		case <-first:
			count1++
		case <-second:
			count2++
		}
	}
	// select is blocking, if nothing is active to read from any channel.

	fmt.Printf("c1Count: %v\nc2Count: %v\n", count1, count2)

}

func timeoutSelect() {
	/*
	 * if channel is nil, then it wil not execute and not make panic
	 */

	var c <-chan int
	fmt.Println("entering")
	select {
	case <-c:
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}

	fmt.Println("exit")
}

func defaultSelect() {
	start := time.Now()
	var c1, c2 <-chan int
	fmt.Printf("entering\n")
	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after %v\n", time.Since(start))
	}
	fmt.Printf("exit\n")
}

/*
 * 10-20k iterations around 1 nanosecond
 */
func forSelectLoop() {
	done := make(chan int)
	go func() {
		time.Sleep(1 * time.Nanosecond)
		done <- 1
	}()

	var cycle int

LOOP:
	for {
		select {
		case <-done:
			break LOOP // break always break nearest enclosing function, in our case this is select.
		default:
		}
		// simulating work
		cycle++
	}

	fmt.Println(cycle)
}
