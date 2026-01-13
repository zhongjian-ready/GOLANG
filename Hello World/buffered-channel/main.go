package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 3)

	fmt.Println("buffered channel created")
	
	go func() {
		defer fmt.Println("go func finished")

		for i := 0; i < 5; i++ { // send 5 values, if buffer is full, it will block
			ch <- i
			fmt.Printf("go func sent %d\n", i)
		}
	}()

	time.Sleep(1 * time.Second) // wait for goroutine to start sending

	for i := 0; i < 3; i++ {
		val := <-ch
		fmt.Printf("main received %d\n", val)
	}

	fmt.Println("main finished")
}
			

