package main

import (
	"fmt"
	"time"

	"go.uber.org/ratelimit"
)

// uber的库里面没看到bucket设置?
func main() {
	rl := ratelimit.New(100) // 1秒生成100个令牌
	time.Sleep(2 * time.Second)
	prev := time.Now()

	for i := 0; i < 10; i++ {
		now := rl.Take()              // 获取一个令牌
		fmt.Println(i, now.Sub(prev)) // 计算与上一个令牌的间隔
		prev = now
	}
}
