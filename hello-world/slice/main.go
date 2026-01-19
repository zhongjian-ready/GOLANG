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

	x := make([]int, 0, 10)
	x = append(x, 1, 2, 3)
	y := append(x, 4)
	z := append(x, 5)
	fmt.Println("x =", x) // x 不受 y 和 z 影响 x = [1 2 3]
	fmt.Println("y =", y) // y = [1 2 3 5] // y 和 z 共享 x 的底层数组
	fmt.Println("z =", z) // z = [1 2 3 5] // y 和 z 共享 x 的底层数组
}
