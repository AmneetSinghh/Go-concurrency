package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	op := 2
	switch op {
	case 1:
		practice()
	case 2:
		contextPractice()
	default:
	}
}

func practice() {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Simulate work
	go func() {
		time.Sleep(1 * time.Second) // Simulate a delay longer than the timeout
		fmt.Println("Goroutine finished work")
	}()

	// Wait for the context to be done
	fmt.Println("waiting")
	<-ctx.Done() // waits till context cancelled.
	fmt.Println("ends")
	// Check why the context was done
	if err := ctx.Err(); err != nil {
		switch err {
		case context.Canceled:
			fmt.Println("Context was canceled")
		case context.DeadlineExceeded:
			fmt.Println("Context deadline was exceeded")
		}
	} else {
		fmt.Println("fuck")
	}
}

func contextPractice() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // always cancel when function ends

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(ctx); err != nil {
			fmt.Printf("cannot print greeetings: %v\n", err)
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if err := printFarewell(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
		}
	}()

	wg.Wait()
}

func printGreeting(ctx context.Context) error {
	greeting, err := genGreetings(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell(ctx context.Context) error {
	farewell, err := genFarewell(ctx)
	fmt.Printf("%s here!\n", farewell)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func genGreetings(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(ctx context.Context) (string, error) {
	fmt.Printf("%s genFarewell!\n")

	switch locale, err := locale(ctx); {
	case err != nil:
		fmt.Println("farewell bro")
		return "", err
	case locale == "EN/US":
		fmt.Println("farewell printed this")
		return "Goodbye", nil
	}
	fmt.Printf("%s here farewell!\n")
	return "", fmt.Errorf("unsupported locale")
}

func locale(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done(): // when context cancelled, this done channel is closed.
		fmt.Println("this runs, parent context deadline: 1s")
		return "", ctx.Err() // it can only be off two types.
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil
}
