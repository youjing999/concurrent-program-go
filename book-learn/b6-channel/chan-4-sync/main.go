package main

import (
	"fmt"
	"time"
)

func chanSync1(ch chan string) {
	fmt.Println("goroutine 1")
	ch <- "flag"
}

func chanSync2(ch chan string) {
	msg := <-ch
	fmt.Println("goroutine 2")
	fmt.Println("msg from chanSync 1", msg)
}

func main() {
	ch := make(chan string)

	go chanSync2(ch)
	go chanSync1(ch)

	time.Sleep(time.Second)
	fmt.Println("done")

}
