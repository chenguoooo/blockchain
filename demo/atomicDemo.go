package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	// counter是所有goroutine都要增加的变量
	counter int64
	// wg用来等待程序的结束
	wg sync.WaitGroup
	// mutex 用来定义一段代码临界区
	mutex sync.Mutex
)

func main() {
	// 计数加2，表示要等待两个goroutine
	wg.Add(2)
	// 创建两个goroutine
	go incCounter(1)
	go incCounter(2)
	// 等待goroutine结束
	wg.Wait()
	fmt.Println("Final Counter:", counter)
}

// incCounter增加包里counter变量的值
func incCounter(id int) {
	// 延时调用，在函数退出时调用Done来通知main函数工作已经完成
	defer wg.Done()

	for count := 0; count < 2; count++ {
		// 捕获counter的值
		value := counter
		// 当前goroutine从线程退出，并放回到队列
		runtime.Gosched()
		// 增加本地value变量的值
		value++
		// 将该值保存回counter
		counter = value
	}
}

// incCounter增加包里counter变量的值
func incCounter1(id int) {
	defer wg.Done()

	for count := 0; count < 2; count++ {
		// 安全的对counter加1
		atomic.AddInt64(&counter, 1)
		// 当前goroutine从线程退出，并放回到队列
		runtime.Gosched()
	}
}

// incCounter增加包里counter变量的值
func incCounter2(id int) {
	// 延时调用，在函数退出时调用Done来通知main函数工作已经完成
	defer wg.Done()

	for count := 0; count < 2; count++ {
		// 同一时刻只允许一个goroutine进入这个临界区
		mutex.Lock()
		{
			// 捕获counter的值
			value := counter
			// 当前goroutine从线程退出，并放回到队列
			runtime.Gosched()
			// 增加本地value变量的值
			value++
			// 将该值保存回counter
			counter = value
		}
		// 释放锁，允许其他正在等待的goroutine进入临界区
		mutex.Unlock()
	}
}
