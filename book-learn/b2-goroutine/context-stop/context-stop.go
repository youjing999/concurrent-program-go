package main

import (
	"context"
	"fmt"
	"time"
)

func printNum(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("协程被取消")
			return
		default:
			fmt.Printf("子协程的数字：%d\n", i)
			time.Sleep(500 * time.Millisecond)
		}
	}
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go printNum(ctx, 10)

	//3秒后取消协程
	time.Sleep(3 * time.Second)
	cancel()

	//等待协程结束
	time.Sleep(1 * time.Second)
	fmt.Println("出协程退出")
}
