package main

import "fmt"

type st struct {
	index int64
}
func main() {
	s:= []*st{&st{
		index: 1,
	}}
	s= nil
	for _, num := range s {
		fmt.Println(num)
	}
}
