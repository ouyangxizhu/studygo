package main

import "fmt"

type testStruct struct {
	index int64
}

func main() {
	st := &testStruct{
		index: 0,
	}
	if st.index == 0 {
		st.index = 90
	}else{
		st.index = 100
	}
	fmt.Println(st.index)
}
