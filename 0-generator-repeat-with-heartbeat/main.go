package main

import (
	"fmt"
	"time"
)

type any = interface{}

func repeatWithHeartbeat(done <-chan any, values ...any) (<-chan any, <-chan any) {
	heartbeat := make(chan any)
	valStream := make(chan any)
	go func() {
		defer close(heartbeat)
		defer close(valStream)
		for {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
			for _, v := range values {
				select {
				case <-done:
					fmt.Println("generator hit done")
					return
				case valStream <- v: // send v to the valueStream
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()
	return heartbeat, valStream
}

func main() {
	done := make(chan any)
	time.AfterFunc(20*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := repeatWithHeartbeat(done, 1, 2, 3)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Println(r)
		case <-time.After(timeout):
			fmt.Println("goroutine is unhealthy!")
			return
		}
	}
}
