package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	done := false

	go func() {
		time.Sleep(time.Second * 2)
		mu.Lock()
		fmt.Println("子协程处理")
		done = true
		cond.Signal()
		mu.Unlock()
	}()

	mu.Lock()
	if !done {
		cond.Wait()
	}
	fmt.Println("协程等待结束")
	mu.Unlock()

}
