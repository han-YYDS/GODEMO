package main

import "fmt"

// init 函数在包被导入时,自行执行
func init() {
	fmt.Println("mypackage 初始化")
}

func main() {
	fmt.Println("main ")
}
