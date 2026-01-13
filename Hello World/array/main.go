package main

import "fmt"

func main() {
	var arr [5]int = [5]int{1, 2, 3, 4, 5}

	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])
	}

	for index, value := range arr {
		fmt.Printf("index=%d, value=%d\n", index, value)
	}

	for _, value := range arr { // _省略变量
		fmt.Printf("value=%d\n", value)
	}

	fmt.Printf("arr types = %T\n", arr)
}