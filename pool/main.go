package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 用来统计实例真正创建的次数
var numCalcsCreated int32

// 创建实例的函数
func createBuffer() interface{} {
	// 这里要注意下，非常重要的一点。这里必须使用原子加，不然有并发问题；
	atomic.AddInt32(&numCalcsCreated, 1)
	buffer := make([]byte, 1024)
	return &buffer
}

// Pool 池里的元素随时可能释放掉，释放策略完全由 runtime 内部管理；
// Get 获取到的元素对象可能是刚创建的，也可能是之前创建好 cache 住的。使用者无法区分；
// Pool 池里面的元素个数你无法知道；
// sync.Pool 本质用途是增加临时对象的重用率，减少 GC 负担

func main() {
	// 创建实例
	bufferPool := &sync.Pool{
		New: createBuffer,
	}

	// 多 goroutine 并发测试
	numWorkers := 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			// 申请一个 buffer 实例
			// buffer := bufferPool.Get()  // 如果采用Pool.Get, 则会减少新对象的构造次数

			buffer := createBuffer()
			_ = buffer.(*[]byte)
			// 释放一个 buffer 实例, 放回pull中
			defer bufferPool.Put(buffer)
		}()
	}
	wg.Wait()
	fmt.Printf("%d buffer objects were created.\n", numCalcsCreated)
}
