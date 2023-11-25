package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	pool := make(chan int, 3)
	for i := 0; i < 3; i++ {
		pool <- i
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		n := <-pool
		fmt.Println("goroutine 1 acquired resource", n)
		// 模拟执行任务
		fmt.Println("goroutine 1 is working...")
		time.Sleep(time.Second)
		wg.Done()
		// 将资源归还到资源池中
		pool <- n
		fmt.Println("goroutine 1 released resource", n)
	}()
	// 协程2 从资源池中获取资源并执行任务
	go func() {
		n := <-pool
		fmt.Println("goroutine 2 acquired resource", n)
		// 模拟执行任务
		fmt.Println("goroutine 2 is working...")
		time.Sleep(time.Second)
		wg.Done()
		// 将资源归还到资源池中
		pool <- n
		fmt.Println("goroutine 2 released resource", n)
	}()

	// 协程3 从资源池中获取资源并执行任务
	go func() {
		n := <-pool
		fmt.Println("goroutine 3 acquired resource", n)
		// 模拟执行任务
		fmt.Println("goroutine 3 is working...")
		time.Sleep(time.Second)
		wg.Done()
		// 将资源归还到资源池中
		pool <- n
		fmt.Println("goroutine 3 released resource", n)
	}()

	// 等待所有协程执行完毕
	wg.Wait()
	fmt.Println("All goroutines have finished")
}
