package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"sync"
)

func main() {
	// 创建一个通道，用于接收中断信号
	sigChan := make(chan os.Signal, 1)
	// 监听中断信号（SIGINT）和终止信号（SIGTERM）
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 创建一个通道，用于示例中的阻塞操作
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)

	// 在单独的goroutine中启动阻塞操作
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ch:
				// 处理接收到的数据
				fmt.Println("Received data")
			case <-sigChan:
				// 收到中断信号，可以在这里进行清理工作
				fmt.Println("Interrupt received, shutting down...")
				return
			}
		}
	}()

	// 发送数据到通道，模拟阻塞操作
	ch <- 1


	wg.Wait()
}

