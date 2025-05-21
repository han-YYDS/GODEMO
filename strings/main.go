package main

import (
	"fmt"
	"strings"
)

func TestCut() {
	str := "abc#"
	before, after, found := strings.Cut(str, "#")
	fmt.Println("before: ", before)
	fmt.Println("after: ", after)
	fmt.Println("found: ", found)
}

func main() {
	TestCut()
}
