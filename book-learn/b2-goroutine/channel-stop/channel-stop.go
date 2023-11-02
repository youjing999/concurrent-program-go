package main

import (
	"fmt"
	"time"
)

func printNum(stop <-chan struct{}, n int) {
	for i := 0; i < n; i++ {
		select {
		case <-stop:
			fmt.Println("协程被关闭")
			return
		default:
			fmt.Printf("子协程的数字：%d\n", i)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	stopChan := make(chan struct{})
	go printNum(stopChan, 10)
	time.Sleep(3 * time.Second)
	close(stopChan)

	time.Sleep(time.Second * 1)
	fmt.Println("主协程退出")
}
