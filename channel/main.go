package main

import (
	"fmt"
	"time"
)

// 20250226:
// 函数参数中的只读channel和只写channel
// func ReadFullMetaAsync(fname string, msg chan<- *EventData) {}

// ---------------------------------------------     CHANNEL BUFFER  ------------------------------------------------------------------------------------------------
func TestChannelSize() { // 有缓冲的channel
	ch := make(chan int, 2)

	go func() {
		// time.Sleep(time.Second)
		for val := range ch { // channel ->
			fmt.Println("output: ", val)
			time.Sleep(2 * time.Second)
		}
	}()

	// send data

	fmt.Println("Sending 1...")
	ch <- 1

	fmt.Println("Sending 2...")
	ch <- 2

	fmt.Println("Sending 3...")
	ch <- 3 // 会在这里阻塞一秒, 等待消费端消费

	fmt.Println("Sending 4...") //
	ch <- 4

	fmt.Println("Sending 5...")
	ch <- 5

	time.Sleep(5 * time.Second)
}

// ---------------------------------------------     CHANNEL BUFFER  ------------------------------------------------------------------------------------------------

// ---------------------------------------------     CHANNEL NIL  ------------------------------------------------------------------------------------------------
func TestNil() {
	// chan *int
	ch := make(chan int, 1)

	ch <- 1
	// would error
	// ch <- nil
}

func TestNilptr() {
	// chan *int
	ch := make(chan *int, 1)
	// 对于 ptr channel 是能够在其中传输nil的
	ch <- nil

	select {
	case <-ch:
		fmt.Println("channel")
	default:
		fmt.Println("default")
	}
}

// ---------------------------------------------     CHANNEL NIL  ------------------------------------------------------------------------------------------------

func main() {
	// TestNilptr()
	TestChannelSize()
}
