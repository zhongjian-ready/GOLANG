package main

import (
	"fmt"
)

func foo(p int) {
	p = 100
	fmt.Println("p =", p)
}

func foo1(p *int)  {
	*p = 200
	fmt.Println("p =", *p)
}

func main() {
	a := 10
	foo(a)
	fmt.Println("a =", a)
	foo1(&a)
	fmt.Println("a =", a)

	// & 取地址符 * 取值符
}