package main

import (
	"encoding/json"
	"fmt"
)

type Movie struct {
	Title  string `json:"title"`
	Year   int    `json:"year"`
	Rating float32 `json:"rating"`
	Actors []string `json:"actors"`
}

func main() {
	m := Movie{
		Title:  "Inception",
		Year:   2010,
		Rating: 8.8,
		Actors: []string{"Leonardo DiCaprio", "Joseph Gordon-Levitt", "Ellen Page"},
	}

	// 编码的过程 结构体 --> JSON字符串
	data, err := json.Marshal(m) 
	if err != nil {
		fmt.Println("json.Marshal error =", err)
		return
	}
	fmt.Printf("data = %s\n", data)

	// 解码的过程 JSON字符串 --> 结构体
	var m2 Movie
	err = json.Unmarshal(data, &m2)
	if err != nil {
		fmt.Println("json.Unmarshal error =", err)
		return
	}
	fmt.Printf("m2 = %+v\n", m2)
}