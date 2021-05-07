package main

import (
	"fmt"
	"sync"
)

const (
	goroutineNum = 4
)
func main() {
	var wg sync.WaitGroup //不用初始化，可以直接用
	// wg.Add(goroutineNum) // 或在这一次性加完
	for i := 0; i < goroutineNum; i++ {
		wg.Add(1) //每次循环（启动一个goroutine）加1，或在上面提前一次性加完
		go func(num int) {
			defer wg.Done()
			fmt.Println(num)
		}(i)
	}
	wg.Wait() //阻塞等待直到所有的4个goroutine执行完
	fmt.Println("finish")
	//执行结果，数字顺序不固定
	/*
	3
	1
	2
	0
	finish
	*/
}
