package main

import (
	"fmt"
	"sync"
	"time"
)

type any = interface{}

func repeat(done <-chan any, values ...any) <-chan any {
	valueStream := make(chan any)
	go func() {
		defer close(valueStream) // run after return hits on the close of done
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

func fanIn(done <-chan any, channels ...<-chan any) <-chan any {
	var wg sync.WaitGroup
	multiplexedStream := make(chan any)

	multiplex := func(c <-chan any) {
		defer wg.Done() // run after return hits on the close of done
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i: // send i to the multiplexedStream
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	// wait on the wait group and close the multiplexed stream
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func main() {
	done := make(chan any)
	go func() {
		time.Sleep(10 * time.Second)
		close(done)
	}()

	aStream := repeat(done, "a")
	bStream := repeat(done, "b")
	cStream := repeat(done, "c")

	charStream := fanIn(done, aStream, bStream, cStream)

	go func() {
		for char := range charStream {
			fmt.Printf("%v ", char)
		}
	}()

	<-done
	fmt.Printf("\n")
}
