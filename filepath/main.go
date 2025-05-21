package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// TODO: 20250226: filepath.Walk遍历目录
// -------------------------------------------------- FILEPATH WALK ------------------------------------------------------------------------------

func collectMetaFiles(dir string) ([]string, error) {
	var metaFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 检查文件是否以 .meta 结尾
		if !info.IsDir() && filepath.Ext(path) == ".db" {
			metaFiles = append(metaFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return metaFiles, nil
}

func TestWalk() {
	// 假设我们要遍历的目录是 "exampleDir"
	dir := "/home/ubuntu20/esmd/node0"
	// dir := "/home/ubuntu20/esmd/node0/data/db/t/meta"
	metaFiles, err := collectMetaFiles(dir)
	if err != nil {
		fmt.Printf("Error collecting .meta files: %v\n", err)
		return
	}

	for _, file := range metaFiles {
		fmt.Println(file)
	}
}

// -------------------------------------------------- FILEPATH WALK ------------------------------------------------------------------------------

func main() {
	TestWalk()
}
