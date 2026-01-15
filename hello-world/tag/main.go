package main

import (
	"fmt"
	"reflect"
)

type resume struct {
	Name string `json:"name" doc:"姓名"`
	Age  int    `json:"age"`
}

func findTag (str interface{}) {
	t := reflect.TypeOf(str).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("field name: %s, json tag: %s, doc tag: %s\n", field.Name, field.Tag.Get("json"), field.Tag.Get("doc"))
	}
}

func main() {
	var r resume
	findTag(&r)
}