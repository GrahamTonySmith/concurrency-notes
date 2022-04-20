package main

import (
	"fmt"
	"time"
)

func repeateGenerator(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream) // run after return hits on the close
		for {
			for _, v := range values {
				select {
				case <-done:
					fmt.Println("generator hit done")
					return
				case valueStream <- v: // send v to the valueStream
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()
	return valueStream
}

func main() {
	done := make(chan interface{})
	intGenerator := repeateGenerator(done, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9)

	go func() {
		time.Sleep(20 * time.Second)
		close(done)
	}()

	for i := range intGenerator {
		fmt.Println(i)
	}
}
