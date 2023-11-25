package main

import (
	"fmt"
	"sync"
)

func main() {
	var rwLock sync.RWMutex
	var data int

	wg := sync.WaitGroup{}
	wg.Add(5)

	go func() {
		rwLock.Lock()
		defer rwLock.Unlock()

		data = 23
		fmt.Println("write operation completed")
		wg.Done()
	}()

	for i := 0; i < 4; i++ {
		go func() {
			rwLock.RLock()
			defer rwLock.RUnlock()

			fmt.Println("Read operation: ", data)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("all operations have finished:", data)
}

func deadFunc() {
	var mu1, mu2 sync.Mutex
	wg := sync.WaitGroup{}
	wg.Add(2)
	// 协程1 先获取mu1再获取mu2
	go func() {
		mu1.Lock()
		fmt.Println("goroutine 1 m1 Lock")
		mu2.Lock()
		fmt.Println("goroutine 1 m2 Lock")
		mu2.Unlock()
		fmt.Println("goroutine 1 m2 unlock")
		mu1.Unlock()
		fmt.Println("goroutine 1 m1 unlock")
		wg.Done()
	}()

	// 协程2 先获取mu2再获取mu1
	go func() {
		mu2.Lock()
		fmt.Println("goroutine 2 m2 lock")
		mu1.Lock()
		fmt.Println("goroutine 2 m1 lock")
		mu1.Unlock()
		fmt.Println("goroutine 2 m1 unlock")
		mu2.Unlock()
		fmt.Println("goroutine 2 m2 unlock")
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("al goroutines have finished")
}
