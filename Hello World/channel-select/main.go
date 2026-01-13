package main

import (
	"fmt"
	"time"
)

func fibonacci(ch, quit chan int) {
	x, y := 1, 1

	for {
		select {
		case ch <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("fibonacci go func finished")
			return
		}
	}
}

func main() {
	ch := make(chan int)
	quit := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("main received %d from channel\n", <-ch)
		}
		quit <- 0
	}()

	fibonacci(ch, quit)

	time.Sleep(1 * time.Second) // wait for goroutine to finish
}