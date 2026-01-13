package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	// use go create an anonymous function that formal parameters is empty and has no return value
	// create a new goroutine
	go func() {
		defer fmt.Println("go func defer called")
		// anonymous function body
		func() {
			defer fmt.Println("inner func defer called")
			runtime.Goexit() // terminate the current goroutine
			fmt.Println("inner func called")
		}()
		
		fmt.Println("go func called")
	}()

	// create another goroutine with parameters and return value, this function will return a bool value but we can not get it in main goroutine
	go func (s string, n int) bool {
		fmt.Printf("go func with params called: s=%s, n=%d\n", s, n)
		return true
	}("hello", 42)

	// infinite loop to keep main goroutine alive
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("main goroutine is running")
	}
}

