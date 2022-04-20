package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	sayHello := func() {
		defer wg.Done()
		fmt.Println("Hello, World!")
	}
	wg.Add(1)
	go sayHello() // fork point
	wg.Wait()     // join point
}
