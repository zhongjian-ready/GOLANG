package main

import "fmt"

/*
arrays in Go are fixed-size sequences of elements of a single type. They are defined using the syntax [n]T,
where n is the number of elements and T is the type of each element.
1. fixed size
2. same type
3. indexable
4. contiguous in memory
*/
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