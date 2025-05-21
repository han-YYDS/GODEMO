package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Fname  string
	Offset int64
	Head   string
	Type   string
}

func showTool() {
	fmt.Printf("\t%s <command> arguments\n\n", filepath.Base(os.Args[0]))
	fmt.Println("The commands are:")
	fmt.Println("\t-f, --file=name\t\t filename.")
	fmt.Println("\t-o, --offset=offset\t\t filename.")
	fmt.Println("\t-h, --head=bytes\t\t head bytes.")
	fmt.Println("\t-t, --type=[meta、ack]#\t meta or ack file")
	fmt.Println("")
}

// 定义命令行参数
func ParseConf() *Config {
	conf := &Config{}
	flag.StringVar(&conf.Fname, "f", "", "")
	flag.Int64Var(&conf.Offset, "o", -1, "")
	flag.StringVar(&conf.Head, "h", "", "")
	flag.StringVar(&conf.Type, "t", "meta", "")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "\nUsage of %s:\n", filepath.Base(os.Args[0]))
		showTool()
		os.Exit(0)
	}

	flag.Parse()
	fmt.Println("fname: ", conf.Fname)
	return conf
}

func main() {
	tool := ParseConf()

	if len(tool.Head) > 0 {
		// tool.DumpHead()
	} else if len(tool.Fname) > 0 {
		// tool.ParseFile()
	} else {
		fmt.Fprintf(flag.CommandLine.Output(), "\nUsage of %s:\n", filepath.Base(os.Args[0]))
		showTool()
		os.Exit(0)
	}
}
