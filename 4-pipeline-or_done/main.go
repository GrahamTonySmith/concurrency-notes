package main

type any = interface{}

// for when we work with channels and we can not assume
// behaviour about how those channels behave when canceled
// via their done channel
func orDone(done <-chan any, c <-chan any) <-chan any {
	valueStream := make(chan any)
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done: // we called done so return
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valueStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valueStream
}
