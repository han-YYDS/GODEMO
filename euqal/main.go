package main

import (
	"fmt"
)

type MyStruct struct {
    Field1 int
    Field2 string
    // 其他字段...
}

func main() {
    a := MyStruct{Field1: 1, Field2: "test"}
    b := MyStruct{Field1: 1, Field2: "test"}

    if a == b {
        fmt.Println("结构体相等")
    } else {
        fmt.Println("结构体不相等")
    }
}

