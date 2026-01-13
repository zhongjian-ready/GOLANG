package main

import (
	"fmt"
	"time"
)

func newTask()  {
	i := 0
	for {
		i++
		fmt.Printf("new task: i = %d\n", i)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	go newTask() // 开启一个新的 goroutine 执行 newTask 函数

	i := 0
	for {
		i++
		fmt.Printf("main task: i = %d\n", i)
		time.Sleep(1 * time.Second)
	}
}
