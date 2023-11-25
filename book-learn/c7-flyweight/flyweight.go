package main

import (
	"fmt"
	"sync"
)

type User struct {
	Name string
	Id   int
}

func main() {
	pool := sync.Pool{
		New: func() interface{} {
			return new(User)
		},
	}
	// 从对象池中获取一个对象
	user := pool.Get().(*User)
	user.Name = "xz"
	user.Id = 1
	pool.Put(user)

	// 从对象池中获取一个对象
	user2 := pool.Get().(*User)
	fmt.Println(*user2)
}
