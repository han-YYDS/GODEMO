package main

import (
	"fmt"

	"github.com/han-YYDS/mymod"
)

func main() {
	sum := mymod.Add(5, 3)
	diff := mymod.Subtract(5, 3)
	fmt.Printf("5 + 3 = %d\n", sum)
	fmt.Printf("5 - 3 = %d\n", diff)

	product := mymod.Multiply(5, 3)
	fmt.Printf("5 * 3 = %d\n", product)
}
