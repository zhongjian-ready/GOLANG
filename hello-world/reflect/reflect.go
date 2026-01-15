package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Id   int
	Name string
	Age  int
}

func (u User) Call() {
	fmt.Println("User Call method called")
	fmt.Printf("%v \n", u)
}

func main() {
	u := User{
		Id:   1,
		Name: "zhangsan",
		Age:  30,
	}
	DoFieldsAndMethods(u)
}

func DoFieldsAndMethods(input interface{}) {
	inputType := reflect.TypeOf(input)
	fmt.Println("input type =", inputType.Name())

	inputValue := reflect.ValueOf(input)
	fmt.Println("input value =", inputValue)

	for i := 0; i < inputType.NumField(); i++ {
		f := inputType.Field(i)
		v := inputValue.Field(i).Interface()
		fmt.Printf("field name = %s, field value = %v\n", f.Name, v)
		
	}

	for i := 0; i < inputType.NumMethod(); i++ {
		m := inputType.Method(i)
		fmt.Printf("method name = %s\n : %v\n", m.Name, m.Type)
	}


	
}