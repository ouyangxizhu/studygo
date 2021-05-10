package main

import "fmt"

func main() {
	var a interface{} = nil
	m:= a.(map[int]int) //panic: interface conversion: interface {} is nil, not map[int]int
	fmt.Println(m)
}
