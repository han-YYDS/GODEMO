package main

import (
	"fmt"
	"reflect"
)

func test1() {
	var v interface{}
	v = (*int)(nil)       // 内容为nil的int指针
	fmt.Println(v == nil) // false
}

// 刚刚声明出来的 data 和 in 变量，确实是输出结果是 nil，判断结果也是 true。
// 怎么把变量 data 一赋予给变量 in，输出结果依然是 nil，但判定却变成了 false。
func test2() {
	var data *byte
	var in interface{}

	fmt.Println(data, data == nil) // 定义但未赋值的指针 - true
	fmt.Println(in, in == nil)     // 定义未赋值的接口 - true

	in = data
	fmt.Println(in, in == nil) // false, 其值为nil,但其type并不是nil,所以不会判断为nil(interface就是nil type)
}

func test3() {
	var data *byte
	var in interface{}

	in = data
	fmt.Println(IsNil(in))
}

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		// reflect中会对interface做特殊处理
		return vi.IsNil()
	}
	return false
}

func main() {
	// test1()
	// test2()
	test3()
}
