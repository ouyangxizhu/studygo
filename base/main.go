package main

import (
	"fmt"
	"strconv"
)

func main() {
	i := 32
	var in int64 = 43
	i2:= int(in)
	fmt.Printf("%#v\n", i2)

	s := string(i)
	println(s)
	fmt.Printf("%#v\n", s)

	s2 := strconv.Itoa(i)
	fmt.Printf("%#v\n", s2)
}



