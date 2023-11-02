package main

import (
	"fmt"
	"time"
)

func main() {

	chan1 := make(chan string)
	chan2 := make(chan string)

	go func() {
		time.Sleep(time.Second * 2)
		chan1 <- "chan1 msg"
	}()

	go func() {
		time.Sleep(time.Second * 2)
		chan2 <- "chan2 msg"
	}()

	for i := 0; i < 2; i++ {
		select {
		case <-chan1:
			fmt.Println("chan1 msg out")
		case <-chan2:
			fmt.Println("chan2 msg out")
		}

	}
}
