package main

import "time"

func main() {
	// main函数中goroutine如果结束，子goroutine不会再执行
	println("main函数开始")
	go func() {
		time.Sleep(1 * time.Second)
		println("main中子goroutine结束")
	}()
	println("main函数结束")
}
