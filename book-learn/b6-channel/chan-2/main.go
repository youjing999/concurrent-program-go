package main

import "fmt"

func Producer(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
		fmt.Println("send:", i)
	}
	close(ch)
}

func Consumer(ch chan int) {
	for {
		val, ok := <-ch
		if !ok {
			break
		}
		fmt.Println("consume:", val)
	}
}

func main() {
	ch := make(chan int, 2)
	go Producer(ch)
	Consumer(ch)
}
