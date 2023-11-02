package main

import (
	"fmt"
	"sync"
)

var (
	num int
	m   sync.Mutex
)

func increase(wg *sync.WaitGroup, caseNum int) {
	for i := 0; i < 500; i++ {
		m.Lock()
		num += 1
		fmt.Println("第", caseNum, "个case", num)
		m.Unlock()
	}
	wg.Done()
}

func main() {
	num = 0
	wg := sync.WaitGroup{}

	wg.Add(2)
	go increase(&wg, 1)
	go increase(&wg, 2)

	wg.Wait()
	fmt.Println(num)
}
