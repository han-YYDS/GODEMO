package main

import (
	"fmt"
	"reflect"
)

type Value struct {
	Name   string
	Gender string
}

type Value1 struct {
	Name   string
	Gender *string
}

type Value2 struct {
	Name   string
	GoodAt []string
}

func test() {
	// 结构体内容可比较且相等
	v1 := Value{Name: "煎鱼", Gender: "男"}
	v2 := Value{Name: "煎鱼", Gender: "男"}
	if v1 == v2 {
		fmt.Println("脑子进煎鱼了")
	} else {
		fmt.Println("脑子没进煎鱼")
	}

	// 直接判断指针的字面值,其不相等
	v3 := Value1{Name: "煎鱼", Gender: new(string)}
	v4 := Value1{Name: "煎鱼", Gender: new(string)}
	if v3 == v4 {
		fmt.Println("脑子进煎鱼了")
		return
	} else {
		fmt.Println("脑子没进煎鱼")

	}

}

// func test2() {
// 	// 报错: truct containing []string cannot be compared
//  // 具有slice成员的结构体不能直接比较
// 	v5 := Value2{Name: "煎鱼", GoodAt: []string{"炸", "煎", "蒸"}}
// 	v6 := Value2{Name: "煎鱼", GoodAt: []string{"炸", "煎", "蒸"}}
// 	if v5 == v6 {
// 		fmt.Println("脑子进煎鱼了")
// 		return
// 	} else {
// 		fmt.Println("脑子没进煎鱼")
// 	}
// }

type ValueA struct {
	Name string
}

type ValueB struct {
	Name string
}

// func test3() {
// 	v1 := Value1{Name: "煎鱼"}
// 	v2 := Value2{Name: "煎鱼"}
// 	// 报错 : mismatched types Value1 and Value2
// 	// 不同类型的struct不能直接比较
// 	if v1 == v2 {
// 		fmt.Println("脑子进煎鱼了")
// 		return
// 	}

// 	fmt.Println("脑子没进煎鱼")
// }

func testDeepEqual() {

	v3 := Value1{Name: "煎鱼", Gender: new(string)}
	v4 := Value1{Name: "煎鱼", Gender: new(string)}
	if reflect.DeepEqual(v3, v4) {
		fmt.Println("111 脑子进煎鱼了")
	} else {
		fmt.Println("222 脑子没进煎鱼")

	}

	// slice成员
	v5 := Value2{Name: "煎鱼", GoodAt: []string{"炸", "煎", "蒸"}}
	v6 := Value2{Name: "煎鱼", GoodAt: []string{"炸", "煎", "蒸"}}
	if reflect.DeepEqual(v5, v6) {
		fmt.Println("111 脑子进煎鱼了")
	} else {
		fmt.Println("222 脑子没进煎鱼")
	}

	// 不同struct
	v1 := Value1{Name: "煎鱼"}
	v2 := Value2{Name: "煎鱼"}
	if reflect.DeepEqual(v1, v2) {
		fmt.Println("111 脑子进煎鱼了")
	}

	fmt.Println("222 脑子没进煎鱼")
}

func main() {
	testDeepEqual()
}
