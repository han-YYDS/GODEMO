package main

import (
	"fmt"
	"time"
)

func namedReturn() (result int) {
	result = 1
	defer func() {
		time.Sleep(time.Second * 3)
		result++
	}()
	return result
}

func anonymousReturn() int {
	var result int = 1
	defer func() {
		time.Sleep(time.Second * 3)
		result++
	}()
	return result
}

func main() {
	fmt.Println(namedReturn())
	fmt.Println("111 end")
	fmt.Println(anonymousReturn())
	fmt.Println("222 end")
}
