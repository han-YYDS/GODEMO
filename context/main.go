package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("timeout ")
			return
		default:
			fmt.Println("default")
			time.Sleep(500 * time.Millisecond)
		}
	}

}

func TestTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer func() {
		fmt.Println("Defer")
		cancel()
	}()
	go worker(ctx, "111")
	go worker(ctx, "222")
	time.Sleep(3 * time.Second)
	fmt.Println("Main ")
}

func main() {
	TestTimeout()
}
