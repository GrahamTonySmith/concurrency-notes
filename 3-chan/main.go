package main

import (
	"fmt"
	"time"
)

func main() {
	stringStream := make(chan string)
	go func() {
		time.Sleep(5 * time.Second)
		stringStream <- "Hello, Channels!"
	}()
	salutation, ok := <-stringStream         // if the chan is closed ok is false
	fmt.Printf("(%v): %v\n", ok, salutation) // <-stringStream blocks until it gets a value
}
