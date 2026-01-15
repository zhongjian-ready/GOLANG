package main

func deferFunc() {
	println("defer func called")
}

func returnFunc() int {
	println("return func called")
	return 42
}

func main() {
	defer deferFunc()
	println("main func called")
	result := returnFunc()
	println("result =", result)
}