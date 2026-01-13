package main

import "fmt"

type Animal interface {
	Speak() string
	Sleep() string
	GetType() string
	GetColor() string
}

type Dog struct {
	animalType string
	color      string
}

// 实现 Animal 接口的方法(need to implement all methods)
func (d Dog) Speak() string {
	return "Woof!"
}

func (d Dog) Sleep() string {
	return "Zzz..."
}

func (d Dog) GetType() string {
	return d.animalType
}

func (d Dog) GetColor() string {
	return d.color
}

// cat
type Cat struct {
	animalType string
	color      string
}

func (c Cat) Speak() string {
	return "Meow!"
}

func (c Cat) Sleep() string {
	return "Zzz..."
}

func (c Cat) GetType() string {
	return c.animalType
}

func (c Cat) GetColor() string {
	return c.color
}

// 实现多态
func makeAnimalSpeak(a Animal) {
	fmt.Printf("The %s (%s) says: %s and sleeps: %s\n", a.GetType(), a.GetColor(), a.Speak(), a.Sleep())
}

func main() {
	dog := Dog{animalType: "Dog", color: "Brown"}
	cat := Cat{animalType: "Cat", color: "Black"}

	makeAnimalSpeak(dog)
	makeAnimalSpeak(cat)
	
}
