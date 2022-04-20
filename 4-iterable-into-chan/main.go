package main

func main() {
	done := make(chan interface{})
	stringStream := make(chan interface{})

	for _, s := range []string{"a", "b", "c"} {
		select {
		case <-done:
			return
		case stringStream <- s:
		}
	}
}
