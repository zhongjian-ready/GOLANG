package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	go func() {
		defer fmt.Println("defer in goroutine called")
		for i := 0; i < 5; i++ {
			ch <- i
			fmt.Printf("goroutine sent %d\n", i)
		}
		close(ch) // if there is no close, main goroutine will block forever
	}() 

	for {
		if data, ok := <-ch; ok {
			fmt.Printf("main received %d\n", data)
		} else {
			break
		}
	}
	
	// we also can use range to receive values until the channel is closed
	/*
	for data := range ch {
		fmt.Printf("main received %d\n", data)
	}
	*/
	fmt.Println("main finished receiving")
	time.Sleep(1 * time.Second) // wait for goroutine to finish
}
