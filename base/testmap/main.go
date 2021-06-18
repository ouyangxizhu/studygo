package main

import "fmt"

func main() {

	//var rankTypeToTopLeftText = map[int64][]string{
	//	1:                []string{"带货榜主播口碑评分较高购买更放心"},
	//}
	//s:= rankTypeToTopLeftText[0]
	//fmt.Println(s == nil)
	//fmt.Printf("%#v\n", s)
	//
	//var inttoint = map[int64]int{
	//	1:                1,
	//}
	//i:= inttoint[0]
	//fmt.Println(i == 0)
	//fmt.Printf("%#v", i)

	//var intToBool = map[int64] bool{
	//	1: true,
	//}
	//fmt.Println(intToBool[0])

	m := make(map[int64]bool, 10)
	fmt.Println(len(m))

}
