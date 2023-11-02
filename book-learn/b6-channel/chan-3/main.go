package main

import (
	"fmt"
	"time"
)

func send(ch chan int) {
	fmt.Println("send the first msg")
	ch <- 1
	fmt.Println("send the second msg")
	ch <- 2
}
func main() {
	ch := make(chan int, 1)
	go send(ch)
	time.Sleep(time.Second * 2)
	fmt.Println("receive 1")
	fmt.Println(<-ch)
	fmt.Println("receive 2")
	fmt.Println(<-ch)
}
