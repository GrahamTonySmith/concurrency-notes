package main

import (
	"fmt"
	"math/rand"
)

type any = interface{}

func doWork(done <-chan any) (<-chan any, <-chan int) {
	heartbeat := make(chan any)
	workStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(workStream)

		for i := 0; i < 10; i++ {
			select {
			case heartbeat <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case workStream <- rand.Intn(10):
			}
		}
	}()
	return heartbeat, workStream
}

func main() {
	done := make(chan any)
	defer close(done)

	heartbeat, results := doWork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Printf("pulse\n")
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}
