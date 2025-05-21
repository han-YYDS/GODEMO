package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	prompt "github.com/c-bata/go-prompt"
)

func executor(t string) {
	if t == "q" {
		os.Exit(1)
	}
	fmt.Println("echo: ", t)
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "new topic", Description: "Create a new topic"},
		{Text: "new group", Description: "Create a new group"},
	}
}

func main() {
	// 创建一个信号通道
	sigs := make(chan os.Signal, 1)
	// 监听os.Interrupt信号（ctrl + c）和syscall.SIGTERM信号
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)

	// 创建一个上下文，用于取消操作
	ctx, cancel := context.WithCancel(context.Background())

	// 启动命令行界面
	go func() {
		p := prompt.New(
			executor,
			completer,
		)
		p.Run()
		cancel()
	}()

	// 监听信号并取消上下文
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		cancel()
	}()

	// 等待上下文取消
	<-ctx.Done()
}
