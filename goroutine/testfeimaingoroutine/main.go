package main

import "time"

func main() {
	// main函数如果不结束，子函数即使结束，子函数中的goroutine也会执行
	println("main函数开始")
	notMainFunc()
	time.Sleep(3 * time.Second)
	println("main函数结束")
}

func notMainFunc() {
	println("notMainFunc开始")
	go func() {
		time.Sleep(1 * time.Second)
		println("notMainFunc中子goroutine结束")
	}()
	println("notMainFunc结束")
}
