package main

import "fmt"

func main() {	// This is a placeholder for the main function.
	
str := "Go语言"

// 按 byte 统计（UTF-8 编码）
// len(str) = 2 (Go) + 3 (语) + 3 (言) = 8 bytes
fmt.Println(len(str)) 

// 按 rune 统计（Unicode 字符）
// len([]rune(str)) = 2 (Go) + 1 (语) + 1 (言) = 4 runes
fmt.Println(len([]rune(str))) 

}