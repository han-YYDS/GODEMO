package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

//-------------------------------------- READ DIR ----------------------------------------------

func TestReadDir() {
	// 读取当前目录下的所有文件和目录
	files, err := ioutil.ReadDir("..")
	if err != nil {
		log.Fatal(err)
	}

	// 遍历文件和目录
	for _, file := range files {
		// 使用IsDir()判断是否为目录=
		if file.IsDir() {
			fmt.Printf("Directory: %s\n", file.Name())
		} else {
			fmt.Printf("File: %s\n", file.Name())
		}
	}
}

//-------------------------------------- READ DIR ----------------------------------------------

func main() {
	TestReadDir()
}
