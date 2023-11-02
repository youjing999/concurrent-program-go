package main

import (
	"fmt"
	"sync"
)

var num = 0
var mu sync.Mutex

func main() {
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		for i := 0; i < 10000; i++ {
			mu.Lock()
			num++
			mu.Unlock()
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			mu.Lock()
			num--
			mu.Unlock()
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("num = ", num)
}
