package main

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"golang.org/x/time/rate"
) // 需要import的rate库，其它import暂时忽略

// 漏桶: 通过限制consume_max速率来实现限流
// 令牌桶: 相当于存在一个资源堆积的逻辑,比如在系统启动时, 其资源充沛, 所以此时消费速率应当高于阈值

// 生成0->X的数据集
func generateData(num int) []any {
	var data []any
	for i := 0; i < num; i++ {
		data = append(data, i)
	}
	return data
}

// 处理数据，数字*10
func process(obj any) (any, error) {
	integer, ok := obj.(int)
	if !ok {
		return nil, errors.New("invalid integer")
	}
	time.Sleep(1)
	nextInteger := integer * 10
	if integer%99 == 0 {
		return nextInteger, errors.New("not a happy number")
	}
	return nextInteger, nil
}

func TestRate(t *testing.T) {
	limit := rate.Limit(100) // QPS：50
	burst := 25              // 桶容量25
	limiter := rate.NewLimiter(limit, burst)
	size := 500 // 数据量500

	data := generateData(size)
	var wg sync.WaitGroup
	startTime := time.Now()

	// 模拟500条请求
	for i, item := range data {
		wg.Add(1)
		go func(idx int, obj any) {
			defer wg.Done()
			// 拿到令牌
			if err := limiter.Wait(context.Background()); err != nil {
				t.Logf("[%d] [EXCEPTION] wait err: %v", idx, err)
			}
			// 执行业务逻辑
			processed, err := process(obj)
			if err != nil {
				t.Logf("[%d] [ERROR] processed: %v, err: %v", idx, processed, err)
			} else {
				t.Logf("[%d] [OK] processed: %v", idx, processed)
			}
		}(i, item)
	}
	wg.Wait()
	endTime := time.Now()

	// qps为50, 处理500条数据, 应该在10s左右完成
	t.Logf("start: %v, end: %v, seconds: %v", startTime, endTime, endTime.Sub(startTime).Seconds())
}
