package main

import (
	"fmt"
	"math"
)

func TestRound() {
	a := 64
	b := 17
	c := 15
	// 四舍五入
	avg1 := int(math.Round(float64(a) / float64(b)))
	avg2 := int(math.Round(float64(a) / float64(c)))
	fmt.Println(avg1, " ", avg2)
}

func main() {
	TestRound()
}
