package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("start")
	tick := time.Tick(time.Second)

	for t := range tick {
		fmt.Printf("%v\n", t)
	}
}
