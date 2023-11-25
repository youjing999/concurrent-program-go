package main

import (
	"fmt"
	"time"
)

func main() {

	c := make(chan int)

	go func() {
		time.Sleep(3 * time.Second)
	}()

	select {
	case res := <-c:
		fmt.Println(res)
	case <-time.After(time.Second * 2):
		fmt.Println("超时")

	}
}
