package main

import "encoding/json"

func main() {
	m := make(map[int]int)
	m[1] = 1
	marshal, _ := json.Marshal(m)
	println(string(marshal))

	mStr := make(map[string]int)
	mStr["1"] = 1
	marshalStr, _ := json.Marshal(mStr)
	println(string(marshalStr))
}
