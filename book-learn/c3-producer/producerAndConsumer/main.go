package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- int) {
	for i := 0; i < 100; i++ {
		ch <- i
		fmt.Print(i)
	}
	close(ch)
}

func consumer(ch <-chan int) {
	for i := range ch {
		fmt.Println("message from channel, no", i)
	}
}

func main() {
	ch := make(chan int, 10)
	go producer(ch)
	consumer(ch)
	time.Sleep(time.Second)

}
