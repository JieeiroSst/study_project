package main

import "fmt"

type callback func()
type fn func(int)

func main() {
	barriers := Barrier(10)

	function := func(n int) {
		// generate multiple coroutine which print "Hi from i" and "Bye from i"
		for i := range barriers {
			go func(i int, ch chan bool) {
				for j := 0; j < n; j++ {
					// have to print all "Hi from i" before "Bye from i"
					fmt.Println("Hi from", i)
					ch <- false
					Wait(ch)

					fmt.Println("Bye from", i)
					ch <- true
					Wait(ch)
				}
			}(i, barriers[i])
		}
	}

	callback := func() {
		// wait for all "Hi from i" printed
		for i := range barriers {
			<- barriers[i]
		}

		PulseAll(barriers)

		// wait for all "Bye from i" printed
		for i := range barriers {
			<- barriers[i]
		}

		PulseAll(barriers)
	}

	Sync(10, function, callback)
}

func Barrier(i int) []chan bool {
	// return slice of #i bool channel(s)
	channels := make([]chan bool, i)
	for i := range channels {
		channels[i] = make(chan bool)
	}
	return channels
}

func Sync(n int, function fn, callback callback) {
	function(n)
	for i := 0; i < n; i++ {
		callback()
	}
}

func PulseAll(barriers []chan bool) {
	// send signal to all channel
	for i := range barriers {
		barriers[i] <- true
	}
}

func Wait(ch chan bool) {
	<- ch
}