package main

import (
	"fmt"
	"time"
)

func main() {
	names := []string{"Eric", "Harry", "Robert", "Mark"}
	for _, name := range names {
		go func(name string) {
			fmt.Printf("Hello %s\n", names)
		}(name)
	}
	time.Sleep(time.Millisecond)
}
