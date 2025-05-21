package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("程序开始运行")

	// 使用 time.After 函数在10秒后返回一个时间通道
	timeout := time.After(5 * time.Second)

	// 模拟程序的一些操作
	for {
		select {
		case <-timeout:
			// 当时间到了，从 timeout 通道接收到数据，程序退出
			fmt.Println("10秒时间已到，程序退出")
			return
		default:
			// 这里可以添加程序的具体逻辑
			// 为了避免CPU占用过高，让程序休眠一小段时间
			time.Sleep(100 * time.Millisecond)
			fmt.Println("111")
		}
	}
}
