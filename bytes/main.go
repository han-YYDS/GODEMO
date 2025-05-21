package main

import (
	"bytes"
	"fmt"
)

func TestEqual() {
	a := []byte("Hello, World!")
	b := []byte("Hello, World!")
	c := []byte("Goodbye, World!")

	// 使用 bytes.Equal 判断
	fmt.Println(bytes.Equal(a, b)) // 输出: true
	fmt.Println(bytes.Equal(a, c)) // 输出: false
}

func main() {
	TestEqual()
}
