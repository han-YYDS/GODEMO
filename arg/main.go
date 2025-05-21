package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/erikdubbelboer/gspt"
)

func ParseConf() {
	var file string
	var demo string
	flag.StringVar(&file, "file", "", "")
	flag.StringVar(&demo, "demo", "", "")
	flag.Parse()
	fmt.Println("arg file: ", file)
	fmt.Println("arg demo: ", demo)
}

func main() {
	gspt.SetProcTitle("some title")

	ParseConf()

	for {
		time.Sleep(1 * time.Second)
	}
}
