package main

import (
	"fmt"
	"sync"
	"time"
)

func GoroutineWG() {

	wg := sync.WaitGroup{}
	printOdd := func() {
		defer wg.Done()
		for i := 1; i <= 10; i += 2 {
			fmt.Println(i)
			time.Sleep(100 * time.Millisecond)
		}
	}

	printEven := func() {
		defer wg.Done()
		for i := 2; i <= 10; i += 2 {
			fmt.Println(i)
			time.Sleep(100 * time.Millisecond)
		}
	}
	wg.Add(2)
	go printOdd()
	go printEven()
	wg.Wait()
	fmt.Println("after group")
}

func main() {
	GoroutineWG()
}
