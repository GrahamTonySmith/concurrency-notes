package main

import (
	"fmt"
	"time"
)

type any = interface{}

func send(c chan int) {
	fmt.Println("sending an int to channel c started")
	c <- 1
	fmt.Println("sending an int to channel c finished") // we will never hit this
}

func main() {
	var c chan int
	go send(c)
	time.Sleep(1 * time.Second)
	<-c
}
