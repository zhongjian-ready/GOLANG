package main

import (
	"fmt"
)

func main() {
	var map1 map[string]int

	// 忽略这个警告 
	if map1 == nil {
		fmt.Println("map1 is nil")
		map1 = make(map[string]int)
	}
	map1["a"] = 10
	map1["b"] = 20
	fmt.Println("map1 =", map1)

	value, ok := map1["a"]
	if ok {
		fmt.Println("key a exists, value =", value)
	} else {
		fmt.Println("key a does not exist")
	}

	delete(map1, "b")
	fmt.Println("map1 after deleting key b =", map1)

	cityMap := map[string]string{
		"Beijing": "100000",
		"Shanghai": "200000",
		"Guangzhou": "510000",
	}
	fmt.Println("cityMap =", cityMap)

	for city, code := range cityMap {
		fmt.Printf("City: %s, Postal Code: %s\n", city, code)
	}

	// 复制map
	map2 := make(map[string]int)
	for k, v := range map1 {
		map2[k] = v
	}
	fmt.Println("map2 =", map2)
}
