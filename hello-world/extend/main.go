package main

import "fmt"

type Human struct {
	name string
	age  int
}

func (h Human) Show() {
	fmt.Println("name =", h.name)
	fmt.Println("age =", h.age)
}

func (h *Human) Eat() {
	fmt.Println(h.name, "is eating")
}

func (h *Human) Sleep() {
	fmt.Println(h.name, "is sleeping")
}

type Student struct {
	Human
	school string
}

func (s Student) Sleep() {
	fmt.Println(s.name, "is sleeping at", s.school)
}

func main() {
	s := Student{
		Human: Human{
			name: "zhangsan",
			age:  20,
		},
		school: "No.1 Middle School",
	}
	s.Show()
	s.Eat()
	s.Sleep()
}
