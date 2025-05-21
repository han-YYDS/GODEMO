package main

import (
	"fmt"
	"log"
	"os"
)

var (
	filename = "demo"
)

func TestFileMode() {}

// 在文件不存在时进行创建
func TestCreate() {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func TestAppend() {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write([]byte("APPEND"))
}

func TestOpen() {

}

func TestRead() {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	fmt.Println(string(buffer[:n]))
}

func TestReadFile() {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}

func main() {
	TestCreate()
	TestAppend()
	// TestRead()
	TestReadFile()
}
