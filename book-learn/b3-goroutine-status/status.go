package main

import (
	"fmt"
	"time"
)

func A() {
	fmt.Println("New: 协程已被创建但还未开始执行任务")
}

func B() {
	fmt.Println("Runnable: 协程正在执行任务")
	time.Sleep(time.Second)
}

func main() {
	//新建
	go A()

	//运行
	go B()

	ch := make(chan bool)
	go func() {
		fmt.Println("Blocked: 协程因为等待channel接收数据而被暂停执行")
		<-ch
	}()

	// 死亡状态
	go func() {
		fmt.Println("Dead: 协程执行完成")
	}()

	time.Sleep(time.Second * 2)
}
