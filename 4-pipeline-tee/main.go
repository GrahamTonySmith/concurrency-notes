package main

import (
	"fmt"
	"time"
)

type any = interface{}

func repeat(done <-chan any, values ...any) <-chan any {
	valueStream := make(chan any)
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

func take(done <-chan any, valueStream <-chan any, num int) <-chan any {
	takeStream := make(chan any)
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				fmt.Println("take got done")
				return
			case takeStream <- <-valueStream: // send the value reseived from valueStream to takeStream
			}
		}
	}()
	return takeStream
}

// we utalised the fact that sending and receiving to and from a nil channel will block
func tee(done <-chan any, in <-chan any) (<-chan any, <-chan any) {
	outOne := make(chan any)
	outTwo := make(chan any)
	go func() {
		defer close(outOne)
		defer close(outTwo)
		for val := range in {
			var outOne, outTwo = outOne, outTwo // shadow the out channels
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case outOne <- val: // after this has been selected we set shadow of outOne to nil so it blocks forever
					outOne = nil
				case outTwo <- val: // after this has been selected we set shadow of outTwo to nil so it blocks forever
					outTwo = nil
				}
			}
		}
	}()
	return outOne, outTwo
}

func main() {
	done := make(chan any)
	defer close(done)

	outOne, outTwo := tee(done, take(done, repeat(done, 1, 2, 3), 5))

	for valOne := range outOne {
		fmt.Printf("outOne: %v, outTwo: %v\n", valOne, <-outTwo)
	}
}
