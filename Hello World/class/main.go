package main

import "fmt"

// Hero 结构体 首字母大写表示对外可见，小写表示包内可见
type Hero struct {
	name string
	age  int
}

func (this Hero) Show() {
	fmt.Println("name =", this.name)
	fmt.Println("age =", this.age)
}

func (this *Hero) GetName() string {
	return this.name
}

func (this *Hero) SetName(name string) {
	this.name = name
}

func main() {
	h := Hero{name: "zhangsan", age: 30}
	h.GetName()
	h.SetName("lisi")
	h.GetName()
}
