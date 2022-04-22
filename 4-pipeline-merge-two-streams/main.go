package main

import (
	"fmt"
	"time"
)

type any = interface{}

func merge(done <-chan any, inOne <-chan int, inTwo <-chan int) <-chan int {
	valueStream := make(chan int)
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				fmt.Println("got a done... returning")
				return
			case value, ok := <-inOne:
				if !ok {
					fmt.Println("inOne is closed")
					inOne = nil
				} else {
					valueStream <- value
				}
			case value, ok := <-inTwo:
				if !ok {
					fmt.Println("inTwo is closed")
					inTwo = nil
				} else {
					valueStream <- value
				}
			default:
				if inOne == nil && inTwo == nil {
					fmt.Println("inOne and inTwo are closed... returning")
					return
				}
			}
		}
	}()
	return valueStream
}

func nValue(value int, n int) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for i := 0; i < n; i++ {
			intStream <- value
			time.Sleep(1 * time.Second)
		}
	}()
	return intStream
}

func main() {
	done := make(chan any)
	time.AfterFunc(15*time.Second, func() { close(done) })

	ones := nValue(1, 5)
	twos := nValue(2, 10)

	for elem := range merge(done, ones, twos) {
		fmt.Println(elem)
	}
}
