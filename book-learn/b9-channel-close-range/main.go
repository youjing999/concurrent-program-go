package main

import (
	"fmt"
)

func main() {
	channel := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			channel <- i
			fmt.Println(i, "in")
		}
		close(channel)
	}()

	for i := range channel {
		fmt.Println(i)
	}
}
