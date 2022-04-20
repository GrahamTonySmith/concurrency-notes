package main

import (
	"fmt"
	"sync"
)

// use close(chan) to signal to n goroutines to start
func main() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin // block until begin is closed
			fmt.Printf("%v has begun\n", i)
		}(i)
	}
	fmt.Println("Unblocking goroutines...")
	close(begin) // which will send the non ok value val, ok <- readChan
	wg.Wait()
}
