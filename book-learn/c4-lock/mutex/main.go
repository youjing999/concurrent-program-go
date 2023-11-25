package main

import (
	"fmt"
	"sync"
	"time"
)

var count = 0
var mu sync.Mutex

func addNum() {
	mu.Lock()
	count++
	mu.Unlock()
}

func main() {
	for i := 1; i <= 10; i++ {
		go addNum()
	}
	time.Sleep(time.Second)
	fmt.Println(count)
}
