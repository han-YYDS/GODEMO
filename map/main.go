package main

import (
	"fmt"
	"sync"
)

// 20240224: LoadOrStore
// --------------------------------------------  TEST LoadOrStore  ----------------------------------------------------------------------------------------------------------------
// 1. 该函数属于 sync.Map, 原生map并不是线程安全的

func TestLoadOrStore() {

}

func TestSyncMap() {
	var m sync.Map
	// 1. 写入
	m.Store("qcrao", 18)
	m.Store("stefno", 20)

	// 2. 读取
	age, _ := m.Load("qcrao")
	fmt.Println(age.(int))

	// 3. 遍历
	m.Range(func(key, value interface{}) bool {
		name := key.(string)
		age := value.(int)
		fmt.Println(name, age)
		return true
	})

	// 4. 删除
	m.Delete("qcrao")
	age, ok := m.Load("qcrao")
	fmt.Println(age, ok)
 
	// 5. 读取或写入, 由于已存在, 所以store不会执行
	age, exist := m.LoadOrStore("stefno", 100)
	fmt.Println(age, exist)
}

//--------------------------------------------  TEST LoadOrStore  ----------------------------------------------------------------------------------------------------------------

// 20240224: sync.Map和map[int]string 的区别
// 20240224: map.Range中的return true和return false

func main() {
	TestSyncMap()
}
