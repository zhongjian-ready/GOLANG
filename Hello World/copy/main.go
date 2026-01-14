package main

// in Golang. there is no Deep Copy  or Shallow Copy concept like other languages. It's also can be said there is all deep copy
// but we can achieve similar behavior using different techniques.

import "fmt"

func main() {
	// Example of array copy (value type)
	arr1 := [3]int{1, 2, 3}
	arr2 := arr1 // This creates a copy of arr1
	arr2[0] = 10
	fmt.Println("arr1:", arr1) // Output: arr1: [1 2 3]
	fmt.Println("arr2:", arr2) // Output: arr2: [10 2 3]

	// Example of slice copy (reference type)
	slice1 := []int{1, 2, 3}
	slice2 := slice1 // This creates a reference to slice1
	slice2[0] = 10
	fmt.Println("slice1:", slice1) // Output: slice1: [10 2 3]
	fmt.Println("slice2:", slice2) // Output: slice2: [10 2 3]

	// To create a true copy of a slice, use the built-in copy function
	slice3 := make([]int, len(slice1))
	copy(slice3, slice1)
	slice3[0] = 20
	fmt.Println("slice1:", slice1) // Output: slice1: [10 2 3]
	fmt.Println("slice3:", slice3) // Output: slice3: [20 2 3]
}
