package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	// time.Sleep(1 * time.Second)
	var a int
	var b = 10
	// b = "12"
	fmt.Println("a = ", a)
	fmt.Println("b = ", b)
	fmt.Printf("%T\n", a)
	
	var c float64 = 3.14
	fmt.Println("c = ", c)
	fmt.Printf("%T\n", c)
}