package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan interface{})

	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
	for {
		select {
		case <-done:
			fmt.Printf("Did %d cycles for work before done was called\n", workCounter)
			return
		default:
			workCounter++
			time.Sleep(1 * time.Second)
		}
	}
}
