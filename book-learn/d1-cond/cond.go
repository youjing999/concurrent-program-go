package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	cond := sync.NewCond(&mu)
	done := false

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mu.Lock()
			for !done {
				fmt.Printf("goroutine %d is waitting\n", i)
				cond.Wait()
				fmt.Printf("goroutine %d was awaked\n", i)
			}
			mu.Unlock()
		}(i)
	}

	time.Sleep(time.Second * 2)
	mu.Lock()
	done = true
	fmt.Println("broadcast")
	cond.Broadcast()
	mu.Unlock()

	wg.Wait()
	fmt.Println("finish")
}
