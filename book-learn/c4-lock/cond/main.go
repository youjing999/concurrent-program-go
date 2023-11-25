package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	ready := false

	for i := 0; i < 5; i++ {
		go func(i int) {
			mu.Lock()
			if !ready {
				cond.Wait()
			}
			fmt.Printf("worked %d is working\n", i)
			mu.Unlock()
		}(i)
	}

	time.Sleep(time.Second)
	mu.Lock()
	ready = true
	cond.Broadcast()
	mu.Unlock()

	time.Sleep(time.Second)
}
