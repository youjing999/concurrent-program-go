package main

import (
	"fmt"
	"sync"
)

func main() {

	count := 0
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			mu.Lock()
			count++
			mu.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(count)
}
