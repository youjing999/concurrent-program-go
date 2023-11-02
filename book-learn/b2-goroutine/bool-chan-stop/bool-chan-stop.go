package main

import (
	"fmt"
	"time"
)

func worker(stopChan chan bool) {
	for {
		select {
		case <-stopChan:
			fmt.Println("stop...")
			return
		default:
			fmt.Println("working")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
func main() {
	stopChan := make(chan bool)
	go worker(stopChan)

	time.Sleep(2 * time.Second)
	stopChan <- false
	time.Sleep(time.Second)
}
