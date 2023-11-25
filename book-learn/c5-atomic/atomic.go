package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	var count int32

	for i := 0; i < 10000; i++ {
		go func() {
			atomic.AddInt32(&count, 1)
		}()

	}
	time.Sleep(time.Second * 3)
	fmt.Println(count)
}
