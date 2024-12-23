package main

import (
	"fmt"
	"time"
)

func worker(stopChan <-chan struct{}) {
	for {
		select {
		case <-stopChan:
			fmt.Println("Worker stopping")
			return
		default:
			// Do some work
			fmt.Println("Working...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	stopChan := make(chan struct{})

	go worker(stopChan)

	time.Sleep(2 * time.Second)

	// Signal the worker to stop
	close(stopChan)

	time.Sleep(1 * time.Second)
	fmt.Println("Main function finished")
}
