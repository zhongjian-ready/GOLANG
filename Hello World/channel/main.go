package main

import "fmt"

func main() {
	// define a channel
	ch := make(chan int)

	go func() {
		defer fmt.Println("go func called")
		fmt.Println("go func started")
		// send a value to the channel
		ch <- 42
		fmt.Println("go func finished sending")
	}()

	// receive a value from the channel， blocking until a value is available
	val := <-ch
	fmt.Println("main received value:", val)
}