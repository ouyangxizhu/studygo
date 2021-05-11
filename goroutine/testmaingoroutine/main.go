package main

import (
	"fmt"
	"time"
)

func main() {
	// main函数中goroutine如果结束，子goroutine不会再执行
	fmt.Println("main函数开始")
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("main中子goroutine结束")
	}()
	fmt.Println("main函数结束")
}
/*
main函数开始
main函数结束
 */
