package main

import "fmt"

func main() {
	s := []int{10, 20, 30, 40, 50}
	fmt.Println(s[1:4])
	fmt.Println(s[:3])
	fmt.Println(s[2:])
	fmt.Println(s[:])
	fmt.Println("len = ", len(s))

	numbers := make([]int, 3, 5)
	fmt.Println(numbers)
	fmt.Println("len =", len(numbers))
	fmt.Println("cap =", cap(numbers))

	numbers = append(numbers, 10)
	fmt.Println(numbers)
	fmt.Println("len =", len(numbers))
	fmt.Println("cap =", cap(numbers))



	// 使用固定次数展示切片扩容，避免使用动态 cap(numbers) 导致死循环
	for i := 0; i < 10; i++ {
		numbers = append(numbers, i*10)
		fmt.Printf("len = %d, cap = %d, numbers = %v\n", len(numbers), cap(numbers), numbers)
	}

	s1 := make([]int, 5)
	copy(s1, s)
	fmt.Println("s1 =", s1)
}