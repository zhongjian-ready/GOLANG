package main

import "fmt"

type Hunman struct {
	name string
	age  int
}

func (h Hunman) Show() {
	fmt.Println("name =", h.name)
	fmt.Println("age =", h.age)
}

func (h *Hunman) Eat() {
	fmt.Println(h.name, "is eating")
}

func (h *Hunman) Sleep() {
	fmt.Println(h.name, "is sleeping")
}

type   struct {
	Hunman
	school string
}

func (s Student) Sleep() {
	fmt.Println(s.name, "is sleeping at", s.school)
}

func main() {
	s := Student{
		Hunman: Hunman{
			name: "zhangsan",
			age:  20,
		},
		school: "No.1 Middle School",
	}
	s.Show()
	s.Eat()
	s.Sleep()
}
