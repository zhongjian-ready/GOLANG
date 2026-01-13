package main

import "log"

func foo(a string, b int)(string, int){
	return a, b
	
}

 
func main() {
	s, n := foo("hello", 42)
	println(s)
	println(n)
	log.Println("标准日志")
}