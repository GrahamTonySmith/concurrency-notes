package main

import (
	"fmt"
	"sync"
)

func main() {
	type Button struct {
		Clicked *sync.Cond
	}

	button := Button{
		Clicked: sync.NewCond(&sync.Mutex{}),
	}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup // instnatiated at zero value
		goroutineRunning.Add(1)             // subscribe will not exit until the go routine is running
		go func() {
			goroutineRunning.Done() // signal we are in go routine
			c.L.Lock()              // take lock because c.Wait will unlock and suspend
			defer c.L.Unlock()      // unlock afer c.Wait because when c.Wait is released by broadcast it locks again
			c.Wait()                // locked, c.Wait, unlock and suspend, broadcast, wake up, lock, continue
			fn()                    // cb
			// finally unlock
		}()
		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximising window.")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying anoying dialog box!")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})

	button.Clicked.Broadcast()
	clickRegistered.Wait()
}
