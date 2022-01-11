package main

import "sync"

var arrOnce sync.Once
var arr []int

// getArr retrieves arr, lazily initializing on first call. Double-checked
// locking is implemented with the sync.Once library function. The first
// goroutine to win the race to call Do() will initialize the array, while
// others will block until Do() has completed. After Do has run, only a
// single atomic comparison will be required to get the array.
func getArr() []int {
	arrOnce.Do(func() {
		arr = []int{0, 1, 2}
	})
	return arr
}

func main() {
	// thanks to double-checked locking, two goroutines attempting to getArr()
	// will not cause double-initialization
	go getArr()
	go getArr()
}