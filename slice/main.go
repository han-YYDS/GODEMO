package main

import "fmt"

func TestCopy() {
	arr1 := []int{1, 2, 3, 4, 5}
	arr2 := arr1[0:3]
	fmt.Println("arr1:", arr1)
	fmt.Println("arr2:", arr2)

	arr1[0] = 3
	fmt.Println("arr1:", arr1)
	fmt.Println("arr2:", arr2)
}

func TestDeepCopy() {
	arr1 := []int{1, 2, 3, 4, 5}
	// copy得预分配空间
	arr2 := make([]int, 5)
	copy(arr2, arr1[0:3])

	fmt.Println("arr1:", arr1)
	fmt.Println("arr2:", arr2)

	arr1[0] = 3
	fmt.Println("arr1:", arr1)
	fmt.Println("arr2:", arr2)
}

func TestAppendCopy() {
	arr1 := []int{1, 2, 3, 4, 5}
	var arr2 []int
	arr2 = append(arr2, arr1[0:3]...)

	fmt.Println("arr1:", arr1)
	fmt.Println("arr2:", arr2)

	arr1[0] = 3
	fmt.Println("arr1:", arr1)
	fmt.Println("arr2:", arr2)
}

// 20250306: 对于浅拷贝之后, 对slice2进行修改会如何影响到slice1?\
// 1. 原来append没有新开辟一个地址, 而是在append的第一个参数上向后扩容并返回
func TestCopy2() {
	slice1 := []int{1, 2, 3}
	slice2 := slice1[1:2]
	fmt.Println("slice1:", slice1)
	fmt.Println("slice2:", slice2)

	// 对于浅拷贝之后的结果, append会被拷贝源产生影响
	slice2 = append(slice2, 4)
	fmt.Println("\n -----------")
	fmt.Println("slice1:", slice1)
	fmt.Println("slice2:", slice2)

	slice2 = append(slice2, 5)
	fmt.Println("\n -----------")
	fmt.Println("slice1:", slice1)
	fmt.Println("slice2:", slice2)

	// 然后在append一定次数之后, slice1没有再受到影响, 原因是发生了扩容, 导致slice1指向了新的空间
	slice2 = append(slice2, 6)
	fmt.Println("\n -----------")
	fmt.Println("slice1:", slice1)
	fmt.Println("slice2:", slice2)
}

func main() {
	// TestDeepCopy()

	// TestAppendCopy()

	TestCopy2()
}
