package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const defaultMaxFileSize = 1 << 30        // 假设文件最大为 1G
const defaultMemMapSize = 128 * (1 << 20) // 假设映射的内存大小为 128M
type Demo struct {
	file    *os.File
	data    *[defaultMaxFileSize]byte
	dataRef []byte
}

func _assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf(msg, v...))
	}
}

func (demo *Demo) mmap() {
	b, err := syscall.Mmap(int(demo.file.Fd()), 0, defaultMemMapSize, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	_assert(err == nil, "failed to mmap", err)
	demo.dataRef = b
	demo.data = (*[defaultMaxFileSize]byte)(unsafe.Pointer(&b[0]))
}
func (demo *Demo) grow(size int64) {
	if info, _ := demo.file.Stat(); info.Size() >= size {
		return
	}
	_assert(demo.file.Truncate(size) == nil, "failed to truncate")
}
func (demo *Demo) munmap() {
	_assert(syscall.Munmap(demo.dataRef) == nil, "failed to munmap")
	demo.data = nil
	demo.dataRef = nil
}

func main() {
	_ = os.Remove("tmp.txt")
	f, _ := os.OpenFile("tmp.txt", os.O_CREATE|os.O_RDWR, 0644)
	demo := &Demo{file: f}
	demo.grow(1)
	demo.mmap()
	defer demo.munmap()
	msg := "hello geektutu!"
	demo.grow(int64(len(msg) * 2))
	for i, v := range msg {
		demo.data[2*i] = byte(v)
		demo.data[2*i+1] = byte(' ')
	}
}
