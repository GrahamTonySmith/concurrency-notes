package main

import "fmt"

// convert discrete values into a stream (chan)
func generator(done <-chan interface{}, intergers ...int) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for _, i := range intergers {
			select {
			// case must send or receive
			case <-done:
				fmt.Printf("generator got done and returning\n")
				return
			case intStream <- i: // put into stream
			}
		}
	}()
	return intStream // closed by defer
}

func multiply(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
	multipliedStream := make(chan int)
	go func() {
		defer close(multipliedStream)
		for i := range intStream {
			select {
			case <-done:
				fmt.Printf("multiplier got done and returning\n")
				return
			case multipliedStream <- i * multiplier: // put into stream
			}
		}
	}()
	return multipliedStream // closed
}

func add(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
	addStream := make(chan int)
	go func() {
		defer close(addStream)
		for i := range intStream {
			select {
			case <-done:
				fmt.Printf("add got done and returning\n")
				return
			case addStream <- i + additive: // put into stream
			}
		}
	}()
	return addStream
}

func main() {
	done := make(chan interface{})
	intStream := generator(done, 1, 2, 3, 4)

	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}
