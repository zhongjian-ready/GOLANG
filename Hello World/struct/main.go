package main

import (
	"fmt"
)

type myint int

type Book struct {
	title  string
	author string
}

func changeBook(b *Book) {
	b.title = "C语言"
	b.author = "www.c.org"
}

func main() {
	var a myint = 10
	fmt.Println("a =", a)
	fmt.Printf("%T\n", a)

	var b Book
	b.title = "Go语言"
	b.author = "www.golang.org"
	fmt.Println("b =", b)
	fmt.Printf("%T\n", b)
	changeBook(&b)
	fmt.Println("b =", b)
}
