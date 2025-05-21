package main

import (
	"fmt"
	"time"
)

func main() {
	total := 100 // 总进度数

	for i := 0; i <= total; i++ {
		// 打印进度条
		printProgressBar(i, total)
		time.Sleep(50 * time.Millisecond) // 模拟工作
	}

	fmt.Println("\n完成!") // 换行，结束输出
}

// printProgressBar 打印动态进度条
func printProgressBar(current, total int) {
	width := 50 // 进度条的宽度
	progress := int(float64(current) / float64(total) * float64(width))

	// 通过使用 \r 来回到行的开头
	bar := fmt.Sprintf(
		"\r[%s%s] %d%%",
		string(repeat('#', progress)),            // 已完成部分
		string(repeat(' ', width-progress)),      // 未完成部分
		int(float64(current)/float64(total)*100), // 百分比
	)
	fmt.Print(bar)
}

// repeat 生成指定字符的重复序列
func repeat(char rune, count int) []rune {
	result := make([]rune, count)
	for i := range result {
		result[i] = char
	}
	return result
}
