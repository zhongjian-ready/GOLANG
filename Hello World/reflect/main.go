package main

import (
	"fmt"
	"reflect"
)

func reflectNumber(v interface{}) {
	fmt.Println("type:", reflect.TypeOf(v))
	fmt.Println("value:", reflect.ValueOf(v))

	value := reflect.ValueOf(v)
	kind := value.Kind()	
	if kind == reflect.Int {
		n := value.Int()
		fmt.Println("int value:", n)
	} else if kind == reflect.Float64 {
		n := value.Float()
		fmt.Println("float64 value:", n)
	}
}

func main_old() {
	var num float64 = 3.14
	reflectNumber(num)
}
