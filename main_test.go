package main

import (
	"sync"
	"testing"
)

func TestBenchmarkMutex(t *testing.T) {
	var m sync.Mutex
	for n := 0; n < 10000; n++ {
		m.Lock()
		m.Unlock()
	}
}