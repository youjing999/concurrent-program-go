package main

import "fmt"

func main() {

	ch := make(chan int)

	go func() {
		for i := 1; i <= 30; i++ {
			ch <- i
		}
		close(ch)
	}()

	go func() {
		for {
			num, ok := <-ch
			if ok {
				fmt.Println("receive", num)
			} else {
				fmt.Println("channel was closed")
				break
			}
		}

	}()

	fmt.Scanln()
	fmt.Println("main goroutine was done")
}
