package main

import (
	"fmt"
	"strconv"
)

func main() {
	var i string = ""
	num, _ := strconv.ParseInt(i, 10, 64)
	fmt.Println(num)

}
