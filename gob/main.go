package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type User struct {
	Username string
	Password string
}

// encode
func test1() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer) // 创建编码器

	user := User{Username: "gopher", Password: "golang"}
	err := encoder.Encode(user) // 编码
	if err != nil {
		fmt.Println("Error encoding:", err)
		return
	}

	fmt.Println("Serialized data:", buffer.Bytes())
}

// decode
func test2() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	user := User{Username: "gopher", Password: "golang"}
	err := encoder.Encode(user)
	if err != nil {
		fmt.Println("Error encoding:", err)
		return
	}

	var decodedUser User
	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&decodedUser)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}

	fmt.Println("Decoded user:", decodedUser)
}

type Person struct {
	Name   string
	Age    int
	Action interface{}
}

type Run struct {
	Speed int
}

func test3() {
	var dao bytes.Buffer
	encoder := gob.NewEncoder(&dao)
	decoder := gob.NewDecoder(&dao)
	p := Person{Name: "chen", Age: 18, Action: Run{80}}
	err := encoder.Encode(&p)
	if err != nil {
		panic(err)
	}
	fmt.Println(dao.String())
	var d Person
	err = decoder.Decode(&d)
	if err != nil {
		panic(err)
	}
	fmt.Println(d)
}

func main() {
	test2()
}
