package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

// 20250224: atomic.LoadInt32
// ----------------------------------------  TEST LOAD  -------------------------------------------------------------------------
// 通过在 go run 后面加 -race来检测是否有数据竞争
// load: 原子态的加载
func TestLoadInt32() {
	var value int32 = 200
	result := atomic.LoadInt32(&value)
	fmt.Printf("Current value: %d\n", result)
}

func TestLoadRace() {
	count := 0
	var wg sync.WaitGroup
	wg.Add(2)
	go incrementor(&count, &wg)
	go reader(&count, &wg)
	wg.Wait()

	fmt.Println("ret: ", count)
}

func incrementor(count *int, wg *sync.WaitGroup) {
	for i := 0; i < 1000000; i++ {
		(*count)++ // 非原子操作
	}
	wg.Done()
}

func reader(count *int, wg *sync.WaitGroup) {
	ret := 0
	for i := 0; i < 100000; i++ {
		ret += *count // 非原子读取
	}
	fmt.Println(ret)
	wg.Done()
}

func TestLoadAtomic() {
	var count int32 = 0
	var wg sync.WaitGroup
	wg.Add(2)
	go atomic_incrementor(&count, &wg)
	go atomic_reader(&count, &wg)
	wg.Wait()

	fmt.Println("ret: ", count)
}

func atomic_incrementor(count *int32, wg *sync.WaitGroup) {
	for i := 0; i < 1000; i++ {
		atomic.AddInt32(count, 1) // 原子操作
	}
	wg.Done()
}

func atomic_reader(count *int32, wg *sync.WaitGroup) {
	for i := 0; i < 1000; i++ {
		atomic.LoadInt32(count) // 原子读取
	}
	wg.Done()
}

//----------------------------------------  TEST LOAD  -------------------------------------------------------------------------

// 20250224: StoreInt32
// ----------------------------------------  TEST STORE  -------------------------------------------------------------------------
// 原子态的写数据
func atomicUpdater(sharedVar *int32, val int32, wg *sync.WaitGroup) {
	atomic.AddInt32(sharedVar, 1)
	// atomic.StoreInt32(sharedVar, val)
	wg.Done()
}

// 使用普通赋值操作更新变量的函数
func regularUpdater(sharedVar *int32, val int32, wg *sync.WaitGroup) {
	*sharedVar += 1
	wg.Done()
}

func TestStoreAtomic() {
	var sharedVar int32
	var wg sync.WaitGroup

	// 设置 CPU 核心数，以最大化并发效果
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 用1000个协程来更新一个数
	// 使用 atomic.Store32 更新变量
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go atomicUpdater(&sharedVar, int32(i), &wg)
	}

	// 等待所有 atomic 更新完成
	wg.Wait()
	fmt.Printf("After atomic updates: %d\n", sharedVar)

	// 重置共享变量
	sharedVar = 0

	// 使用普通赋值操作更新变量
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go regularUpdater(&sharedVar, int32(i), &wg)
	}

	// 等待所有普通更新完成
	wg.Wait()
	fmt.Printf("After regular updates: %d\n", sharedVar)

	// 注意：由于数据竞争，regularUpdater 的结果可能不是最后一个赋值操作的值
}

// ----------------------------------------  TEST STORE  -------------------------------------------------------------------------

func main() {
	// TestLoadInt32()
	// TestLoadRace()
	// TestLoadAtomic()
	TestStoreAtomic()
}
