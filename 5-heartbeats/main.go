package main

import (
	"fmt"
	"time"
)

type any = interface{}

func doWork(done <-chan any, pulseInterval time.Duration) (<-chan any, <-chan time.Time) {
	heartbeat := make(chan any) // set up a heartbeat and results channel to return
	results := make(chan time.Time)
	go func() {
		defer close(heartbeat)
		defer close(results)

		pulse := time.Tick(pulseInterval)       // set the heartbeat t opulse at the pulse interval
		workGen := time.Tick(2 * pulseInterval) // use a ticker to simulate work so we can see heartbeats

		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}: // struct with no properties has no address
			default: // we include a default because we can not guarentie that heartbeat chan will be received from and thus might be blocked
			}
		}

		sendResult := func(r time.Time) {
			for {
				select {
				case <-done:
					return
				case <-pulse: // just like with done channels any time we send or receive we need to handle heartbeats
					sendPulse()
				case results <- r:
					return
				}
			}
		}

		for {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}
	}()
	return heartbeat, results
}

func main() {
	done := make(chan any)
	time.AfterFunc(10*time.Second, func() { close(done) }) // wait 10 seconds and close 1

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat: // select from heartbeat, we should get it every timeout / 2, if not then we have a problem
			if ok == false {
				return
			}
			fmt.Printf("pulse\n")
		case r, ok := <-results: // 5 select from the results chan
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout): // select from time.After if the timeout is hit before we get a heartbeat
			fmt.Printf("timedout")
			return
		}
	}
}
