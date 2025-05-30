package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("demo.txt", os.O_CREATE|os.O_APPEND, 0644)
	fmt.Println("openfile: ", err)
	defer file.Close()
	log.SetOutput(file)

	log.Println("demo demo demo")
	file.Sync()
}
