package main

import (
	"fmt"
	"time"
)

func main() {

	ch := make(chan int)

	go func() {
		fmt.Println("in goroutine,wait channel")
		<-ch
		fmt.Println("channel wait success")
	}()

	time.Sleep(time.Second * 2)

	fmt.Println("main goroutine，set num to channel")
	ch <- 1
	time.Sleep(time.Second)
}
