package main

import (
	"fmt"
	"reflect"
)

type user struct {
	name string `export:"namee"`
	age  int    `export:"agee"`
	// 可以有多个标签
	demo string `tag1:"demo" tag2:"demo"`
}

func main() {
	u := user{"tom", 11, "d"}

	t := reflect.TypeOf(u)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		exporttag := field.Tag.Get("export")
		fmt.Println("Field: ", field.Name, exporttag)
	}
}
