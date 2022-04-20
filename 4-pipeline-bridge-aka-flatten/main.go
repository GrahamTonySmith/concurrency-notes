package main

import "fmt"

type any = interface{}

// bridge a stream of streams in to a stream
func bridge(done <-chan any, chanStream <-chan <-chan any) <-chan any {
	valueStream := make(chan any)
	go func() {
		defer close(valueStream)
		for {
			var stream <-chan any // a nil channel
			select {              // blocks until a case can be executed
			case <-done:
				return // and close
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			}
			for val := range stream {
				select {
				case <-done:
				case valueStream <- val: // send the value to the valueStream
				}
			}
		}
	}()
	return valueStream
}

func genVals() <-chan <-chan any {
	chanStream := make(chan (<-chan any))
	go func() {
		defer close(chanStream)
		for i := 0; i < 10; i++ {
			stream := make(chan any, 1)
			stream <- i
			close(stream)
			chanStream <- stream
		}
	}()
	return chanStream
}

func main() {
	done := make(chan any)
	defer close(done)

	for v := range bridge(done, genVals()) {
		fmt.Printf("%v ", v)
	}
	fmt.Printf("\n")
}
