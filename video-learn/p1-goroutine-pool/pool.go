package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/panjf2000/ants/v2"
)

func GoroutineAnts() {
	//统计当前存在的goroutine的数据
	go func() {
		for {
			fmt.Println("the num of goroutine:", runtime.NumGoroutine())
			time.Sleep(500 * time.Millisecond)
		}
	}()

	//初始化协程池 goroutine pool
	size := 1024
	pool, err := ants.NewPool(size)
	if err != nil {
		log.Fatalln(err)
	}
	defer pool.Release()

	// 利用 pool 调度需要并发的大量goroutine
	for {
		//向pool中提交一个执行的goroutine
		err := pool.Submit(func() {
			v := make([]int, 1024)
			_ = v
			fmt.Println("in goroutine")
			time.Sleep(10 * time.Second)
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
}
