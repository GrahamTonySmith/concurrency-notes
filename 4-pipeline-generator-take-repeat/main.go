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

func main() {
	done := make(chan any)
	defer close(done)

	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}
	fmt.Printf("\n")
}
