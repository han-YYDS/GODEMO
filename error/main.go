package main

import (
	"errors"
	"fmt"
)
if err != nil {
	t.Fatal(err)
}
  
if err.Error() != "grammar rules not match" {
	t.Fatal(err)
}



// 判断error类型
func TestErrorEQ() {
	err1 := fmt.Errorf("ERROR")
	err2 := fmt.Errorf("ERROR")
	fmt.Println(err1.Error() == err2.Error()) // 判断报错信息是否相同
	fmt.Println(errors.Is(err1, err2))        // 判断底层是否是同一实例
	fmt.Println(errors.Is(err1, err1))        // 判断底层是否是同一实例
}

var errDemo = errors.New("my err")

func TestErrorIs() {
	err := fmt.Errorf("hello %w", errDemo)

	fmt.Printf("myErr:%s , err:%s \n", errDemo, err)

	fmt.Println("使用 == 的结果:", err == errDemo)
	fmt.Println("使用 errors.Is(err, myErr) 的结果:", errors.Is(err, errDemo))
	fmt.Println("使用 errors.Is(myErr, err) 的结果:", errors.Is(errDemo, err))
}

func main() {
	TestErrorEQ()
	TestErrorIs()
}
