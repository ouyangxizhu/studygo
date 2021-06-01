package main

import "fmt"

func main() {

	var rankTypeToTopLeftText = map[int64][]string{
		1:                []string{"带货榜主播口碑评分较高购买更放心"},
	}
	s:= rankTypeToTopLeftText[0]
	fmt.Println(s == nil)
	fmt.Printf("%#v", s)
}
