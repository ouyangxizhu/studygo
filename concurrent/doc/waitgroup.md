[TOC]

# 介绍

## 官方介绍

> // A WaitGroup waits for a collection of goroutines to finish.

> // The main goroutine calls Add to set the number of

> // goroutines to wait for. Then each of the goroutines

> // runs and calls Done when finished. At the same time,

> // Wait can be used to block until all goroutines have finished.

> //

> // A WaitGroup must not be copied after first use.

大致翻译就是waitgroup用作主goroutine等待一组goroutine执行结束之后再执行

注意这里的最后一行，第一次使用后不允许复制使用。即使在函数中传值也要用指针传值（因为go中函数的传参数都是值传递）

源码中主要设计到的三个概念：counter、waiter和semaphore

counter： 当前还未执行结束的goroutine计数器

waiter : 等待goroutine-group结束的goroutine数量，即有多少个等候者

semaphore: 信号量

信号量是Unix系统提供的一种保护共享资源的机制，用于防止多个线程同时访问某个资源。

可简单理解为信号量为一个数值：

当信号量>0时，表示资源可用，获取信号量时系统自动将信号量减1；

当信号量=0时，表示资源暂不可用，获取信号量时，当前线程会进入睡眠，当信号量为正时被唤醒。


## 常用场景

主goroutine调用Add函数设置需要等待的goroutine的数量（可以多次设置）,当每个goroutine执行完成后调用Done函数(一般是在goroutine内调用Done()函数，将counter减1)，主goroutine调用Wait函数用于阻塞等待直到该组中的所有goroutine都执行完成。

看一个小demo

```go
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
```

这里不用new初始化是因为waitgroup是结构体类型，并且sync.WaitGroup的成员中没有需要通过`new`或者`make`进行初始化的，即声明的时候就赋值相应的零值了

# 源码分析

版本：go version go1.15.10 darwin/amd64(运行go version 可以得到)

## 结构体定义

```go
type WaitGroup struct {
	noCopy noCopy

	// 64-bit value: high 32 bits are counter, low 32 bits are waiter count.
	// 64-bit atomic operations require 64-bit alignment, but 32-bit
	// compilers do not ensure it. So we allocate 12 bytes and then use
	// the aligned 8 bytes in them as state, and the other 4 as storage
	// for the sema.
	state1 [3]uint32
}
```

- noCopy 这个属性使WaitGroup结构体不允许拷贝使用,只能用指针传递（避免复制使用的一个技巧，可以告诉vet工具违反了复制使用的规则）//禁止拷贝，如果有拷贝构建不会报错，可以用go vet 或 go tool vet 检测是否有拷贝错误，

- state1 声明为[3]uint32，一个uint32是32位，即4个字节；数组长度是3，即一共12个字节，共96位。再看一下state1的注释 

  1）64bit(8bytes)的值分成两段，高32bit是计数值，低32bit是waiter的计数。（另外32bit用作信号量）

  2）64位原子操作需要64位对齐。但是32位的编译器不支持。

  3）开辟了一个12字节的数组。64位对齐的8字节用作存储state（state包含计数值-高32位和waiter的计数-低32位），其余的4字节（32位）存储sema（信号量），即在不同的架构中代表的含义不一样

即WaitGroup结构体包含了一个noCopy 的辅助字段和一个具有复合意义的state1字段

## state()方法

```go
// state returns pointers to the state and sema fields stored within wg.state1.
// unsafe.Pointer其实就是类似C的void *，在golang中是用于各种指针相互转换的桥梁。
// uintptr是golang的内置类型，是能存储指针的整型，uintptr的底层类型是int，它和unsafe.Pointer可相互转换。
// uintptr和unsafe.Pointer的区别就是：unsafe.Pointer只是单纯的通用指针类型，用于转换不同类型指针，它不可以参与指针运算；
// 而uintptr是用于指针运算的，GC 不把 uintptr 当指针，也就是说 uintptr 无法持有对象，uintptr类型的目标会被回收。
// state()函数可以获取到wg.state1数组中元素组成的二进制对应的十进制的值
func (wg *WaitGroup) state() (statep *uint64, semap *uint32) {
	if uintptr(unsafe.Pointer(&wg.state1))%8 == 0 {
    // 如果地址是64bit（8字节）对齐的，数组前两个元素（8个字节）做state，后一个元素做信号量
		return (*uint64)(unsafe.Pointer(&wg.state1)), &wg.state1[2]
	} else {
    // 如果地址是32bit对齐的，数组后两个元素用来做state，它可以用来做64bit的原子操作，第一个元素32bit用来做信号量
    // 若不是,说明是4字节对齐,则后移4个字节后,这样必为8字节对齐,然后取后面8个字节作为*uint64类型
		return (*uint64)(unsafe.Pointer(&wg.state1[1])), &wg.state1[0]
	}
}
```

state()方法用于得到state的地址和信号量的地址

因为对 64 位整数的原子操作要求整数的地址是 64 位对齐的，所以针对 64 位和 32 位环境的 state 字段的组成是不一样的。

