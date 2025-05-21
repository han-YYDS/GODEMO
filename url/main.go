package main

import (
	"fmt"
	"net/url"
)

func TestParseAbsolute() {
	rawURL := "http://user:pass@example.com:8080/path?query=1#fragment"
	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Scheme:", u.Scheme)     // "http" 协议
	fmt.Println("User:", u.User)         // "user:pass" user信息
	fmt.Println("Host:", u.Host)         // "example.com:8080" socket
	fmt.Println("Path:", u.Path)         // "/path" 域名之后的urlpath
	fmt.Println("RawQuery:", u.RawQuery) // "query=1"
	fmt.Println("Fragment:", u.Fragment) // "fragment"
}

func TestParseRelative() {
	rawURL := "/path/to/resource?query=1#fragment"
	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Scheme:", u.Scheme)     // ""
	fmt.Println("Host:", u.Host)         // ""
	fmt.Println("Path:", u.Path)         // "/path/to/resource"
	fmt.Println("RawQuery:", u.RawQuery) // "query=1"
	fmt.Println("Fragment:", u.Fragment) // "fragment"
}

func main() {
	// TestParseAbsolute()
	TestParseRelative()
}
