package main

import (
	"fmt"
	"sync"
	"time"
)

var count = 0
var rwLock sync.RWMutex
var wg sync.WaitGroup

func read(i int) {
	defer wg.Done()
	rwLock.RLock()
	defer rwLock.RUnlock()
	fmt.Printf("read:%v goroutine:%d\n", i, count)
}

func write() {
	defer wg.Done()
	rwLock.Lock()
	defer rwLock.Unlock()
	count++
	fmt.Println("write:", count)
}

func main() {
	for i := 0; i < 10; i++ {
		wg.Add(1)
		time.Sleep(time.Millisecond * 1000)
		go read(i)
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go write()
	}
	wg.Wait()
	fmt.Println("final")
}