在 64 位环境下，state1 的第一个元素是 waiter 数，第二个元素是 WaitGroup 的计数值，第三个元素是信号量。

![64位state1值分布](/Users/bytedance/go/src/github.com/ouyangxizhu/studygo/concurrent/doc/waitgroup.assets/367c0ea5ead347acc6cf779554d9727c.png)

在 32 位环境下，如果 state1 不是 64 位对齐的地址，那么 state1 的第一个元素是信号量，后两个元素分别是 waiter 数和计数值。

![在这里插入图片描述](https://img-blog.csdnimg.cn/img_convert/334a529b22c44f4a5a77ebfffe9ecf48.png#pic_center)

接下里，我们一一看 Add 方法、 Done 方法、 Wait 方法的实现原理。

## Add()方法

**Add方法实现思路：**

Add方法主要操作的state1字段中计数值部分。当Add方法被调用时，首先会将delta参数值左移32位(计数值在高32位)，然后内部通过原子操作将这个值加到计数值上。需要注意的是，delta的取值范围可正可负，因为调用Done()方法时，内部通过Add(-1)方法实现的。

**代码实现如下**：

```go
// Add adds delta, which may be negative, to the WaitGroup counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
//
// Note that calls with a positive delta that occur when the counter is zero
// must happen before a Wait. Calls with a negative delta, or calls with a
// positive delta that start when the counter is greater than zero, may happen
// at any time.
// Typically this means the calls to Add should execute before the statement
// creating the goroutine or other event to be waited for.
// If a WaitGroup is reused to wait for several independent sets of events,
// new Add calls must happen after all previous Wait calls have returned.
// See the WaitGroup example.
func (wg *WaitGroup) Add(delta int) {
  // statep表示wait数和计数值
  // 低32位表示wait数，高32位表示计数值
	statep, semap := wg.state()
  //数据竞态检测，默认是false,开启消耗cpu性能 ,先不管
	if race.Enabled {
		_ = *statep // trigger nil deref early
		if delta < 0 {
			// Synchronize decrements with Wait.
			race.ReleaseMerge(unsafe.Pointer(wg))
		}
		race.Disable()
		defer race.Enable()
	}
  // uint64(delta)<<32 将delta左移32位
  // 因为高32位表示计数值，所以将delta左移32，增加到计数值上
	state := atomic.AddUint64(statep, uint64(delta)<<32)
  // 当前计数值
	v := int32(state >> 32)
  // 阻塞在检查点的wait数
	w := uint32(state)
  //数据竞态检测,先不管
	if race.Enabled && delta > 0 && v == int32(delta) {
		// The first increment must be synchronized with Wait.
		// Need to model this as a read, because there can be
		// several concurrent wg.counter transitions from 0.
		race.Read(unsafe.Pointer(semap))
	}
  // 计数器小于0 报错
	if v < 0 {
		panic("sync: negative WaitGroup counter")
	}
  	// waiter值不为0,累加后的counter值和delta相等,说明Wait()方法没有在Add()方法之后调用,触发panic,因为正确的做法是先Add()后Wait()

	if w != 0 && delta > 0 && v == int32(delta) {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
  // 如果等待为0或者计数器大于0 意味着没有等待或者还有读锁 不需要唤醒goroutine则返回 add操作完毕
  //Add()添加正常返回
	//1.counter > 0,说明还不需要释放信号量，可以直接返回
	//2. waiter  = 0 ,说明没有等待的goroutine，也不需要释放信号量，可以直接返回

	if v > 0 || w == 0 {
		return
	}
	// This goroutine has set counter to 0 when waiters > 0.
	// Now there can't be concurrent mutations of state:
	// - Adds must not happen concurrently with Wait,
	// - Wait does not increment waiters if it sees counter == 0.
	// Still do a cheap sanity check to detect WaitGroup misuse.
  //下面是 counter == 0 并且 waiter > 0的情况
	//现在若原state和新的state不等，则有以下两种可能
	//1. Add 和 Wait方法同时调用
	//2. counter已经为0，但waiter值有增加，这种情况永远不会触发信号量了
	// 以上两种情况都是错误的，所以触发异常
	//注：state := atomic.AddUint64(statep, uint64(delta)<<32)  这一步调用之后，state和*statep的值应该是相等的，除非有以上两种情况发生
  // 当等待计数器> 0时，而goroutine设置为0。
    // 此时不可能有同时发生的状态突变:
    // - 增加不能与等待同时发生，
    // - 如果计数器counter == 0，不再增加等待计数器
	if *statep != state {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
  // Reset waiters count to 0.
	//将waiter 和 counter都置为0
	*statep = 0
  // 如果计数值v为0并且waiter的数量w不为0，那么state的值就是waiter的数量
  // 将waiter的数量设置为0，因为计数值v也是0,所以它们俩的组合*statep直接设置为0即可。此时需要并唤醒所有的waiter
  // 唤醒所有等待的线程
  //原子递减信号量，并通知等待的goroutine
	for ; w != 0; w-- {
    // 目的是作为一个简单的wakeup原语，以供同步使用。true为唤醒排在等待队列的第一个goroutine
		runtime_Semrelease(semap, false, 0)
	}
}
```

## Done()方法

内部就是调用Add(-1)方法，这里就不细讲了。

```go
// Done decrements the WaitGroup counter by one.
// Done方法实际就是计数器减1
func (wg *WaitGroup) Done() {
	wg.Add(-1)
}
```

##  Wait()方法

**wait实现思路：**

不断检查state值。如果其中的计数值为零，则说明所有的子goroutine已全部执行完毕，调用者不必等待，直接返回。如果计数值大于零，说明此时还有任务没有完成，那么调用者变成等待者，需要加入wait队列，并且阻塞自己。

执行阻塞，直到所有的WaitGroup数量变成0。

**代码实现如下：**

```go
// Wait blocks until the WaitGroup counter is zero.
//调用Wait方法会阻塞当前调用的goroutine直到 counter的值为0
//也会增加waiter的值
func (wg *WaitGroup) Wait() {
	statep, semap := wg.state()
	if race.Enabled {
		_ = *statep // trigger nil deref early
		race.Disable()
	}
  //一直等待，直到无需等待或信号量触发，才返回
	for {
		state := atomic.LoadUint64(statep)
    // 将state右移32位，表示当前计数值
		v := int32(state >> 32)
    // w表示waiter等待值
		w := uint32(state)
    //如果counter值为0，则说明所有goroutine都退出了，无需等待，直接退出
		if v == 0 {
			// Counter is 0, no need to wait.
			if race.Enabled {
				race.Enable()
				race.Acquire(unsafe.Pointer(wg))
			}
      // 如果当前计数值为零，表示当前子goroutine已全部执行完毕，则直接返回
			return
		}
		// Increment waiters count.
    // 否则使用原子操作将state值加一。
    // 添加等待数量 如果cas失败则重新获取状态 避免计数有错
    //原子增加waiter的值，CAS方法，外面for循环会一直尝试，保证多个goroutine同时调用Wait()也能正常累加waiter
		if atomic.CompareAndSwapUint64(statep, state, state+1) {
			if race.Enabled && w == 0 {
				// Wait must be synchronized with the first Add.
				// Need to model this is as a write to race with the read in Add.
				// As a consequence, can do the write only for the first waiter,
				// otherwise concurrent Waits will race with each other.
				race.Write(unsafe.Pointer(semap))
			}
      // 阻塞休眠等待
      // 阻塞goroutine 等待唤醒
      //一直等待信号量sema，直到信号量触发，
      // 目的是作为一个简单的sleep原语，以供同步使用
			runtime_Semacquire(semap)
      //从上面的Add()方法看到，触发信号量之前会将seatep置为0(即counter和waiter都置为0)，所以此时应该也为0
			//如果不为0，说明WaitGroup此时又执行了Add()或者Wait()操作，所以会触发panic
			if *statep != 0 {
				panic("sync: WaitGroup is reused before previous Wait has returned")
			}
			if race.Enabled {
				race.Enable()
				race.Acquire(unsafe.Pointer(wg))
			}
      // 被唤醒，不再阻塞，返回
			return
		}
	}
}

```

# 注意点
1.Add()必须在Wait()前调用

2.Add()设置的值必须与实际等待的goroutine个数一致，如果设置的值大于实际的goroutine数量，可能会一直阻塞。如果小于会触发panic

3.WaitGroup不可拷贝，可以通过指针传递，否则很容易造成BUG

 

以下为值拷贝引起的Bug示例

demo1：因为值拷贝引起的死锁

```go
func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0 ; i < 5 ; i++ {
	 	test(wg)
	}
	wg.Wait()
}

func test(wg sync.WaitGroup) {
	go func() {
		fmt.Println("hello")
		wg.Done()
	}()
}
```




demo2:因为值拷贝引起的不会阻塞等待现象

```go
func main() {
	var wg sync.WaitGroup
	for i := 0 ; i < 5 ; i++ {
	 	test(wg)
	}
	wg.Wait()
}

func test(wg sync.WaitGroup) {
	go func() {
		wg.Add(1)
		fmt.Println("hello")
		time.Sleep(time.Second*5)
		wg.Done()
	}()
}
```




demo3:因为值拷贝引发的panic

```go
type person struct {
	wg sync.WaitGroup
}

func (t *person) say()  {
	go func() {
		fmt.Println("say Hello!")
		time.Sleep(time.Second*5)
		t.wg.Done()
	}()
}

func main() {
	var wg sync.WaitGroup
	t := person{wg:wg}
	wg.Add(5)
	for  i := 0 ; i< 5 ;i++ {
		t.say()
	}
	wg.Wait()
}

​																																																																																																			

# TODO

将race相关代码去掉分析一下

https://www.linkinstar.wiki/2020/03/15/golang/source-code/sync-waitgroup-source-code/

https://www.cnblogs.com/ricklz/p/14496612.html

https://xumc.github.io/blog/2019/11/13/waitgroup

https://www.jianshu.com/p/774587ddf25c?utm_campaign=maleskine&utm_content=note&utm_medium=seo_notes&utm_source=recommendation

