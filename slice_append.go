package main

import (
	"fmt"
)

func main() {
	a := map[string]int{"one": 1, "two": 2, "three": 3, "four": 4}

	var b []string // 定义顺序
	b = append(b, "one", "two", "three", "four")

	for k, v := range a { // 无序循环
		fmt.Printf("%v : %v, ", k, v)
	}

	fmt.Println()

	for _, element := range b { // 按定义的顺序循环
		fmt.Printf("%v : %v, ", element, a[element])
	}
}
