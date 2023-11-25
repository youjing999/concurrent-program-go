package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			once.Do(func() {
				fmt.Println("only one")
			})
		}(i)

		fmt.Println("goroutine ", i)
		wg.Done()
	}
	wg.Wait()
	fmt.Println("finish")
}
