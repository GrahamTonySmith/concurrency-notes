package main

/**
	1. Give the pool a New variable that is thread safe when called
	2. Make no assumptions about the state of the object you get back
	3. Make sure to call put when you're finished otherwise pool is usless
	4. Objects in the pool must be roughly uniform in makeup
**/
import (
	"fmt"
	"sync"
)

func main() {
	var numCalcsCreated int
	calcPool := sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}

	// warm the pool
	for i := 0; i > 5; i++ {
		calcPool.Put(calcPool.New())
	}

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)
		}()
	}
	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}
