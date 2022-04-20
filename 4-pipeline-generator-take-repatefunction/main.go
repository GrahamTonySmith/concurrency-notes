package main

import (
	"fmt"
	"math/rand"
)

type any = interface{}

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

func repeatFn(done <-chan any, fn func() any) <-chan any {
	valueStream := make(chan any)
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn(): // send the return val of fn into the valueStream
			}
		}
	}()
	return valueStream
}

func main() {
	done := make(chan any)
	defer close(done)

	rand := func() any { return rand.Int() }

	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}
}
