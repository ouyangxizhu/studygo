package main

import (
	"fmt"
	"time"
)

func main() {
	// main函数如果不结束，父goroutine即使结束，父goroutine中的子goroutine也会执行
	fmt.Println("main函数开始")
	go notMainFunc()
	time.Sleep(3 * time.Second)
	fmt.Println("main函数结束")
}

func notMainFunc() {
	fmt.Println("notMainFunc开始")
	go func() {
		time.Sleep(1 * time.Second)
		println("notMainFunc中子goroutine结束")
	}()
	fmt.Println("notMainFunc结束")
}
/*
main函数开始
notMainFunc开始
notMainFunc结束
notMainFunc中子goroutine结束
main函数结束
*/
