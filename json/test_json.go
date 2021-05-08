package main

import (
	"encoding/json"
	"fmt"
)

// 结构体里面的属性必须大写，否则序列化后的结果是空的
// 结构体是不是指针没关系，但是json.Unmarshal传入的必须是指针
type stPtr struct {
	Str  string      `json:"str"`
	I    int         `json:"i"`
	Stru *StIn        `json:"stru"`
	M    map[int]int `json:"m"`
}
type st struct {
	Str  string      `json:"str"`
	I    int         `json:"i"`
	Stru StIn        `json:"stru"`
	M    map[int]int `json:"m"`
}

type stLowLetter struct {
	str string `json:"str"`
	i   int    `json:"i"`
}

type StIn struct {
	In int `json:"in"`
}



func main() {
	//fmt.Println("======属性小写结构体=====")
	//lowLetterStruct := stLowLetter{
	//	str: "str",
	//	i:   12,
	//}
	//fmt.Println("序列化前", lowLetterStruct)//序列化前 {str 12}
	//stLowLetterMarshal, err := json.Marshal(lowLetterStruct)
	//fmt.Println("非指针序列化结果", string(stLowLetterMarshal))//非指针序列化结果 {}
	//
	//stLowLetterPtrMarshal, err := json.Marshal(&lowLetterStruct)
	//fmt.Println("指针结构体序列化结果", string(stLowLetterPtrMarshal))//指针结构体序列化结果 {}
	
	//fmt.Println("======属性大写结构体=====")
	//fmt.Println("======属性结构体非指针=====")
	//s := st{
	//	Str: "str",
	//	I:   12,
	//	Stru: StIn{
	//		2,
	//	},
	//	M: map[int]int{
	//		1:1,
	//	},
	//}
	//fmt.Println("序列化前", s)
	//fmt.Println("======结构体序列化=====")
	//stMarshal, _ := json.Marshal(s)
	//stMarshalStr := string(stMarshal)
	//fmt.Println("非指针序列化结果", stMarshalStr)
	//fmt.Println("======反序列化=====")
	//item0 := st{}
	//json.Unmarshal([]byte(stMarshalStr), item0)
	//fmt.Println("Unmarshal非指针序列化结果", item0)
	//
	//item0Ptr := &st{}
	//json.Unmarshal([]byte(stMarshalStr), item0Ptr)
	//fmt.Println("Unmarshal指针序列化结果", *item0Ptr)


	//fmt.Println("======属性结构体是指针=====")
	//sPtr := stPtr{
	//	Str: "str",
	//	I:   12,
	//	Stru: &StIn{
	//		2,
	//	},
	//	M: map[int]int{
	//		1:1,
	//	},
	//}
	//fmt.Println("序列化前", sPtr)
	//fmt.Println("======结构体序列化=====")
	//stMarshalPtr, err := json.Marshal(sPtr)
	//fmt.Println("非指针序列化结果", string(stMarshalPtr))
	//
	//stPtrMarshal, err := json.Marshal(&sPtr)
	//fmt.Println("指针结构体序列化结果", string(stPtrMarshal))
	//
	//fmt.Println("======map中值非指针结构=====")
	//m := make(map[int]stPtr)
	//m[1] = sPtr
	//fmt.Println("序列化前结果", m)
	//marshal, err := json.Marshal(m)
	//if err != nil {
	//	fmt.Println("序列化错误", err)
	//	return
	//}
	//marshalStr := string(marshal)
	//fmt.Println("序列化结果", marshalStr)
	//
	//fmt.Println("======map中值指针结构体=====")
	//m2 := make(map[int]*stPtr)
	//m2[1] = &sPtr
	//fmt.Println("序列化前结果", m2)
	//marshal2, _ := json.Marshal(m2)
	//marshalStr2 := string(marshal2)
	//fmt.Println("序列化结果", marshalStr2)
	//
	//fmt.Println("======map中值非指针结构体=====")
	//m3 := make(map[int]stPtr)
	//m3[1] = sPtr
	//fmt.Println("序列化前结果", m3)
	//marshal3, _ := json.Marshal(m3)
	//marshalStr3 := string(marshal3)
	//fmt.Println("序列化结果", marshalStr3)
	//fmt.Println("map中值非指针结构体序列化结果和map中值指针结构体结果是否相等:", marshalStr3 == marshalStr2)
	//
	//fmt.Println("======反序列化======")
	//itemMap:= make(map[int]stPtr)
	//json.Unmarshal([]byte(marshalStr3), itemMap)
	//fmt.Println("非指针Unmarshal", itemMap)
	//
	//json.Unmarshal([]byte(marshalStr3), &itemMap)
	//fmt.Println("入参指针非指针Unmarshal", itemMap)
	//
	//
	//itemMapPtr:= make(map[int]*stPtr)
	//json.Unmarshal([]byte(marshalStr3), &itemMapPtr)
	//fmt.Println("指针Unmarshal", itemMapPtr)
	//
	//fmt.Println("======map中值指针结构体=====")
	//m3Ptr := make(map[int]*stPtr)
	//m3Ptr[1] = &sPtr
	//fmt.Println("序列化前结果", m3Ptr)
	//marshal3Ptr, _ := json.Marshal(m3Ptr)
	//marshalStr3Ptr := string(marshal3Ptr)
	//fmt.Println("序列化结果", marshalStr3Ptr)
	//
	//fmt.Println("======反序列化======")
	//itemMapPtrPtr:= make(map[int]*stPtr)
	//json.Unmarshal([]byte(marshalStr3Ptr), &itemMapPtrPtr)
	//fmt.Println("指针Unmarshal", itemMapPtrPtr)
	//
	//itemMapPtrPtr1:= make(map[int]*stPtr)
	//json.Unmarshal([]byte(marshalStr3Ptr), &itemMapPtrPtr1)
	//fmt.Println("指针Unmarshal", *itemMapPtrPtr1[1].Stru)

	i:= 1
	marshal, err := json.Marshal(i)
	if err != nil {
		fmt.Println(err)
		return
	}
	println(string(marshal))

	bytes, err := json.Marshal(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	println(string(bytes))

	itemMap:=make(map[int]int)
	err = json.Unmarshal(bytes, &itemMap)

	itemMap1:=-1
	err = json.Unmarshal(bytes, &itemMap1)
	println(itemMap1)

	itemMap2:="fadf"
	err = json.Unmarshal(bytes, &itemMap2)
	println(itemMap2)

	itemMap3:="fadf"
	str2, err := json.Marshal("itemMap3")
	err = json.Unmarshal([]byte(str2), &itemMap3)
	println(itemMap3)

	itemMap4:=1
	str4, err := json.Marshal(itemMap4)
	err = json.Unmarshal(str4, &itemMap4)
	println(itemMap4)
}
